<script setup>
import { computed, ref } from 'vue'
import { blackboardStates, channelTabs, messageTimeline, taskLinks } from '../content/communicationMock'

const activeChannel = ref('all')

const filteredMessages = computed(() => {
  if (activeChannel.value === 'all') return messageTimeline
  return messageTimeline.filter((msg) => msg.channel === activeChannel.value)
})

function boardTagClass(status) {
  if (status === '已锁定') return 'locked'
  if (status === '已同步') return 'synced'
  if (status === '更新中') return 'updating'
  return 'pending'
}

function taskClass(status) {
  if (status === '已完成') return 'done'
  if (status === '进行中') return 'doing'
  return 'todo'
}
</script>

<template>
  <section class="panel com-page">
    <div class="panel-head">
      <h2>多 Agent 通信页面</h2>
      <div class="panel-actions">
        <button class="btn ghost" type="button">导出通信日志</button>
      </div>
    </div>
    <p class="panel-intro">黑板共享状态、消息时间线与任务关联统一展示，帮助快速定位协同冲突与上下游依赖。</p>

    <div class="com-grid">
      <article class="com-card">
        <h3>黑板共享状态</h3>
        <div class="blackboard-list">
          <section v-for="item in blackboardStates" :key="item.id" class="board-item">
            <p class="board-key">{{ item.key }}</p>
            <p class="board-value">{{ item.value }}</p>
            <div class="board-meta">
              <span>{{ item.owner }}</span>
              <span>{{ item.taskId }}</span>
              <span class="board-tag" :class="boardTagClass(item.status)">{{ item.status }}</span>
            </div>
          </section>
        </div>
      </article>

      <article class="com-card">
        <h3>消息时间线</h3>
        <div class="channel-tabs">
          <button
            v-for="tab in channelTabs"
            :key="tab.id"
            type="button"
            :class="{ active: activeChannel === tab.id }"
            @click="activeChannel = tab.id"
          >
            {{ tab.label }}
          </button>
        </div>
        <div class="timeline-list">
          <section v-for="msg in filteredMessages" :key="msg.id" class="timeline-item">
            <div class="timeline-top">
              <span class="sender">{{ msg.sender }}</span>
              <span class="time">{{ msg.time }}</span>
            </div>
            <p class="timeline-content">{{ msg.content }}</p>
            <div class="timeline-meta">
              <span>目标：{{ msg.target }}</span>
              <span>任务：{{ msg.taskId }}</span>
            </div>
          </section>
          <p v-if="filteredMessages.length === 0" class="empty">当前通道暂无消息</p>
        </div>
      </article>

      <article class="com-card">
        <h3>任务关联视图</h3>
        <div class="task-link-list">
          <section v-for="task in taskLinks" :key="task.id" class="task-link-item">
            <div class="task-top">
              <p>{{ task.taskId }} · {{ task.title }}</p>
              <span class="task-status" :class="taskClass(task.status)">{{ task.status }}</span>
            </div>
            <p class="task-lines">关联消息：{{ task.messages.join(', ') }}</p>
            <p class="task-lines">关联 Agent：{{ task.relatedAgents.join(' / ') }}</p>
          </section>
        </div>
      </article>
    </div>
  </section>
</template>

<style scoped>
.com-page {
  display: grid;
  gap: 12px;
  min-height: 0;
}

.com-grid {
  min-height: 0;
  display: grid;
  grid-template-columns: 1fr 1fr 0.9fr;
  gap: 12px;
}

.com-card {
  border-radius: 14px;
  background: rgba(16, 14, 11, 0.88);
  padding: 12px;
  min-height: 0;
  display: grid;
  grid-template-rows: auto 1fr;
  gap: 10px;
}

h3 {
  margin: 0;
  color: #f2dc9f;
  font-size: 22px;
}

.blackboard-list,
.timeline-list,
.task-link-list {
  min-height: 0;
  overflow-y: auto;
  display: grid;
  align-content: start;
  gap: 8px;
  padding-right: 4px;
}

.board-item,
.timeline-item,
.task-link-item {
  border-radius: 12px;
  background: rgba(24, 20, 16, 0.9);
  padding: 10px;
}

.board-key {
  margin: 0;
  font-family: Consolas, 'Courier New', monospace;
  color: #e8cf8f;
  font-size: 12px;
}

.board-value {
  margin: 6px 0;
  font-size: 14px;
}

.board-meta {
  display: inline-flex;
  gap: 8px;
  flex-wrap: wrap;
  color: var(--muted);
  font-size: 12px;
}

.board-tag {
  border-radius: 999px;
  padding: 2px 8px;
}

.board-tag.locked {
  background: rgba(212, 175, 55, 0.2);
}

.board-tag.synced {
  background: rgba(132, 194, 157, 0.2);
}

.board-tag.updating {
  background: rgba(121, 168, 198, 0.2);
}

.board-tag.pending {
  background: rgba(184, 171, 144, 0.2);
}

.channel-tabs {
  display: inline-flex;
  gap: 8px;
  flex-wrap: wrap;
}

.channel-tabs button {
  border: none;
  border-radius: 999px;
  background: rgba(26, 22, 17, 0.9);
  color: var(--muted);
  padding: 6px 10px;
  font-size: 12px;
  cursor: pointer;
}

.channel-tabs button.active {
  color: #1a140c;
  background: linear-gradient(145deg, #8a6b24, var(--accent-2));
}

.timeline-top,
.timeline-meta,
.task-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}

.timeline-content {
  margin: 6px 0;
  font-size: 14px;
}

.timeline-meta,
.time,
.task-lines {
  color: var(--muted);
  font-size: 12px;
}

.sender {
  font-weight: 700;
}

.task-top p {
  margin: 0;
  font-size: 14px;
  font-weight: 700;
}

.task-status {
  border-radius: 999px;
  padding: 2px 8px;
  font-size: 12px;
}

.task-status.done {
  background: rgba(132, 194, 157, 0.2);
}

.task-status.doing {
  background: rgba(121, 168, 198, 0.2);
}

.task-status.todo {
  background: rgba(184, 171, 144, 0.2);
}

.task-lines {
  margin: 5px 0 0;
}

.empty {
  color: var(--muted);
  font-size: 13px;
  margin: 0;
}

@media (max-width: 1360px) {
  .com-grid {
    grid-template-columns: 1fr;
  }
}
</style>
