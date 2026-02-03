# 部署脚本使用指南

## 概述

`deploy.sh` 是一个自动化部署脚本，用于将 Kratos 微服务部署到生产环境。它支持多服务部署、自动备份、健康检查和回滚功能。

## 功能特性

- ✅ **多服务支持**: 支持部署 `api` 和 `admin` 服务
- ✅ **自动备份**: 部署前自动备份当前版本
- ✅ **健康检查**: 部署后自动进行健康检查
- ✅ **自动回滚**: 健康检查失败时自动回滚
- ✅ **版本管理**: 保留多个历史版本备份
- ✅ **Supervisor集成**: 自动配置和管理 Supervisor
- ✅ **灵活配置**: 支持环境变量和命令行参数

## 快速开始

### 基本用法

```bash
# 部署默认服务 (api)
./scripts/deploy.sh

# 部署指定服务
./scripts/deploy.sh api
./scripts/deploy.sh admin

# 部署所有服务
./scripts/deploy.sh --all
```

### 常用选项

```bash
# 跳过代码拉取（使用当前代码）
./scripts/deploy.sh --skip-fetch api

# 跳过测试（快速部署）
./scripts/deploy.sh --skip-tests api

# 查看服务状态
./scripts/deploy.sh --status

# 回滚服务
./scripts/deploy.sh --rollback api

# 查看帮助
./scripts/deploy.sh --help
```

## 环境变量配置

脚本支持通过环境变量进行配置：

```bash
# 应用部署路径
export APP_PATH="/opt/kratos-kit"

# 环境类型
export APP_ENV="production"

# Supervisor配置目录
export SUPERVISOR_CONF_DIR="/etc/supervisor/conf.d"

# 日志目录
export LOG_DIR="/var/log/kratos-kit"

# 备份目录
export BACKUP_DIR="/opt/kratos-kit/backups"

# 运行用户
export DEPLOY_USER="www-data"

# 最大备份数量
export MAX_BACKUPS=5
```

## 部署流程

1. **检查依赖**: 验证必要的命令（go, git, supervisorctl, curl）
2. **拉取代码**: 从 Git 仓库拉取最新代码（可选跳过）
3. **环境检查**: 检查并设置环境变量
4. **运行测试**: 运行单元测试（可选跳过）
5. **构建应用**: 编译 Go 二进制文件
6. **备份当前版本**: 备份现有的二进制文件
7. **停止服务**: 停止当前运行的服务
8. **安装新版本**: 复制新的二进制文件
9. **更新配置**: 更新 Supervisor 配置
10. **启动服务**: 启动新版本服务
11. **健康检查**: 验证服务是否正常运行
12. **自动回滚**: 如果健康检查失败，自动回滚到上一个版本

## 服务配置

脚本通过关联数组配置每个服务的参数：

- **端口映射**: `SERVICE_PORTS`
  - `api`: 9000
  - `admin`: 9001

- **命令路径**: `SERVICE_CMDS`
  - `api`: `cmd/api/main.go`
  - `admin`: `cmd/admin/main.go`

- **配置文件**: `SERVICE_CONFIGS`
  - `api`: `configs/config.yaml`
  - `admin`: `configs/config_admin.local.yaml`

## Supervisor 配置

脚本会自动生成 Supervisor 配置文件，位置：`/etc/supervisor/conf.d/kratos-{service}.conf`

配置示例：
```ini
[program:kratos-api]
directory=/opt/kratos-kit
command=/opt/kratos-kit/bin/api -f /opt/kratos-kit/configs/config.yaml
autostart=true
autorestart=true
startretries=3
startsecs=5
stderr_logfile=/var/log/kratos-kit/api.err.log
stdout_logfile=/var/log/kratos-kit/api.out.log
environment=APP_ENV="production",APP_PORT="9000"
user=www-data
```

## 备份管理

- 备份文件命名格式: `{service}_{timestamp}`
- 备份位置: `$BACKUP_DIR`
- 自动清理: 只保留最新的 N 个备份（默认 5 个）

## 健康检查

部署后会自动进行健康检查：
- 检查 URL: `http://localhost:{port}/health`
- 最大重试次数: 10 次
- 重试间隔: 3 秒
- 超时时间: 5 秒

如果健康检查失败，脚本会自动回滚到上一个版本。

## 回滚操作

### 自动回滚
部署过程中如果健康检查失败，会自动回滚。

### 手动回滚
```bash
./scripts/deploy.sh --rollback api
```

回滚流程：
1. 查找最新的备份文件
2. 停止当前服务
3. 恢复备份文件
4. 启动服务
5. 执行健康检查

## WSL 环境支持

脚本自动检测 WSL 环境并提供相应的安装指导。

### WSL 环境准备

1. **安装 Go**
```bash
# 方法1: 使用包管理器（可能版本较旧）
sudo apt update
sudo apt install -y golang-go

# 方法2: 下载最新版本（推荐）
wget https://go.dev/dl/go1.22.3.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
```

