// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package app

import (
	"context"

	"log-system-backend/application/log-api/api/internal/svc"
	"log-system-backend/application/log-api/api/internal/types"
	"log-system-backend/common/ctxutils"
	"log-system-backend/common/rpc/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAppsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAppsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAppsLogic {
	return &ListAppsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListAppsLogic) ListApps(req *types.ListAppsReq) (resp *types.ListAppsResp, err error) {
	userIdStr, err := ctxutils.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.AuthRpc.ListApps(l.ctx, &auth.ListAppsRequest{
		UserId: userIdStr,
	})
	if err != nil {
		return nil, err
	}

	var apps []types.AppInfo
	for _, app := range rpcResp.Apps {
		apps = append(apps, types.AppInfo{
			AppId:       app.AppId,
			AppCode:     app.AppCode,
			AppName:     app.AppName,
			Description: app.Description,
			AppSecret:   app.AppSecret,
		})
	}

	return &types.ListAppsResp{
		Apps: apps,
	}, nil
}
