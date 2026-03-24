# 多智能体平台任务追踪表

状态图标说明：`✅ 已完成` `🟨 进行中` `⬜ 未开始`

基线来源：`deep-research-report.md`  
记录起点：当前已交付前端为“产品介绍页（落地页）”。
主架构文档：`多智能体协作架构设计.md`

## 1）前端页面与 UI 任务（按业务架构拆分）

| ID | 页面 / 模块 | UI 任务 | 验收标准 | 状态 |
|---|---|---|---|---|
| FE-PAGE-001 | 落地页 / 产品介绍 | 实现单页介绍（Hero、架构、能力、API、路线图、风险、页脚） | 章节完整且可读 | ✅ |
| FE-PAGE-002 | 工作台框架 | 设计应用壳层（左侧导航 + 顶栏 + 内容区） | 支持多页面路由与响应式布局 | ✅ |
| FE-PAGE-016 | 全局项目库导航 | 左侧导航改为项目库，右侧先进入项目总览/项目详情入口，再进入各业务页 | 支持项目点击进入详情，并跳转群聊编排/编排执行/上下文记忆等模块 | ✅ |
| FE-PAGE-003 | Agent Compiler 页面（已并入群聊编排） | 原编译输入与 IR 产出流程合并到群聊编排页 | 独立编译器路由下线，流程并入群聊步骤 | ✅ |
| FE-PAGE-004 | Agent Assembler 页面 | 设计 IR -> AgentGraph 组装 UI | 展示角色模板、技能、工具、记忆组合视图 | ✅ |
| FE-PAGE-013 | 需求输入入口页面（已并入群聊页） | 取消独立需求入口，合并到群聊编排页面顶部 | 在群聊页可直接输入需求并重开确认流程 | ✅ |
| FE-PAGE-014 | 群聊编排页面 | 左对话右信息 + 结果步骤化展开（纵向） | 聊天区可滚动到底，步骤一/二可展开，Agent方块化 | ✅ |
| FE-PAGE-015 | Agent 详情页面 | 实现 Agent 资产详情视图 | 展示 `ROLE.md`、`SKILL.md` 与工具图标列表 | ✅ |
| FE-PAGE-005 | Orchestrator 页面 | 设计编排与执行 UI | 可视化 DAG/FSM/Planner 任务状态与流转 | ✅ |
| FE-PAGE-006 | 上下文与记忆页面 | 设计 Global/Task/Agent 记忆 UI | 支持“项目库侧栏可点击切换项目详情”与上下文槽位、RAG 检索展示 | ✅ |
| FE-PAGE-007 | 多 Agent 通信页面 | 设计黑板与消息流 UI | 展示共享状态、消息时间线、任务关联 | ✅ |
| FE-PAGE-008 | 调试与自愈页面 | 设计异常、重试、回滚 UI | 可查看失败任务、恢复策略与执行结果 | ✅ |
| FE-PAGE-009 | 安全与租户页面 | 设计权限边界与租户隔离 UI | 展示权限范围、沙箱策略、租户视图 | ⬜ |
| FE-PAGE-010 | 可观测性页面 | 设计指标/日志/追踪仪表盘 UI | 展示关键卡片、趋势图、追踪详情 | ⬜ |
| FE-PAGE-011 | API Explorer 页面 | 设计 API 请求/响应演示 UI | 覆盖 `/compile`、`/assembleAgents`、`/executeTask` Mock 流程 | ⬜ |
| FE-PAGE-012 | 配置与策略页面 | 设计模型/参数/策略管理 UI | 支持策略配置与合规说明展示 | ⬜ |

## 2）前端交互任务

