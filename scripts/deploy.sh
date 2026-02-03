#!/bin/bash

# 设置错误时退出
set -e

# 定义颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

log_debug() {
    echo -e "${BLUE}[DEBUG]${NC} $1"
}

# 配置变量（可通过环境变量覆盖）
SERVICES=("api" "admin")  # 支持的服务列表
DEFAULT_SERVICE="api"     # 默认部署的服务
APP_PATH="${APP_PATH:-/opt/kratos-kit}"  # 应用部署路径
SUPERVISOR_CONF_DIR="${SUPERVISOR_CONF_DIR:-/etc/supervisor/conf.d}"  # Supervisor配置目录
LOG_DIR="${LOG_DIR:-/var/log/kratos-kit}"  # 日志目录
BACKUP_DIR="${BACKUP_DIR:-${APP_PATH}/backups}"  # 备份目录
HEALTH_CHECK_URL="${HEALTH_CHECK_URL:-http://localhost:9000/health}"  # 健康检查URL
HEALTH_CHECK_TIMEOUT="${HEALTH_CHECK_TIMEOUT:-30}"  # 健康检查超时时间（秒）
MAX_BACKUPS="${MAX_BACKUPS:-5}"  # 保留的最大备份数量
DEPLOY_USER="${DEPLOY_USER:-www-data}"  # 运行用户

# 服务配置映射
declare -A SERVICE_PORTS
SERVICE_PORTS["api"]="9000"
SERVICE_PORTS["admin"]="9001"

declare -A SERVICE_CMDS
SERVICE_CMDS["api"]="cmd/api/main.go"
SERVICE_CMDS["admin"]="cmd/admin/main.go"

declare -A SERVICE_CONFIGS
SERVICE_CONFIGS["api"]="configs/config.yaml"
SERVICE_CONFIGS["admin"]="configs/config_admin.local.yaml"

# 检测是否在 WSL 环境
is_wsl() {
    if [ -f /proc/version ]; then
        grep -qi microsoft /proc/version && return 0
    fi
    [ -n "$WSL_DISTRO_NAME" ] && return 0
    return 1
}

# 显示安装指导
show_install_guide() {
    local cmd=$1
    local is_wsl_env=$(is_wsl && echo "true" || echo "false")
    
    log_error "$cmd is required but not installed."
    echo ""
    
    case $cmd in
        go)
            log_info "To install Go:"
            if [ "$is_wsl_env" = "true" ]; then
                echo "  Ubuntu/Debian: sudo apt update && sudo apt install -y golang-go"
                echo "  Or download from: https://go.dev/dl/"
                echo "  After installation, add to PATH: export PATH=\$PATH:/usr/local/go/bin"
            else
                echo "  Download from: https://go.dev/dl/"
                echo "  Or use package manager:"
                echo "    Ubuntu/Debian: sudo apt install golang-go"
                echo "    CentOS/RHEL: sudo yum install golang"
                echo "    macOS: brew install go"
            fi
            ;;
        git)
            log_info "To install Git:"
            if [ "$is_wsl_env" = "true" ]; then
                echo "  sudo apt update && sudo apt install -y git"
            else
                echo "  Ubuntu/Debian: sudo apt install git"
                echo "  CentOS/RHEL: sudo yum install git"
                echo "  macOS: brew install git"
            fi
            ;;
        supervisorctl)
            log_info "To install Supervisor:"
            if [ "$is_wsl_env" = "true" ]; then
                echo "  sudo apt update && sudo apt install -y supervisor"
            else
                echo "  Ubuntu/Debian: sudo apt install supervisor"
                echo "  CentOS/RHEL: sudo yum install supervisor"
            fi
            echo "  Start service: sudo systemctl start supervisor"
            echo "  Enable service: sudo systemctl enable supervisor"
            ;;
        curl)
            log_info "To install curl:"
            if [ "$is_wsl_env" = "true" ]; then
                echo "  sudo apt update && sudo apt install -y curl"
            else
                echo "  Ubuntu/Debian: sudo apt install curl"
                echo "  CentOS/RHEL: sudo yum install curl"
            fi
            ;;
    esac
    echo ""
}

