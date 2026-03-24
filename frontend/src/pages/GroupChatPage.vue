<script setup>
import { computed, nextTick, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  appendConversation,
  assembleAgents,
  collabFlowState,
  finalizeIR,
  loadProjectContext,
  setSelectedAgent,
  submitRequirement,
} from '../content/collabFlowMock'

const router = useRouter()
const route = useRoute()
const inputText = ref('')
const messageListRef = ref(null)

const centerIrOpen = ref(true)
const generationVisible = ref(false)
const generationRunning = ref(false)
const generationSteps = ref([
  { id: 'confirm', label: '需求确认中', status: 'pending' },
  { id: 'ir', label: '生成 IR 中', status: 'pending' },
  { id: 'template', label: '预定义模板匹配中', status: 'pending' },
  { id: 'agent', label: '生成 Agents 实例图中', status: 'pending' },
  { id: 'done', label: '生成完成，即将前往组装器页面', status: 'pending' },
])

function syncRequirement() {
  const normalized = inputText.value.trim().replace(/^需求[:：]\s*/, '')
  if (!normalized) return
  const projectId = String(route.params.projectId || 'proj-001')
  submitRequirement(projectId, normalized).catch((error) => {
    console.error(error)
  })
  inputText.value = ''
}

function sendMessage() {
  const trimmed = inputText.value.trim()
  if (!trimmed) return

  if (/^需求[:：]/.test(trimmed)) {
    syncRequirement()
    return
  }

  const projectId = String(route.params.projectId || 'proj-001')
  appendConversation(projectId, trimmed).catch((error) => {
    console.error(error)
  })
  inputText.value = ''
}

function produceIR() {
  runGenerationFlow()
}

function openAgent(agentId) {
  setSelectedAgent(agentId)
  const projectId = String(route.params.projectId || 'proj-001')
  router.push(`/studio/projects/${projectId}/agent/${agentId}`)
}

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

function resetGenerationSteps() {
  generationSteps.value = generationSteps.value.map((step) => ({ ...step, status: 'pending' }))
}

async function runGenerationFlow() {
  if (generationRunning.value) return

  const projectId = String(route.params.projectId || 'proj-001')
  generationRunning.value = true
  generationVisible.value = true
  resetGenerationSteps()

  try {
    for (const step of generationSteps.value) {
      step.status = 'doing'
      if (step.id === 'ir') {
        collabFlowState.irReady = false
      }
      await sleep(750)
      if (step.id === 'ir') {
        await finalizeIR(projectId)
      }
      if (step.id === 'agent') {
        await assembleAgents(projectId)
      }
      step.status = 'done'
    }

    await sleep(450)
    router.push(`/studio/projects/${projectId}/assembler`)
  } catch (error) {
    console.error(error)
  } finally {
    generationRunning.value = false
  }
}

watch(
  () => collabFlowState.conversation.length,
  async () => {
    await nextTick()
    if (messageListRef.value) {
      messageListRef.value.scrollTop = messageListRef.value.scrollHeight
    }
  },
  { immediate: true },
)

watch(
  () => route.params.projectId,
  (projectId) => {
    const pid = String(projectId || 'proj-001')
    loadProjectContext(pid).catch((error) => {
      console.error(error)
    })
  },
  { immediate: true },
)
</script>

