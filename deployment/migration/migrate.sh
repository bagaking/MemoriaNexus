#!/bin/bash
# File: deployment/migration/migrate.sh

set -e

# 使用Docker Compose设置的环境变量
DB_USER="root"
DB_PASSWORD="${MYSQL_ROOT_PASSWORD:-password}"
DB_NAME="${MYSQL_DATABASE:-memorianexus}"

# SQL文件目录
MIGRATIONS_PATH="/migration"

MYSQL_CMD="mysql -u $DB_USER -p$DB_PASSWORD $DB_NAME"

# 函数：执行向上迁移
run_migration_up() {
    for file in `ls $MIGRATIONS_PATH/*_up.sql | sort -V`; do
        echo "Applying migration: $file"
        $MYSQL_CMD < "$file"
    done
}

# 函数：执行向下迁移
run_migration_down() {
    for file in `ls $MIGRATIONS_PATH/*_down.sql | sort -Vr`; do
        echo "Reverting migration: $file"
        $MYSQL_CMD < "$file"
    done
}

# 主命令处理
command=$1

case $command in
    up)
        echo "Running migrations up..."
        run_migration_up
        ;;
    down)
        echo "Running migrations down..."
        run_migration_down
        ;;
    *)
        echo "Unknown command: $command"
        exit 1
        ;;
esac