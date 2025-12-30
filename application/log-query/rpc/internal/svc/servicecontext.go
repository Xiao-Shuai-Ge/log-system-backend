package svc

import (
	"log-system-backend/application/log-query/internal/repository"
	"log-system-backend/application/log-query/internal/service"
	"log-system-backend/application/log-query/rpc/internal/config"
	"log-system-backend/common/es"
)

type ServiceContext struct {
	Config       config.Config
	QueryService service.QueryService
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
		Config:       c,
		QueryService: service.NewQueryService(repo),
	}
}
