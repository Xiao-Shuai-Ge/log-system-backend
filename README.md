# Distributed Log System (Log-System)

这是一个基于 **Go-Zero** 微服务框架构建的分布式日志查询演示系统。系统旨在模拟高性能日志采集、存储及实时查询场景。

## 🏗 架构设计

### 1. 逻辑架构
系统采用分层解耦架构，确保业务逻辑与通信协议（HTTP/gRPC）分离：

- **Entry Layer (入口层)**: `api/` 和 `rpc/` 目录。仅负责请求分发、参数校验及协议转换，不包含任何业务逻辑。
- **Service Layer (业务层)**: `internal/service/`。核心业务逻辑实现，供 API 和 RPC 共同调用。
- **Repository Layer (仓储层)**: `internal/repository/`。负责与 Elasticsearch、Redis 等底层存储交互，屏蔽存储细节。
- **Common Layer (公共层)**: `common/`。提供全局中间件、错误码、工具函数等。

### 2. 技术栈
- **框架**: [Go-Zero](https://github.com/zeromicro/go-zero) (微服务脚手架)
- **存储**: Elasticsearch (日志索引与搜索)
- **通信**: gRPC & RESTful API
- **部署**: Docker & Docker Compose
- **追踪**: OpenTelemetry / Jaeger (可选)

---

## 📂 项目结构

```text
.
├── application/                # 微服务集合
│   ├── log-ingester/           # 日志接收服务 (负责接收并写入日志)
│   │   ├── api/                # HTTP 入口 (go-zero 生成)
│   │   ├── rpc/                # gRPC 入口 (go-zero 生成)
│   │   ├── internal/
│   │   │   ├── service/        # 核心业务逻辑 (手工编写)
│   │   │   ├── repository/     # 数据持久化逻辑 (ES 操作)
│   │   │   └── config/         # 配置文件定义
│   │   └── etc/                # 配置文件 (yaml)
│   └── log-query/              # 日志查询服务 (负责搜索日志)
├── common/                     # 公共组件库
│   ├── errorx/                 # 业务错误定义
│   ├── middleware/             # 统一中间件
│   └── utils/                  # 工具类
├── deploy/                     # 部署相关
│   ├── docker-compose.yml      # 一键启动环境
│   └── sql/                    # (如有) 数据库初始化脚本
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
4. **配置隔离**：每个微服务维护自己的 `etc/*.yaml`，公共配置可提取到 `common`。
