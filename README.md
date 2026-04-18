# w2work3

西二任务三，算一个尝试。
由于是一个多月前写的，现在回头看感觉烂的没边，尝试补救了一下。

## 部署

```bash
docker compose -f docs/docker-compose.yml up -d
go mod tidy
go run .
```

## 接口文档

OpenAPI 文件在：

- `docs/w2work3.openapi.json`

## 项目架构

```text
.
├── main.go
├── config.yaml
├── docs
│   ├── docker-compose.yml
│   └── w2work3.openapi.json
└── internal
    ├── apperr
    │   └── apperr.go
    ├── config
    │   └── config.go
    ├── constant
    │   └── constant.go
    ├── handler
    │   ├── user.go
    │   └── todo.go
    ├── infra
    │   └── db
    │       └── db.go
    ├── middleware
    │   └── jwtauth.go
    ├── model
    │   ├── user.go
    │   └── todo.go
    ├── repository
    │   ├── user.go
    │   └── todo.go
    ├── service
    │   ├── ports.go
    │   ├── user.go
    │   └── todo.go
    └── utils
        ├── jwt
        │   └── jwt.go
        └── password
            └── password.go
```
