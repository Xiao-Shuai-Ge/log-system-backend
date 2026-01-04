package logic

import (
	"context"

	"log-system-backend/application/user-auth/rpc/auth"
	"log-system-backend/application/user-auth/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyAppSecretLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyAppSecretLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyAppSecretLogic {
	return &VerifyAppSecretLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *VerifyAppSecretLogic) VerifyAppSecret(in *auth.VerifyAppSecretRequest) (*auth.VerifyAppSecretResponse, error) {
	if in.AppId == "" || in.AppSecret == "" {
		return &auth.VerifyAppSecretResponse{IsValid: false}, nil
	}

	// First try to find by ID (UUID)
	app, err := l.svcCtx.AppRepo.FindOne(l.ctx, in.AppId)
	if err != nil {
		// If not found by ID, try by AppCode as "AppId" is often used interchangeably with ClientID
		appByCode, errCode := l.svcCtx.AppRepo.FindOneByAppCode(l.ctx, in.AppId)
		if errCode != nil {
			logx.Errorf("VerifyAppSecret: app not found by id or code: %s", in.AppId)
			return &auth.VerifyAppSecretResponse{IsValid: false}, nil
		}
		app = appByCode
	}

	if app.AppSecret == in.AppSecret {
		return &auth.VerifyAppSecretResponse{
			IsValid: true,
			AppName: app.AppName,
		}, nil
	}

	return &auth.VerifyAppSecretResponse{IsValid: false}, nil
}