<template>
  <section class="chat-page">
    <div class="chat-head">
      <h2>项目：{{ collabFlowState.projectName }}</h2>
      <button class="btn" type="button" :disabled="generationRunning" @click="produceIR">
        {{ generationRunning ? '生成流程执行中...' : collabFlowState.irReady ? '重新生成 IR' : '确认并生成 IR' }}
      </button>
    </div>

    <div class="chat-layout">
      <article class="col conversation-col">
        <h3>群聊对话（你 + 分析Agent）</h3>

        <section v-if="generationVisible" class="generation-flow">
          <p class="flow-title">生成流程</p>
          <div class="flow-steps">
            <article v-for="step in generationSteps" :key="step.id" class="flow-step" :class="step.status">
              <span class="flow-dot"></span>
              <p>{{ step.label }}</p>
              <small>{{ step.status === 'done' ? '已完成' : step.status === 'doing' ? '进行中' : '等待中' }}</small>
            </article>
          </div>
        </section>

        <div ref="messageListRef" class="message-list">
          <article v-for="msg in collabFlowState.conversation" :key="msg.id" class="msg" :class="msg.speaker">
            <p class="msg-title">{{ msg.title }}</p>
            <p class="msg-body">{{ msg.text }}</p>
          </article>
        </div>

        <div class="chat-input-row">
          <input v-model="inputText" type="text" placeholder="输入需求或补充信息（需求请以“需求:”开头）" @keydown.enter="sendMessage" />
          <button type="button" @click="sendMessage">发送</button>
        </div>
      </article>

      <article class="col center-col">
        <h3>群聊信息 / Agent状态</h3>
        <div class="agent-wall">
          <button v-for="agent in collabFlowState.agents" :key="agent.id" type="button" class="agent-square" @click="openAgent(agent.id)">
            <span>{{ agent.name }}</span>
            <small>{{ agent.state }}</small>
          </button>
        </div>

        <section class="mini-result-card">
          <div class="mini-head">
            <h4>IR产出成果</h4>
            <button class="mini-toggle" type="button" @click="centerIrOpen = !centerIrOpen">{{ centerIrOpen ? '收起' : '展开' }}</button>
          </div>
          <div v-if="centerIrOpen" class="mini-body">
            <p class="step-text">第一步：结构化 IR 文件</p>
            <div class="result-file compact">
              <div class="file-icon"></div>
              <div>
                <p class="result-main">{{ collabFlowState.irFileName }}</p>
                <p class="result-sub">{{ collabFlowState.irSource }}</p>
              </div>
            </div>
          </div>
        </section>
      </article>
    </div>
  </section>
</template>

<style scoped>
.chat-page {
  width: 100%;
  height: 100%;
  min-height: 0;
  display: grid;
  grid-template-rows: auto 1fr;
  gap: 12px;
}

.chat-head {
  border-radius: 14px;
  padding: 10px 14px;
  background: rgba(20, 18, 14, 0.75);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.chat-head h2 {
  margin: 0;
  font-size: 30px;
}

.chat-layout {
  height: 100%;
  min-height: 0;
  display: grid;
  grid-template-columns: 7fr 3fr;
  gap: 12px;
}

.col {
  min-height: 0;
  border-radius: 16px;
  background: rgba(21, 18, 14, 0.8);
  padding: 12px;
  overflow: hidden;
}

.conversation-col {
  display: grid;
  grid-template-rows: auto 1fr auto;
  gap: 10px;
}

.center-col {
  display: grid;
  grid-template-rows: auto auto auto;
  gap: 10px;
  align-content: start;
  overflow-y: auto;
  padding-right: 6px;
}

.center-col::-webkit-scrollbar {
  width: 8px;
}

.center-col::-webkit-scrollbar-thumb {
  border-radius: 999px;
  background: rgba(212, 175, 55, 0.35);
}

h3,
h4 {
  margin: 0;
  color: #f2dc9f;
}

h3 {
  font-size: 24px;
}

h4 {
  font-size: 32px;
}

.message-list {
  min-height: 0;
  overflow-y: auto;
  padding-right: 4px;
  display: grid;
  align-content: start;
  gap: 10px;
}

.message-list::-webkit-scrollbar {
  width: 8px;
}

.message-list::-webkit-scrollbar-thumb {
  border-radius: 999px;
  background: rgba(212, 175, 55, 0.35);
}

.msg {
  max-width: 86%;
  border-radius: 14px;
  padding: 10px 12px;
  background: rgba(30, 25, 18, 0.84);
}

.msg.user {
  justify-self: end;
  background: rgba(60, 50, 35, 0.78);
}

.msg.agent {
  justify-self: start;
}

.msg-title {
  margin: 0;
  font-weight: 700;
  color: #f4e3b8;
}

.msg-body {
  margin: 6px 0 0;
  font-size: 17px;
}

.chat-input-row {
  border-radius: 12px;
  background: rgba(14, 12, 10, 0.9);
  padding: 8px;
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 8px;
}

.chat-input-row input {
  width: 100%;
  border: none;
  outline: none;
  background: transparent;
  color: var(--text);
  padding: 8px 10px;
  font-size: 16px;
}

.chat-input-row button {
  border: none;
  border-radius: 10px;
  background: linear-gradient(145deg, #8a6b24, var(--accent-2));
  color: #1a140c;
  padding: 0 14px;
  font-weight: 700;
  cursor: pointer;
}

.agent-wall {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
}

.agent-wall {
  align-content: start;
  grid-auto-rows: 112px;
}

.agent-square {
  border: none;
  border-radius: 16px;
  background: rgba(16, 14, 12, 0.95);
  color: var(--text);
  display: grid;
  place-content: center;
  gap: 4px;
  text-align: center;
  cursor: pointer;
  padding: 8px;
}

.agent-square span {
  font-size: 15px;
  font-weight: 700;
}

.agent-square small {
  color: var(--muted);
  font-size: 12px;
}

.mini-result-card {
  border-radius: 12px;
  background: rgba(22, 19, 15, 0.88);
  padding: 10px;
  display: grid;
  gap: 8px;
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.06);
}

.mini-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}

