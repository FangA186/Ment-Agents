export const projectLibrary = [
  { id: 'proj-001', name: '营销活动自动化系统', status: '进行中', updatedAt: '2026-03-23 11:20' },
  { id: 'proj-002', name: '内容分发增长系统', status: '规划中', updatedAt: '2026-03-23 10:42' },
  { id: 'proj-003', name: '售后工单自动分派', status: '进行中', updatedAt: '2026-03-23 09:58' },
  { id: 'proj-004', name: '多渠道投放联动', status: '已归档', updatedAt: '2026-03-22 18:14' },
]

const memoryByProject = {
  'proj-001': {
    layers: [
      {
        id: 'global',
        name: 'Global Memory',
        desc: '全局策略、合规规则、通用知识',
        slots: [
          {
            id: 'g-01',
            key: 'policy.audit',
            value: '所有执行链路必须可追踪、可回滚',
            updatedAt: '2026-03-23 09:21:05',
            source: '审计中心',
          },
          {
            id: 'g-02',
            key: 'domain.lexicon',
            value: '营销术语词典（v3）',
            updatedAt: '2026-03-23 09:22:44',
            source: '知识库',
          },
        ],
      },
      {
        id: 'task',
        name: 'Task Memory',
        desc: '任务上下文、阶段目标、中间产物',
        slots: [
          {
            id: 't-01',
            key: 'task-02.plan',
            value: '执行链路拆解为 3 个可回滚步骤',
            updatedAt: '2026-03-23 10:04:11',
            source: '规划Agent',
          },
          {
            id: 't-02',
            key: 'task-02.constraints',
            value: '失败率 > 8% 触发回滚',
            updatedAt: '2026-03-23 10:05:32',
            source: '审计Agent',
          },
        ],
      },
      {
        id: 'agent',
        name: 'Agent Memory',
        desc: 'Agent 个体偏好、历史反馈、局部经验',
        slots: [
          {
            id: 'a-01',
            key: 'planner.preference',
            value: '偏向分阶段任务而非并行爆发',
            updatedAt: '2026-03-23 10:11:08',
            source: '规划Agent',
          },
          {
            id: 'a-02',
            key: 'executor.retry_pattern',
            value: '优先重试工具调用，最多 2 次',
            updatedAt: '2026-03-23 10:11:56',
            source: '执行Agent',
          },
        ],
      },
    ],
    rag: [
      {
        id: 'r-01',
        query: '人工审核入口约束',
        hit: 'policy.audit',
        score: 0.92,
        snippet: '所有执行链路必须保留人工审核入口，且可回滚。',
      },
      {
        id: 'r-02',
        query: '失败阈值回滚策略',
        hit: 'task-02.constraints',
        score: 0.88,
        snippet: '失败率 > 8% 时触发自动回滚并通知审计Agent。',
      },
      {
        id: 'r-03',
        query: '执行重试策略',
        hit: 'executor.retry_pattern',
        score: 0.84,
        snippet: '工具调用异常优先重试，两次失败进入降级通道。',
      },
    ],
  },
  'proj-002': {
    layers: [
      {
        id: 'global',
        name: 'Global Memory',
        desc: '全局策略、合规规则、通用知识',
        slots: [
          {
            id: 'g2-01',
            key: 'policy.content_quality',
            value: '敏感词命中即拦截发布',
            updatedAt: '2026-03-23 10:12:19',
            source: '内容审计中心',
          },
        ],
      },
      {
        id: 'task',
        name: 'Task Memory',
        desc: '任务上下文、阶段目标、中间产物',
        slots: [
          {
            id: 't2-01',
            key: 'task-11.schedule',
            value: '按平台时区进行错峰发布',
            updatedAt: '2026-03-23 10:18:42',
            source: '调度Agent',
          },
        ],
      },
      {
        id: 'agent',
        name: 'Agent Memory',
        desc: 'Agent 个体偏好、历史反馈、局部经验',
        slots: [
          {
            id: 'a2-01',
            key: 'publisher.best_slot',
            value: '短视频平台晚间 20:00-22:00 转化更高',
            updatedAt: '2026-03-23 10:22:01',
            source: '发布Agent',
          },
        ],
      },
    ],
    rag: [
      {
        id: 'r2-01',
        query: '内容发布时间段',
        hit: 'publisher.best_slot',
        score: 0.86,
        snippet: '短视频平台晚间 20:00-22:00 转化更高。',
      },
    ],
  },
  'proj-003': {
    layers: [
      {
        id: 'global',
        name: 'Global Memory',
        desc: '全局策略、合规规则、通用知识',
        slots: [
          {
            id: 'g3-01',
            key: 'sla.default',
            value: 'P1工单 15 分钟内响应',
            updatedAt: '2026-03-23 09:40:03',
            source: '客服运营',
          },
        ],
      },
      {
        id: 'task',
        name: 'Task Memory',
        desc: '任务上下文、阶段目标、中间产物',
        slots: [
          {
            id: 't3-01',
            key: 'task-24.routing_rule',
            value: '按故障类型 + 地域双重路由',
            updatedAt: '2026-03-23 09:46:55',
            source: '分派Agent',
          },
        ],
      },
      {
        id: 'agent',
        name: 'Agent Memory',
        desc: 'Agent 个体偏好、历史反馈、局部经验',
        slots: [
          {
            id: 'a3-01',
            key: 'dispatcher.fallback',
            value: '主责任组满载时自动转入备份组',
            updatedAt: '2026-03-23 09:51:13',
            source: '调度Agent',
          },
        ],
      },
    ],
    rag: [
      {
        id: 'r3-01',
        query: '高优先工单SLA',
        hit: 'sla.default',
        score: 0.91,
        snippet: 'P1工单 15 分钟内响应。',
      },
    ],
  },
  'proj-004': {
    layers: [
      {
        id: 'global',
        name: 'Global Memory',
        desc: '全局策略、合规规则、通用知识',
        slots: [
          {
            id: 'g4-01',
            key: 'budget.guard',
            value: '单渠道预算波动不得超过 20%',
            updatedAt: '2026-03-22 17:11:36',
            source: '风控中心',
          },
        ],
      },
      {
        id: 'task',
        name: 'Task Memory',
        desc: '任务上下文、阶段目标、中间产物',
        slots: [
          {
            id: 't4-01',
            key: 'task-31.sync_window',
            value: '按小时同步渠道数据',
            updatedAt: '2026-03-22 17:25:08',
            source: '同步Agent',
          },
        ],
      },
      {
        id: 'agent',
        name: 'Agent Memory',
        desc: 'Agent 个体偏好、历史反馈、局部经验',
        slots: [
          {
            id: 'a4-01',
            key: 'optimizer.safe_mode',
            value: '波动异常时自动进入保守策略',
            updatedAt: '2026-03-22 17:28:19',
            source: '优化Agent',
          },
        ],
      },
    ],
    rag: [
      {
        id: 'r4-01',
        query: '预算防护规则',
        hit: 'budget.guard',
        score: 0.9,
        snippet: '单渠道预算波动不得超过 20%。',
      },
    ],
  },
}

export function getProjectMemory(projectId) {
  return memoryByProject[projectId] || memoryByProject['proj-001']
}
