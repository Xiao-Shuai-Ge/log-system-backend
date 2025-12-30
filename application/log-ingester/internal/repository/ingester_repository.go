package repository

import "context"

// LogRepository 定义了日志存储的接口
type LogRepository interface {
	Save(ctx context.Context, data map[string]interface{}) error
}

// esRepository 是 Elasticsearch 的实现
type esRepository struct {
	// 这里未来可以持有 ES 的 client
}

func NewESRepository() LogRepository {
	return &esRepository{}
}

func (r *esRepository) Save(ctx context.Context, data map[string]interface{}) error {
	// 实现具体的 ES 写入逻辑
	return nil
}
