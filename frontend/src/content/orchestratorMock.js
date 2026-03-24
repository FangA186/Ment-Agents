export const modeTabs = [
  { id: 'dag', label: 'DAG 编排' },
  { id: 'fsm', label: 'FSM 状态机' },
  { id: 'planner', label: 'Planner Loop' },
]

export const orchestrationByMode = {
  dag: {
    intro: '按依赖拓扑顺序推进任务，适合确定性链路。',
    nodes: [
      { id: 'N1', title: '需求确认', owner: '分析Agent', status: 'done', dependsOn: '-' },
      { id: 'N2', title: '计划拆解', owner: '规划Agent', status: 'doing', dependsOn: 'N1' },
      { id: 'N3', title: '执行分发', owner: '执行Agent', status: 'todo', dependsOn: 'N2' },
      { id: 'N4', title: '质量复核', owner: '评估Agent', status: 'todo', dependsOn: 'N3' },
    ],
    logs: [
      '10:31:12 N1 -> done',
      '10:31:31 N2 -> doing',
      '10:31:44 等待 N2 完成后激活 N3',
    ],
  },
  fsm: {
    intro: '按状态迁移驱动流程，适合多分支与回退控制。',
    nodes: [
      { id: 'S1', title: '草案态 Draft', owner: '规划Agent', status: 'done', dependsOn: '-' },
      { id: 'S2', title: '审核态 Review', owner: '审计Agent', status: 'doing', dependsOn: 'S1' },
      { id: 'S3', title: '执行态 Execute', owner: '执行Agent', status: 'todo', dependsOn: 'S2' },
      { id: 'S4', title: '回滚态 Rollback', owner: '执行Agent', status: 'todo', dependsOn: 'S3' },
    ],
    logs: [
      '10:32:08 状态迁移 Draft -> Review',
      '10:32:26 审核中，等待风险项确认',
      '10:32:45 若失败率>8%，触发 Rollback',
    ],
  },
  planner: {
    intro: '动态规划 + 执行反馈闭环，适合开放型任务。',
    nodes: [
      { id: 'P1', title: '生成计划', owner: '规划Agent', status: 'done', dependsOn: '-' },
      { id: 'P2', title: '执行计划', owner: '执行Agent', status: 'doing', dependsOn: 'P1' },
      { id: 'P3', title: '评估反馈', owner: '评估Agent', status: 'todo', dependsOn: 'P2' },
      { id: 'P4', title: '迭代修正', owner: '规划Agent', status: 'todo', dependsOn: 'P3' },
    ],
    logs: [
      '10:33:02 Planner 生成第1版计划',
      '10:33:29 执行阶段反馈：策略A点击率偏低',
      '10:33:56 进入评估并准备下一轮修正',
    ],
  },
}
