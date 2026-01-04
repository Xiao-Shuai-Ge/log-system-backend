package logic

import (
	"context"

	"log-system-backend/application/auth/rpc/auth"
	"log-system-backend/application/auth/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAppsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAppsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAppsLogic {
	return &ListAppsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListAppsLogic) ListApps(in *auth.ListAppsRequest) (*auth.ListAppsResponse, error) {
	apps, err := l.svcCtx.AppService.ListUserApps(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	var appInfos []*auth.AppInfo
	for _, app := range apps {
		appInfos = append(appInfos, &auth.AppInfo{
			AppId:       app.ID,
			AppCode:     app.AppCode,
			AppName:     app.AppName,
			Description: app.Description,
		})
	}

	return &auth.ListAppsResponse{
		Apps: appInfos,
	}, nil
}