# 检查必要的命令是否存在
check_requirements() {
    local requirements=("go" "git" "supervisorctl" "curl")
    local missing=()
    
    log_info "Checking requirements..."
    
    # 检测环境
    if is_wsl; then
        log_info "Detected WSL environment"
    fi
    
    # 检查每个依赖
    for cmd in "${requirements[@]}"; do
        if ! command -v "$cmd" &> /dev/null; then
            missing+=("$cmd")
        else
            # 显示版本信息并验证
            case $cmd in
                go)
                    local go_version=$(go version 2>/dev/null | awk '{print $3}')
                    log_debug "Found $cmd: $go_version"
                    
                    # 检查 Go 版本（需要 1.22+）
                    local go_major=$(go version 2>/dev/null | sed -n 's/.*go\([0-9]*\)\.\([0-9]*\).*/\1/p')
                    local go_minor=$(go version 2>/dev/null | sed -n 's/.*go\([0-9]*\)\.\([0-9]*\).*/\2/p')
                    
                    if [ -n "$go_major" ] && [ -n "$go_minor" ]; then
                        if [ "$go_major" -lt 1 ] || ([ "$go_major" -eq 1 ] && [ "$go_minor" -lt 22 ]); then
                            log_warn "Go version $go_version is too old. This project requires Go 1.22+"
                            log_info "Please upgrade Go: https://go.dev/dl/"
                        else
                            log_debug "Go version check passed"
                        fi
                    fi
                    ;;
                git)
                    local git_version=$(git --version 2>/dev/null)
                    log_debug "Found $cmd: $git_version"
                    ;;
                supervisorctl)
                    # 检查 Supervisor 服务是否运行
                    if ! systemctl is-active --quiet supervisor 2>/dev/null && ! service supervisor status >/dev/null 2>&1; then
                        log_warn "Supervisor service may not be running"
                        log_info "Start it with: sudo systemctl start supervisor (systemd) or sudo service supervisor start (sysvinit)"
                    fi
                    log_debug "Found $cmd"
                    ;;
                *)
                    log_debug "Found $cmd"
                    ;;
            esac
        fi
    done
    
    # 如果有缺失的依赖，显示错误和安装指导
    if [ ${#missing[@]} -gt 0 ]; then
        log_error "Missing required dependencies: ${missing[*]}"
        echo ""
        for cmd in "${missing[@]}"; do
            show_install_guide "$cmd"
        done
        exit 1
    fi
    
    log_info "All requirements satisfied"
}

# 验证服务名称
validate_service() {
    local service=$1
    if [[ ! " ${SERVICES[@]} " =~ " ${service} " ]]; then
        log_error "Invalid service: $service. Available services: ${SERVICES[*]}"
        exit 1
    fi
}

# 拉取最新代码
fetch_latest_code() {
    local skip_fetch=${1:-false}
    
    if [ "$skip_fetch" = "true" ]; then
        log_warn "Skipping code fetch (--skip-fetch flag set)"
        return 0
    fi
    
    log_info "Fetching latest code..."
    
    # 检查是否有未提交的更改
    if ! git diff --quiet HEAD 2>/dev/null; then
        log_warn "There are uncommitted changes."
        # 检查是否在交互式终端
        if [ -t 0 ]; then
            read -p "Do you want to stash them? (y/n) " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                git stash push -m "Auto-stashed before deployment $(date +%Y%m%d_%H%M%S)"
                log_info "Changes stashed"
            else
                log_error "Please commit or stash changes first."
                exit 1
            fi
        else
            # 非交互式环境，自动stash
            log_info "Non-interactive terminal, auto-stashing changes..."
            git stash push -m "Auto-stashed before deployment $(date +%Y%m%d_%H%M%S)"
            log_info "Changes stashed"
        fi
    fi
    
    # 获取当前分支
    local current_branch=$(git symbolic-ref --short HEAD 2>/dev/null || echo "main")
    
    # 拉取最新代码
    git fetch origin || {
        log_error "Failed to fetch from origin"
        exit 1
    }
    
    # 检查是否有更新
    local local_commit=$(git rev-parse HEAD)
    local remote_commit=$(git rev-parse origin/$current_branch 2>/dev/null || echo "")
    
    if [ -z "$remote_commit" ]; then
        log_warn "Remote branch origin/$current_branch not found, using current HEAD"
        return 0
    fi
    
    if [ "$local_commit" = "$remote_commit" ]; then
        log_info "Already up to date with origin/$current_branch"
        return 0
    fi
    
    # 重置到远程分支最新状态
    log_info "Updating to latest code on branch $current_branch..."
    git reset --hard origin/$current_branch || {
        log_error "Failed to reset to origin/$current_branch"
        exit 1
    }
    
    log_info "Successfully updated to latest code on branch $current_branch"
}

# 检查环境变量
check_env() {
    local service=$1
    
    if [ -z "$APP_ENV" ]; then
        export APP_ENV="production"
        log_warn "APP_ENV not set, using default: production"
    fi
    
    # 根据服务设置端口
    if [ -z "$APP_PORT" ]; then
        export APP_PORT="${SERVICE_PORTS[$service]}"
        log_info "APP_PORT not set, using service default: $APP_PORT"
    fi

    # 创建必要的目录
    mkdir -p "$APP_PATH/bin"
    mkdir -p "$BACKUP_DIR"
    mkdir -p "$LOG_DIR"
    
    log_info "Environment check completed"
    log_debug "APP_ENV: $APP_ENV"
    log_debug "APP_PORT: $APP_PORT"
    log_debug "APP_PATH: $APP_PATH"
}

# 构建应用
build_app() {
    local service=$1
    local skip_tests=${2:-false}
    local binary_name="${service}"
    local cmd_path="${SERVICE_CMDS[$service]}"
    
    if [ -z "$cmd_path" ]; then
        log_error "No command path configured for service: $service"
        exit 1
    fi
    
    log_info "Building $service service..."
    
    # 清理之前的构建
    go clean -cache
    
    # 运行测试（可选）
    if [ "$skip_tests" != "true" ]; then
        log_info "Running tests..."
        if ! go test ./... -short; then
            log_warn "Some tests failed, but continuing build..."
        fi
    else
        log_warn "Skipping tests (--skip-tests flag set)"
    fi
    
    # 获取版本信息
    local version=$(git describe --tags --always 2>/dev/null || echo "unknown")
    local commit=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    local build_time=$(date +%Y-%m-%d_%H:%M:%S)
    
    # 构建应用
    log_info "Compiling $service binary..."
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
        -ldflags "-s -w -X main.Version=$version -X main.CommitSHA=$commit -X main.BuildTime=$build_time" \
        -o "bin/${binary_name}" \
        "$cmd_path"
    
    if [ $? -eq 0 ]; then
        log_info "Build successful: bin/${binary_name}"
        log_debug "Version: $version"
        log_debug "Commit: $commit"
        log_debug "Build time: $build_time"
    else
        log_error "Build failed"
        exit 1
    fi
}

# 更新supervisor配置
update_supervisor_conf() {
    local service=$1
    local program_name="kratos-${service}"
    local supervisor_conf="${SUPERVISOR_CONF_DIR}/${program_name}.conf"
    local binary_path="${APP_PATH}/bin/${service}"
    local config_path="${APP_PATH}/${SERVICE_CONFIGS[$service]}"
    local port="${SERVICE_PORTS[$service]}"
    
    log_info "Updating supervisor configuration for $service..."
    
    # 创建Supervisor配置目录（如果不存在）
    sudo mkdir -p "$SUPERVISOR_CONF_DIR"
    
    # 生成Supervisor配置
    log_info "Creating/updating supervisor configuration..."
    sudo tee "$supervisor_conf" > /dev/null <<EOF
[program:${program_name}]
directory=${APP_PATH}
command=${binary_path} -f ${config_path}
autostart=true
autorestart=true
startretries=3
startsecs=5
stopwaitsecs=10
stderr_logfile=${LOG_DIR}/${service}.err.log
stdout_logfile=${LOG_DIR}/${service}.out.log
logfile_maxbytes=50MB
logfile_backups=10
environment=APP_ENV="${APP_ENV}",APP_PORT="${port}"
user=${DEPLOY_USER}
priority=999
EOF
    
    log_info "Supervisor configuration updated: $supervisor_conf"
}

# 健康检查
health_check() {
    local service=$1
    local port="${SERVICE_PORTS[$service]}"
    local url="http://localhost:${port}/health"
    local max_attempts=10
    local attempt=0
    local interval=3
    
    log_info "Performing health check for $service on $url..."
    
    while [ $attempt -lt $max_attempts ]; do
        if curl -f -s -m 5 "$url" > /dev/null 2>&1; then
            log_info "Health check passed for $service"
            return 0
        fi
        
        attempt=$((attempt + 1))
        if [ $attempt -lt $max_attempts ]; then
            log_debug "Health check attempt $attempt/$max_attempts failed, retrying in ${interval}s..."
            sleep $interval
        fi
    done
    
    log_error "Health check failed for $service after $max_attempts attempts"
    return 1
}

# 备份当前版本
backup_current_version() {
    local service=$1
    local binary_path="${APP_PATH}/bin/${service}"
    local timestamp=$(date +%Y%m%d_%H%M%S)
    local backup_file="${BACKUP_DIR}/${service}_${timestamp}"
    
    if [ -f "$binary_path" ]; then
        log_info "Backing up current version..."
        mkdir -p "$BACKUP_DIR"
        cp "$binary_path" "$backup_file"
        log_info "Backup created: $backup_file"
        
        # 清理旧备份（只保留最新的N个）
        local backup_count=$(ls -1t "${BACKUP_DIR}/${service}_"* 2>/dev/null | wc -l)
        if [ "$backup_count" -gt "$MAX_BACKUPS" ]; then
            log_info "Cleaning up old backups (keeping latest $MAX_BACKUPS)..."
            ls -1t "${BACKUP_DIR}/${service}_"* 2>/dev/null | tail -n +$((MAX_BACKUPS + 1)) | xargs rm -f
        fi
        
        echo "$backup_file"  # 返回备份文件路径
    else
        log_warn "No existing binary to backup"
        echo ""
    fi
}

# 部署应用
deploy() {
    local service=$1
    local env="$APP_ENV"
    local program_name="kratos-${service}"
    local binary_name="${service}"
    local source_binary="bin/${binary_name}"
    local target_binary="${APP_PATH}/bin/${binary_name}"
    
    log_info "Deploying $service to $env environment..."
    
    # 检查源二进制文件是否存在
    if [ ! -f "$source_binary" ]; then
        log_error "Source binary not found: $source_binary"
        log_error "Please run build first"
        exit 1
    fi
    
    # 备份当前版本
    local backup_file=$(backup_current_version "$service")
    
    # 停止服务
    log_info "Stopping service $service..."
    if sudo supervisorctl status "$program_name" &>/dev/null; then
        sudo supervisorctl stop "$program_name" || {
            log_warn "Failed to stop service (may not be running)"
        }
    else
        log_info "Service not running, skipping stop"
    fi
    
    # 复制新的二进制文件
    log_info "Installing new binary..."
    sudo cp "$source_binary" "$target_binary"
    sudo chmod +x "$target_binary"
    sudo chown "$DEPLOY_USER:$DEPLOY_USER" "$target_binary"
    
    # 复制配置文件（如果存在）
    local config_source="${SERVICE_CONFIGS[$service]}"
    if [ -f "$config_source" ]; then
        log_info "Copying configuration file..."
        sudo mkdir -p "$APP_PATH/configs"
        sudo cp "$config_source" "$APP_PATH/configs/"
    fi
    
    # 更新supervisor配置
    update_supervisor_conf "$service"
    
    # 重新加载supervisor配置
    log_info "Reloading supervisor configuration..."
    sudo supervisorctl reread
    sudo supervisorctl update
    
    # 启动服务
    log_info "Starting service $service..."
    sudo supervisorctl start "$program_name" || {
        log_error "Failed to start service"
        # 如果启动失败，尝试恢复备份
        if [ -n "$backup_file" ] && [ -f "$backup_file" ]; then
            log_warn "Restoring backup..."
            sudo cp "$backup_file" "$target_binary"
            sudo supervisorctl start "$program_name"
        fi
        exit 1
    }
    
    # 等待服务启动
    sleep 3
    
    # 检查服务状态
    log_info "Checking service status..."
    sudo supervisorctl status "$program_name"
    
    # 健康检查
    if ! health_check "$service"; then
        log_error "Health check failed after deployment"
        # 如果健康检查失败，可以选择回滚
        if [ -n "$backup_file" ] && [ -f "$backup_file" ]; then
            log_warn "Rolling back due to health check failure..."
            sudo supervisorctl stop "$program_name"
            sudo cp "$backup_file" "$target_binary"
            sudo supervisorctl start "$program_name"
            log_info "Rolled back to previous version"
        fi
        exit 1
    fi
    
    log_info "Deployment completed successfully for $service"
}

# 显示使用说明
show_usage() {
    cat <<EOF
Usage: $0 [OPTIONS] [SERVICE]

Deploy Kratos services using Supervisor

Arguments:
  SERVICE              Service to deploy (api|admin). Default: $DEFAULT_SERVICE

Options:
  -h, --help           Show this help message
  -s, --skip-fetch     Skip fetching latest code from git
  -t, --skip-tests     Skip running tests before build
  -a, --all            Deploy all services
  --list               List available services
  --status             Show status of all services
  --rollback SERVICE   Rollback service to previous version

Environment Variables:
  APP_PATH             Application deployment path (default: $APP_PATH)
  APP_ENV              Environment (default: production)
  APP_PORT             Application port (auto-detected by service)
  SUPERVISOR_CONF_DIR  Supervisor config directory (default: $SUPERVISOR_CONF_DIR)
  LOG_DIR              Log directory (default: $LOG_DIR)
  BACKUP_DIR           Backup directory (default: $BACKUP_DIR)
  DEPLOY_USER          User to run service (default: $DEPLOY_USER)
  MAX_BACKUPS          Maximum number of backups to keep (default: $MAX_BACKUPS)

Examples:
  $0 api                # Deploy api service
  $0 --all              # Deploy all services
  $0 --status           # Show service status
  $0 --rollback api     # Rollback api service

EOF
}

# 列出服务状态
show_status() {
    log_info "Service Status:"
    echo ""
    for service in "${SERVICES[@]}"; do
        local program_name="kratos-${service}"
        echo "=== $service ==="
        if sudo supervisorctl status "$program_name" &>/dev/null; then
            sudo supervisorctl status "$program_name"
        else
            echo "  Not configured in supervisor"
        fi
        echo ""
    done
}

# 回滚服务
rollback_service() {
    local service=$1
    local program_name="kratos-${service}"
    local target_binary="${APP_PATH}/bin/${service}"
    
    log_info "Rolling back $service..."
    
    # 查找最新的备份
    local latest_backup=$(ls -1t "${BACKUP_DIR}/${service}_"* 2>/dev/null | head -n 1)
    
    if [ -z "$latest_backup" ] || [ ! -f "$latest_backup" ]; then
        log_error "No backup found for $service"
        exit 1
    fi
    
    log_info "Found backup: $latest_backup"
    
    # 停止服务
    log_info "Stopping service..."
    sudo supervisorctl stop "$program_name" || true
    
    # 恢复备份
    log_info "Restoring backup..."
    sudo cp "$latest_backup" "$target_binary"
    sudo chmod +x "$target_binary"
    sudo chown "$DEPLOY_USER:$DEPLOY_USER" "$target_binary"
    
    # 启动服务
    log_info "Starting service..."
    sudo supervisorctl start "$program_name"
    
    # 等待并检查
    sleep 3
    sudo supervisorctl status "$program_name"
    
    if health_check "$service"; then
        log_info "Rollback completed successfully"
    else
        log_warn "Rollback completed but health check failed"
    fi
}

# 主函数
main() {
    local service=""
    local skip_fetch=false
    local skip_tests=false
    local deploy_all=false
    local show_status_only=false
    local rollback_service_name=""
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_usage
                exit 0
                ;;
            -s|--skip-fetch)
                skip_fetch=true
                shift
                ;;
            -t|--skip-tests)
                skip_tests=true
                shift
                ;;
            -a|--all)
                deploy_all=true
                shift
                ;;
            --list)
                echo "Available services: ${SERVICES[*]}"
                exit 0
                ;;
            --status)
                show_status_only=true
                shift
                ;;
            --rollback)
                if [ -z "$2" ]; then
                    log_error "Service name required for rollback"
                    show_usage
                    exit 1
                fi
                rollback_service_name="$2"
                shift 2
                ;;
            api|admin)
                if [ -z "$service" ]; then
                    service="$1"
                else
                    log_error "Multiple services specified. Use --all to deploy all services."
                    exit 1
                fi
                shift
                ;;
            *)
                log_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done
    
    # 处理特殊命令
    if [ "$show_status_only" = "true" ]; then
        show_status
        exit 0
    fi
    
    if [ -n "$rollback_service_name" ]; then
        validate_service "$rollback_service_name"
        log_info "Checking sudo privileges..."
        if ! sudo -v 2>/dev/null; then
            log_error "Sudo privileges are required"
            exit 1
        fi
        rollback_service "$rollback_service_name"
        exit 0
    fi
    
    # 确定要部署的服务
    if [ "$deploy_all" = "true" ]; then
        service="all"
    elif [ -z "$service" ]; then
        service="$DEFAULT_SERVICE"
        log_info "No service specified, using default: $service"
    fi
    
    # 验证服务
    if [ "$service" != "all" ]; then
        validate_service "$service"
    fi
    
    # 检查是否具有sudo权限
    log_info "Checking sudo privileges..."
    if ! sudo -v 2>/dev/null; then
        log_error "Sudo privileges are required for deployment"
        echo ""
        log_info "Please ensure:"
        echo "  1. You have sudo access (run: sudo -v to test)"
        echo "  2. Your user is in the sudoers group"
        if is_wsl; then
            echo "  3. In WSL, you may need to set up sudo password"
            echo "     Run: sudo passwd root (if needed)"
        fi
        exit 1
    fi
    log_info "Sudo privileges confirmed"
    
    # 检查依赖
    check_requirements
    
    # 拉取最新代码
    fetch_latest_code "$skip_fetch"
    
    # 部署服务
    if [ "$service" = "all" ]; then
        log_info "Deploying all services..."
        for svc in "${SERVICES[@]}"; do
            log_info "========================================="
            log_info "Deploying service: $svc"
            log_info "========================================="
            
            check_env "$svc"
            build_app "$svc" "$skip_tests"
            deploy "$svc"
            
            log_info "Service $svc deployed successfully"
            echo ""
        done
        log_info "All services deployed successfully!"
    else
        check_env "$service"
        build_app "$service" "$skip_tests"
        deploy "$service"
        log_info "Deployment completed successfully!"
    fi
}

# 运行主函数
main "$@"
