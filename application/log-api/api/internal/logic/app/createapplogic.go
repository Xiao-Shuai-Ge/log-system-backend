// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package app

import (
	"context"
	"encoding/json"
	"fmt"

	"log-system-backend/application/log-api/api/internal/svc"
	"log-system-backend/application/log-api/api/internal/types"
	"log-system-backend/common/rpc/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAppLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAppLogic {
	return &CreateAppLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAppLogic) CreateApp(req *types.CreateAppReq) (resp *types.CreateAppResp, err error) {
	userId := l.ctx.Value("userId")
	var userIdStr string
	if v, ok := userId.(string); ok {
		userIdStr = v
	} else if v, ok := userId.(json.Number); ok {
		userIdStr = v.String()
	} else {
		return nil, fmt.Errorf("invalid user id in context")
	}

	rpcResp, err := l.svcCtx.AuthRpc.CreateApp(l.ctx, &auth.CreateAppRequest{
		AppCode:     req.AppCode,
		AppName:     req.AppName,
		Description: req.Description,
		UserId:      userIdStr,
	})
	if err != nil {
		return nil, err
	}

	return &types.CreateAppResp{
		AppId: rpcResp.AppId,
	}, nil
}