2. **安装 Supervisor**
```bash
sudo apt update
sudo apt install -y supervisor
sudo systemctl start supervisor
sudo systemctl enable supervisor
```

3. **配置 Sudo 权限**
```bash
# 如果 sudo 需要密码，确保已配置
sudo passwd root  # 设置 root 密码（如果需要）
```

4. **验证安装**
```bash
go version      # 应该显示 Go 1.22+
git --version
supervisorctl --version
curl --version
```

## 故障排查

### 依赖未安装

如果脚本提示缺少依赖，会显示详细的安装指导：

```bash
# 示例：Go 未安装
[ERROR] go is required but not installed.

To install Go:
  Ubuntu/Debian: sudo apt update && sudo apt install -y golang-go
  Or download from: https://go.dev/dl/
```

### Go 版本过低

脚本会自动检查 Go 版本，需要 Go 1.22+：

```bash
# 如果版本过低，会显示警告
[WARN] Go version go1.20 is too old. This project requires Go 1.22+
```

### Sudo 权限问题

如果遇到 sudo 权限问题：

```bash
# 测试 sudo 权限
sudo -v

# 如果失败，检查用户是否在 sudoers 组
groups

# 添加用户到 sudoers（需要 root 权限）
sudo usermod -aG sudo $USER
# 然后重新登录
```

### 服务启动失败
```bash
# 查看 Supervisor 状态
sudo supervisorctl status kratos-api

# 查看日志
tail -f /var/log/kratos-kit/api.err.log
tail -f /var/log/kratos-kit/api.out.log

# 检查 Supervisor 服务是否运行
sudo systemctl status supervisor
```

### 健康检查失败
1. 检查服务是否正常启动
   ```bash
   sudo supervisorctl status kratos-api
   ```

2. 检查端口是否被占用
   ```bash
   sudo netstat -tlnp | grep 9000
   # 或
   sudo ss -tlnp | grep 9000
   ```

3. 检查配置文件是否正确
   ```bash
   cat /opt/kratos-kit/configs/config.yaml
   ```

4. 查看应用日志
   ```bash
   tail -f /var/log/kratos-kit/api.out.log
   ```

### 权限问题
确保：
- 脚本有执行权限: `chmod +x scripts/deploy.sh`
- 用户有 sudo 权限
- 部署路径有写权限
  ```bash
  sudo mkdir -p /opt/kratos-kit
  sudo chown -R $USER:$USER /opt/kratos-kit
  ```

### WSL 特定问题

1. **Supervisor 服务未运行**
   ```bash
   # WSL 中可能需要手动启动
   sudo service supervisor start
   # 或
   sudo /etc/init.d/supervisor start
   ```

2. **PATH 环境变量问题**
   ```bash
   # 确保 Go 在 PATH 中
   echo $PATH | grep go
   # 如果不在，添加到 ~/.bashrc
   export PATH=$PATH:/usr/local/go/bin
   ```

3. **文件权限问题**
   ```bash
   # WSL 中文件权限可能需要特殊处理
   sudo chmod +x /opt/kratos-kit/bin/api
   ```

## 最佳实践

1. **生产环境部署前**:
   - 在测试环境验证部署脚本
   - 确保所有环境变量已正确设置
   - 检查备份目录有足够空间

2. **部署时**:
   - 使用 `--skip-tests` 快速部署（如果测试已在 CI/CD 中运行）
   - 监控部署日志
   - 部署后验证服务功能

3. **维护**:
   - 定期清理旧备份
   - 监控磁盘空间
   - 保留重要版本的备份

## 示例场景

### 场景1: 首次部署
```bash
# 设置环境变量
export APP_PATH="/opt/kratos-kit"
export APP_ENV="production"

# 部署服务
./scripts/deploy.sh api
```

### 场景2: 快速热更新
```bash
# 跳过测试和代码拉取，快速部署
./scripts/deploy.sh --skip-fetch --skip-tests api
```

### 场景3: 批量部署
```bash
# 部署所有服务
./scripts/deploy.sh --all
```

### 场景4: 紧急回滚
```bash
# 快速回滚到上一个版本
./scripts/deploy.sh --rollback api
```

## 注意事项

1. ⚠️ 部署脚本会使用 `git reset --hard`，确保重要更改已提交
2. ⚠️ 生产环境部署前务必在测试环境验证
3. ⚠️ 确保有足够的磁盘空间存储备份
4. ⚠️ 健康检查失败会自动回滚，请确保备份可用
5. ⚠️ Supervisor 配置会覆盖现有配置，请备份重要配置

## 相关文件

- `scripts/deploy.sh`: 部署脚本
- `Makefile`: 构建命令
- `configs/config.yaml`: API 服务配置
- `configs/config_admin.local.yaml`: Admin 服务配置

## 支持

如有问题或建议，请提交 Issue 或联系开发团队。

