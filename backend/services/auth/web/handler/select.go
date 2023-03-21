package handler

import (
	"context"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-pipedrive/pkg/log"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/auth/web/core/domain"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/auth/web/core/port"
	pclient "github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/client"
	"go-micro.dev/v4/client"
	"golang.org/x/sync/singleflight"
)

var group singleflight.Group

type UserSelectHandler struct {
	service       port.UserAccessService
	client        client.Client
	pipedriveAuth pclient.PipedriveAuthClient
	logger        log.Logger
}

func NewUserSelectHandler(
	service port.UserAccessService,
	client client.Client,
	pipedriveAuth pclient.PipedriveAuthClient,
	logger log.Logger,
) UserSelectHandler {
	return UserSelectHandler{
		service:       service,
		client:        client,
		pipedriveAuth: pipedriveAuth,
		logger:        logger,
	}
}

func (u UserSelectHandler) GetUser(ctx context.Context, uid *string, res *domain.UserAccess) error {
	user, err, _ := group.Do(*uid, func() (interface{}, error) {
		user, err := u.service.GetUser(ctx, *uid)
		if err != nil {
			u.logger.Errorf("could not get user with id: %s. Reason: %s", *uid, err.Error())
			return nil, err
		}

		if user.ExpiresAt <= time.Now().UnixMilli() {
			u.logger.Debug("user token has expired. Trying to refresh!")
			token, terr := u.pipedriveAuth.RefreshAccessToken(ctx, user.RefreshToken)
			if terr != nil {
				u.logger.Errorf("could not refresh user's %s token. Reason: %s", *uid, terr.Error())
				return nil, terr
			}

			u.logger.Debugf("user's %s token has been refreshed", *uid)
			access := domain.UserAccess{
				ID:           user.ID,
				AccessToken:  token.AccessToken,
				RefreshToken: token.RefreshToken,
				TokenType:    token.TokenType,
				Scope:        token.Scope,
				ApiDomain:    token.ApiDomain,
				ExpiresAt:    time.Now().Local().Add(time.Second * time.Duration(token.ExpiresIn-700)).UnixMilli(),
			}

			_, err := u.service.UpdateUser(ctx, access)
			if err != nil {
				u.logger.Debugf("could not persist a new user's %s token. Reason: %s. Sending a fallback message!", *uid, err.Error())
				return nil, err
			}

			u.logger.Debugf("user's %s token has been updated", *uid)
			return access, nil
		}

		return user, nil
	})

	if usr, ok := user.(domain.UserAccess); ok {
		*res = usr
		return nil
	}

	return err
}
