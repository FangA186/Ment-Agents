<script setup>
import { computed, reactive, ref } from 'vue'
import { debugFilters, debugRuns, recoveryPlaybooks } from '../content/debugSelfHealMock'

const runs = reactive(debugRuns.map((run) => ({ ...run, timeline: run.timeline.map((item) => ({ ...item })) })))
const activeFilter = ref('all')
const selectedRunId = ref(runs[0]?.id || '')
const actionMessage = ref('已加载最近一次异常处置上下文。')

const filteredRuns = computed(() => {
  if (activeFilter.value === 'all') return runs
  return runs.filter((run) => run.status === activeFilter.value)
})

const selectedRun = computed(() => filteredRuns.value.find((run) => run.id === selectedRunId.value) || runs.find((run) => run.id === selectedRunId.value) || runs[0])
const activePlaybook = computed(() => recoveryPlaybooks.find((item) => item.id === selectedRun.value?.strategyId) || recoveryPlaybooks[0])

const summaryCards = computed(() => {
  const blocked = runs.filter((run) => run.status === 'blocked').length
  const recovering = runs.filter((run) => run.status === 'recovering').length
  const healed = runs.filter((run) => run.status === 'healed').length
  const rolledBack = runs.filter((run) => run.status === 'rolled_back').length

  return [
    { label: '待人工接管', value: `${blocked} 个`, tone: 'blocked' },
    { label: '自动恢复中', value: `${recovering} 个`, tone: 'recovering' },
    { label: '已恢复', value: `${healed} 个`, tone: 'healed' },
    { label: '已回滚', value: `${rolledBack} 个`, tone: 'rollback' },
  ]
})

function statusLabel(status) {
  if (status === 'recovering') return '恢复中'
  if (status === 'blocked') return '待人工接管'
  if (status === 'healed') return '已恢复'
  if (status === 'rolled_back') return '已回滚'
  return '待处理'
}

function statusClass(status) {
  if (status === 'recovering') return 'recovering'
  if (status === 'blocked') return 'blocked'
  if (status === 'healed') return 'healed'
  if (status === 'rolled_back') return 'rollback'
  return 'pending'
}

function severityClass(severity) {
  if (severity === 'P1') return 'sev-p1'
  if (severity === 'P2') return 'sev-p2'
  return 'sev-p3'
}

function timelineClass(type) {
  if (type === 'error') return 'error'
  if (type === 'recovery') return 'recovery'
  if (type === 'guard') return 'guard'
  if (type === 'done') return 'done'
  return 'analysis'
}

function selectRun(runId) {
  selectedRunId.value = runId
  actionMessage.value = '已切换到新的异常 run，恢复策略面板已同步更新。'
}

function prependEvent(type, title, detail) {
  selectedRun.value.timeline.unshift({
    id: `${selectedRun.value.id}-${Date.now()}`,
    time: new Date().toLocaleTimeString('zh-CN', { hour12: false }),
    type,
    title,
    detail,
  })
}

function triggerRetry() {
  if (!selectedRun.value) return
  selectedRun.value.retryCount += 1
  selectedRun.value.status = 'recovering'
  selectedRun.value.updatedAt = '刚刚'
  selectedRun.value.result = {
    headline: `已触发第 ${selectedRun.value.retryCount} 次重试`,
    detail: '保留当前上下文与审计链路，优先恢复主执行路径。',
    evidence: ['retry mode=guarded', `retry count=${selectedRun.value.retryCount}`, 'context locked=true'],
  }
  prependEvent('recovery', '触发受控重试', '重新执行当前节点，并保留原始上下文快照。')
  actionMessage.value = `最近动作：已对 ${selectedRun.value.title} 触发第 ${selectedRun.value.retryCount} 次重试。`
}

function executeRollback() {
  if (!selectedRun.value || !selectedRun.value.canRollback) return
  selectedRun.value.status = 'rolled_back'
  selectedRun.value.canRollback = false
  selectedRun.value.updatedAt = '刚刚'
  selectedRun.value.result = {
    headline: '已执行回滚并恢复到最近稳定检查点',
    detail: '回滚完成后保留问题快照，等待后续复盘。',
    evidence: ['checkpoint=current-last-green', 'downstream reset=true', 'audit trail persisted'],
  }
  prependEvent('guard', '执行回滚', '撤销异常节点并恢复至最近绿色检查点。')
  actionMessage.value = `最近动作：${selectedRun.value.title} 已执行回滚。`
}

