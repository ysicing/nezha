# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 在此代码库中工作时提供指导。

## 项目概览

哪吒监控是一个可自托管的服务器监控系统，在统一代码库中包含两个主要组件：

- **Dashboard** (`cmd/dashboard`)：基于 Web 的监控服务端，包含 gRPC 服务（需要 CGO 支持 SQLite）
- **Agent** (`cmd/agent`)：轻量级监控客户端（纯 Go，无 CGO 依赖）

模块名称：`github.com/naiba/nezha`
Go 版本：1.24.0

## 开发命令

### 构建

**Dashboard（需要 CGO 和交叉编译工具链）：**
```bash
# 本地构建（使用本机编译器）
CGO_ENABLED=1 go build -o nezha-dashboard ./cmd/dashboard

# Linux AMD64 静态链接构建
CGO_ENABLED=1 CC=x86_64-linux-gnu-gcc go build \
  -ldflags="-s -w -extldflags '-static -fpic'" \
  -o dashboard-linux-amd64 ./cmd/dashboard
```

**Agent（无需 CGO）：**
```bash
CGO_ENABLED=0 go build -o nezha-agent ./cmd/agent
```

**生产环境构建：**
```bash
# 使用 GoReleaser 构建所有平台
goreleaser build --snapshot --clean
```

### 运行

**Dashboard（开发环境）：**
```bash
# 使用 docker-compose（推荐）
docker-compose up

# 直接运行（80 端口用于 HTTP，5555 端口用于 gRPC）
./nezha-dashboard
```

**Agent：**
```bash
# 连接到 dashboard
./nezha-agent -s <dashboard_地址:端口> -p <客户端密钥>
```

### 测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./model
go test ./pkg/monitor
go test ./cmd/agent

# 详细输出模式
go test -v ./...

# 带覆盖率统计
go test -cover ./...
```

### Protocol Buffers

修改 `proto/nezha.proto` 后执行：
```bash
bash script/proto.sh
```

该脚本会生成 Go 代码（`proto/nezha.pb.go` 和 `proto/nezha_grpc.pb.go`）。

## 架构要点

### 通信流程

1. **Agent → Dashboard**：gRPC 双向流
   - `ReportSystemState()`：周期性指标上报（CPU、内存、磁盘、网络）
   - `ReportSystemInfo()`：连接时上报主机信息
   - `RequestTask()`：长生命周期流，接收来自 dashboard 的任务
   - `IOStream()`：类 WebSocket 的双向流，用于终端/文件操作
   - `LookupGeoIP()`：GeoIP 查询服务

2. **Dashboard 内部**：
   - HTTP (Gin)：Web UI + RESTful API
   - gRPC Server：Agent 通信
   - Singleton 服务：在 `service/singleton/singleton.go` 中初始化的后台工作进程

### Singleton 服务

`service/singleton/` 中的关键后台服务：
- `AlertSentinel`：实时告警监控（状态变化触发）
- `ServiceSentinel`：周期性 HTTP/TCP/Ping 监控（基于 cron）
- `CronTask`：在 agent 上执行计划任务
- `Notification`：告警推送到多个通知渠道
- `DDNS`：动态 DNS 更新

所有 singleton 在 dashboard 启动时初始化一次。它们维护长生命周期的 goroutine，并通过数据库（GORM + SQLite）共享状态。

### 数据库访问模式

- ORM：GORM 配合 SQLite（`data/sqlite.db`）
- 模型定义位于 `model/` 目录
- 数据库句柄通过 `singleton.DB` 访问
- 多记录操作使用事务

### Agent 任务执行

任务流程：Dashboard → `RequestTask()` 流 → Agent 执行 → `ReportTask()` 报告给 Dashboard

任务类型包括：
- 命令执行
- 指标上报触发
- 终端 I/O 操作
- 文件操作

### RPC 认证

Dashboard 和 Agent 使用共享密钥进行认证：
- Dashboard 通过 `service/rpc/auth.go` 验证 agent 连接
- Agent 在连接元数据中提供客户端密钥
- 每个服务器在数据库中有唯一的 ID 和密钥对

### 前端主题

`resource/template/theme-*/` 中的多主题支持：
- 主题选择存储在配置中
- 模板使用 Go 的 `html/template`，共享组件位于 `resource/template/component/`
- 静态资源通过 `resource/resource.go` 中的 `//go:embed` 嵌入

