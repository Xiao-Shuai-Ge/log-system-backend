package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string `gorm:"primaryKey;type:char(36)"`
	Username  string `gorm:"uniqueIndex;type:varchar(255);not null"`
	Password  string `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	// 多对多关系
	Apps []*App `gorm:"many2many:user_apps;"`
}

// BeforeCreate GORM 钩子，在插入前生成 UUID
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return
}

type UserRepository interface {
	Insert(ctx context.Context, user *User) error
	FindOneByUsername(ctx context.Context, username string) (*User, error)
}

type mysqlUserRepository struct {
	db *gorm.DB
}

// NewMysqlUserRepository 创建 MySQL 用户仓储实例并自动迁移表结构
func NewMysqlUserRepository(db *gorm.DB) UserRepository {
	// 自动迁移
	db.AutoMigrate(&User{})
	return &mysqlUserRepository{
		db: db,
	}
}

// Insert 插入新用户
func (r *mysqlUserRepository) Insert(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// FindOneByUsername 根据用户名查找用户
func (r *mysqlUserRepository) FindOneByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

var ErrNotFound = gorm.ErrRecordNotFound
