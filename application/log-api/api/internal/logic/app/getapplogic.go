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

type GetAppLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppLogic {
	return &GetAppLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAppLogic) GetApp(req *types.GetAppReq) (resp *types.GetAppResp, err error) {
	rpcResp, err := l.svcCtx.AuthRpc.GetApp(l.ctx, &auth.GetAppRequest{
		AppId: req.AppId,
	})
	if err != nil {
		return nil, err
	}

	return &types.GetAppResp{
		AppId:       rpcResp.AppId,
		AppCode:     rpcResp.AppCode,
		AppName:     rpcResp.AppName,
		Description: rpcResp.Description,
		AppSecret:   rpcResp.AppSecret,
		CreatedAt:   rpcResp.CreatedAt,
	}, nil
}
