package service

import (
	"context"
	"log-system-backend/application/log-query/internal/repository"
)

type QueryService interface {
	SearchLog(ctx context.Context, source, keyword string, metadata map[string]string, page, pageSize int64) (*repository.SearchResult, error)
}

type queryService struct {
	repo repository.LogRepository
}

// NewQueryService 创建日志查询服务实例
func NewQueryService(repo repository.LogRepository) QueryService {
	return &queryService{
		repo: repo,
	}
}

// SearchLog 调用仓储层执行日志搜索
func (s *queryService) SearchLog(ctx context.Context, source, keyword string, metadata map[string]string, page, pageSize int64) (*repository.SearchResult, error) {
	return s.repo.Search(ctx, repository.SearchQuery{
		Source:   source,
		Keyword:  keyword,
		Metadata: metadata,
		Page:     page,
		PageSize: pageSize,
	})
}
