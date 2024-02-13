.
├── cmd
│   └── main.go
├── config
│   ├── app.dev.yaml
│   ├── app.prod.yaml
│   ├── log.dev.yaml
│   └── log.prod.yaml
├── deployment
│   ├── migration
│   │   ├── 01.passport_down.sql
│   │   ├── 01.passport_up.sql
│   │   ├── 02.profile_down.sql
│   │   ├── 02.profile_up.sql
│   │   └── migrate.sh
│   ├── Dockerfile
│   └── ci_cd.yaml
├── development
│   ├── dial
│   │   ├── auth.http
│   │   ├── http-client.env.json
│   │   └── profile.http
│   └── docker-compose.yml
├── doc
│   ├── API_SPEC.md
│   ├── CODE_STRUCTURE.md
│   ├── CONTRIBUTING.md
│   ├── DEPENDENCIES.md
│   └── PROJECT_STRUCTURE.md
├── internal
│   ├── repository
│   │   ├── cache.go
│   │   ├── dc.go
│   │   └── rds.go
│   ├── tests
│   │   ├── integration
│   │   │   └── api_test.go
│   │   └── unit
│   │       └── calculator_test.go
│   └── util
│       └── util.go
├── pkg
│   ├── auth
│   │   ├── jwt.go
│   │   ├── jwtmw.go
│   │   └── oauth.go
│   └── memcurve
│       ├── calculator.go
│       └── curvemodel.go
├── script
│   └── setup_project.sh
├── src
│   ├── app
│   │   └── gw
│   │       ├── error_handler.go
│   │       └── routes.go
│   ├── core
│   │   ├── analytics
│   │   │   ├── reporter.go
│   │   │   └── types.go
│   │   ├── interfaces
│   │   │   └── port.go
│   │   ├── reminder
│   │   │   ├── service.go
│   │   │   └── types.go
│   │   ├── review
│   │   │   ├── scheduler.go
│   │   │   └── session.go
│   │   └── handlers.go
│   ├── iam
│   │   ├── passport
│   │   │   ├── model
│   │   │   │   ├── repo.go
│   │   │   │   └── user.go
│   │   │   ├── init.go
│   │   │   ├── login.go
│   │   │   ├── register.go
│   │   │   └── util.go
│   │   ├── session
│   │   │   ├── longterm.go
│   │   │   └── shortterm.go
│   │   └── handlers.go
│   └── profile
│       ├── model
│       │   ├── profile.go
│       │   └── repo.go
│       ├── PROFILE_CLASSES.mermaid
│       ├── basic.go
│       └── init.go
├── LICENSE
├── MODULES.puml
├── Makefile
├── README.md
├── go.mod
└── go.sum
