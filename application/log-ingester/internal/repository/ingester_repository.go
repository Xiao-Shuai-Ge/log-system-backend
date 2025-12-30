package repository

import (
	"context"
	"fmt"

	"log-system-backend/common/es"

	"github.com/elastic/go-elasticsearch/v8"
)

// LogRepository 定义了日志存储的接口
type LogRepository interface {
	Save(ctx context.Context, data map[string]interface{}) error
}

// esRepository 是 Elasticsearch 的实现
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

func (r *esRepository) Save(ctx context.Context, data map[string]interface{}) error {
	if r.client == nil {
		return fmt.Errorf("es client is nil")
	}
	if r.index == "" {
		return fmt.Errorf("es index is empty")
	}

	resp, err := es.IndexJSON(ctx, r.client, r.index, data, r.client.Index.WithRefresh("wait_for"))
	if err != nil {
		return err
	}
	_, err = es.ReadBody(resp)
	return err
}
