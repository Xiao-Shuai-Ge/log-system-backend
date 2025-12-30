package config

import (
	"log-system-backend/common/es"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	ES    es.Config `json:",optional"`
	Index string    `json:",optional"`
}
