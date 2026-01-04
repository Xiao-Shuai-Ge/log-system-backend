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
	// Many-to-Many relationship
	Apps []*App `gorm:"many2many:user_apps;"`
}

// BeforeCreate is a GORM hook that generates a UUID for the user before insertion
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

func NewMysqlUserRepository(db *gorm.DB) UserRepository {
	// Auto migrate
	db.AutoMigrate(&User{})
	return &mysqlUserRepository{
		db: db,
	}
}

func (r *mysqlUserRepository) Insert(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

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
