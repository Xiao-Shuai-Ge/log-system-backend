package service

import (
	"context"
	"log-system-backend/application/log-ingester/internal/repository"
)

// IngesterService 定义了日志接收的核心业务逻辑接口
type IngesterService interface {
	WriteLog(ctx context.Context, data map[string]interface{}) error
}

type ingesterService struct {
	repo repository.LogRepository
}

func NewIngesterService(repo repository.LogRepository) IngesterService {
	return &ingesterService{
		repo: repo,
	}
}

func (s *ingesterService) WriteLog(ctx context.Context, data map[string]interface{}) error {
	// 这里可以添加业务逻辑，比如：
	// 1. 数据清洗
	// 2. 格式化转换
	// 3. 权限校验等
	return s.repo.Save(ctx, data)
}
