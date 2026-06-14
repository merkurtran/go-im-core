# go-im-core
IM


# IM-Core-Backend

## 项目结构

```
IM-Core-Backend/
├── cmd/                        # 程序入口
│   └── api/
│       └── main.go             # 启动 HTTP 服务
├── internal/                   # 私有代码，不允许外部导入
│   ├── config/                 # 配置定义和加载
│   ├── domain/                 # 核心业务（实体、接口）
│   │   ├── user/               # 用户领域
│   │   │   ├── model.go        # User 结构体
│   │   │   └── repository.go   # User 存储接口（interface）
│   │   └── message/            # 消息领域
│   │       ├── model.go
│   │       └── repository.go
│   ├── service/                # 业务逻辑层
│   │   ├── user_service.go
│   │   └── message_service.go
│   ├── repository/             # 数据访问层（实现 domain 的接口）
│   │   ├── mongo/
│   │   │   ├── user_repo.go
│   │   │   └── message_repo.go
│   │   └── redis/
│   │       └── presence_repo.go
│   ├── handler/                # HTTP 处理器（原 controller）
│   │   ├── auth_handler.go
│   │   └── message_handler.go
│   ├── middleware/             # 中间件
│   │   ├── auth.go
│   │   └── logger.go
│   └── websocket/              # WebSocket 管理
│       ├── hub.go              # 连接管理中心
│       └── client.go           # 单个连接
├── pkg/                        # 可复用的公共库
│   ├── jwt/
│   ├── validator/
│   └── response/
├── api/                        # API 契约（可选）
│   └── proto/                  # 如果以后用 gRPC
├── scripts/                    # 脚本
├── deployments/                # Docker, k8s 配置
├── docs/                       # 文档
├── go.mod
├── go.sum
├── Makefile
└── README.md
```