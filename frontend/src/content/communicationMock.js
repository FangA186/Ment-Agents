export const channelTabs = [
  { id: 'all', label: '全部消息' },
  { id: 'blackboard', label: '黑板同步' },
  { id: 'task', label: '任务协同' },
  { id: 'alert', label: '异常告警' },
]

export const blackboardStates = [
  {
    id: 'bb-01',
    key: 'campaign.goal',
    value: '营销活动自动化系统',
    owner: '分析Agent',
    status: '已锁定',
    taskId: 'task-01',
  },
  {
    id: 'bb-02',
    key: 'campaign.constraints',
    value: '保留人工审核入口 / 支持回滚',
    owner: '审计Agent',
    status: '已同步',
    taskId: 'task-01',
  },
  {
    id: 'bb-03',
    key: 'execution.plan',
    value: '拆解 3 阶段任务并绑定执行人',
    owner: '规划Agent',
    status: '更新中',
    taskId: 'task-02',
  },
  {
    id: 'bb-04',
    key: 'quality.metrics',
    value: '准确率 >= 92%，可解释性 >= 85%',
    owner: '评估Agent',
    status: '待确认',
    taskId: 'task-03',
  },
]

export const messageTimeline = [
  {
    id: 'm-001',
    channel: 'blackboard',
    sender: '分析Agent',
    target: 'Blackboard',
    content: '写入 campaign.goal，并广播至规划/执行链路',
    time: '10:21:12',
    taskId: 'task-01',
  },
  {
    id: 'm-002',
    channel: 'task',
    sender: '规划Agent',
    target: '执行Agent',
    content: '已下发 task-02 子任务，请确认工具依赖',
    time: '10:21:44',
    taskId: 'task-02',
  },
  {
    id: 'm-003',
    channel: 'alert',
    sender: '审计Agent',
    target: '执行Agent',
    content: '检测到风险：回滚策略未附带触发阈值',
    time: '10:22:03',
    taskId: 'task-02',
  },
  {
    id: 'm-004',
    channel: 'blackboard',
    sender: '执行Agent',
    target: 'Blackboard',
    content: '回滚阈值已补充：失败率>8% 自动回退',
    time: '10:22:27',
    taskId: 'task-02',
  },
  {
    id: 'm-005',
    channel: 'task',
    sender: '评估Agent',
    target: '规划Agent',
    content: '请求补充评估指标采样窗口（当前缺失）',
    time: '10:23:15',
    taskId: 'task-03',
  },
]

export const taskLinks = [
  {
    id: 'link-01',
    taskId: 'task-01',
    title: '需求边界确认',
    messages: ['m-001'],
    relatedAgents: ['分析Agent', '审计Agent'],
    status: '已完成',
  },
  {
    id: 'link-02',
    taskId: 'task-02',
    title: '执行链路组装',
    messages: ['m-002', 'm-003', 'm-004'],
    relatedAgents: ['规划Agent', '执行Agent', '审计Agent'],
    status: '进行中',
  },
  {
    id: 'link-03',
    taskId: 'task-03',
    title: '质量评估闭环',
    messages: ['m-005'],
    relatedAgents: ['评估Agent', '规划Agent'],
    status: '待开始',
  },
]
