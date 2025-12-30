// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"log-system-backend/application/log-api/api/internal/svc"
	"log-system-backend/application/log-api/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type WriteLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWriteLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WriteLogLogic {
	return &WriteLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WriteLogLogic) WriteLog(req *types.WriteLogReq) (resp *types.WriteLogResp, err error) {
	err = l.svcCtx.LogApiService.WriteLog(l.ctx, req.Source, req.Level, req.Content)
	if err != nil {
		l.Errorf("Failed to write log: %v", err)
		return nil, err
	}

	return &types.WriteLogResp{
		Ok: true,
	}, nil
}
