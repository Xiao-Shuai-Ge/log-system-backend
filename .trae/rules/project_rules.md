# 需求

构建一个分布式的日志查询演示系统。该系统由多个微服务组成，能够模拟接收、存储和查询来自不同服务器或应用的日志数据。

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

## 核心框架/库

微服务框架：使用 Go Zero 的服务发现、RPC、配置管理等微服务所需的基础组件。API、RPC 均要使用 Go Zero 提供的工具生成。
RPC协议：使用 gRPC（高效）或JSON over HTTP（简单）进行服务间通信。
数据存储：使用Elasticsearch（ES）存储和索引日志（便于关键词查询）。
部署：使用 Docker 容器化每个微服务，并通过 docker-compose.yml 文件在开发环境一键编排部署所有组件。
系统结构与服务划分：
- **Log API (BFF/Gateway)**：作为系统唯一的 RESTful 入口，负责聚合各个 RPC 服务的数据并处理用户鉴权。
- **Log Ingester (RPC Service)**：接收日志写入请求，专注于高并发持久化到 Elasticsearch。
- **Log Query (RPC Service)**：提供强大的日志检索接口，屏蔽 ES 查询复杂度。
- **User Auth (RPC Service)**：负责用户管理、登录鉴权及权限控制。
- **Web UI**：Vue 前端，仅与 Log API 通信。

关键技术点：
- **聚合层设计**：Log API 严禁包含复杂的存储层逻辑，仅负责调用 RPC 服务进行业务编排。
- **服务隔离**：内部 RPC 服务不暴露 HTTP 端口，增强安全性。
- **高效通信**：BFF 与领域服务之间统一使用 gRPC 协议。

拓展功能：
为日志接收端增加身份认证，防止恶意写入。实现对日志数据的简单统计图表展示（如错误日志数量随时间变化）。
