package logic

import (
	"context"

	"log-system-backend/application/auth/rpc/auth"
	"log-system-backend/application/auth/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAppLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAppLogic {
	return &DeleteAppLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteAppLogic) DeleteApp(in *auth.DeleteAppRequest) (*auth.DeleteAppResponse, error) {
	err := l.svcCtx.AppService.DeleteApp(l.ctx, in.AppId)
	if err != nil {
		return nil, err
	}

	return &auth.DeleteAppResponse{}, nil
}
