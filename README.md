

# cronin 服务器任务管理平台

## 介绍
cronin 是一个基于 Go 语言开发的服务器任务管理平台，支持定时任务、流水线任务及 Webhook 接收任务，适用于多种场景。平台提供了可视化界面，方便用户进行任务配置、监控、日志查看等操作。

## 特点
- 多任务类型支持：HTTP、RPC、Shell �icker、SQL、Jenkins、Git �icker �icker。
- 完善的权限控制：角色权限设置、用户管理。
- 支持流水线任务：多个任务组合，实现复杂业务逻辑。
- Webhook 接收任务：通过 Webhook 触发任务执行。
- 任务日志追踪：提供日志与分布式追踪功能。
- 支持多种数据源：MySQL、ClickHouse、SQLite、Git、Jenkins、Gitee/GitHub API。
- 丰富的前端组件：基于 Element UI 构建，提供良好的交互体验。

## 功能
### 任务类型
- **定时任务**：支持周期性任务，可设置 Cron 表达式。
- **HTTP 任务**：定时调用 HTTP 接口。
- **RPC 任务**：支持 gRPC �icker。
- **Shell 任务**：执行本地 Shell 命令。
- **SQL 任务**：执行 SQL 语句，支持 MySQL 和 ClickHouse。
- **Jenkins 任务**：触发 Jenkins 构建。
- **Git 任务**：支持 Gitee/GitHub API 操作 PR、文件更新、代码审查等。

### 流水线任务
- 多任务组合执行。
- 支持任务失败处理策略（跳过、停止、继续）。
- 提供日志追踪功能。

### Webhook 接收任务
- 通过 Webhook 触发任务。
- 支持多种 Webhook 源（如 Tapd、GitHub、Gitee）。
- 可配置任务触发规则。

### 日志与追踪
- 任务执行日志记录。
- 分布式追踪支持，记录任务执行链路。
- 支持日志检索、过滤、删除。

### 系统设置
- �icker 管理：配置 Jenkins、Git、SQL �icker。
- 全局变量设置：定义全局变量供任务调用。
- 环境配置：设置不同环境的任务配置。
- 消息模板配置：自定义任务执行后的消息通知模板。
- 权限管理：角色权限分配，支持不同角色访问不同资源。
- 用户管理：用户增删改查、密码修改、状态控制。

### 前端组件
- 基于 Element UI 的前端组件。
- 提供任务配置表单、日志展示、Webhook 配置等组件。
- 支持任务状态显示、日志追踪、参数替换等交互功能。

## 功能预览
- 任务配置页面：支持多种任务类型配置，包括 HTTP、Shell、SQL、Jenkins 等。
- 流水线任务页面：可视化配置多个任务的执行顺序与失败策略。
- Webhook 接收任务页面：配置 Webhook URL、触发条件与任务参数。
- 任务日志页面：展示任务执行日志、状态、耗时等信息。
- 任务追踪页面：展示任务执行链路，支持链路详情查看。
- 系统设置页面：配置 Jenkins、Git、SQL 源，设置全局变量与环境信息。
- 用户管理页面：管理用户权限、状态、密码、账号等。
- 权限控制页面：分配角色权限，控制不同角色的访问范围。

## 文档
- [手册](README.md#手册)
- [安装](README.md#安装)
- [捐助与支持](README.md#捐助与支持)
- [参与贡献](README.md#参与贡献)

## 手册
- [任务设置](work/config_set.md)
- [流水线设置](work/pipeline_set.md)
- [接收设置](work/receive_set.md)
- [消息模板设置](work/message_template_set.md)
- [链接设置设置](work/source_set.md)
- [用户设置](work/user_set.md)

## 安装
### 一、获取程序包
- 从源码编译：
```sh
git clone https://gitee.com/cronin/cronin.git
cd cronin
make build
```
- 或使用 Docker：
```sh
docker pull cronin:latest
```

### 二、完善配置
- 配置数据库连接信息（MySQL/SQLite）。
- 配置 HTTP、gRPC、Jenkins、Git �icker 信息。
- 配置全局变量与环境信息。

### 三、运行
- 本地运行：
```sh
./cronin
```
- Docker 运行：
```sh
docker run -d -p 9003:9003 cronin:latest
```

## 捐助与支持
如需支持 cronin 的持续开发，可通过以下方式捐助：
- 支付宝
- 微信
- 银行转账

## 参与贡献
欢迎参与 cronin 的开发与文档完善：
- Fork 仓库并提交 PR。
- 参与 issue 讨论。
- 提交文档优化建议。

## 前端组件
- 任务配置组件：`web/components/config.js`, `web/components/config_form.js`
- 流水线组件：`web/components/pipeline.js`, `web/components/pipeline_form.js`
- Webhook 配置组件：`web/components/receive.js`, `web/components/receive_form.js`
- 消息模板组件：`web/components/message_template.js`
- 用户管理组件：`web/components/user.js`
- 权限管理组件：`web/components/role.js`
- �icker �

## 后端服务
### 服务启动
- HTTP 服务：`internal/server/http.go`
- 路由配置：`internal/server/router_*.go`
- 中间件：`internal/server/http_middleware.go`

### 任务调度
- 定时任务调度：`internal/biz/cron.go`
- 流水线任务调度：`internal/biz/cron_pipeline.go`
- Webhook 接收任务：`internal/biz/job_receive.go`

### 任务执行
- HTTP 任务：`internal/biz/job_config.go`
- gRPC 任务：`internal/biz/job_config.go`
- Shell 任务：`internal/biz/job_config.go`
- SQL 任务：`internal/biz/job_config_sql.go`
- Jenkins 任务：`internal/biz/job_config_jenkins.go`
- Git 任务：`internal/biz/job_config_git.go`

### 数据模型
- 任务配置：`internal/models/cron_config.go`
- 流水线配置：`internal/models/cron_pipeline.go`
- Webhook 配置：`internal/models/cron_receive.go`
- 日志记录：`internal/models/cron_log.go`
- 分布式追踪：`internal/models/cron_log_span.go`, `internal/models/cron_log_span_index.go`
- 系统设置：`internal/models/cron_setting.go`
- 用户与角色：`internal/models/cron_user.go`, `internal/models/cron_auth_role.go`

## 依赖库
- Gin：用于构建 HTTP 服务。
- GORM：用于数据库操作。
- Gitee/GitHub API：用于 Git 相关操作。
- gRPCurl：用于 gRPC 接口调用。
- Redis：缓存任务状态。
- MySQL/ClickHouse：任务执行数据存储。
- JWT：用户鉴权。
- SSE：实时推送任务状态。

## 协议与授权
- 使用 MIT 授权协议。
- 项目代码遵循开源规范，欢迎社区参与贡献与改进。