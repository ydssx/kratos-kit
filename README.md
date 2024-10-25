# github.com/ydssx/kratos-kit

基于Kratos框架构建的后端项目。


## 技术栈

- Kratos: Go服务框架
- Gorm: 一个出色的ORM库
- Redis: 用于缓存
- MySQL: 用于持久化存储
- Asynq: 一个简单的分布式任务队列

## 项目结构
```bash
github.com/ydssx/kratos-kit/
├── api/                # API定义
├── cmd/                # 应用程序入口点
├── internal/
│   ├── biz/            # 业务逻辑层
│   ├── data/           # 数据访问层
│   ├── server/         # HTTP、gRPC和Asynq服务器定义
│   ├── service/        # 服务接口实现
│   └── job/            # 定时任务和队列任务处理函数
├── configs/            # 配置文件
├── docs/               # 文档
├── models/             # 数据库模型定义
├── pkg/                # 公共包，例如日志、工具等
├── third_party/        # 第三方包
├── scripts/            # 脚本文件
└── Makefile            # 项目管理命令

```

## 快速开始

### 先决条件

- [Go 1.22+](https://go.dev/dl/go1.22.3.windows-amd64.msi)
- [Git](https://github.com/git-for-windows/git/releases/download/v2.45.1.windows.1/Git-2.45.1-64-bit.exe)
- [Redis](https://github.com/tporadowski/redis/releases/download/v5.0.14.1/Redis-x64-5.0.14.1.msi)
- MySQL
- make 工具 (Windows系统推荐使用scoop安装，详见[官方文档](https://scoop.sh/),执行`scoop install make`)

### 本地运行

1. 克隆代码库

2. 安装工具

```bash
make init
```

3. 下载依赖

```bash
go mod tidy
```

 根据需要修改configs/config.test.yaml配置。

4. 运行项目

```bash
make run
```

## 使用说明
- 构建项目：`make build`
- 生成proto定义代码与swagger文档：`make gen`
- 生成依赖注入代码：`make wire`
- 生成数据库模型代码：`make gorm-gen`


## 相关文档和资源
- [Kratos 官方文档](https://go-kratos.dev/docs/)
- [wire 官方文档](https://github.com/google/wire)
- [protobuf 官方文档](https://protobuf.dev/)
- [protoc-gen-validate 官方文档](https://github.com/bufbuild/protoc-gen-validate)