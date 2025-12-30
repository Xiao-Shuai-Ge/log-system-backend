package logic

import (
	"context"

	"log-system-backend/application/log-query/rpc/internal/svc"
	"log-system-backend/application/log-query/rpc/types/query"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/types/known/structpb"
)

type SearchLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogLogic {
	return &SearchLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchLogLogic) SearchLog(in *query.SearchLogReq) (*query.SearchLogResp, error) {
	res, err := l.svcCtx.QueryService.SearchLog(l.ctx, in.Source, in.Keyword, in.Metadata, in.Page, in.PageSize)
	if err != nil {
		return nil, err
	}

	logs := make([]*query.LogItem, 0, len(res.Logs))
	for _, logMap := range res.Logs {
		id, _ := logMap["_id"].(string)
		source, _ := logMap["source"].(string)

		// Timestamp handling: priority to @timestamp
		timestamp, ok := logMap["@timestamp"].(string)
		if !ok {
			// fallback to timestamp
			if ts, ok := logMap["timestamp"].(string); ok {
				timestamp = ts
			}
		}

		content, err := structpb.NewStruct(logMap)
		if err != nil {
			l.Logger.Errorf("Failed to convert log map to struct: %v", err)
			continue
		}

		logs = append(logs, &query.LogItem{
			Id:        id,
			Source:    source,
			Timestamp: timestamp,
			Content:   content,
		})
	}

	return &query.SearchLogResp{
		Logs:  logs,
		Total: res.Total,
	}, nil
}