| ID | 交互项 | 任务描述 | 验收标准 | 状态 |
|---|---|---|---|---|
| FE-INT-001 | 语言切换 | 实现中英文页面级切换 | 页面内容同步切换 | ✅ |
| FE-INT-002 | 锚点导航 | 实现章节锚点跳转 | 顶部导航可准确定位到章节 | ✅ |
| FE-INT-003 | 轻量动效 | 加入弱动效（滚动/入场） | 动效克制且可被无障碍设置降级 | ✅ |
| FE-INT-004 | 任务流模拟 | 实现分步 Mock 任务流程交互 | 用户可在 UI 中完成一次完整模拟流程 | ✅ |
| FE-INT-005 | 状态反馈 | 补齐加载/成功/失败/空态反馈 | 核心面板都有明确状态反馈 | ✅ |
| FE-INT-006 | 跨页导航 | 增加路由切换与激活态高亮 | 当前页面状态可清晰识别 | ✅ |
| FE-INT-007 | 需求到群聊跳转 | 输入需求后进入群聊并延续上下文 | 可从需求页 -> 群聊页 -> Agent详情页连贯操作 | ✅ |
| FE-INT-008 | 步骤化结果交互 | 结果区改为“第一步/第二步”纵向可展开 | 每一步可展开收起，第二步承载 Agent 方块入口 | ✅ |
| FE-INT-009 | 四栏舞台布局 | 群聊页改为“对话区 + 状态区 + IR成果区 + Agent成果区” | 桌面端四栏稳定、各栏内部独立滚动 | ✅ |

## 3）前端全局风格系统任务

| ID | 样式范围 | 任务描述 | 验收标准 | 状态 |
|---|---|---|---|---|
| FE-STYLE-001 | 设计令牌 | 定义颜色、间距、圆角、阴影、字体等令牌 | 令牌集中管理并被复用 | ✅ |
| FE-STYLE-002 | 全局视觉方向 | 建立极简文档风视觉基线 | 当前页面视觉一致性通过 | ✅ |
| FE-STYLE-003 | 组件样式规范 | 建立可复用卡片/表格/标签/按钮/表单样式 | 跨页面可直接复用 | ⬜ |
| FE-STYLE-004 | 主题策略 | 定义并实现明确的主题策略 | 主题行为可测试且稳定 | ⬜ |
| FE-STYLE-005 | 响应式规范 | 定义断点和多布局规则 | 目标视口无重叠、无截断 | ⬜ |

## 4）Mock 数据优先任务

| ID | 数据层 | 任务描述 | 验收标准 | 状态 |
|---|---|---|---|---|
| FE-MOCK-001 | 落地页内容 Mock | 建立结构化双语内容映射 | 页面内容全部来自 Mock 配置 | ✅ |
| FE-MOCK-002 | 工作流 Mock Schema | 定义 compile/assemble/execute 的 Mock 数据结构 | 可支撑全流程 UI 模拟 | ✅ |
| FE-MOCK-003 | 仪表盘 Mock 数据集 | 准备指标/日志/追踪数据集 | 无后端也可完整渲染 | ⬜ |
| FE-MOCK-004 | 异常场景 Mock | 加入超时/失败/部分成功场景 | 异常分支可被完整测试 | ⬜ |

## 5）前端自检任务

| ID | 自检项 | 标准 | 状态 |
|---|---|---|---|
| FE-QA-001 | 交互可用性 | 点击、切换、锚点、滚动正常 | ✅ |
| FE-QA-002 | Mock 渲染完整性 | 无空白块、无关键字段缺失 | ✅ |
| FE-QA-003 | 响应式可读性 | 桌面端与移动端基础可读通过 | ✅ |
| FE-QA-004 | 构建验证 | `npm run build` 成功 | ✅ |
| FE-QA-005 | 多页面冒烟检查 | 所有业务页面具备基础冒烟测试 | ✅ |

## 6）后端 / 平台任务（路线图保留）

