package logic

import (
	"context"

	"log-system-backend/application/user-auth/rpc/auth"
	"log-system-backend/application/user-auth/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyAppAccessLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyAppAccessLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyAppAccessLogic {
	return &VerifyAppAccessLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *VerifyAppAccessLogic) VerifyAppAccess(in *auth.VerifyAppAccessRequest) (*auth.VerifyAppAccessResponse, error) {
	hasAccess, err := l.svcCtx.AppService.VerifyUserAccess(l.ctx, in.UserId, in.AppCode)
	if err != nil {
		return nil, err
	}

	return &auth.VerifyAppAccessResponse{
		HasAccess: hasAccess,
	}, nil
}
