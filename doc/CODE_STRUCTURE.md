# Memoria Nexus - Code Structure

Memoria Nexus is designed with modularity and a clean architecture in mind. Below is the directory structure within the project, highlighting how our codebase is organized to separate concerns, making it easier to manage and extend.

## Directory Structure

```plaintext
.
├── cmd                             # 主应用入口点所在的目录
│   └── main.go                     # 主程序入口文件，负责应用启动和配置加载
├── config                          # 存放配置文件的目录
│   ├── app.dev.yaml                # 应用程序的主要配置文件 (开发环境)
│   ├── app.prod.yaml               # 应用程序的主要配置文件 (生产环境)
│   ├── log.dev.yaml                # 用于配置日志管理的配置文件 (开发环境)
│   └── log.prod.yaml               # 用于配置日志管理的配置文件 (生产环境)
├── deployment                      # 部署相关的脚本和配置文件
│   ├── Dockerfile                  # 用于构建Docker镜像的Dockerfile
│   ├── db_migration.sh             # 数据库迁移脚本
│   └── ci_cd.yaml                  # 持续集成和部署(CI/CD)的配置文件
├── doc                             # 存放文档的目录
│   ├── CODE_STRUCTURE.md           # 代码结构说明文档
│   ├── API_SPEC.md                 # API规格说明文档（OpenAPI/Swagger文件等）
│   └── DEPENDENCIES.md             # 项目依赖说明文档
├── pkg                             # 可由外部应用程序使用的库代码
│   ├── memcurve                    # 存储记忆曲线算法相关代码
│   │   ├── calculator.go           # 计算复习间隔的算法实现
│   │   └── curvemodel.go           # 记忆曲线模型的定义
│   └── auth                        # 认证库，用于用户认证和授权
│       ├── jwt.go                  # JWT认证实现
│       └── oauth.go                # OAuth认证实现
├── script                          # 各类脚本
├── src                             # 源代码目录
│   ├── app                         # 应用层代码
│   │   ├── gw                      # API接口层代码
│   │   │   ├── middleware.go       # API中间件
│   │   │   ├── error_handler.go    # API错误处理逻辑
│   │   │   └── routes.go           # 路由定义
│   │   └── static                  # 静态文件和资源目录
│   ├── core                        # 领域核心业务逻辑层
│   │   ├── handlers.go             # API请求处理函数
│   │   ├── review                  # 复习功能的领域逻辑
│   │   │   ├── scheduler.go        # 复习计划调度逻辑
│   │   │   └── session.go          # 复习会话逻辑
│   │   ├── reminder                # 提醒功能的领域逻辑
│   │   │   ├── service.go          # 提醒服务实现
│   │   │   └── types.go            # 提醒相关的类型定义
│   │   ├── analytics               # 分析功能的领域逻辑
│   │   │   ├── reporter.go         # 报告服务实现
│   │   │   └── types.go            # 分析相关的类型定义
│   │   └── interfaces              # 定义内部和外部通信的接口
│   │       └── port.go             # 实体与外部设备或系统通信接口
│   └── profile                     # 账户管理的代码
│       ├── handlers.go             # API请求处理函数
│       ├── passport                # 账户登录等
│       │   ├── account.go          # 账户管理逻辑
│       │   └── login.go            # 登录和登出相关逻辑
│       └── session                 # 连接
│           ├── longterm.go         # 长效 session 管理逻辑
│           └── shortterm.go        # 短效 session 管理逻辑
├── internal                        # 私有应用程序和库代码
│   ├── repository                  # 数据持久化层代码 (将在对应的应用层被初始化)
│   │   ├── rds.go                  # ORM层的逻辑，适配MySQL等SQL类型数据库
│   │   ├── cache.go                # Redis缓存逻辑
│   │   └── dc.go                   # 分布式配置中心逻辑
│   ├── utils                       # 内部使用的一些工具函数
│   │   └── utils.go                # 工具函数的实现
│   └── tests                       # 自动化测试代码
│       ├── unit                    # 单元测试
│       │   └── calculator_test.go  # 记忆曲线计算器的单元测试
│       └── integration             # 集成测试
│           └── api_test.go         # API接口的集成测试
├── .gitignore                      # Git忽略文件
├── LICENSE                         # 项目的许可证文件
└── README.md                       # 项目说明文件
```