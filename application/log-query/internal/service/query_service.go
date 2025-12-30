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

func NewQueryService(repo repository.LogRepository) QueryService {
	return &queryService{
		repo: repo,
	}
}

func (s *queryService) SearchLog(ctx context.Context, source, keyword string, metadata map[string]string, page, pageSize int64) (*repository.SearchResult, error) {
	return s.repo.Search(ctx, repository.SearchQuery{
		Source:   source,
		Keyword:  keyword,
		Metadata: metadata,
		Page:     page,
		PageSize: pageSize,
	})
}