function markRecovered() {
  if (!selectedRun.value) return
  selectedRun.value.status = 'healed'
  selectedRun.value.updatedAt = '刚刚'
  selectedRun.value.result = {
    headline: '异常已闭环恢复',
    detail: '当前 run 已恢复绿色状态，并生成复盘摘要。',
    evidence: ['status=green', 'post-check passed=true', 'handoff report ready'],
  }
  prependEvent('done', '恢复确认完成', '已通过恢复后检查，run 重新回到绿色状态。')
  actionMessage.value = `最近动作：${selectedRun.value.title} 已标记为恢复完成。`
}
</script>

<template>
  <section class="panel debug-page">
    <div class="panel-head">
      <h2>调试与自愈页面</h2>
      <div class="panel-actions">
        <span class="stage-pill error">{{ statusLabel(selectedRun.status) }}</span>
        <button class="btn ghost" type="button">导出复盘摘要</button>
      </div>
    </div>

    <p class="panel-intro">聚焦异常 run、恢复策略、回滚轨迹与恢复结果，帮助你快速定位失败任务并执行自愈动作。</p>

    <div class="summary-grid">
      <article v-for="card in summaryCards" :key="card.label" class="summary-card" :class="card.tone">
        <p>{{ card.label }}</p>
        <strong>{{ card.value }}</strong>
      </article>
    </div>

    <div class="filter-row">
      <button
        v-for="filter in debugFilters"
        :key="filter.id"
        type="button"
        class="filter-chip"
        :class="{ active: filter.id === activeFilter }"
        @click="activeFilter = filter.id"
      >
        {{ filter.label }}
      </button>
    </div>

    <div class="debug-grid">
      <article class="debug-card run-card">
        <div class="debug-card-head">
          <h3>异常 run 列表</h3>
          <span>{{ filteredRuns.length }} 条</span>
        </div>
        <div class="run-list">
          <button
            v-for="run in filteredRuns"
            :key="run.id"
            type="button"
            class="run-item"
            :class="{ active: run.id === selectedRun.id }"
            @click="selectRun(run.id)"
          >
            <div class="run-top">
              <strong>{{ run.title }}</strong>
              <span class="severity-pill" :class="severityClass(run.severity)">{{ run.severity }}</span>
            </div>
            <p class="run-summary">{{ run.summary }}</p>
            <div class="run-meta">
              <span>{{ run.phase }}</span>
              <span>{{ run.owner }}</span>
              <span class="status-pill" :class="statusClass(run.status)">{{ statusLabel(run.status) }}</span>
            </div>
          </button>
        </div>
      </article>

      <article class="debug-card detail-card">
        <div class="debug-card-head">
          <h3>{{ selectedRun.title }}</h3>
          <span>{{ selectedRun.id }} · {{ selectedRun.updatedAt }}</span>
        </div>

        <div class="detail-strip">
          <article class="detail-box">
            <p class="detail-label">触发条件</p>
            <strong>{{ selectedRun.trigger }}</strong>
          </article>
          <article class="detail-box">
            <p class="detail-label">影响范围</p>
            <strong>{{ selectedRun.impact }}</strong>
          </article>
          <article class="detail-box">
            <p class="detail-label">当前策略</p>
            <strong>{{ activePlaybook.name }}</strong>
          </article>
        </div>

        <div class="timeline-panel">
          <div class="timeline-line"></div>
          <article v-for="event in selectedRun.timeline" :key="event.id" class="timeline-event" :class="timelineClass(event.type)">
            <div class="event-dot"></div>
            <div class="event-body">
              <div class="event-top">
                <strong>{{ event.title }}</strong>
                <span>{{ event.time }}</span>
              </div>
              <p>{{ event.detail }}</p>
            </div>
          </article>
        </div>
      </article>

      <div class="right-stack">
        <article class="debug-card strategy-card">
          <div class="debug-card-head">
            <h3>恢复策略矩阵</h3>
            <span>{{ activePlaybook.id }}</span>
          </div>
          <div class="playbook-list">
            <section
              v-for="playbook in recoveryPlaybooks"
              :key="playbook.id"
              class="playbook-item"
              :class="{ active: playbook.id === activePlaybook.id }"
            >
              <div class="playbook-head">
                <strong>{{ playbook.name }}</strong>
                <span>{{ playbook.owner }}</span>
              </div>
              <p class="playbook-trigger">{{ playbook.trigger }}</p>
              <p class="playbook-label">执行步骤</p>
              <ul>
                <li v-for="step in playbook.steps" :key="step">{{ step }}</li>
              </ul>
              <p class="playbook-label">护栏</p>
              <ul>
                <li v-for="guard in playbook.guardrails" :key="guard">{{ guard }}</li>
              </ul>
            </section>
          </div>
        </article>

        <article class="debug-card action-card">
          <div class="debug-card-head">
            <h3>恢复控制台</h3>
            <span>retry={{ selectedRun.retryCount }}</span>
          </div>
          <div class="action-grid">
            <button class="btn" type="button" @click="triggerRetry">触发重试</button>
            <button class="btn ghost" type="button" :disabled="!selectedRun.canRollback" @click="executeRollback">执行回滚</button>
            <button class="btn ghost" type="button" @click="markRecovered">标记恢复完成</button>
          </div>
          <p class="action-message">{{ actionMessage }}</p>
        </article>

        <article class="debug-card result-card">
          <div class="debug-card-head">
            <h3>执行结果</h3>
            <span>{{ statusLabel(selectedRun.status) }}</span>
          </div>
          <div class="result-box">
            <p class="result-headline">{{ selectedRun.result.headline }}</p>
            <p class="result-detail">{{ selectedRun.result.detail }}</p>
            <div class="metrics-grid">
              <section v-for="metric in selectedRun.metrics" :key="metric.label" class="metric-box">
                <p>{{ metric.label }}</p>
                <strong>{{ metric.value }}</strong>
              </section>
            </div>
            <div class="evidence-list">
              <p class="playbook-label">关键证据</p>
              <ul>
                <li v-for="evidence in selectedRun.result.evidence" :key="evidence">{{ evidence }}</li>
              </ul>
            </div>
          </div>
        </article>
      </div>
    </div>
  </section>