.mini-toggle {
  border: none;
  background: transparent;
  color: #f2dc9f;
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
}

.mini-body {
  display: grid;
  gap: 8px;
}

.step-text {
  margin: 0;
  font-size: 14px;
  font-weight: 700;
}

.result-file {
  border-radius: 12px;
  background: rgba(13, 12, 10, 0.92);
  padding: 10px;
  display: grid;
  grid-template-columns: 58px 1fr;
  gap: 10px;
  align-items: center;
}

.result-file.compact {
  grid-template-columns: 46px 1fr;
}

.file-icon {
  width: 50px;
  height: 58px;
  border: 1px solid rgba(212, 175, 55, 0.3);
  border-radius: 4px;
  background: linear-gradient(180deg, rgba(70, 60, 44, 0.5), rgba(40, 35, 26, 0.8));
  position: relative;
}

.result-file.compact .file-icon {
  width: 40px;
  height: 50px;
}

.file-icon::after {
  content: '';
  position: absolute;
  right: -1px;
  top: -1px;
  width: 16px;
  height: 16px;
  clip-path: polygon(0 0, 100% 100%, 100% 0);
  background: rgba(128, 109, 72, 0.8);
}

.result-main {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
}

.result-sub {
  margin: 4px 0 0;
  font-size: 13px;
  color: var(--muted);
}

.generation-flow {
  border-radius: 12px;
  background: rgba(13, 12, 10, 0.9);
  padding: 10px;
  display: grid;
  gap: 8px;
}

.flow-title {
  margin: 0;
  color: #f2dc9f;
  font-size: 14px;
  font-weight: 700;
}

.flow-steps {
  display: grid;
  gap: 8px;
}

.flow-step {
  border-radius: 10px;
  background: rgba(24, 20, 16, 0.9);
  padding: 8px 10px;
  display: grid;
  grid-template-columns: auto 1fr auto;
  align-items: center;
  gap: 8px;
}

.flow-dot {
  width: 9px;
  height: 9px;
  border-radius: 999px;
  background: rgba(184, 171, 144, 0.45);
}

.flow-step p {
  margin: 0;
  font-size: 13px;
}

.flow-step small {
  color: var(--muted);
  font-size: 12px;
}

.flow-step.doing .flow-dot {
  background: rgba(121, 168, 198, 0.8);
}

.flow-step.done .flow-dot {
  background: rgba(132, 194, 157, 0.8);
}

@media (max-width: 1680px) {
  .chat-layout {
    grid-template-columns: 7fr 3fr;
  }
}

@media (max-width: 1200px) {
  .chat-page {
    height: auto;
  }

  .chat-layout {
    grid-template-columns: 1fr;
    grid-template-rows: auto;
    height: auto;
  }

  .col {
    min-height: 420px;
  }
}
</style>

