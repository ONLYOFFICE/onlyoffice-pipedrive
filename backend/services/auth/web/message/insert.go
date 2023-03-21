package message

import (
	"context"

	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/auth/web/core/domain"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/auth/web/core/port"
	"github.com/mitchellh/mapstructure"
)

type InsertMessageHandler struct {
	service port.UserAccessService
}

func BuildInsertMessageHandler(service port.UserAccessService) InsertMessageHandler {
	return InsertMessageHandler{
		service: service,
	}
}

func (i InsertMessageHandler) GetHandler() func(context.Context, interface{}) error {
	return func(ctx context.Context, payload interface{}) error {
		var user domain.UserAccess
		if err := mapstructure.Decode(payload, &user); err != nil {
			return err
		}
		_, err := i.service.UpdateUser(ctx, user)
		return err
	}
}
