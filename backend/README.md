# 后端接口说明（Go + Docker + MySQL）

当前后端已使用 Go + MySQL 持久化，并实现前端联调所需接口。

## 一键启动（推荐）

```bash
cd backend
docker compose up -d --build
```

服务地址：

- API: `http://127.0.0.1:8080`
- MySQL: `127.0.0.1:3306`（root/root）

## 编译器配置（BE-001）

后端已实现 `BE-001 Agent Compiler`，支持从配置文件选择调用通道：

- `local_api`：调用本地解析服务
- `model_api`：调用模型 API（OpenAI 兼容 chat completions）

配置文件：`backend/config/config.yaml`

```yaml
compiler:
  mode: local_api # local_api | model_api
  timeout_seconds: 20
  local_api:
    endpoint: http://127.0.0.1:9090/parse
    api_key: ""
  model_api:
    endpoint: https://api.openai.com/v1/chat/completions
    api_key: ""
    model: gpt-4o-mini
```

可选环境变量覆盖：

- `CONFIG_FILE`
- `COMPILER_MODE`
- `COMPILER_LOCAL_API_ENDPOINT`
- `COMPILER_LOCAL_API_KEY`
- `COMPILER_MODEL_API_ENDPOINT`
- `COMPILER_MODEL_API_KEY`
- `COMPILER_MODEL_NAME`
- `COMPILER_TIMEOUT_SECONDS`

## 本地直接启动 API

前提：本地已有可连接 MySQL。

```bash
cd backend
set MYSQL_DSN=root:root@tcp(127.0.0.1:3306)/ment_agents?charset=utf8mb4&parseTime=True&loc=Local
go run ./cmd/api
```

## 已实现接口

1. 项目库
- `GET /api/projects`
- `POST /api/projects`
- `GET /api/projects/:id`

2. 对话
- `GET /api/projects/:id/chat`
- `POST /api/projects/:id/chat`

3. 编译（产出 IR）
- `POST /api/projects/:id/compile`

4. 组装（产出 AgentGraph）
- `POST /api/projects/:id/assemble`

5. 产物读取
- `GET /api/projects/:id/artifacts`

## 数据持久化说明

- 使用 GORM + MySQL 8.0。
- 启动时自动迁移表结构：
  - `projects`
  - `chat_messages`
  - `project_artifacts`
- 首次启动会自动写入一个种子项目 `proj-001`。
