/**
 * @typedef {'todo' | 'in_progress' | 'done'} ConceptStatus
 */

/**
 * @typedef {{
 *   id: string
 *   name: string
 *   keywords: string[]
 *   imagePaths: string[]
 *   strengths: string[]
 *   risks: string[]
 *   recommendedScene: string
 *   status: ConceptStatus
 * }} DesignConcept
 */

/**
 * @typedef {{
 *   taskId: string
 *   title: string
 *   status: 'pending' | 'running' | 'done' | 'error'
 *   owner: string
 * }} WorkbenchTask
 */

/**
 * @typedef {{
 *   inputPrompt: string
 *   irPreview: string
 *   logs: string[]
 *   nextAction: string
 * }} CompilePanelMockState
 */

/**
 * @type {DesignConcept[]}
 */
export const designConcepts = [
  {
    id: 'A',
    name: '模块解构式工作台',
    keywords: ['非对称分栏', '流程胶囊', '高对比信息层'],
    imagePaths: [
      'reports/ui-proposals/concept-a-workbench.svg',
      'reports/ui-proposals/concept-a-compiler.svg',
    ],
    strengths: ['流程导向清晰', '状态反馈直观', '适合快速上手'],
    risks: ['在低分辨率下需要严格控制信息密度'],
    recommendedScene: '强调任务链路与实时状态反馈',
    status: 'done',
  },
  {
    id: 'B',
    name: '时间轴编排式工作台',
    keywords: ['任务流舞台', '状态层叠', '节奏导览'],
    imagePaths: [
      'reports/ui-proposals/concept-b-workbench.svg',
      'reports/ui-proposals/concept-b-compiler.svg',
    ],
    strengths: ['阶段推进感强', '汇报展示友好', '流程状态明显'],
    risks: ['并行任务扩展时需要额外可视化规则'],
    recommendedScene: '强调编排时序与阶段推进',
    status: 'done',
  },
  {
    id: 'C',
    name: '图谱操作台式工作台',
    keywords: ['节点画布', '语义连线', '侧边编译面板'],
    imagePaths: [
      'reports/ui-proposals/concept-c-workbench.svg',
      'reports/ui-proposals/concept-c-compiler.svg',
    ],
    strengths: ['关系表达力强', '适合复杂依赖', '可视化辨识度高'],
    risks: ['首次使用存在学习成本'],
    recommendedScene: '强调 Agent 关系网和策略联动',
    status: 'done',
  },
]

/**
 * @type {WorkbenchTask[]}
 */
export const workbenchMockState = [
  { taskId: 'T1', title: '解析需求并生成 IR', status: 'running', owner: 'Compiler Agent' },
  { taskId: 'T2', title: '组装 Agent 模板链', status: 'pending', owner: 'Assembler Agent' },
  { taskId: 'T3', title: '执行任务并回收指标', status: 'pending', owner: 'Orchestrator Agent' },
]

/**
 * @type {CompilePanelMockState}
 */
export const compilePanelMockState = {
  inputPrompt: '请构建一个自动化营销多 Agent 流程，包含调研、文案、投放和复盘。',
  irPreview: '{\"goal\":\"自动化营销\",\"tasks\":[{\"id\":\"t1\",\"role\":\"analysis\"},{\"id\":\"t2\",\"role\":\"copywriting\"}]}',
  logs: ['[10:21:08] 意图识别完成', '[10:21:09] IR 结构生成完成', '[10:21:10] 等待 Agent 组装'],
  nextAction: '触发组装并进入执行预演',
}

