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

// SearchLog 处理来自 RPC 的日志搜索请求，将结果转换为 Protobuf 响应格式
func (l *SearchLogLogic) SearchLog(in *query.SearchLogReq) (*query.SearchLogResp, error) {
	res, err := l.svcCtx.QueryService.SearchLog(l.ctx, in.Source, in.Keyword, in.Metadata, in.Page, in.PageSize)
	if err != nil {
		return nil, err
	}

	logs := make([]*query.LogItem, 0, len(res.Logs))
	for _, logMap := range res.Logs {
		id, _ := logMap["_id"].(string)
		source, _ := logMap["source"].(string)

		// 时间戳处理：优先使用 @timestamp
		timestamp, ok := logMap["@timestamp"].(string)
		if !ok {
			// 备选方案使用 timestamp 字段
			if ts, ok := logMap["timestamp"].(string); ok {
				timestamp = ts
			}
		}

		content, err := structpb.NewStruct(logMap)
		if err != nil {
			l.Logger.Errorf("将日志 Map 转换为 Protobuf Struct 失败: %v", err)
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
