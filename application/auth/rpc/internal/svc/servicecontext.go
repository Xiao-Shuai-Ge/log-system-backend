package svc

import (
	"log-system-backend/application/auth/internal/repository"
	"log-system-backend/application/auth/internal/service"
	"log-system-backend/application/auth/rpc/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	AuthService service.AuthService
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.DataSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Auto migrate is handled in repository
	repo := repository.NewMysqlUserRepository(db)
	authService := service.NewAuthService(repo, c.JwtAuth.AccessSecret, c.JwtAuth.AccessExpire)

	return &ServiceContext{
		Config:      c,
		AuthService: authService,
	}
}
