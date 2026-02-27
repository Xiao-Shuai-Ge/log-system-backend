package repository

import (
	"context"
	"errors"
	"time"

	"log-system-backend/common/cryptox"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type App struct {
	ID          string `gorm:"primaryKey;type:char(36)"`
	AppCode     string `gorm:"uniqueIndex;type:varchar(50);not null;comment:应用唯一标识码"`
	AppName     string `gorm:"type:varchar(255);not null;comment:应用名称"`
	AppSecret   string `gorm:"type:varchar(64);not null;comment:用于 API 访问的应用密钥"`
	Description string `gorm:"type:text;comment:应用描述"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	// 多对多关系
	Users []*User `gorm:"many2many:user_apps;"`
}

// BeforeCreate GORM 钩子，在插入前生成 UUID 和应用密钥
func (a *App) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	if a.AppSecret == "" {
		secret, err := cryptox.GenerateRandomString(32) // 32 字节 = 64 位十六进制字符
		if err != nil {
			return err
		}
		a.AppSecret = secret
	}
	return
}

type AppRepository interface {
	Insert(ctx context.Context, app *App) error
	FindOne(ctx context.Context, id string) (*App, error)
	FindOneByAppCode(ctx context.Context, appCode string) (*App, error)
	Update(ctx context.Context, app *App) error
	Delete(ctx context.Context, id string) error
	ListByUserID(ctx context.Context, userID string) ([]*App, error)
	AssignUser(ctx context.Context, appID, userID string) error
	RemoveUser(ctx context.Context, appID, userID string) error
	VerifyUserAccess(ctx context.Context, userID, appCode string) (bool, error)
}

type mysqlAppRepository struct {
	db *gorm.DB
}

// NewMysqlAppRepository 创建 MySQL 应用仓储实例并自动迁移表结构
func NewMysqlAppRepository(db *gorm.DB) AppRepository {
	// 自动迁移
	db.AutoMigrate(&App{})
	return &mysqlAppRepository{
		db: db,
	}
}

// Insert 插入新应用
func (r *mysqlAppRepository) Insert(ctx context.Context, app *App) error {
	return r.db.WithContext(ctx).Create(app).Error
}

// FindOne 根据 ID 查找应用
func (r *mysqlAppRepository) FindOne(ctx context.Context, id string) (*App, error) {
	var app App
	err := r.db.WithContext(ctx).First(&app, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &app, nil
}

// FindOneByAppCode 根据 AppCode 查找应用
func (r *mysqlAppRepository) FindOneByAppCode(ctx context.Context, appCode string) (*App, error) {
	var app App
	err := r.db.WithContext(ctx).First(&app, "app_code = ?", appCode).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &app, nil
}

// Update 更新应用信息
func (r *mysqlAppRepository) Update(ctx context.Context, app *App) error {
	return r.db.WithContext(ctx).Save(app).Error
}

// Delete 根据 ID 删除应用
func (r *mysqlAppRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&App{}, "id = ?", id).Error
}

// ListByUserID 获取用户有权限访问的所有应用列表
func (r *mysqlAppRepository) ListByUserID(ctx context.Context, userID string) ([]*App, error) {
	var user User
	// 预加载用户的应用列表
	err := r.db.WithContext(ctx).Preload("Apps").First(&user, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user.Apps, nil
}

// AssignUser 为应用分配用户访问权限
func (r *mysqlAppRepository) AssignUser(ctx context.Context, appID, userID string) error {
	// 使用关联模式
	// 简单方法：先找到应用和用户，然后追加关联
	var app App
	if err := r.db.WithContext(ctx).First(&app, "id = ?", appID).Error; err != nil {
		return err
	}
	var user User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", userID).Error; err != nil {
		return err
	}

	return r.db.WithContext(ctx).Model(&app).Association("Users").Append(&user)
}

// RemoveUser 移除用户的应用访问权限
func (r *mysqlAppRepository) RemoveUser(ctx context.Context, appID, userID string) error {
	var app App
	if err := r.db.WithContext(ctx).First(&app, "id = ?", appID).Error; err != nil {
		return err
	}
	var user User
	// 通常只需要 ID 即可移除，但 GORM 需要对象或 ID
	user.ID = userID
	return r.db.WithContext(ctx).Model(&app).Association("Users").Delete(&user)
}

// VerifyUserAccess 验证用户是否具有指定应用的访问权限
func (r *mysqlAppRepository) VerifyUserAccess(ctx context.Context, userID, appCode string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("user_apps").
		Joins("JOIN apps ON apps.id = user_apps.app_id").
		Where("user_apps.user_id = ? AND apps.app_code = ?", userID, appCode).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
