/**
 * @typedef {'zh' | 'en'} Locale
 */

/**
 * @typedef {{
 *   id: string
 *   label: string
 * }} NavItem
 */

/**
 * @typedef {{
 *   brand: string
 *   nav: NavItem[]
 *   hero: { kicker: string, title: string, summary: string, highlights: string[] }
 *   architecture: { title: string, intro: string, flowTitle: string, flow: string[], bundleTitle: string, bundle: string[] }
 *   capabilities: { title: string, intro: string, blocks: { title: string, items: string[] }[] }
 *   api: { title: string, intro: string, table: { method: string, endpoint: string, request: string, response: string }, rows: { method: string, endpoint: string, request: string, response: string }[] }
 *   roadmap: { title: string, intro: string, phases: { stage: string, name: string, goal: string }[] }
 *   risk: { title: string, intro: string, items: { title: string, detail: string }[] }
 *   footer: { reference: string, notice: string }
 * }} SiteCopy
 */

/**
 * @type {Record<Locale, SiteCopy>}
 */
export const siteContent = {
  zh: {
    brand: 'Meta-Agent Studio',
    nav: [
      { id: 'architecture', label: '架构总览' },
      { id: 'capabilities', label: '核心能力' },
      { id: 'api', label: 'API 示例' },
      { id: 'roadmap', label: '路线图' },
      { id: 'risk', label: '风险控制' },
    ],
    hero: {
      kicker: '自动化定制多 Agent 平台',
      title: '把自然语言需求，编译为可执行的 Agent 系统',
      summary:
        '本页面基于 deep-research-report.md 提炼，聚焦可读、可信、可演示的产品化表达。首版采用静态单页，展示平台核心设计、接口示例与实施路线。',
      highlights: ['NL -> IR -> AgentGraph', 'DAG / FSM / Planner Loop', 'RAG 记忆与黑板通信', '安全隔离与可观测'],
    },
    architecture: {
      title: '架构总览',
      intro: '系统以 Agent Compiler 为中心，把用户意图转换为可调度、可追踪、可复用的执行图。',
      flowTitle: '核心流程',
      flow: [
        '用户自然语言输入',
        'NL 解析与意图识别',
        'IR 结构化生成',
        'Assembler 组装 AgentGraph',
        'Orchestrator 调度执行',
        '反馈与自优化闭环',
      ],
      bundleTitle: 'Agent 构成',
      bundle: ['角色模板（Role Template）', '技能模块（Skill）', '工具接口（Tool）', '状态与长期记忆（Memory）'],
    },
    capabilities: {
      title: '核心能力',
      intro: '平台能力覆盖编排、记忆、通信、调试、安全与运维观测，支持从 MVP 向生产化演进。',
      blocks: [
        {
          title: '编排模式',
          items: ['DAG：流程稳定、调试友好', 'FSM：状态驱动、适合复杂交互', 'Planner Loop：动态规划、自治调度'],
        },
        {
          title: '上下文与记忆',
          items: ['Global / Task / Agent 三级上下文', 'RAG 向量检索增强注入', 'Memory Slots 结构化更新'],
        },
        {
          title: '多 Agent 通信',
          items: ['黑板模式共享状态', '发布-订阅降低耦合', '统一消息协议与审计追踪'],
        },
        {
          title: '自愈与安全',
          items: ['Debug Agent 异常检测与纠错', 'ACK + 重试 + 熔断容错', '多租户隔离与工具沙箱化'],
        },
      ],
    },
    api: {
      title: 'API 示例（展示用）',
      intro: '以下接口来自报告中的简化示例，仅用于落地页展示，不发起真实请求。',
      table: {
        method: '方法',
        endpoint: '接口',
        request: '请求体',
        response: '响应体',
      },
      rows: [
        {
          method: 'POST',
          endpoint: '/compile',
          request: '{"query":"开发电商网站"}',
          response:
            '{"ir":{"goal":"开发电商网站","tasks":[{"id":"t1","type":"analysis"},{"id":"t2","type":"backend"},{"id":"t3","type":"frontend"}]}}',
        },
        {
          method: 'POST',
          endpoint: '/assembleAgents',
          request: '{"ir":{...}}',
          response: '{"agents":[{"id":"A1","role":"业务分析师"},{"id":"A2","role":"后端工程师"},{"id":"A3","role":"前端工程师"}]}',
        },
        {
          method: 'POST',
          endpoint: '/executeTask',
          request: '{"agentId":"A2","taskId":"t2","input":"..."}',
          response: '{"status":"success","output":"数据库搭建完成"}',
        },
      ],
    },
    roadmap: {
      title: '路线图与实施节奏',
      intro: '报告建议从 MVP 验证核心机制，再逐步扩展自动化、智能化与规模化能力。',
      phases: [
        { stage: '阶段 1', name: 'PoC / MVP', goal: '固定模板 Agent + 简单 DAG + 手工技能集成' },
        { stage: '阶段 2', name: '自动化 Agent', goal: 'LLM 自动拆解任务并生成 Agent，提升任务分配自动化' },
        { stage: '阶段 3', name: '智能优化', goal: '完善 IR 版本管理与 Planner 调度，实现混合执行模式' },
        { stage: '阶段 4', name: '生产化规模', goal: '构建 Agent Store、多租户支持与高并发能力' },
      ],
    },
    risk: {
      title: '风险与缓解',
      intro: '平台建设需同时关注技术正确性、性能、合规与安全边界。',
      items: [
        { title: '技术风险', detail: '应对 LLM 幻觉和多 Agent 一致性问题，依赖 Debug Agent 与测试覆盖。' },
        { title: '性能风险', detail: '通过上下文压缩、RAG 检索、缓存与水平扩展缓解延迟和吞吐瓶颈。' },
        { title: '安全风险', detail: '执行环境沙箱化、权限最小化、输入校验与漏洞扫描是基础要求。' },
        { title: '合规风险', detail: '遵循隐私与版权法规，保留审计链路并提供用户可控的数据保留策略。' },
      ],
    },
    footer: {
      reference: '参考来源：deep-research-report.md（已做结构化提炼）。',
      notice: '免责声明：本页为产品说明与架构展示，不构成法律、合规或工程交付承诺。',
    },
  },
  en: {
    brand: 'Meta-Agent Studio',
    nav: [
      { id: 'architecture', label: 'Architecture' },
      { id: 'capabilities', label: 'Capabilities' },
      { id: 'api', label: 'API Examples' },
      { id: 'roadmap', label: 'Roadmap' },
      { id: 'risk', label: 'Risk Control' },
    ],
    hero: {
      kicker: 'Automated Multi-Agent Platform',
      title: 'Compile natural language goals into executable agent systems',
      summary:
        'This page is distilled from deep-research-report.md and focuses on a readable, trustworthy, and demo-ready product narrative. The first release is a static single-page experience.',
      highlights: ['NL -> IR -> AgentGraph', 'DAG / FSM / Planner Loop', 'RAG memory + blackboard communication', 'Security isolation + observability'],
    },
    architecture: {
      title: 'Architecture Overview',
      intro: 'The platform centers on an Agent Compiler that turns user intent into schedulable, traceable, and reusable execution graphs.',
      flowTitle: 'Core Flow',
      flow: [
        'User natural-language input',
        'Intent parsing and extraction',
        'Structured IR generation',
        'Assembler builds AgentGraph',
        'Orchestrator schedules execution',
        'Feedback and self-optimization loop',
      ],
      bundleTitle: 'Agent Composition',
      bundle: ['Role Template', 'Skill Modules', 'Tool Interfaces', 'State and Long-Term Memory'],
    },
    capabilities: {
      title: 'Core Capabilities',
      intro:
        'The platform combines orchestration, memory, communication, debugging, security, and observability to evolve from MVP to production scale.',
      blocks: [
        {
          title: 'Orchestration Modes',
          items: [
            'DAG: stable flow and easy debugging',
            'FSM: state-driven control for complex tasks',
            'Planner Loop: dynamic planning and autonomous dispatch',
          ],
        },
        {
          title: 'Context & Memory',
          items: ['Three-layer context: Global / Task / Agent', 'RAG-based retrieval injection', 'Memory Slots for structured updates'],
        },
        {
          title: 'Multi-Agent Communication',
          items: ['Shared blackboard state model', 'Publish-subscribe decoupling', 'Unified message schema and audit trail'],
        },
        {
          title: 'Self-Healing & Security',
          items: ['Debug Agent for anomaly detection', 'ACK + retry + circuit-breaker resilience', 'Tenant isolation and tool sandboxing'],
        },
      ],
    },
    api: {
      title: 'API Examples (Display Only)',
      intro: 'The endpoints below are simplified examples from the report. They are displayed only and are not called by this page.',
      table: {
        method: 'Method',
        endpoint: 'Endpoint',
        request: 'Request',
        response: 'Response',
      },
      rows: [
        {
          method: 'POST',
          endpoint: '/compile',
          request: '{"query":"Build an e-commerce website"}',
          response:
            '{"ir":{"goal":"Build an e-commerce website","tasks":[{"id":"t1","type":"analysis"},{"id":"t2","type":"backend"},{"id":"t3","type":"frontend"}]}}',
        },
        {
          method: 'POST',
          endpoint: '/assembleAgents',
          request: '{"ir":{...}}',
          response: '{"agents":[{"id":"A1","role":"Business Analyst"},{"id":"A2","role":"Backend Engineer"},{"id":"A3","role":"Frontend Engineer"}]}',
        },
        {
          method: 'POST',
          endpoint: '/executeTask',
          request: '{"agentId":"A2","taskId":"t2","input":"..."}',
          response: '{"status":"success","output":"Database setup completed"}',
        },
      ],
    },
    roadmap: {
      title: 'Roadmap & Delivery Rhythm',
      intro: 'The report recommends validating core mechanics in MVP first, then expanding automation, intelligence, and production scale.',
      phases: [
        { stage: 'Phase 1', name: 'PoC / MVP', goal: 'Fixed-role agents, simple DAG orchestration, and manual skill integration' },
        { stage: 'Phase 2', name: 'Automated Agents', goal: 'LLM-driven task decomposition and automatic agent generation' },
        { stage: 'Phase 3', name: 'Intelligent Optimization', goal: 'IR versioning and planner-based hybrid execution' },
        { stage: 'Phase 4', name: 'Production Scale', goal: 'Agent Store, multi-tenant support, and high-concurrency operation' },
      ],
    },
    risk: {
      title: 'Risks & Mitigation',
      intro: 'Building this platform requires balancing technical quality, performance, security boundaries, and compliance.',
      items: [
        { title: 'Technical Risk', detail: 'Address hallucinations and cross-agent consistency with Debug Agent workflows and test coverage.' },
        { title: 'Performance Risk', detail: 'Use context compression, RAG retrieval, caching, and horizontal scaling to reduce latency bottlenecks.' },
        { title: 'Security Risk', detail: 'Sandbox execution, least-privilege access, input validation, and continuous vulnerability scanning are mandatory.' },
        { title: 'Compliance Risk', detail: 'Respect privacy and IP requirements with auditable trails and user-controlled data retention policies.' },
      ],
    },
    footer: {
      reference: 'Source: deep-research-report.md (structured and productized).',
      notice: 'Disclaimer: This page is a product and architecture overview, not a legal or delivery commitment.',
    },
  },
}
