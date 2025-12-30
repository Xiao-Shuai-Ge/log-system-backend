// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"encoding/json"

	"log-system-backend/application/log-api/api/internal/svc"
	"log-system-backend/application/log-api/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogLogic {
	return &SearchLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchLogLogic) SearchLog(req *types.SearchLogReq) (resp *types.SearchLogResp, err error) {
	var metadata map[string]string
	if req.Metadata != "" {
		if err := json.Unmarshal([]byte(req.Metadata), &metadata); err != nil {
			return nil, err
		}
	}

	rpcResp, err := l.svcCtx.LogApiService.SearchLog(l.ctx, req.Source, req.Keyword, metadata, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	logs := make([]types.LogItem, 0, len(rpcResp.Logs))
	for _, logItem := range rpcResp.Logs {
		logs = append(logs, types.LogItem{
			Id:        logItem.Id,
			Source:    logItem.Source,
			Timestamp: logItem.Timestamp,
			Content:   logItem.Content.AsMap(),
		})
	}

	return &types.SearchLogResp{
		Logs:  logs,
		Total: rpcResp.Total,
	}, nil
}
