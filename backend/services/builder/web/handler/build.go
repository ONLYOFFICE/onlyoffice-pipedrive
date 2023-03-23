package handler

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	plog "github.com/ONLYOFFICE/onlyoffice-pipedrive/pkg/log"
	pclient "github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/client"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/client/model"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/constants"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/crypto"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/request"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/response"
	"github.com/mileusna/useragent"
	"go-micro.dev/v4/client"
	"golang.org/x/sync/singleflight"
)

var _ErrNoSettingsFound = errors.New("could not find document server settings")
var _ErrOperationTimeout = errors.New("operation timeout")

type ConfigHandler struct {
	namespace  string
	logger     plog.Logger
	client     client.Client
	apiClient  pclient.PipedriveApiClient
	jwtManager crypto.JwtManager
	gatewayURL string
	group      singleflight.Group
}

func NewConfigHandler(
	namespace string,
	logger plog.Logger,
	client client.Client,
	jwtManager crypto.JwtManager,
	gatewayURL string,
) ConfigHandler {
	return ConfigHandler{
		namespace:  namespace,
		logger:     logger,
		client:     client,
		apiClient:  pclient.NewPipedriveApiClient(),
		jwtManager: jwtManager,
		gatewayURL: gatewayURL,
	}
}

func (c ConfigHandler) processConfig(user response.UserResponse, req request.BuildConfigRequest, ctx context.Context) (response.BuildConfigResponse, error) {
	var config response.BuildConfigResponse
	var wg sync.WaitGroup
	usrChan := make(chan model.User, 1)
	settingsChan := make(chan response.DocSettingsResponse, 1)
	errorsChan := make(chan error, 2)

	go func() {
		wg.Add(1)
		defer wg.Done()
		u, err := c.apiClient.GetMe(ctx, model.Token{
			AccessToken:  user.AccessToken,
			RefreshToken: user.RefreshToken,
			TokenType:    user.TokenType,
			Scope:        user.Scope,
			ApiDomain:    user.ApiDomain,
		})

		if err != nil {
			c.logger.Debugf("could not get pipedrive user: %s", err.Error())
			errorsChan <- err
			return
		}

		c.logger.Debugf("populating pipedrive user %d channel", u.ID)
		usrChan <- u
		c.logger.Debugf("successfully populated pipedrive channel")
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		var docs response.DocSettingsResponse
		if err := c.client.Call(ctx, c.client.NewRequest(fmt.Sprintf("%s:settings", c.namespace), "SettingsSelectHandler.GetSettings", fmt.Sprint(req.CID)), &docs); err != nil {
			c.logger.Debugf("could not document server settings: %s", err.Error())
			errorsChan <- err
			return
		}

		if docs.DocAddress == "" || docs.DocSecret == "" {
			c.logger.Debugf("no settings found")
			errorsChan <- _ErrNoSettingsFound
			return
		}

		c.logger.Debugf("populating document server %d settings channel", req.CID)
		settingsChan <- docs
		c.logger.Debugf("successfully populated document server settings channel")
	}()

	c.logger.Debugf("waiting for goroutines to finish execution")
	wg.Wait()
	c.logger.Debugf("goroutines have finished the execution")

	select {
	case err := <-errorsChan:
		return config, err
	case <-ctx.Done():
		return config, _ErrOperationTimeout
	default:
		c.logger.Debugf("select default")
	}

	usr := <-usrChan
	settings := <-settingsChan
	t := "desktop"
	ua := useragent.Parse(req.UserAgent)

	if ua.Mobile || ua.Tablet {
		t = "mobile"
	}

	downloadToken := request.PipedriveTokenContext{
		UID: usr.ID,
		CID: usr.CompanyID,
	}

	downloadToken.IssuedAt = 0
	downloadToken.ExpiresAt = time.Now().Add(4 * time.Minute).UnixMilli()
	tkn, _ := c.jwtManager.Sign(settings.DocSecret, downloadToken)

	config = response.BuildConfigResponse{
		Document: response.Document{
			Key:   req.DocKey,
			Title: req.Filename,
			URL:   fmt.Sprintf("%s/download?cid=%d&fid=%s&token=%s", c.gatewayURL, usr.CompanyID, req.FileID, tkn),
		},
		EditorConfig: response.EditorConfig{
			User: response.User{
				ID:   fmt.Sprint(usr.ID + usr.CompanyID),
				Name: usr.Name,
			},
			CallbackURL: fmt.Sprintf(
				"%s/callback?cid=%d&did=%s&fid=%s&filename=%s",
				c.gatewayURL, usr.CompanyID, req.Deal, req.FileID, req.Filename,
			),
			Customization: response.Customization{
				Goback: response.Goback{
					RequestClose: true,
				},
				Plugins:       false,
				HideRightMenu: true,
			},
			Lang: usr.Language.Lang,
		},
		Type:      t,
		ServerURL: settings.DocAddress,
	}

	if strings.TrimSpace(req.Filename) != "" {
		ext := strings.ReplaceAll(filepath.Ext(req.Filename), ".", "")
		fileType, err := constants.GetFileType(ext)
		if err != nil {
			return config, err
		}
		config.Document.FileType = ext
		config.Document.Permissions = response.Permissions{
			Edit:                 constants.IsExtensionEditable(ext),
			Comment:              true,
			Download:             true,
			Print:                false,
			Review:               false,
			Copy:                 true,
			ModifyContentControl: true,
			ModifyFilter:         true,
		}
		config.DocumentType = fileType
	}

	token, err := c.jwtManager.Sign(settings.DocSecret, config)
	if err != nil {
		c.logger.Debugf("could not sign document server config: %s", err.Error())
		return config, err
	}

	config.Token = token
	return config, nil
}

func (c ConfigHandler) BuildConfig(ctx context.Context, payload request.BuildConfigRequest, res *response.BuildConfigResponse) error {
	c.logger.Debugf("processing a docs config: %s", payload.Filename)

	config, err, _ := c.group.Do(fmt.Sprint(payload.UID+payload.CID), func() (interface{}, error) {
		req := c.client.NewRequest(
			fmt.Sprintf("%s:auth", c.namespace), "UserSelectHandler.GetUser",
			fmt.Sprint(payload.UID+payload.CID),
		)

		var ures response.UserResponse
		if err := c.client.Call(ctx, req, &ures); err != nil {
			c.logger.Debugf("could not get user %d access info: %s", payload.UID+payload.CID, err.Error())
			return nil, err
		}

		config, err := c.processConfig(ures, payload, ctx)
		if err != nil {
			return nil, err
		}

		return config, nil
	})

	if cfg, ok := config.(response.BuildConfigResponse); ok {
		*res = cfg
		return nil
	}

	return err
}
