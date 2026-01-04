package logic

import (
	"context"

	"log-system-backend/application/auth/rpc/auth"
	"log-system-backend/application/auth/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAppLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAppLogic {
	return &CreateAppLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// App management
func (l *CreateAppLogic) CreateApp(in *auth.CreateAppRequest) (*auth.CreateAppResponse, error) {
	appID, err := l.svcCtx.AppService.CreateApp(l.ctx, in.AppCode, in.AppName, in.Description, in.UserId)
	if err != nil {
		return nil, err
	}

	return &auth.CreateAppResponse{
		AppId: appID,
	}, nil
}