</template>

<style scoped>
.debug-page {
  display: grid;
  gap: 12px;
  min-height: 0;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 10px;
}

.summary-card {
  border-radius: 14px;
  padding: 12px 14px;
  background: rgba(16, 14, 11, 0.9);
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.04);
}

.summary-card p,
.metric-box p {
  margin: 0;
  color: var(--muted);
  font-size: 12px;
}

.summary-card strong,
.metric-box strong {
  display: block;
  margin-top: 6px;
  font-size: 26px;
}

.summary-card.blocked {
  background: linear-gradient(145deg, rgba(93, 29, 29, 0.3), rgba(16, 14, 11, 0.92));
}

.summary-card.recovering {
  background: linear-gradient(145deg, rgba(44, 79, 99, 0.28), rgba(16, 14, 11, 0.92));
}

.summary-card.healed {
  background: linear-gradient(145deg, rgba(43, 92, 64, 0.28), rgba(16, 14, 11, 0.92));
}

.summary-card.rollback {
  background: linear-gradient(145deg, rgba(92, 71, 35, 0.28), rgba(16, 14, 11, 0.92));
}

.filter-row {
  display: inline-flex;
  gap: 8px;
  flex-wrap: wrap;
}

.filter-chip {
  border: none;
  border-radius: 999px;
  padding: 7px 12px;
  color: var(--muted);
  background: rgba(20, 18, 14, 0.84);
  cursor: pointer;
}

.filter-chip.active {
  color: #1a140c;
  background: linear-gradient(145deg, #8a6b24, var(--accent-2));
  font-weight: 700;
}

.debug-grid {
  min-height: 0;
  display: grid;
  grid-template-columns: 0.95fr 1.25fr 1fr;
  gap: 12px;
}

.right-stack {
  min-height: 0;
  display: grid;
  grid-template-rows: 1.15fr auto 1fr;
  gap: 12px;
}

.debug-card {
  border-radius: 14px;
  background: rgba(16, 14, 11, 0.9);
  padding: 12px;
  min-height: 0;
  display: grid;
  grid-template-rows: auto 1fr;
  gap: 10px;
}

.debug-card-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}

.debug-card-head h3 {
  margin: 0;
  color: #f2dc9f;
  font-size: 21px;
}

.debug-card-head span {
  color: var(--muted);
  font-size: 12px;
}

.run-list,
.playbook-list {
  min-height: 0;
  overflow-y: auto;
  display: grid;
  align-content: start;
  gap: 8px;
}

.run-item,
.playbook-item,
.detail-box,
.metric-box {
  border: 1px solid transparent;
  border-radius: 12px;
  background: rgba(24, 20, 16, 0.92);
}

