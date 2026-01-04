# Distributed Log System (Log-System)

这是一个基于 **Go-Zero** 微服务框架构建的分布式日志查询演示系统。系统旨在模拟高性能日志采集、存储及实时查询场景。

## 启动命令行

### 1. Docker Compose 配置项启动（Elasticsearch、Etcd、Kibana）

```bash
# 在项目根目录下执行
docker-compose -f deploy/docker-compose.yml up -d
```

### 2. 微服务命令行启动

```bash
go run application/log-api/api/log.go
go run application/log-ingester/rpc/ingester.go
go run application/log-query/rpc/query.go
go run application/user-auth/rpc/auth.go

```

## 🏗 架构设计

### 1. 逻辑架构 (BFF + 微服务模式)
系统采用 **BFF (Backend For Frontend)** 架构，将外部接口与内部领域服务完全解耦：

- **Gateway Layer (API 网关/BFF)**: 统一的 RESTful API 入口，负责请求路由、用户鉴权、数据聚合及简单的业务编排。
- **Domain Service Layer (领域服务层)**: 核心业务逻辑实现，仅提供 gRPC 接口，不直接暴露给前端。
    - **Log Ingester**: 专职日志的高并发写入与存储分发。
    - **Log Query**: 专职海量日志的复杂检索与聚合分析。
    - **User Auth**: 专职用户权限管理、Token 签发与校验。
- **Common Layer (公共层)**: 提供全局错误处理、中间件、工具函数等。

### 2. 技术栈
- **框架**: [Go-Zero](https://github.com/zeromicro/go-zero) (微服务脚手架)
- **存储**: Elasticsearch (日志索引与搜索)
- **通信**: gRPC (内部服务间) & RESTful API (对外暴露)
- **部署**: Docker & Docker Compose
- **追踪**: OpenTelemetry / Jaeger (可选)

---

## 📂 项目结构

```text
.
├── application/                # 微服务集合
│   ├── log-api/                # API 网关 (BFF)，唯一的 HTTP 入口
│   │   └── api/                # RESTful 接口定义与聚合逻辑
│   ├── log-ingester/           # 日志接收服务 (RPC)
│   │   └── rpc/                # gRPC 接口，负责 ES 写入
│   ├── log-query/              # 日志查询服务 (RPC)
│   │   └── rpc/                # gRPC 接口，负责 ES 检索
│   └── user-auth/              # 用户认证服务 (RPC)
│       └── rpc/                # gRPC 接口，负责鉴权逻辑
├── common/                     # 公共组件库
│   ├── errorx/                 # 业务错误定义
│   ├── middleware/             # 统一中间件
│   └── rpc/                    # gRPC 客户端
├── deploy/                     # 部署相关
│   ├── docker-compose.yml      # 一键启动环境
│   └── sql/                    # (如有) 数据库初始化脚本
├── etc/                        # 配置文件 (yaml)
└── README.md
```

---

## 🚀 快速开始

### 1. 环境准备
确保已安装以下工具：
- Go 1.20+
- Docker & Docker Compose
- [goctl](https://go-zero.dev/docs/tasks/installation/goctl) (go-zero 命令行工具)

### 2. 启动基础环境
```bash
cd deploy
docker-compose up -d
```

### 3. 生成代码 (开发参考)
```bash
# 生成 API 代码
goctl api go -api *.api -dir .

# 生成 RPC 代码
goctl rpc protoc *.proto --go_out=. --go-grpc_out=. --zrpc_out=.
```

---

## 🛠 开发规范

1. **逻辑下沉**：严禁在 `api/internal/logic` 或 `rpc/internal/logic` 中直接编写复杂的业务逻辑。
2. **接口先行**：`repository` 层应定义接口，方便进行单元测试和存储介质切换。
3. **错误处理**：使用 `common/errorx` 定义的业务错误码，避免在代码中到处硬编码错误信息。
4. **配置隔离**：配置文件统一放在项目根目录 `etc/*.yaml`，并建议按服务入口命名（例如 `logingester-api.yaml`、`logingester-rpc.yaml`）。
