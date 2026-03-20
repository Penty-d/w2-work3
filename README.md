# w2work3

这是一个基于 Go、Hertz、Gorm 和 PostgreSQL 实现的 Todo API 项目，主要完成用户注册登录、JWT 鉴权，以及待办事项的增删改查、分页查询和按状态批量处理。

## 写在前面

本文档由ai编写

## 项目结构

```text
w2-work3
├── main.go
├── config.yaml
├── docker-compose.yml
├── api.json
├── go.mod
├── internal
│   ├── config
│   │   └── config.go
│   ├── handler
│   │   └── handler.go
│   ├── middleware
│   │   └── jwtauth.go
│   ├── model
│   │   ├── user.go
│   │   └── todo.go
│   ├── repository
│   │   └── repository.go
│   ├── service
│   │   └── service.go
│   └── utils
│       ├── jwt
│       │   └── jwt.go
│       └── password
│           └── password.go
└── 框架.txt
```

## 各模块说明

### `main.go`

项目入口文件，负责：

- 加载配置
- 初始化数据库
- 创建 repository、service、handler
- 注册路由和中间件
- 启动 Hertz 服务

### `config.yaml`

项目配置文件，包含：

- 服务端口
- PostgreSQL 连接信息
- JWT 密钥与过期时间
- Redis 预留配置

### `internal/config`

负责读取并解析配置文件，将配置映射到结构体中供其他模块使用。

### `internal/model`

定义项目使用的数据模型：

- `User`：用户信息
- `Todo`：待办事项信息
- `TodoQueryConditions`：Todo 查询条件

### `internal/repository`

数据访问层，负责直接操作数据库，包括：

- 用户的创建、查询、删除
- Todo 的创建、查询、更新、删除

这一层主要处理 SQL/Gorm 相关逻辑。

### `internal/service`

业务逻辑层，负责：

- 参数合法性校验
- 用户注册、登录、删除
- Todo 的新增、查询、更新、删除
- 批量更新状态和批量删除等业务处理

这一层连接 handler 和 repository，是项目的核心业务层。

### `internal/handler`

HTTP 接口处理层，负责：

- 接收请求
- 解析 JSON 参数
- 调用 service 层
- 返回统一 JSON 响应

### `internal/middleware`

中间件层，目前主要实现 JWT 鉴权：

- 解析请求头中的 Bearer Token
- 校验 token 合法性
- 将当前用户信息写入上下文

### `internal/utils`

工具层，提供通用能力：

- `jwt`：生成和解析 JWT
- `password`：密码哈希与校验

### `api.json`

项目的 OpenAPI 接口文档文件，记录了接口路径、请求方法、请求体示例、响应示例和鉴权方式。

## 项目架构

本项目采用三层架构：

1. `handler`：处理 HTTP 请求和响应
2. `service`：处理业务逻辑
3. `repository`：处理数据库访问

请求的大致流程为：

客户端请求 -> 中间件鉴权 -> handler -> service -> repository -> 数据库

## 已实现的主要接口

### 用户模块

- `POST /api/v1/user/signup`：用户注册
- `POST /api/v1/user/login`：用户登录
- `DELETE /api/v1/user/delete`：删除用户

### Todo 模块

- `POST /api/v1/todo/add`：新增 Todo
- `GET /api/v1/todo`：查询 Todo 列表
- `PATCH /api/v1/todo`：更新单个 Todo
- `PATCH /api/v1/todo/status`：批量更新 Todo 状态
- `DELETE /api/v1/todo`：删除指定 Todo
- `DELETE /api/v1/todo/status`：按状态删除 Todo
- `DELETE /api/v1/todo/all`：删除当前用户全部 Todo

## 运行方式

### 1. 启动 PostgreSQL

```bash
docker compose up -d
```

### 2. 启动项目

```bash
go run .
```

默认服务地址为：

```text
http://0.0.0.0:8080
```

## 技术栈

- Go
- Hertz
- Gorm
- PostgreSQL
- JWT
- Viper
