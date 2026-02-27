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

// NewESRepository 创建一个基于 Elasticsearch 的日志仓储实例
func NewESRepository(client *elasticsearch.Client, index string) LogRepository {
	return &esRepository{
		client: client,
		index:  index,
	}
}

// Save 将单条日志数据保存到 Elasticsearch
func (r *esRepository) Save(ctx context.Context, data map[string]interface{}) error {
	if r.client == nil {
		return fmt.Errorf("ES 客户端为空")
	}
	if r.index == "" {
		return fmt.Errorf("ES 索引名为空")
	}

	resp, err := es.IndexJSON(ctx, r.client, r.index, data)
	if err != nil {
		return err
	}
	_, err = es.ReadBody(resp)
	return err
}
