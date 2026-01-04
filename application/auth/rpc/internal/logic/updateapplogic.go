package logic

import (
	"context"

	"log-system-backend/application/auth/rpc/auth"
	"log-system-backend/application/auth/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAppLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAppLogic {
	return &UpdateAppLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAppLogic) UpdateApp(in *auth.UpdateAppRequest) (*auth.UpdateAppResponse, error) {
	err := l.svcCtx.AppService.UpdateApp(l.ctx, in.AppId, in.AppName, in.Description)
	if err != nil {
		return nil, err
	}

	return &auth.UpdateAppResponse{}, nil
}
