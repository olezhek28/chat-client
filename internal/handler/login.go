package handler

import (
	"context"

	"github.com/olezhek28/chat-client/internal/model"
)

func (h *Handler) Login(ctx context.Context, info *model.AuthInfo) error {
	refreshToken, err := h.authClient.GetRefreshToken(ctx, info)
	if err != nil {
		return err
	}

	accessToken, err := h.authClient.GetAccessToken(ctx, refreshToken)
	if err != nil {
		return err
	}

	err = h.redisClient.Set(model.UsernameKey, info.Username, 0)
	if err != nil {
		return err
	}

	err = h.redisClient.Set(model.PasswordKey, info.Password, 0)
	if err != nil {
		return err
	}

	err = h.redisClient.Set(model.AccessTokenKey, accessToken, 0)
	if err != nil {
		return err
	}

	err = h.redisClient.Set(model.RefreshTokenKey, refreshToken, 0)
	if err != nil {
		return err
	}

	return nil
}
