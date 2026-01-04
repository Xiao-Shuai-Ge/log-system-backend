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

type UpdateAppLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAppLogic {
	return &UpdateAppLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAppLogic) UpdateApp(req *types.UpdateAppReq) (resp *types.UpdateAppResp, err error) {
	_, err = l.svcCtx.AuthRpc.UpdateApp(l.ctx, &auth.UpdateAppRequest{
		AppId:       req.AppId,
		AppName:     req.AppName,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	return &types.UpdateAppResp{}, nil
}
