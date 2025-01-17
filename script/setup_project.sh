#!/bin/bash

# 创建目录结构
mkdir -p cmd \
         config \
         deployment \
         doc \
         pkg/memcurve \
         pkg/auth \
         src/app/gw \
         src/app/static \
         src/core/handlers \
         src/core/review \
         src/core/reminder \
         src/core/analytic \
         src/core/interfaces \
         src/profile/handlers \
         src/profile/passport \
         src/profile/session \
         internal/repository \
         internal/utils \
         internal/tests/unit \
         internal/tests/integration

# 定义创建Go文件的函数
create_go_file() {
    local file=$1
    local package=$2
    local description=$3

    # 检查文件是否存在
    if [[ -e "$file" ]]; then
        echo "File $file exist, skip creating."
    else
        # 创建文件并添加内容
        cat <<EOF > "$file"
// Package $package provides $description.
//
// This file was generated by setup_project.sh script.

package $package

// TODO: Implement $description.
EOF
        echo "Created $file"
    fi
}

# Go文件列表
declare -A go_files=(
    [cmd/memorial_nexus.go]="main:The application entrypoint"
    [pkg/memcurve/calculator.go]="memcurve:Ebbinghaus forgetting curve calculations"
    [pkg/memcurve/curvemodel.go]="memcurve:Data models of the memory curve"
    [pkg/auth/jwt.go]="auth:JWT-based authentication logic"
    [pkg/auth/oauth.go]="auth:OAuth-based authentication logic"
    [src/app/gw/middleware.go]="gw:HTTP middleware components"
    [src/app/gw/error_handler.go]="gw:HTTP error handling"
    [src/app/gw/routes.go]="gw:Route definitions and mappings"
    [src/core/handlers.go]="handlers:API request handlers"
    [src/core/review/scheduler.go]="review:Scheduling logic for reviews"
    [src/core/review/session.go]="review:Review session management logic"
    [src/core/reminder/service.go]="reminder:Reminder service implementation"
    [src/core/reminder/types.go]="reminder:Types related to reminders"
    [src/core/analytic/reporter.go]="analytics:Analytics reporting logic"
    [src/core/analytic/types.go]="analytics:Types related to analytics"
    [src/core/interfaces/port.go]="interfaces:Port definitions for internal and external communication"
    [src/profile/handlers.go]="handlers:User profile handling logic"
    [src/profile/passport/model/account.go]="passport:Account management logic"
    [src/profile/passport/model/repo.go]="passport:Account management logic"
    [src/profile/passport/init.go]="passport:init the passport service"
    [src/profile/passport/register.go]="passport:Handler of account register"
    [src/profile/passport/login.go]="passport:Handler of user login"
    [src/profile/session/longterm.go]="session:Long-term session management logic"
    [src/profile/session/shortterm.go]="session:Short-term session management logic"
    [internal/repository/rds.go]="repository:ORM-related operations, adapting SQL type databases"
    [internal/repository/cache.go]="repository:Redis cache logic implementation"
    [internal/repository/dc.go]="repository:Distributed configuration center logic implementation"
    [internal/utils/utils.go]="utils:Utility functions for internal use"
    [internal/tests/unit/calculator_test.go]="unit:Unit tests for the memory curve calculator"
    [internal/tests/integration/api_test.go]="integration:Integration tests for API endpoints"
)

# 使用create_go_file函数创建Go文件
for file in "${!go_files[@]}"; do
    IFS=':' read -r package description <<< "${go_files[$file]}"
    create_go_file "$file" "$package" "$description"
done

# 创建其他重要的非Go文件列表
other_files=(
    config/app.dev.yaml
    config/app.prod.yaml
    config/log.dev.yaml
    config/log.prod.yaml
    deployment/Dockerfile
    deployment/db_migration.sh
    deployment/ci_cd.yaml
    doc/API_SPEC.md
    doc/DEPENDENCIES.md
)

# 检查其他重要文件
for file in "${other_files[@]}"; do
    if [[ -e "$file" ]]; then
      echo "File $file exist, skip creating."
    else
      touch "$file"
      echo "Created $file"
    fi
done

echo "Project structure and files creation process completed."