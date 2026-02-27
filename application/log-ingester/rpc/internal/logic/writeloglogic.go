package logic

import (
	"context"

	"log-system-backend/application/log-ingester/rpc/internal/svc"
	"log-system-backend/application/log-ingester/rpc/types/ingester"
	"log-system-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type WriteLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWriteLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WriteLogLogic {
	return &WriteLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// WriteLog 处理来自 RPC 的日志写入请求，将 Protobuf Struct 转换为 Map 并调用服务层
func (l *WriteLogLogic) WriteLog(in *ingester.WriteLogReq) (*ingester.WriteLogResp, error) {
	if in == nil || in.Data == nil {
		return nil, errorx.NewCodeError(errorx.CodeParamError, "日志数据不能为空")
	}

	data := in.Data.AsMap()
	if len(data) == 0 {
		return nil, errorx.NewCodeError(errorx.CodeParamError, "日志数据内容为空")
	}

	if err := l.svcCtx.Ingester.WriteLog(l.ctx, data); err != nil {
		return nil, err
	}

	return &ingester.WriteLogResp{
		Result: "ok",
	}, nil
}
