<script setup>
import { computed, onBeforeUnmount, ref } from 'vue'
import { modeTabs, orchestrationByMode } from '../content/orchestratorMock'

const activeMode = ref('dag')
const isAutoPlaying = ref(false)
const rows = ref(orchestrationByMode[activeMode.value].nodes.map((n) => ({ ...n })))
const logs = ref([...orchestrationByMode[activeMode.value].logs])

let autoTimer = null

const current = computed(() => orchestrationByMode[activeMode.value])

const links = computed(() =>
  rows.value
    .filter((row) => row.dependsOn && row.dependsOn !== '-')
    .map((row) => ({ from: row.dependsOn, to: row.id })),
)

function switchMode(modeId) {
  stopAuto()
  activeMode.value = modeId
  rows.value = orchestrationByMode[modeId].nodes.map((n) => ({ ...n }))
  logs.value = [...orchestrationByMode[modeId].logs]
}

function statusLabel(status) {
  if (status === 'done') return '已完成'
  if (status === 'doing') return '执行中'
  return '待执行'
}

function statusClass(status) {
  if (status === 'done') return 'done'
  if (status === 'doing') return 'doing'
  return 'todo'
}

function resetFlow() {
  stopAuto()
  rows.value = orchestrationByMode[activeMode.value].nodes.map((n) => ({ ...n }))
  logs.value = [...orchestrationByMode[activeMode.value].logs]
}

function runStep(appendAutoTag = false) {
  const currentDoing = rows.value.findIndex((x) => x.status === 'doing')
  const nextTodo = rows.value.findIndex((x) => x.status === 'todo')

  if (currentDoing === -1 && nextTodo === -1) {
    stopAuto()
    logs.value = [`${timeNow()} 已到达终态，无可推进节点`, ...logs.value].slice(0, 12)
    return
  }

  if (currentDoing >= 0) {
    rows.value[currentDoing].status = 'done'
  }

  if (nextTodo >= 0) {
    rows.value[nextTodo].status = 'doing'
  }

  const tag = appendAutoTag ? '（自动）' : '（手动）'
  logs.value = [`${timeNow()} 推进一步 ${tag}`, ...logs.value].slice(0, 12)
}

function toggleAuto() {
  if (isAutoPlaying.value) {
    stopAuto()
    return
  }

  isAutoPlaying.value = true
  logs.value = [`${timeNow()} 自动播放已启动`, ...logs.value].slice(0, 12)
  autoTimer = setInterval(() => runStep(true), 1200)
}

function stopAuto() {
  if (autoTimer) {
    clearInterval(autoTimer)
    autoTimer = null
  }
  if (isAutoPlaying.value) {
    isAutoPlaying.value = false
    logs.value = [`${timeNow()} 自动播放已暂停`, ...logs.value].slice(0, 12)
  }
}

function timeNow() {
  return new Date().toLocaleTimeString('zh-CN', { hour12: false })
}

onBeforeUnmount(() => {
  stopAuto()
})
</script>

