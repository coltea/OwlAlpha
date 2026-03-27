# OwlAlpha 部署说明

## 1. 部署目标

OwlAlpha 的第一版部署目标是让用户可以在本地电脑、家庭服务器或小型云服务器中，以较低门槛完成完整系统部署。系统默认使用 Docker Compose 管理所有核心组件，并通过本地卷持久化数据。

## 2. 部署方式

推荐使用 Docker Compose 部署以下服务：

- `db`：业务数据库
- `redis`：缓存与任务辅助
- `api`：GoFrame 后端 API 服务
- `nginx`：统一入口、静态资源分发与反向代理，同时对外提供前端访问入口

## 3. 推荐部署拓扑

```text
                +-------------------+
                |       nginx       |
                | 静态资源 / 统一入口 |
                +---------+---------+
                          |
                          |
                     +----v----+
                     |   api   |
                     +----+----+
                          |
              +-----------+-----------+
              |                       |
        +-----v-----+           +-----v-----+
        |    db     |           |   redis   |
        +-----------+           +-----------+

               +----------------------+
               | OpenAI Compatible API|
               |   via custom URL     |
               +----------------------+
```

## 4. 环境要求

### 最低建议

- Docker 24+
- Docker Compose v2
- 2 CPU
- 4 GB 内存
- 20 GB 可用磁盘

### 推荐配置

- 4 CPU
- 8 GB 内存
- SSD 存储

## 5. 数据持久化设计

建议将以下数据通过卷进行持久化：

- PostgreSQL 数据目录
- Redis 数据目录
- 报告快照目录
- 系统日志目录
- 上传文件或缓存目录

## 6. 环境变量规划

### 6.1 基础环境变量

- `APP_NAME`：应用名称
- `APP_ENV`：运行环境，如 `dev`、`test`、`prod`
- `TZ`：时区，建议 `Asia/Shanghai`
- `SERVER_PORT`：后端服务端口

### 6.2 数据库相关

- `DB_HOST`
- `DB_PORT`
- `DB_NAME`
- `DB_USER`
- `DB_PASSWORD`

### 6.3 Redis 相关

- `REDIS_HOST`
- `REDIS_PORT`
- `REDIS_PASSWORD`
- `REDIS_DB`

### 6.4 鉴权相关

- `JWT_SECRET`
- `JWT_EXPIRE_HOURS`

### 6.5 模型相关

- `LLM_PROVIDER`
- `OPENAI_BASE_URL`
- `OPENAI_API_KEY`
- `OPENAI_MODEL`

### 6.6 数据源相关

- `MARKET_PROVIDER`
- `NEWS_PROVIDER`
- `TUSHARE_TOKEN` 或其他 A 股数据源配置

### 6.7 任务调度相关

- `SCHEDULER_ENABLED`
- `SCHEDULER_TIMEZONE`
- `ANALYSIS_CONCURRENCY`
- `ANALYSIS_RETRY_MAX`

## 7. Compose 设计建议

第一版 `docker-compose.yml` 建议支持以下特性：

- 所有服务位于同一个默认网络中。
- Nginx 对外暴露统一端口。
- `db` 和 `redis` 只在内部网络暴露。
- 使用 `.env` 文件注入环境变量。

## 8. 部署步骤建议

### 第一步：准备环境

- 安装 Docker 与 Docker Compose。
- 克隆项目代码。
- 复制环境变量模板并填写必要配置。

### 第二步：初始化配置

- 设置数据库连接信息。
- 设置 Redis 连接信息。
- 设置 JWT 密钥。
- 设置 OpenAI 兼容接口的 `Base URL`、模型名和 API Key。
- 设置 A 股数据源配置。

### 第三步：启动服务

- 执行 `docker compose up -d` 拉起基础服务。

### 第四步：检查状态

- 检查后端健康接口。
- 检查管理后台是否可访问。
- 检查数据库连接、Redis 连接和 OpenAI 兼容接口是否正常。

## 9. 运维建议

- 为 PostgreSQL 做定期备份。
- 为报告和日志目录做卷级别备份。
- 控制日志输出级别，避免无界增长。
- 为模型调用和任务执行增加超时与重试限制。
- 在公网部署时配置 HTTPS 与安全访问策略。

## 10. 安全建议

- 所有密钥通过环境变量注入，不写死在镜像中。
- 对外部署时限制数据库和 Redis 直连暴露。
- 配置强随机 JWT 密钥。
- 默认不在日志中记录完整 API Key 和敏感请求体。
- 后续可增加配置变更审计和登录审计。

## 11. 升级策略建议

- 升级前备份数据库与本地卷。
- 先更新镜像或代码，再执行数据库迁移。
- 检查配置项是否有新增或变更。
- 重启服务并验证关键链路。

## 12. 未来增强项

- 支持开发环境和生产环境的多套 Compose 配置。
- 支持独立 worker 服务处理批量分析任务。
- 支持外部对象存储和更细粒度的日志采集。
