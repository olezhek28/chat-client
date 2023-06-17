package interceptor

import (
	"context"
	"time"

	"github.com/olezhek28/chat-client/internal/client/grpc/auth"
	"github.com/olezhek28/chat-client/internal/client/redis"
	"github.com/olezhek28/chat-client/internal/model"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthInterceptor struct {
	authClient  auth.Client
	redisClient redis.Client
}

func NewAuthInterceptor(authClient auth.Client, redisClient redis.Client) *AuthInterceptor {
	return &AuthInterceptor{
		authClient:  authClient,
		redisClient: redisClient,
	}
}

func (i *AuthInterceptor) Run(refreshTokenPeriod time.Duration, accessTokenPeriod time.Duration) {
	go func() {
		t := time.NewTicker(refreshTokenPeriod)
		ctx := context.Background()

		for _ = range t.C {
			username, err := i.redisClient.Get(model.UsernameKey)
			if err != nil {
				log.Error().Err(err).Msg("failed to get username from redis")
				continue
			}

			password, err := i.redisClient.Get(model.PasswordKey)
			if err != nil {
				log.Error().Err(err).Msg("failed to get password from redis")
				continue
			}

			refreshToken, err := i.authClient.GetRefreshToken(ctx, &model.AuthInfo{
				Username: username,
				Password: password,
			})
			if err != nil {
				log.Error().Err(err).Msg("failed to get refresh token")
				continue
			}

			err = i.redisClient.Set(model.RefreshTokenKey, refreshToken, 0)
			if err != nil {
				log.Error().Err(err).Msg("failed to set refresh token to redis")
				continue
			}

			log.Info().Msg("refresh token has been updated")
		}
	}()

	go func() {
		t := time.NewTicker(accessTokenPeriod)
		ctx := context.Background()

		for _ = range t.C {
			refreshToken, err := i.redisClient.Get(model.RefreshTokenKey)
			if err != nil {
				log.Error().Err(err).Msg("failed to get refresh token from redis")
				continue
			}

			accessToken, err := i.authClient.GetAccessToken(ctx, refreshToken)
			if err != nil {
				log.Error().Err(err).Msg("failed to get access token")
				continue
			}

			err = i.redisClient.Set(model.AccessTokenKey, accessToken, 0)
			if err != nil {
				log.Error().Err(err).Msg("failed to set access token to redis")
				continue
			}

			log.Info().Msg("access token has been updated")
		}
	}()
}

func (i *AuthInterceptor) Unary(ctx context.Context, method string, req interface{}, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	accessToken, err := i.redisClient.Get(model.AccessTokenKey)
	if err != nil {
		return err
	}

	md := metadata.New(map[string]string{"Authorization": "Bearer " + accessToken})
	ctx = metadata.NewOutgoingContext(ctx, md)

	return invoker(ctx, method, req, reply, cc, opts...)
}