| ID | 模块 | 任务描述 | 验收标准 | 状态 |
|---|---|---|---|---|
| BE-001 | Agent Compiler | NL 解析器（意图识别/实体抽取）已实现，`compile` 接口返回 `intent/entities/tasks/constraints/tools` | 用户输入可转为结构化请求 | ✅ |
| BE-002 | IR 层 | IR Schema v1 设计与校验 | 含 `goal/tasks/constraints/tools` 核心字段 | ⬜ |
| BE-003 | IR 层 | IR 版本管理机制 | 支持语义化版本与兼容策略 | ⬜ |
| ORC-001 | 编排引擎 | DAG 执行模式 | 按依赖顺序执行任务图 | ⬜ |
| ORC-002 | 编排引擎 | FSM 执行模式 | 支持状态迁移驱动流程 | ⬜ |
| ORC-003 | 编排引擎 | Planner Loop 模式 | 支持动态规划与任务分发 | ⬜ |
| AGT-001 | Agent 组装 | 角色模板库 | 支持按角色实例化 Agent | ⬜ |
| AGT-002 | Agent 组装 | Skill 接口标准化 | 统一输入/输出与权限元数据 | ⬜ |
| AGT-003 | Agent 组装 | Tool 接口标准化 | 统一协议并可审计调用 | ⬜ |
| MEM-001 | 记忆系统 | Global/Task/Agent 分层上下文 | 支持分层读写 | ⬜ |
| MEM-002 | 记忆系统 | RAG 检索接入 | 支持向量检索并注入上下文 | ⬜ |
| MEM-003 | 记忆系统 | Memory Slots 更新机制 | 支持按槽位定向更新 | ⬜ |
| COM-001 | 通信系统 | 黑板模式共享状态 | Agent 可通过黑板协同 | ⬜ |
| COM-002 | 通信系统 | 消息协议与审计元数据统一 | 含消息 ID/时间戳/任务 ID | ⬜ |
| DBG-001 | 调试系统 | Debug Agent 异常检测 | 可识别超时/失败/无响应 | ⬜ |
| DBG-002 | 调试系统 | 重试/回滚/降级策略 | 恢复策略可执行 | ⬜ |
| SEC-001 | 安全体系 | 多租户隔离 | 执行与数据边界隔离 | ⬜ |
| SEC-002 | 安全体系 | 工具调用沙箱化 | 外部调用在受限环境运行 | ⬜ |
| SEC-003 | 安全体系 | 审计与隐私策略 | 满足可追溯与合规要求 | ⬜ |
| OBS-001 | 可观测性 | 指标/日志/追踪三支柱 | 关键观测数据可查询 | ⬜ |
| OPS-001 | 工程化 | CI/CD 与测试流水线 | 构建/单测/集测可运行 | ⬜ |
| API-001 | 平台 API | 项目库/对话/编译/组装/产物 5 组接口已落地（MySQL 持久化），前端已完成基础对接 | 接口可调用并返回结构化结果 | ✅ |
| DATA-001 | 数据存储 | 关系库/向量库/缓存分层落地 | 满足配置、记忆、热数据场景 | ⬜ |   

## 7）架构文档与 Prompt 体系任务

| ID | 模块 | 任务描述 | 验收标准 | 状态 |
|---|---|---|---|---|
| DOC-001 | 架构文档 | 新增《多智能体协作架构设计》主文档 | 文档覆盖主链路、模块分层、状态机、接口与存储 | ✅ |
| DOC-002 | 架构文档 | 与 `deep-research-report.md` 对齐校验 | 术语、阶段、主链路一致，无范式冲突 | ✅ |
| DOC-003 | 架构文档 | 补齐端到端时序图、模块图、状态机图 | 至少包含 3 类图，且可映射当前系统 | ✅ |
| PROMPT-001 | Prompt 体系 | 基础协作协议 Prompt | 明确通用约束、输出规范、禁止项 | ✅ |
| PROMPT-002 | Prompt 体系 | 需求澄清 Prompt | 包含模板、校验器、失败规则、示例 | ✅ |
| PROMPT-003 | Prompt 体系 | IR 生成 Prompt | 包含模板、校验器、失败规则、示例 | ✅ |
| PROMPT-004 | Prompt 体系 | Agent 组装 Prompt | 包含模板、校验器、失败规则、示例 | ✅ |
| PROMPT-005 | Prompt 体系 | 编排执行 Prompt | 包含模板、校验器、失败规则、示例 | ✅ |
| PROMPT-006 | Prompt 体系 | Reflection / 审计 Prompt | 包含模板、校验器、失败规则、示例 | ✅ |
| PROMPT-007 | Prompt 体系 | 最终交付与确认 Prompt | 包含模板、校验器、失败规则、示例 | ✅ |
| PROMPT-008 | Prompt 体系 | Prompt 资产文件化与后端配置挂载 | `backend/config/prompts/` 独立文件与 `prompt_assets` 配置可直接引用 | ✅ |
| PROMPT-009 | Prompt 体系 | Prompt 运行时接入后端主链路 | 启动时完成 PromptLoader 装载，`chat` 使用澄清 Prompt，`compile` 使用编译 Prompt | ✅ |
| PROMPT-010 | Prompt 体系 | Validator Prompt 运行时二次校验接入 | 澄清结果与 IR 结果在运行时完成契约校验，不通过时自动回退兜底结果 | ✅ |
| ORC-004 | 编排架构 | Planner 主控状态机 | 文档与任务流中已明确状态推进规则 | ⬜ |
| ORC-005 | 编排架构 | 返工与人工闸门机制 | 明确 2 轮返工上限与 IR/最终确认闸门 | ⬜ |
| ORC-006 | 编排架构 | 编排日志回放与前端展示 | 支持多轮 Envelope 回放与 UI 映射 | ⬜ |

