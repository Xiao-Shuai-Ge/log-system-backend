// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"log-system-backend/application/log-api/api/internal/config"
	"log-system-backend/application/log-api/api/internal/middleware"
	"log-system-backend/application/log-api/internal/service"
	"log-system-backend/common/rpc/auth"
	"log-system-backend/common/rpc/logingester"
	"log-system-backend/common/rpc/logquery"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config            config.Config
	LogApiService     service.LogApiService
	LogQueryRpc       logquery.LogQuery
	AuthRpc           auth.Auth
	AppAuthMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	logQueryRpc := logquery.NewLogQuery(zrpc.MustNewClient(c.LogQueryRpc))
	authRpc := auth.NewAuth(zrpc.MustNewClient(c.AuthRpc))
	return &ServiceContext{
		Config:            c,
		LogApiService:     service.NewLogApiService(logingester.NewLogIngester(zrpc.MustNewClient(c.LogIngesterRpc)), logQueryRpc, authRpc),
		LogQueryRpc:       logQueryRpc,
		AuthRpc:           authRpc,
		AppAuthMiddleware: middleware.NewAppAuthMiddleware(authRpc).Handle,
	}
}
