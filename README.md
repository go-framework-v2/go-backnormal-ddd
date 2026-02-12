# go-backnormal-ddd

go-backnormal框架，结合ddd落地实现模板

## go-backnormal+DDD

在自研go-backnormal web框架基础上，使用DDD领域驱动设计方法论落地实现

## 目录结构说明

```
src/
├── api/                          # 用户接口层：路由与 HTTP 入口
│   ├── route.go                  # 路由注册入口
│   └── route_login.go            # 登录相关路由及 Handler（调用 application/identity）
├── cmd/                          # 启动入口与配置
│   ├── main.go
│   └── config_*.yaml
├── config/                       # 配置加载（viper 等）
├── cons/                         # 全局常量（如 auth_type、is_valid）
├── internal/                     # 内部代码（DDD 分层）
│   ├── application/              # 应用层：用例编排、DTO
│   │   └── identity/             # 身份限界上下文（登录等）
│   │       ├── login_service.go  # 登录服务定义与构造函数
│   │       ├── guest_login.go    # 游客登录用例
│   │       ├── wechat_login.go   # 微信登录用例
│   │       ├── ali_mobile_login.go
│   │       ├── send_mobile_code.go
│   │       ├── sms_verify_login.go
│   │       └── dto/              # 子包：请求/响应 DTO，按接口拆分
│   │           ├── login_common_dto.go
│   │           ├── guest_login_dto.go
│   │           ├── wechat_login_dto.go
│   │           ├── ali_mobile_login_dto.go
│   │           ├── send_mobile_code_dto.go
│   │           └── sms_verify_login_dto.go
│   ├── domain/                   # 领域层：聚合根、值对象、仓储接口
│   │   └── identity/             # 身份限界上下文
│   │       ├── app/              # 应用聚合子包
│   │       │   ├── app.go        # 聚合根 App、RestoreApp
│   │       │   ├── app_id.go     # 值对象 AppID
│   │       │   ├── app_repository.go
│   │       │   └── app_wechat_config.go  # 值对象（隶属 App）
│   │       └── user/             # 用户聚合子包
│   │           ├── user.go       # 聚合根 User、RestoreUser
│   │           ├── user_id.go
│   │           └── user_repository.go
│   └── infrastructure/           # 基础设施层：持久化实现
│       └── persistence/
│           └── mysql/
│               ├── model/        # 表对应 PO（biz_app、biz_user、biz_sms_code 等）
│               └── repository/   # 仓储实现（实现 domain 的接口）
│                   ├── bizApp_repository_impl.go
│                   ├── bizUser_repository_impl.go
│                   └── bizSmsCode_repository_impl.go
├── middleware/                   # 中间件（JWT、限流等）
├── res/                          # 资源初始化（DB 等）
├── tool/                         # 通用工具（微信/阿里/短信等封装）
└── validate/                     # 参数校验
```

### 分层与约定


| 层级           | 目录                             | 职责                                                                                       |
| -------------- | -------------------------------- | ------------------------------------------------------------------------------------------ |
| **用户接口层** | `api/`                           | 注册路由、解析请求、调用应用服务、返回响应                                                 |
| **应用层**     | `internal/application/identity/` | 用例编排（登录、发码、验码等），依赖领域接口与 DTO；根目录平铺用例，`dto/` 子包放请求/响应 |
| **领域层**     | `internal/domain/identity/`      | 聚合根、值对象、仓储接口，无基础设施依赖；按聚合划分子包`app/`、`user/`                    |
| **基础设施层** | `internal/infrastructure/`       | 仓储实现、PO、DB 访问，实现 domain 中定义的接口                                            |

### 依赖关系

- **api** → application/identity（调用 `identity.NewLoginService`、`loginService.GuestLogin` 等）
- **application/identity** → domain/identity/app、domain/identity/user、infrastructure/repository（依赖领域接口与仓储实现）
- **infrastructure** → domain/identity/app、domain/identity/user（实现 `app.AppRepository`、`user.UserRepository`）
- **domain** 不依赖 application 与 infrastructure

扩展新限界上下文时，在 `application/` 与 `domain/` 下新增同级模块目录（如 `order/`），并在 `api/` 中注册对应路由即可。
