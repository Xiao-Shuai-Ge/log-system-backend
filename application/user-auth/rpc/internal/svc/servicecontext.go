package svc

import (
	"log-system-backend/application/user-auth/internal/repository"
	"log-system-backend/application/user-auth/internal/service"
	"log-system-backend/application/user-auth/rpc/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	AuthService service.AuthService
	AppService  service.AppService
	AppRepo     repository.AppRepository
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.DataSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Auto migrate is handled in repository
	userRepo := repository.NewMysqlUserRepository(db)
	appRepo := repository.NewMysqlAppRepository(db)
	authService := service.NewAuthService(userRepo, c.JwtAuth.AccessSecret, c.JwtAuth.AccessExpire)
	appService := service.NewAppService(appRepo, userRepo)

	return &ServiceContext{
		Config:      c,
		AuthService: authService,
		AppService:  appService,
		AppRepo:     appRepo,
	}
}
