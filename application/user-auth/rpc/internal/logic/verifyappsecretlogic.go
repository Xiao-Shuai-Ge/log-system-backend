package logic

import (
	"context"
	"crypto/subtle"
	"errors"

	"log-system-backend/application/user-auth/internal/repository"
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

// VerifyAppSecret 验证应用 ID (或 AppCode) 与应用密钥是否匹配
func (l *VerifyAppSecretLogic) VerifyAppSecret(in *auth.VerifyAppSecretRequest) (*auth.VerifyAppSecretResponse, error) {
	if in.AppId == "" || in.AppSecret == "" {
		return &auth.VerifyAppSecretResponse{IsValid: false}, nil
	}

	// 首先尝试通过 ID (UUID) 查找应用
	app, err := l.svcCtx.AppRepo.FindOne(l.ctx, in.AppId)
	if err != nil {
		if !errors.Is(err, repository.ErrNotFound) {
			logx.Errorf("VerifyAppSecret: 通过 ID '%s' 查找应用出错: %v", in.AppId, err)
			return nil, err
		}

		// 如果通过 ID 未找到，则尝试通过 AppCode 查找（AppId 有时与 ClientID 混用）
		app, err = l.svcCtx.AppRepo.FindOneByAppCode(l.ctx, in.AppId)
		if err != nil {
			if !errors.Is(err, repository.ErrNotFound) {
				logx.Errorf("VerifyAppSecret: 通过代码 '%s' 查找应用出错: %v", in.AppId, err)
				return nil, err
			}

			logx.Infof("VerifyAppSecret: 通过 ID 或代码均未找到应用: %s", in.AppId)
			return &auth.VerifyAppSecretResponse{IsValid: false}, nil
		}
	}

	// 使用恒定时间比较来防止计时攻击
	if subtle.ConstantTimeCompare([]byte(app.AppSecret), []byte(in.AppSecret)) == 1 {
		return &auth.VerifyAppSecretResponse{
			IsValid: true,
			AppName: app.AppName,
			AppCode: app.AppCode,
		}, nil
	}

	return &auth.VerifyAppSecretResponse{IsValid: false}, nil
}
