import { reactive } from 'vue'
import { apiAssemble, apiCompile, apiGetArtifacts, apiGetChat, apiGetProject, apiSendChat } from '../services/api'

const defaultAgents = [
  { id: 'analyst', name: '分析Agent', state: '已确认', color: 'gold' },
  { id: 'planner', name: '规划Agent', state: '待实例化', color: 'blue' },
  { id: 'executor', name: '执行Agent', state: '待实例化', color: 'green' },
  { id: 'evaluator', name: '评估Agent', state: '待实例化', color: 'gold' },
  { id: 'auditor', name: '审计Agent', state: '待实例化', color: 'gray' },
  { id: 'reporter', name: '汇总Agent', state: '待实例化', color: 'gray' },
]

const agentProfiles = {
  planner: {
    id: 'planner',
    name: '规划Agent',
    summary: '将业务目标拆解为可执行阶段任务',
    style: '稳健 / 可审计 / 可回滚',
    status: '就绪',
    roleFile: 'ROLE.md',
    roleDesc: '定义身份边界与决策原则',
    skillFile: 'SKILL.md',
    skillDesc: '包含输入输出规范与执行策略',
    tools: [
      { icon: '检', name: '检索工具', desc: '用于补全上下文与事实校验' },
      { icon: '排', name: '任务编排器', desc: '生成阶段计划与依赖关系' },
      { icon: '规', name: '规则引擎', desc: '执行策略约束与风控检查' },
    ],
  },
  executor: {
    id: 'executor',
    name: '执行Agent',
    summary: '根据计划自动执行任务并同步状态',
    style: '高效 / 可恢复 / 可追踪',
    status: '待命',
    roleFile: 'ROLE.md',
    roleDesc: '定义执行边界与异常处理权限',
    skillFile: 'SKILL.md',
    skillDesc: '封装执行动作与结果回传规范',
    tools: [
      { icon: '执', name: '执行器', desc: '触发任务与收集执行回执' },
      { icon: '流', name: '流程网关', desc: '管理节点顺序与并行开关' },
      { icon: '志', name: '日志管道', desc: '记录执行日志与链路追踪' },
    ],
  },
  evaluator: {
    id: 'evaluator',
    name: '评估Agent',
    summary: '评估产出质量并给出优化建议',
    style: '客观 / 可解释 / 可复用',
    status: '待命',
    roleFile: 'ROLE.md',
    roleDesc: '定义评分维度、阈值与回退条件',
    skillFile: 'SKILL.md',
    skillDesc: '提供评分算法与优化建议模板',
    tools: [
      { icon: '评', name: '质量评分器', desc: '按指标体系输出评分结果' },
      { icon: '比', name: '版本对比器', desc: '对比当前与历史方案差异' },
      { icon: '档', name: '归档器', desc: '归档评估记录便于复盘' },
    ],
  },
}

function mapMessage(item) {
  const speaker = item.role === 'user' ? 'user' : 'agent'
  return {
    id: item.id,
    speaker,
    title: speaker === 'user' ? '你' : '分析Agent',
    text: item.content,
    createdAt: item.createdAt,
  }
}

export const collabFlowState = reactive({
  requirement: '',
  conversation: [],
  projectName: '项目加载中',
  projectStatus: '规划中',
  irReady: false,
  irFileName: 'IR.json',
  irSource: '待生成',
  selectedAgentId: 'planner',
  agents: [...defaultAgents],
  agentGraph: null,
  loading: false,
  error: '',
})

export async function loadProjectContext(projectId) {
  collabFlowState.loading = true
  collabFlowState.error = ''

  try {
    const [project, messages, artifacts] = await Promise.all([
      apiGetProject(projectId),
      apiGetChat(projectId),
      apiGetArtifacts(projectId),
    ])

    collabFlowState.projectName = project?.name || '未命名项目'
    collabFlowState.projectStatus = project?.status || '规划中'
    collabFlowState.conversation = (messages || []).map(mapMessage)
    collabFlowState.requirement = collabFlowState.conversation.find((item) => item.speaker === 'user')?.text || ''
    collabFlowState.irReady = Boolean(artifacts?.ir)
    collabFlowState.irSource = artifacts?.ir ? `版本 ${artifacts.ir.version} · ${artifacts.ir.generatedAt}` : '待生成'
    collabFlowState.agentGraph = artifacts?.agentGraph || null
  } catch (error) {
    collabFlowState.error = error.message || '加载项目上下文失败'
  } finally {
    collabFlowState.loading = false
  }
}

export async function submitRequirement(projectId, nextRequirement) {
  const text = (nextRequirement || '').trim()
  if (!text) return

  collabFlowState.requirement = text
  collabFlowState.irReady = false
  collabFlowState.irSource = '待重新确认'

  const messages = await apiSendChat(projectId, text)
  collabFlowState.conversation = messages.map(mapMessage)
}

export async function appendConversation(projectId, userInput) {
  const text = (userInput || '').trim()
  if (!text) return

  const messages = await apiSendChat(projectId, text)
  collabFlowState.conversation = messages.map(mapMessage)
}

export async function finalizeIR(projectId) {
  const ir = await apiCompile(projectId)
  collabFlowState.irReady = true
  collabFlowState.irFileName = 'IR.json'
  collabFlowState.irSource = ir ? `版本 ${ir.version} · ${ir.generatedAt}` : '已生成'
}

export async function assembleAgents(projectId) {
  const graph = await apiAssemble(projectId)
  collabFlowState.agentGraph = graph
  return graph
}

export function setSelectedAgent(id) {
  collabFlowState.selectedAgentId = id
}

export function getAgentProfile(id) {
  return agentProfiles[id] || agentProfiles.planner
}
