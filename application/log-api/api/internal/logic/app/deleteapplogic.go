// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package app

import (
	"context"

	"log-system-backend/application/log-api/api/internal/svc"
	"log-system-backend/application/log-api/api/internal/types"
	"log-system-backend/common/rpc/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAppLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAppLogic {
	return &DeleteAppLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteAppLogic) DeleteApp(req *types.DeleteAppReq) (resp *types.DeleteAppResp, err error) {
	_, err = l.svcCtx.AuthRpc.DeleteApp(l.ctx, &auth.DeleteAppRequest{
		AppId: req.AppId,
	})
	if err != nil {
		return nil, err
	}

	return &types.DeleteAppResp{}, nil
}
