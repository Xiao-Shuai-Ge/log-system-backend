package logic

import (
	"context"

	"log-system-backend/application/user-auth/rpc/auth"
	"log-system-backend/application/user-auth/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppLogic {
	return &GetAppLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAppLogic) GetApp(in *auth.GetAppRequest) (*auth.GetAppResponse, error) {
	app, err := l.svcCtx.AppService.GetApp(l.ctx, in.AppId)
	if err != nil {
		return nil, err
	}

	return &auth.GetAppResponse{
		AppId:       app.ID,
		AppCode:     app.AppCode,
		AppName:     app.AppName,
		Description: app.Description,
		CreatedAt:   app.CreatedAt.Unix(),
		AppSecret:   app.AppSecret,
	}, nil
}
