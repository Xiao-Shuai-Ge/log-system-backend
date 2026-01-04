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

func NewAppService(repo repository.AppRepository, userRepo repository.UserRepository) AppService {
	return &appService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *appService) CreateApp(ctx context.Context, appCode, appName, description, userID string) (*repository.App, error) {
	// Check if app code exists
	_, err := s.repo.FindOneByAppCode(ctx, appCode)
	if err == nil {
		return nil, errorx.NewCodeError(errorx.CodeParamError, "app code already exists")
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, errorx.NewCodeError(errorx.CodeInternal, "database error")
	}

	app := &repository.App{
		AppCode:     appCode,
		AppName:     appName,
		Description: description,
	}

	// Insert app
	err = s.repo.Insert(ctx, app)
	if err != nil {
		return nil, errorx.NewCodeError(errorx.CodeInternal, "failed to create app")
	}

	// Assign to user
	if userID != "" {
		err = s.repo.AssignUser(ctx, app.ID, userID)
		if err != nil {
			// Rollback? Or just log error. For now, simple error return.
			// Ideally we should use transaction.
			return app, errorx.NewCodeError(errorx.CodeInternal, "failed to assign user to app")
		}
	}

	return app, nil
}

func (s *appService) UpdateApp(ctx context.Context, appID, appName, description string) error {
	app, err := s.repo.FindOne(ctx, appID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return errorx.NewCodeError(errorx.CodeParamError, "app not found")
		}
		return errorx.NewCodeError(errorx.CodeInternal, "database error")
	}

	app.AppName = appName
	app.Description = description

	err = s.repo.Update(ctx, app)
	if err != nil {
		return errorx.NewCodeError(errorx.CodeInternal, "failed to update app")
	}

	return nil
}

func (s *appService) VerifyUserAccess(ctx context.Context, userID, appCode string) (bool, error) {
	// 1. Get User's apps
	apps, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return false, nil
		}
		return false, errorx.NewCodeError(errorx.CodeInternal, "database error")
	}

	// 2. Check if appCode is in the list
	for _, app := range apps {
		if app.AppCode == appCode {
			return true, nil
		}
	}

	return false, nil
}

func (s *appService) DeleteApp(ctx context.Context, appID string) error {
	// We might want to check if it exists first
	_, err := s.repo.FindOne(ctx, appID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return errorx.NewCodeError(errorx.CodeParamError, "app not found")
		}
		return errorx.NewCodeError(errorx.CodeInternal, "database error")
	}

	err = s.repo.Delete(ctx, appID)
	if err != nil {
		return errorx.NewCodeError(errorx.CodeInternal, "failed to delete app")
	}
	return nil
}

func (s *appService) GetApp(ctx context.Context, appID string) (*repository.App, error) {
	app, err := s.repo.FindOne(ctx, appID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, errorx.NewCodeError(errorx.CodeParamError, "app not found")
		}
		return nil, errorx.NewCodeError(errorx.CodeInternal, "database error")
	}
	return app, nil
}

func (s *appService) ListUserApps(ctx context.Context, userID string) ([]*repository.App, error) {
	apps, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, errorx.NewCodeError(errorx.CodeInternal, "failed to list apps")
	}
	return apps, nil
}
