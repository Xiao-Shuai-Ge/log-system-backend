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

func NewESRepository(client *elasticsearch.Client, index string) LogRepository {
	return &esRepository{
		client: client,
		index:  index,
	}
}

func (r *esRepository) Search(ctx context.Context, q SearchQuery) (*SearchResult, error) {
	if r.client == nil {
		return nil, fmt.Errorf("es client is nil")
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
		// For metadata, we try to use .keyword for exact matching if it's not already specified
		// This is a common pattern in ES for dynamic fields
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
