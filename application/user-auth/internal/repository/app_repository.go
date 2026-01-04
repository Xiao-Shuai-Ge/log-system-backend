package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type App struct {
	ID          string `gorm:"primaryKey;type:char(36)"`
	AppCode     string `gorm:"uniqueIndex;type:varchar(50);not null;comment:Application Identifier Code"`
	AppName     string `gorm:"type:varchar(255);not null;comment:Application Name"`
	AppSecret   string `gorm:"type:varchar(64);not null;comment:Application Secret for API Access"`
	Description string `gorm:"type:text;comment:Application Description"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	// Many-to-Many relationship
	Users []*User `gorm:"many2many:user_apps;"`
}

// BeforeCreate is a GORM hook that generates a UUID for the app before insertion
func (a *App) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	if a.AppSecret == "" {
		a.AppSecret = uuid.New().String() // Simple UUID as secret for now
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
}

type mysqlAppRepository struct {
	db *gorm.DB
}

func NewMysqlAppRepository(db *gorm.DB) AppRepository {
	// Auto migrate
	db.AutoMigrate(&App{})
	return &mysqlAppRepository{
		db: db,
	}
}

func (r *mysqlAppRepository) Insert(ctx context.Context, app *App) error {
	return r.db.WithContext(ctx).Create(app).Error
}

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

func (r *mysqlAppRepository) Update(ctx context.Context, app *App) error {
	return r.db.WithContext(ctx).Save(app).Error
}

func (r *mysqlAppRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&App{}, "id = ?", id).Error
}

func (r *mysqlAppRepository) ListByUserID(ctx context.Context, userID string) ([]*App, error) {
	var user User
	// Preload Apps for the user
	err := r.db.WithContext(ctx).Preload("Apps").First(&user, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user.Apps, nil
}

func (r *mysqlAppRepository) AssignUser(ctx context.Context, appID, userID string) error {
	// We need to use the association mode
	// But simple way is to find App and Append User
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

func (r *mysqlAppRepository) RemoveUser(ctx context.Context, appID, userID string) error {
	var app App
	if err := r.db.WithContext(ctx).First(&app, "id = ?", appID).Error; err != nil {
		return err
	}
	var user User
	// We only need the ID for removal usually, but GORM needs the object or ID
	user.ID = userID
	return r.db.WithContext(ctx).Model(&app).Association("Users").Delete(&user)
}
