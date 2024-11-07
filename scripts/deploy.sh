#!/bin/bash

# 设置错误时退出
set -e

# 定义颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查必要的命令是否存在
check_requirements() {
    local requirements=("go" "git" "supervisorctl")
    
    for cmd in "${requirements[@]}"; do
        if ! command -v "$cmd" &> /dev/null; then
            log_error "$cmd is required but not installed."
            exit 1
        fi
    done
}

# 拉取最新代码
fetch_latest_code() {
    log_info "Fetching latest code..."
    
    # 检查是否有未提交的更改
    if ! git diff --quiet HEAD; then
        log_error "There are uncommitted changes. Please commit or stash them first."
        exit 1
    fi
    
    # 获取当前分支
    local current_branch=$(git symbolic-ref --short HEAD)
    
    # 拉取最新代码
    git fetch origin
    
    # 重置到远程分支最新状态
    git reset --hard origin/$current_branch
    
    log_info "Successfully updated to latest code on branch $current_branch"
}

# 检查环境变量
check_env() {
    if [ -z "$APP_ENV" ]; then
        export APP_ENV="development"
        log_warn "APP_ENV not set, using default: development"
    fi
    
    if [ -z "$APP_PORT" ]; then
        export APP_PORT="8080"
        log_warn "APP_PORT not set, using default: 8080"
    fi

    # 设置应用路径
    if [ -z "$APP_PATH" ]; then
        export APP_PATH="/path/to/your/app"  # 请修改为实际的应用路径
        log_warn "APP_PATH not set, using default: $APP_PATH"
    fi
}

# 构建应用
build_app() {
    log_info "Building application..."
    
    # 清理之前的构建
    go clean
    
    # 运行测试
    log_info "Running tests..."
    go test ./... || {
        log_error "Tests failed"
        exit 1
    }
    
    # 构建应用
    CGO_ENABLED=0 go build -o bin/app cmd/api/main.go
    
    if [ $? -eq 0 ]; then
        log_info "Build successful"
    else
        log_error "Build failed"
        exit 1
    fi
}

# 更新supervisor配置
update_supervisor_conf() {
    local program_name="your-app"  # 请修改为实际的程序名称
    local supervisor_conf="/etc/supervisor/conf.d/${program_name}.conf"
    
    log_info "Updating supervisor configuration..."
    
    # 检查supervisor配置文件是否存在
    if [ ! -f "$supervisor_conf" ]; then
        log_info "Creating new supervisor configuration..."
        sudo tee "$supervisor_conf" > /dev/null <<EOF
[program:${program_name}]
directory=${APP_PATH}
command=${APP_PATH}/bin/app
autostart=true
autorestart=true
stderr_logfile=/var/log/${program_name}.err.log
stdout_logfile=/var/log/${program_name}.out.log
environment=APP_ENV="${APP_ENV}",APP_PORT="${APP_PORT}"
user=www-data
EOF
    fi
}

# 部署应用
deploy() {
    local env="$APP_ENV"
    
    log_info "Deploying to $env environment..."
    
    # 停止服务
    log_info "Stopping service..."
    sudo supervisorctl stop your-app || true
    
    # 备份当前版本
    if [ -f "${APP_PATH}/bin/app" ]; then
        log_info "Backing up current version..."
        cp "${APP_PATH}/bin/app" "${APP_PATH}/bin/app.backup"
    fi
    
    # 复制新的二进制文件
    log_info "Copying new binary..."
    cp bin/app "${APP_PATH}/bin/app"
    chmod +x "${APP_PATH}/bin/app"
    
    # 更新supervisor配置
    update_supervisor_conf
    
    # 重新加载supervisor配置
    log_info "Reloading supervisor configuration..."
    sudo supervisorctl reread
    sudo supervisorctl update
    
    # 启动服务
    log_info "Starting service..."
    sudo supervisorctl start your-app
    
    # 检查服务状态
    log_info "Checking service status..."
    sudo supervisorctl status your-app
}

# 清理函数
cleanup() {
    log_info "Cleaning up..."
    # 如果部署失败，恢复备份
    if [ $? -ne 0 ] && [ -f "${APP_PATH}/bin/app.backup" ]; then
        log_warn "Deployment failed, restoring backup..."
        mv "${APP_PATH}/bin/app.backup" "${APP_PATH}/bin/app"
        sudo supervisorctl restart your-app
    fi
    
    # 删除备份文件
    rm -f "${APP_PATH}/bin/app.backup"
}

# 主函数
main() {
    # 检查是否具有sudo权限
    if ! sudo -v; then
        log_error "Sudo privileges are required for deployment"
        exit 1
    fi
    
    # 检查依赖
    check_requirements
    
    # 拉取最新代码
    fetch_latest_code
    
    # 检查环境变量
    check_env
    
    # 构建应用
    build_app
    
    # 部署应用
    deploy
    
    log_info "Deployment completed successfully!"
}

# 捕获错误并清理
trap cleanup EXIT

# 运行主函数
main "$@"
