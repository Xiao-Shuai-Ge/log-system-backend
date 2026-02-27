package service

import (
	"context"
	"errors"

	"log-system-backend/application/user-auth/internal/repository"
	"log-system-backend/common/errorx"
)

type AppService interface {
	CreateApp(ctx context.Context, appCode, appName, description, userID string) (*repository.App, error)
	UpdateApp(ctx context.Context, appID, appName, description string) error
	DeleteApp(ctx context.Context, appID string) error
	GetApp(ctx context.Context, appID string) (*repository.App, error)
	ListUserApps(ctx context.Context, userID string) ([]*repository.App, error)
	VerifyUserAccess(ctx context.Context, userID, appCode string) (bool, error)
}

type appService struct {
	repo     repository.AppRepository
	userRepo repository.UserRepository
}

// NewAppService 创建应用服务实例
func NewAppService(repo repository.AppRepository, userRepo repository.UserRepository) AppService {
	return &appService{
		repo:     repo,
		userRepo: userRepo,
	}
}

// CreateApp 创建新应用并可选地分配给一个初始用户
func (s *appService) CreateApp(ctx context.Context, appCode, appName, description, userID string) (*repository.App, error) {
	_, err := s.repo.FindOneByAppCode(ctx, appCode)
	if err == nil {
		return nil, errorx.NewCodeError(errorx.CodeParamError, "应用代码已存在")
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, errorx.NewCodeError(errorx.CodeInternal, "数据库错误")
	}

	app := &repository.App{
		AppCode:     appCode,
		AppName:     appName,
		Description: description,
	}

	err = s.repo.Insert(ctx, app)
	if err != nil {
		return nil, errorx.NewCodeError(errorx.CodeInternal, "创建应用失败")
	}

	if userID != "" {
		err = s.repo.AssignUser(ctx, app.ID, userID)
		if err != nil {
			return app, errorx.NewCodeError(errorx.CodeInternal, "分配用户到应用失败")
		}
	}

	return app, nil
}

// UpdateApp 更新应用信息（如名称和描述）
func (s *appService) UpdateApp(ctx context.Context, appID, appName, description string) error {
	app, err := s.repo.FindOne(ctx, appID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return errorx.NewCodeError(errorx.CodeParamError, "应用未找到")
		}
		return errorx.NewCodeError(errorx.CodeInternal, "数据库错误")
	}

	app.AppName = appName
	app.Description = description

	err = s.repo.Update(ctx, app)
	if err != nil {
		return errorx.NewCodeError(errorx.CodeInternal, "更新应用失败")
	}

	return nil
}

// VerifyUserAccess 验证用户是否具有该应用的访问权限
func (s *appService) VerifyUserAccess(ctx context.Context, userID, appCode string) (bool, error) {
	hasAccess, err := s.repo.VerifyUserAccess(ctx, userID, appCode)
	if err != nil {
		return false, errorx.NewCodeError(errorx.CodeInternal, "数据库错误")
	}

	return hasAccess, nil
}

// DeleteApp 删除指定应用
func (s *appService) DeleteApp(ctx context.Context, appID string) error {
	// 检查应用是否存在
	_, err := s.repo.FindOne(ctx, appID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return errorx.NewCodeError(errorx.CodeParamError, "应用未找到")
		}
		return errorx.NewCodeError(errorx.CodeInternal, "数据库错误")
	}

	err = s.repo.Delete(ctx, appID)
	if err != nil {
		return errorx.NewCodeError(errorx.CodeInternal, "删除应用失败")
	}
	return nil
}

// GetApp 获取应用详情
func (s *appService) GetApp(ctx context.Context, appID string) (*repository.App, error) {
	app, err := s.repo.FindOne(ctx, appID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, errorx.NewCodeError(errorx.CodeParamError, "应用未找到")
		}
		return nil, errorx.NewCodeError(errorx.CodeInternal, "数据库错误")
	}
	return app, nil
}

// ListUserApps 获取用户拥有的应用列表
func (s *appService) ListUserApps(ctx context.Context, userID string) ([]*repository.App, error) {
	apps, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, errorx.NewCodeError(errorx.CodeInternal, "列出应用失败")
	}
	return apps, nil
}
