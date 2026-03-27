# OwlAlpha

<p align="center">
  <img src="/images/logo.png" width="300" />
</p>


OwlAlpha 是一个本地优先的 A 股 AI 智能分析系统，面向个人投资者、研究型用户和小型团队，提供可自部署、可扩展、数据默认本地化的智能分析能力。

项目采用前后端分离架构，后端基于 GoFrame，前端以单一 `web/` 管理后台项目承载分析与配置能力，并通过 Docker Compose 提供一键部署能力。V0.1 版本聚焦 A 股每日报告场景，系统聚合 A 股行情、新闻资讯和公司基础信息，通过用户配置的 OpenAI 兼容 API 生成每日分析报告、风险提示与观察要点。

## V0.1 定位

- 只支持 A 股。
- 只提供后台管理端，后台登录后查看分析相关内容。
- 只支持 OpenAI 兼容 API 接入，允许自定义 `Base URL`、模型名和 API Key。
- 核心功能聚焦“每日生成报告”，优先保证主链路可用与可部署。

## 项目目标

- 提供一个支持自部署安装的 A 股智能分析平台。
- 确保业务数据、历史报告与配置默认保存在用户本地环境。
- 支持用户接入自定义 OpenAI 兼容 API。
- 使用清晰的前后端分离架构，为后续市场扩展和功能迭代预留空间。

## 核心能力

- 后台登录：系统后台需要登录后访问分析与配置功能。
- 股票池管理：支持 A 股股票分组、自选池维护、启停分析。
- 行情聚合：采集基础行情、K 线数据并进行本地缓存。
- 新闻分析：聚合与股票相关的新闻资讯，支持去重、时效过滤与上下文整理。
- AI 分析：基于多维上下文生成每日股票分析报告，包含摘要结论、风险提示、观察建议等。
- 报告中心：保存分析快照、模型配置和任务记录，支持按日查看历史报告。
- 任务调度：支持手动分析与定时生成每日报告。
- 本地部署：通过 Docker Compose 拉起 `nginx`、`backend`、`postgres`、`redis` 四个核心容器，其中前端静态文件直接由 `nginx` 容器承载。

## 架构概览

```text
OwlAlpha/
├── backend/        # GoFrame 后端服务
├── web/            # 单一后台前端项目
├── deploy/         # Docker Compose、Nginx、环境变量模板
├── docs/           # 项目文档
└── README.md
```

更详细的设计说明见：

- `docs/product-overview.md`
- `docs/architecture.md`
- `docs/backend-design.md`
- `docs/deployment.md`

## 技术选型

- 后端：Go 1.23+、GoFrame 2.x
- 前端：React、TypeScript、Vite
- 数据库：PostgreSQL
- 缓存与任务辅助：Redis
- AI 模型：OpenAI 兼容接口
- 部署：Docker Compose、Nginx

## 当前阶段

当前仓库处于 `0.1` 版本规划阶段，优先完成以下内容：

- 核心文档与系统设计
- Monorepo 基础目录初始化
- GoFrame 后端骨架
- 单一后台前端骨架
- Docker Compose 部署基线

## 产品边界

`0.1` 版本仅聚焦 A 股每日分析报告场景，暂不包含以下能力：

- 港股和美股支持
- 独立用户端站点
- 本地模型或 Ollama 接入
- 自动下单与券商交易接入
- 复杂回测系统
- 多租户 SaaS 化能力
- 高级组织协同审批流

## 设计原则

- 本地优先：数据默认保存在用户本地数据库和本地卷中。
- 可控部署：所有核心组件均支持自托管运行。
- 分层清晰：后端按控制器、服务、逻辑、数据访问分层，前端以单一后台应用承载分析与配置界面。
- 易于扩展：为后续增加更多分析策略、更多数据源和更多市场留出边界。
- 可审计：任务、报告、配置变更和关键调用链支持追踪。

## 规划中的目录结构

```text
OwlAlpha/
├── backend/
│   ├── api/
│   ├── internal/
│   ├── manifest/
│   ├── resource/
│   └── main.go
├── web/
├── deploy/
│   ├── docker-compose.yml
│   ├── docker/
│   └── env/
├── docs/
└── README.md
```

## 文档清单

- `docs/product-overview.md`：产品定位、目标用户、功能范围与里程碑。
- `docs/architecture.md`：系统架构、模块边界、数据流与部署拓扑。
- `docs/backend-design.md`：GoFrame 后端分层、模块职责与核心接口设计。
- `docs/deployment.md`：Docker Compose 部署说明、环境变量和运维建议。

## 当前已实现的基础骨架

- `backend/`：GoFrame + Gorm 后端基础工程。
- `web/`：React + Vite 单一后台前端骨架。
- `deploy/docker-compose.yml`：PostgreSQL、Redis、Backend、Nginx 组合部署，前端构建产物由 `nginx` 容器直接提供。
- `backend/manifest/sql/001_init.sql`：初始化表结构、默认管理员和示例报告数据。

## 本地启动

### 方式一：Docker Compose

```bash
cp deploy/env/.env.example deploy/env/.env
docker compose -f deploy/docker-compose.yml up --build
```

启动后访问：

- 后台地址：`http://localhost:8080`
- 默认账号：`admin`
- 默认密码：`admin123456`

说明：

- 没有独立 `web` 容器
- 前端代码位于 `web/`，但 Docker 部署时由 `nginx` 镜像构建并托管静态资源

### 方式二：分别启动前后端

后端：

```bash
cd backend
go mod tidy
go run .
```

前端：

```bash
cd web
npm install
npm run dev
```

说明：

- 后端默认读取 `backend/manifest/config/config.yaml`
- 前端默认访问 `/api/v1`
- 如需完整运行，仍建议优先使用 Docker Compose

## 声明

本项目用于学习、研究与辅助分析，不构成任何投资建议。股市有风险，投资需谨慎。
