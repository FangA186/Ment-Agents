export const workbenchTimeline = [
  { id: 'T1', title: '解析需求并生成 IR', owner: 'Compiler Agent', status: 'running', detail: '提取目标、约束与实体。' },
  { id: 'T2', title: '组装 Agent 模板链', owner: 'Assembler Agent', status: 'pending', detail: '绑定角色、技能、工具接口。' },
  { id: 'T3', title: '触发编排执行', owner: 'Orchestrator Agent', status: 'pending', detail: '按时间轴推进任务状态。' },
  { id: 'T4', title: '回收观测与复盘', owner: 'Debug Agent', status: 'pending', detail: '输出指标、日志、追踪信息。' },
]

export const compilerPreset = {
  defaultPrompt: '请构建一个多 Agent 电商增长流程，包含调研、文案、投放和复盘。',
  successIR: `{
  "version": "1.0",
  "goal": "电商增长自动化",
  "tasks": [
    { "id": "t1", "role": "analysis", "status": "done" },
    { "id": "t2", "role": "copywriting", "status": "running" },
    { "id": "t3", "role": "ad_ops", "status": "pending" }
  ],
  "tools": ["web_search", "sql_query", "campaign_api"]
}`,
}

export const assemblerMock = {
  tasks: [
    { id: 't1', name: '市场分析', priority: '高' },
    { id: 't2', name: '文案生成', priority: '中' },
    { id: 't3', name: '投放执行', priority: '中' },
  ],
  agents: [
    {
      taskId: 't1',
      agentId: 'A1',
      roleTemplate: '分析师模板 v1',
      skills: ['检索增强', '数据清洗'],
      tools: ['Web检索', 'SQL查询'],
      memory: {
        work: ['当前轮问题：市场洞察', '输入约束：预算<50k'],
        long: ['历史行业报告向量记忆', '用户偏好：偏数据驱动'],
      },
    },
    {
      taskId: 't2',
      agentId: 'A2',
      roleTemplate: '文案模板 v2',
      skills: ['提示词工程', '风格对齐'],
      tools: ['模板库', '品牌词典'],
      memory: {
        work: ['当前轮目标：输出3版文案', '风格：理性简洁'],
        long: ['品牌语料长期记忆', '历史转化高文案样本'],
      },
    },
    {
      taskId: 't3',
      agentId: 'A3',
      roleTemplate: '运营模板 v1',
      skills: ['渠道投放', '预算控制'],
      tools: ['广告API', '报表中心'],
      memory: {
        work: ['当前轮目标：分渠道投放', 'KPI：CTR > 3%'],
        long: ['历史投放策略库', '异常告警处理记录'],
      },
    },
  ],
  graph: {
    nodes: [
      { id: 'A1', label: '分析师', state: 'ready' },
      { id: 'A2', label: '文案', state: 'running' },
      { id: 'A3', label: '运营', state: 'pending' },
    ],
    edges: [
      { from: 'A1', to: 'A2' },
      { from: 'A2', to: 'A3' },
    ],
  },
}
