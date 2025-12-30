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
	err = l.svcCtx.LogApiService.WriteLog(l.ctx, req.Source, req.Level, req.Content, req.Metadata)
	if err != nil {
		// 这里不再需要重复记录错误日志，因为 service 层或者 handler 已经有记录
		// 直接返回错误，由全局错误处理器处理
		return nil, err
	}

	return &types.WriteLogResp{
		Ok: true,
	}, nil
}
