# Ment-Agents

基于 `deep-research-report.md` 的多智能体平台项目仓库。  
当前阶段优先完成前端页面与交互，采用 Mock 数据先跑通，再逐步推进后端能力。

## 项目目标

- 将自然语言需求编译为可执行的多 Agent 系统（NL -> IR -> AgentGraph）
- 提供可视化的编排、执行、调试与观测体验
- 在保障安全隔离与合规审计的前提下实现可扩展架构

## 当前进展

- 已完成 `frontend/` 落地页（中英切换、锚点导航、基础响应式、Mock 内容渲染）
- 已建立任务追踪文档：`task-tracker.md`
- 已明确默认执行原则：前端优先、Mock 优先、改动后先自检

## 目录结构

```text
.
├─ deep-research-report.md      # 深度研究报告（需求与架构来源）
├─ task-tracker.md              # 任务拆分、状态追踪与执行规则
├─ frontend/                    # Vue + Vite 前端工程
│  ├─ src/
│  │  ├─ App.vue
│  │  ├─ style.css
│  │  └─ content/siteContent.js
│  └─ package.json
└─ README.md
```

## 前端快速启动

```bash
cd frontend
npm install
npm run dev
```

生产构建：

```bash
cd frontend
npm run build
```

## 任务管理

所有任务状态统一维护在 `task-tracker.md`：  
- `✅ 已完成`  
- `🟨 进行中`  
- `⬜ 未开始`

每完成一个功能，立即更新对应任务状态。

## 提交规范

- 提交信息必须使用中文
- 格式：`type(scope): 中文摘要`
- 示例：`feat(frontend): 新增工作台框架与路由导航`

