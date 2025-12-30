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
		// Try exact match on keyword field if possible, otherwise just term
		// We try k.keyword first? No, we don't know if it exists.
		// Let's assume the user passes the correct field name or we append .keyword
		// But for metadata, keys are dynamic. 
		// If we blindly append .keyword, it might fail if mapping doesn't exist.
		// Let's just use the key as provided, but maybe use "match_phrase" or "term".
		// "term" is exact match.
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				k: v,
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
			{"@timestamp": "desc"}, // Assuming @timestamp or just default score
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
