package logic

import (
	"context"
	"fmt"
	"time"

	"log-system-backend/application/log-ingester/rpc/internal/svc"
	"log-system-backend/application/log-ingester/rpc/types/ingester"

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

func (l *WriteLogLogic) WriteLog(in *ingester.WriteLogReq) (*ingester.WriteLogResp, error) {
	if in == nil || in.Data == nil {
		return nil, fmt.Errorf("data is required")
	}

	data := in.Data.AsMap()
	if len(data) == 0 {
		return nil, fmt.Errorf("data is empty")
	}
	if _, ok := data["@timestamp"]; !ok {
		data["@timestamp"] = time.Now().UTC().Format(time.RFC3339Nano)
	}

	if err := l.svcCtx.Ingester.WriteLog(l.ctx, data); err != nil {
		return nil, err
	}

	return &ingester.WriteLogResp{
		Result: "ok",
	}, nil
}
