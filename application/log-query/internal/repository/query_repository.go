package repository

import (
	"context"
	"fmt"
	"log-system-backend/common/es"

	"github.com/elastic/go-elasticsearch/v8"
)

type SearchQuery struct {
	Source   string
	Keyword  string
	Metadata map[string]string
	Page     int64
	PageSize int64
}

type SearchResult struct {
	Logs  []map[string]interface{}
	Total int64
}

type LogRepository interface {
	Search(ctx context.Context, q SearchQuery) (*SearchResult, error)
}

type esRepository struct {
	client *elasticsearch.Client
	index  string
}

// NewESRepository 创建基于 Elasticsearch 的日志查询仓储实例
func NewESRepository(client *elasticsearch.Client, index string) LogRepository {
	return &esRepository{
		client: client,
		index:  index,
	}
}

// Search 根据指定的查询条件从 Elasticsearch 中检索日志
func (r *esRepository) Search(ctx context.Context, q SearchQuery) (*SearchResult, error) {
	if r.client == nil {
		return nil, fmt.Errorf("ES 客户端为空")
	}

	must := []interface{}{}

	if q.Source != "" {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"source.keyword": q.Source,
			},
		})
	}

	if q.Keyword != "" {
		must = append(must, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  q.Keyword,
				"fields": []string{"*"},
			},
		})
	}

	for k, v := range q.Metadata {
		// 对于元数据，我们尝试对特定字段使用 .keyword 进行精确匹配
		// 这是 ES 动态字段的常用处理模式
		field := k
		if k == "level" || k == "source" {
			field = k + ".keyword"
		}
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				field: v,
			},
		})
	}

	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 {
		q.PageSize = 10
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		},
		"from": (q.Page - 1) * q.PageSize,
		"size": q.PageSize,
		"sort": []map[string]interface{}{
			{"@timestamp": map[string]interface{}{"order": "desc", "unmapped_type": "date"}},
			{"_score": "desc"},
		},
	}

	resp, err := es.SearchJSON(ctx, r.client, r.index, query)
	if err != nil {
		return nil, err
	}

	var esResp struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source map[string]interface{} `json:"_source"`
				ID     string                 `json:"_id"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := es.ReadJSON(resp, &esResp); err != nil {
		return nil, err
	}

	logs := make([]map[string]interface{}, 0, len(esResp.Hits.Hits))
	for _, hit := range esResp.Hits.Hits {
		logData := hit.Source
		if logData == nil {
			logData = make(map[string]interface{})
		}
		logData["_id"] = hit.ID
		logs = append(logs, logData)
	}

	return &SearchResult{
		Logs:  logs,
		Total: esResp.Hits.Total.Value,
	}, nil
}