<template>
  <section class="panel orchestrator-page">
    <div class="panel-head">
      <h2>编排执行页面</h2>
      <div class="panel-actions">
        <button class="btn ghost" type="button" @click="resetFlow">重置流程</button>
        <button class="btn ghost" type="button" @click="runStep()">推进一步</button>
        <button class="btn" type="button" @click="toggleAuto">{{ isAutoPlaying ? '暂停自动播放' : '启动自动播放' }}</button>
      </div>
    </div>

    <p class="panel-intro">支持 DAG / FSM / Planner 三种编排模式，展示任务状态、依赖连线与执行日志。</p>

    <div class="mode-tabs">
      <button
        v-for="tab in modeTabs"
        :key="tab.id"
        type="button"
        :class="{ active: tab.id === activeMode }"
        @click="switchMode(tab.id)"
      >
        {{ tab.label }}
      </button>
    </div>

    <article class="mode-intro">{{ current.intro }}</article>

    <div class="orchestrator-grid">
      <article class="orc-card">
        <h3>任务状态流转</h3>
        <div class="flow-table">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>节点</th>
                <th>负责人</th>
                <th>依赖</th>
                <th>状态</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="row in rows" :key="row.id">
                <td><code>{{ row.id }}</code></td>
                <td>{{ row.title }}</td>
                <td>{{ row.owner }}</td>
                <td>{{ row.dependsOn }}</td>
                <td>
                  <span class="status-pill" :class="statusClass(row.status)">
                    {{ statusLabel(row.status) }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </article>

      <div class="right-stack">
        <article class="orc-card">
          <h3>依赖连线（简版）</h3>
          <div class="graph-list">
            <div v-for="row in rows" :key="`node-${row.id}`" class="graph-node" :class="statusClass(row.status)">
              <strong>{{ row.id }}</strong>
              <span>{{ row.title }}</span>
            </div>
            <div v-for="link in links" :key="`${link.from}-${link.to}`" class="graph-edge">{{ link.from }} -> {{ link.to }}</div>
          </div>
        </article>

        <article class="orc-card">
          <h3>执行日志</h3>
          <div class="log-panel">
            <p v-for="line in logs" :key="line">{{ line }}</p>
          </div>
        </article>
      </div>
    </div>
  </section>
</template>

<style scoped>
.orchestrator-page {
  display: grid;
  gap: 12px;
  min-height: 0;
}

.mode-tabs {
  display: inline-flex;
  gap: 8px;
  flex-wrap: wrap;
}

.mode-tabs button {
  border: none;
  border-radius: 999px;
  padding: 7px 12px;
  color: var(--muted);
  background: rgba(20, 18, 14, 0.84);
  cursor: pointer;
}

.mode-tabs button.active {
  color: #1a140c;
  background: linear-gradient(145deg, #8a6b24, var(--accent-2));
  font-weight: 700;
}

.mode-intro {
  border-radius: 12px;
  background: rgba(16, 14, 11, 0.84);
  padding: 10px 12px;
  color: var(--muted);
}

.orchestrator-grid {
  min-height: 0;
  display: grid;
  grid-template-columns: 1.3fr 0.9fr;
  gap: 12px;
}

.right-stack {
  min-height: 0;
  display: grid;
  grid-template-rows: 1fr 1fr;
  gap: 12px;
}

.orc-card {
  border-radius: 14px;
  background: rgba(16, 14, 11, 0.88);
  padding: 12px;
  min-height: 0;
  display: grid;
  grid-template-rows: auto 1fr;
  gap: 8px;
}

.orc-card h3 {
  margin: 0;
  color: #f2dc9f;
  font-size: 21px;
}

.flow-table {
  min-height: 0;
  overflow: auto;
  border-radius: 12px;
}

table {
  width: 100%;
  border-collapse: collapse;
  min-width: 760px;
  background: rgba(10, 9, 8, 0.9);
}

th,
td {
  padding: 10px;
  border-bottom: 1px solid rgba(212, 175, 55, 0.1);
  text-align: left;
  font-size: 13px;
}

code {
  font-family: Consolas, 'Courier New', monospace;
  background: rgba(31, 26, 18, 0.9);
  border-radius: 6px;
  padding: 2px 6px;
}

.status-pill {
  border-radius: 999px;
  padding: 4px 10px;
  font-size: 12px;
}

.status-pill.done {
  background: rgba(132, 194, 157, 0.2);
}

.status-pill.doing {
  background: rgba(121, 168, 198, 0.2);
}

.status-pill.todo {
  background: rgba(184, 171, 144, 0.2);
}

.graph-list,
.log-panel {
  min-height: 0;
  overflow: auto;
  border-radius: 12px;
  background: rgba(10, 9, 8, 0.9);
  padding: 10px;
  display: grid;
  align-content: start;
  gap: 8px;
}

.graph-node {
  border-radius: 10px;
  padding: 8px 10px;
  display: grid;
  gap: 3px;
}

.graph-node strong {
  font-family: Consolas, 'Courier New', monospace;
  font-size: 12px;
}

.graph-node span {
  font-size: 13px;
}

.graph-node.done {
  background: rgba(132, 194, 157, 0.16);
}

.graph-node.doing {
  background: rgba(121, 168, 198, 0.16);
}

.graph-node.todo {
  background: rgba(184, 171, 144, 0.14);
}

.graph-edge {
  border-left: 2px solid rgba(212, 175, 55, 0.32);
  padding: 2px 8px;
  font-family: Consolas, 'Courier New', monospace;
  font-size: 12px;
  color: #e9dcc2;
}

.log-panel p {
  margin: 0;
  font-family: Consolas, 'Courier New', monospace;
  font-size: 12px;
  color: #e9dcc2;
}

@media (max-width: 1280px) {
  .orchestrator-grid {
    grid-template-columns: 1fr;
  }

  .right-stack {
    grid-template-rows: auto auto;
  }
}
</style>