### 国际化

翻译文件位于 `resource/l10n/*.toml`：
- 启动时由 `service/singleton/l10n.go` 加载
- Dashboard 设置中可按用户选择语言
- HTML 中使用模板函数 `Tf` 进行翻译

## 重要构建说明

1. **Dashboard 需要 CGO**，因为依赖 SQLite。交叉编译需要相应的 C 编译器（参见 `.goreleaser.yml` 中的 CC 环境变量配置）。

2. **Agent 必须无 CGO 依赖**，以实现跨平台的最大可移植性。

3. **版本注入**：Dashboard 版本通过 ldflags 注入：
   ```
   -X github.com/naiba/nezha/service/singleton.Version={{.Version}}
   ```
   Agent 版本使用 `main.version` 和 `main.arch`。

4. **Proto 变更**：修改 `.proto` 文件后务必运行 `script/proto.sh`，以确保 agent 的 proto 定义保持同步。

## 配置

Dashboard 配置：`data/config.yaml`（首次运行时自动创建）
Agent 配置：命令行参数（运行 `./nezha-agent --help` 查看）

Dashboard 配置文件示例：
```yaml
debug: false
httpport: 80
grpcport: 5555
language: zh-CN

oauth2:
  type: "github"  # github/gitlab/gitee/gitea
  admin: "username1,username2"
  clientid: "your_oauth_client_id"
  clientsecret: "your_oauth_client_secret"

site:
  brand: "哪吒监控"
  cookiename: "nezha-dashboard"
  theme: "default"
```

关键配置项：
- `httpPort`：Web UI 端口（默认：80）
- `grpcPort`：Agent 通信端口（默认：5555）
- `database`：SQLite 路径（默认：`data/sqlite.db`）
- `siteTitle`/`brand`：Dashboard 品牌名称
- `language`：UI 语言（en-US、zh-CN、zh-TW、es-ES）

## 关键目录说明

```
cmd/
├── dashboard/          # Dashboard 主程序
│   ├── controller/     # HTTP 路由处理器
│   └── rpc/           # gRPC 服务实现
└── agent/             # Agent 主程序

model/                 # 数据模型（GORM）
service/
├── singleton/         # 全局单例服务
└── rpc/              # RPC 服务逻辑层

pkg/                  # 可复用工具包
├── monitor/          # 监控相关工具
├── ddns/            # DDNS 提供商实现
├── geoip/           # GeoIP 查询
├── pty/             # 跨平台终端支持
└── utils/           # 通用工具函数

resource/
├── template/        # HTML 模板（多主题）
├── static/          # CSS/JS/图片
└── l10n/           # 国际化翻译文件

proto/               # Protocol Buffers 定义
script/              # 部署和配置脚本
```

## 常见开发场景

### 添加新的监控指标

1. 修改 `proto/nezha.proto` 中的 `State` 消息
2. 运行 `bash script/proto.sh` 生成代码
3. 在 Agent 端（`cmd/agent`）收集新指标
4. 在 Dashboard 端（`model/server.go`）存储和展示

### 添加新的 RPC 方法

1. 在 `proto/nezha.proto` 中定义新方法
2. 运行 `bash script/proto.sh`
3. 在 `cmd/dashboard/rpc/rpc.go` 实现服务端逻辑
4. 在 Agent 端调用新方法

### 添加新的通知渠道

1. 在 `model/notification.go` 中添加新类型常量
2. 在 `service/singleton/notification.go` 中实现发送逻辑
3. 在 Dashboard UI 中添加配置表单

### 调试技巧

**启用 Debug 日志：**
```yaml
# data/config.yaml
debug: true
```

**查看 gRPC 通信：**
```bash
# 使用 grpcurl 测试
grpcurl -plaintext -d '{"ip":"8.8.8.8"}' localhost:5555 proto.NezhaService/LookupGeoIP
```

**数据库查询：**
```bash
sqlite3 data/sqlite.db "SELECT * FROM servers;"
```
