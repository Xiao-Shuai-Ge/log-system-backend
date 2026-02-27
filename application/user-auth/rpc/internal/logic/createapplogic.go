package logic

import (
	"context"

	"log-system-backend/application/user-auth/rpc/auth"
	"log-system-backend/application/user-auth/rpc/internal/svc"

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

// CreateApp 处理应用创建请求，调用服务层创建应用并返回其 ID 和密钥
func (l *CreateAppLogic) CreateApp(in *auth.CreateAppRequest) (*auth.CreateAppResponse, error) {
	app, err := l.svcCtx.AppService.CreateApp(l.ctx, in.AppCode, in.AppName, in.Description, in.UserId)
	if err != nil {
		return nil, err
	}

	return &auth.CreateAppResponse{
		AppId:     app.ID,
		AppSecret: app.AppSecret,
	}, nil
}
