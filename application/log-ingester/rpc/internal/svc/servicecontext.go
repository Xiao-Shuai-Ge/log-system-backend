package svc

import (
	"log-system-backend/application/log-ingester/internal/repository"
	"log-system-backend/application/log-ingester/internal/service"
	"log-system-backend/application/log-ingester/rpc/internal/config"
	"log-system-backend/common/es"
)

type ServiceContext struct {
	Config   config.Config
	Ingester service.IngesterService
}

func NewServiceContext(c config.Config) *ServiceContext {
	esClient, err := es.NewClient(c.ES)
	if err != nil {
		panic(err)
	}

	index := c.Index
	if index == "" {
		index = "logs"
	}

	repo := repository.NewESRepository(esClient, index)
	return &ServiceContext{
		Config:   c,
		Ingester: service.NewIngesterService(repo),
	}
}