## 8）后端接口前端接入状态（防误判）

| ID | 接口 / 字段 | 后端实现状态 | 前端接入状态 | 备注 |
|---|---|---|---|---|
| API-FE-001 | `GET /api/projects` | ✅ | ✅ | 项目库页面与侧栏已接入 |
| API-FE-002 | `GET /api/projects/:id` | ✅ | ✅ | 项目详情基础信息已接入 |
| API-FE-003 | `GET /api/projects/:id/chat` | ✅ | ✅ | 进入群聊页加载历史消息 |
| API-FE-004 | `POST /api/projects/:id/chat` | ✅ | ✅ | 群聊发送消息已接入 |
| API-FE-005 | `POST /api/projects/:id/compile` | ✅ | ✅（部分） | 已调用并展示 IR 文件信息 |
| API-FE-006 | `POST /api/projects/:id/assemble` | ✅ | ✅（部分） | 已触发组装并在组装器页展示 AgentGraph |
| API-FE-007 | `GET /api/projects/:id/artifacts` | ✅ | ✅（部分） | 已读取 IR 与 AgentGraph |
| API-FE-008 | `compile` 返回 `intent/entities/tasks/constraints/tools` | ✅ | 🟨 | 目前前端仅展示 IR 文件名与来源，尚未完整展示上述结构化字段 |
| API-FE-009 | `compile` 返回 `acceptanceChecklist/completionCriteria/humanGates` | ✅ | ⬜ | 后端已随 Prompt 编译链路落地，前端尚未展示这些结构化字段 |
| API-FE-010 | `POST /api/projects/:id/orchestrate` | ⬜ | ⬜ | 编排执行主入口，当前仅在架构文档中定义 |
| API-FE-011 | `POST /api/projects/:id/confirm` | ⬜ | ⬜ | 用于 `ir_confirm` 与 `final_confirm`，当前待实现 |
| API-FE-012 | `GET /api/projects/:id/orchestration-log` | ⬜ | ⬜ | 用于前端编排回放与多 Agent 通信展示 |

## 9）前端视觉探索阶段（本轮重点）

