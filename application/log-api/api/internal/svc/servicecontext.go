// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"log-system-backend/application/log-api/api/internal/config"
	"log-system-backend/application/log-api/internal/service"
	"log-system-backend/common/rpc/logingester"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config        config.Config
	LogApiService service.LogApiService
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		LogApiService: service.NewLogApiService(logingester.NewLogIngester(zrpc.MustNewClient(c.LogIngesterRpc))),
	}
}