.run-item {
  color: inherit;
  padding: 12px;
  cursor: pointer;
  text-align: left;
}

.run-item.active {
  border-color: rgba(212, 175, 55, 0.3);
  background: linear-gradient(145deg, rgba(212, 175, 55, 0.14), rgba(24, 20, 16, 0.94));
}

.run-top,
.run-meta,
.event-top,
.playbook-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}

.run-summary,
.result-detail,
.playbook-trigger,
.event-body p {
  margin: 6px 0 0;
  color: #ddd2be;
  font-size: 13px;
}

.run-meta {
  margin-top: 10px;
  color: var(--muted);
  font-size: 12px;
  flex-wrap: wrap;
}

.severity-pill,
.status-pill {
  border-radius: 999px;
  padding: 3px 8px;
  font-size: 12px;
}

.sev-p1 {
  background: rgba(156, 48, 48, 0.28);
  color: #ffd2d2;
}

.sev-p2 {
  background: rgba(140, 108, 31, 0.26);
  color: #f6e1a8;
}

.sev-p3 {
  background: rgba(62, 90, 113, 0.24);
  color: #d4ebff;
}

.status-pill.recovering {
  background: rgba(121, 168, 198, 0.2);
}

.status-pill.blocked {
  background: rgba(156, 48, 48, 0.28);
}

.status-pill.healed {
  background: rgba(132, 194, 157, 0.2);
}

.status-pill.rollback {
  background: rgba(212, 175, 55, 0.18);
}

.detail-card {
  grid-template-rows: auto auto 1fr;
}

.detail-strip {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
}

.detail-box,
.metric-box {
  padding: 10px;
}

.detail-label,
.playbook-label {
  margin: 0 0 6px;
  color: var(--muted);
  font-size: 12px;
}

.detail-box strong,
.playbook-item strong {
  font-size: 14px;
  color: #f3e6c3;
}

.timeline-panel {
  position: relative;
  min-height: 0;
  overflow-y: auto;
  display: grid;
  align-content: start;
  gap: 10px;
  padding-left: 22px;
}

.timeline-line {
  position: absolute;
  left: 6px;
  top: 4px;
  bottom: 4px;
  width: 2px;
  background: linear-gradient(180deg, rgba(212, 175, 55, 0.2), rgba(121, 168, 198, 0.3));
}

.timeline-event {
  position: relative;
}

.event-dot {
  position: absolute;
  left: -20px;
  top: 4px;
  width: 10px;
  height: 10px;
  border-radius: 999px;
}

.timeline-event.error .event-dot {
  background: #d87a7a;
}

.timeline-event.analysis .event-dot {
  background: #d4af37;
}

.timeline-event.recovery .event-dot {
  background: #79a8c6;
}

.timeline-event.guard .event-dot {
  background: #f0d680;
}

.timeline-event.done .event-dot {
  background: #84c29d;
}

.event-body {
  border-radius: 12px;
  background: rgba(24, 20, 16, 0.92);
  padding: 10px 12px;
}

.event-top strong {
  font-size: 14px;
}

.event-top span {
  color: var(--muted);
  font-size: 12px;
}

.playbook-item {
  padding: 12px;
}

.playbook-item.active {
  border-color: rgba(212, 175, 55, 0.3);
}

.playbook-item ul,
.evidence-list ul {
  margin: 0;
  padding-left: 18px;
  display: grid;
  gap: 4px;
  color: #ddd2be;
  font-size: 12px;
}

.action-card {
  grid-template-rows: auto auto auto;
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
}

.action-grid .btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.action-message {
  margin: 0;
  border-radius: 12px;
  background: rgba(24, 20, 16, 0.92);
  padding: 10px 12px;
  color: #ddd2be;
  font-size: 13px;
}

.result-box {
  border-radius: 12px;
  background: rgba(24, 20, 16, 0.92);
  padding: 12px;
}

.result-headline {
  margin: 0;
  font-size: 16px;
  color: #f3e6c3;
  font-weight: 700;
}

.metrics-grid {
  margin-top: 12px;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
}

@media (max-width: 1500px) {
  .debug-grid {
    grid-template-columns: 1fr;
  }

  .right-stack {
    grid-template-rows: auto auto auto;
  }
}

@media (max-width: 960px) {
  .summary-grid,
  .detail-strip,
  .metrics-grid,
  .action-grid {
    grid-template-columns: 1fr;
  }
}
</style>