| ID | 任务 | 产出 | 状态 |
|---|---|---|---|
| FE-DESIGN-001 | 视觉探索需求定义（工作台 + 编译器、桌面端、3 套高创新） | 需求边界与评审标准明确 | ✅ |
| FE-DESIGN-002 | 方案 A：模块解构式工作台 | `reports/ui-proposals/concept-a-workbench.svg` `reports/ui-proposals/concept-a-compiler.svg` | ✅ |
| FE-DESIGN-003 | 方案 B：时间轴编排式工作台 | `reports/ui-proposals/concept-b-workbench.svg` `reports/ui-proposals/concept-b-compiler.svg` | ✅ |
| FE-DESIGN-004 | 方案 C：图谱操作台式工作台 | `reports/ui-proposals/concept-c-workbench.svg` `reports/ui-proposals/concept-c-compiler.svg` | ✅ |
| FE-DESIGN-005 | 交互说明与全局风格基线沉淀 | `reports/ui-proposals/方案说明.md` | ✅ |
| FE-DESIGN-006 | 方案选型决策（由你确认） | 已选定方案 B（时间轴编排式）进入实现阶段 | ✅ |
| FE-DESIGN-007 | 选中方案落地任务拆解（组件、路由、Mock 交互） | 进入 FE-PAGE-002/003 实装任务单 | ✅ |

### 方案挂载信息（DesignConcept）

| 方案 | 文件路径 | 设计关键词 | 优点 | 风险 | 推荐场景 | 状态 |
|---|---|---|---|---|---|---|
| A 模块解构式 | `reports/ui-proposals/concept-a-workbench.svg` `reports/ui-proposals/concept-a-compiler.svg` | 非对称分栏、流程胶囊、高对比信息层 | 流程清晰、反馈直观、上手快 | 低分辨率下密度控制要求高 | 强调任务链路与状态反馈 | ✅（备选） |
| B 时间轴编排式 | `reports/ui-proposals/concept-b-workbench.svg` `reports/ui-proposals/concept-b-compiler.svg` | 任务流舞台、状态层叠、节奏导览 | 阶段感强、展示友好、推进明确 | 并行任务扩展需额外规则 | 强调时序推进与阶段管理 | ✅（已选定） |
| C 图谱操作台式 | `reports/ui-proposals/concept-c-workbench.svg` `reports/ui-proposals/concept-c-compiler.svg` | 节点画布、语义连线、侧边控制台 | 关系表达力强、辨识度高 | 初次使用学习成本较高 | 强调关系网与策略联动 | ✅（备选） |

## 默认执行原则（生效中）

1. 优先实现前端页面与交互效果。  
2. 优先使用 Mock 数据跑通前端流程，不被后端联调阻塞。  
3. 每次前端改动完成后，先做前端自检，再进入下一阶段。  
4. 以后新生成的 Markdown（`.md`）文件默认全部使用中文。  
5. 代码提交时，`commit message` 必须使用中文，并遵循规范化提交格式。  
6. UI 设计图默认填充 Mock 数据（字段、状态、指标、示例值）。

## 任务更新规则

1. 任务完成后立即标记为 `✅`。  
2. 进行中的任务标记为 `🟨`，完成后改为 `✅`。  
3. 如有新范围，新增任务行，保持 ID 稳定可追踪。  
4. 每次将某个任务标记为 `✅` 时，必须同步在本文件“任务衔接记录”中写明下一步任务编号与目标。  

## Commit 提交规范（中文）

1. 格式：`type(scope): 中文摘要`  
2. 常用 `type`：`feat` `fix` `docs` `style` `refactor` `test` `chore`  
3. 示例：`feat(frontend): 新增工作台框架与路由导航`  


## 后端接口接入标记规则（新增）

1. 只要后端接口已实现但前端尚未接入，必须在“后端接口前端接入状态（防误判）”中单独列出，并标记为 `🟨` 或 `⬜`。
2. 禁止仅以“后端已实现”判定整体完成，必须同时写明前端接入状态。
3. 对结构化字段（如 `intent/entities/tasks/constraints/tools`）必须单独列项，标明“后端实现状态”与“前端展示状态”。

## 10）任务衔接记录

| 最近完成任务 | 完成说明 | 下一步任务 | 下一步目标 | 状态 |
|---|---|---|---|---|
| FE-PAGE-008 | 已完成调试与自愈页面，实现异常 run 列表、恢复策略矩阵、重试/回滚控制台与执行结果展示，并通过测试构建验证 | FE-PAGE-009 | 进入“安全与租户页面”设计与实现，补齐权限边界、沙箱策略与租户隔离视图 | ⬜ |
