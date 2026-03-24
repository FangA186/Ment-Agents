<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { collabFlowState, getAgentProfile } from '../content/collabFlowMock'

const route = useRoute()
const router = useRouter()

const projectId = computed(() => String(route.params.projectId || 'proj-001'))
const profile = computed(() => getAgentProfile(route.params.id))

function backToChat() {
  router.push(`/studio/projects/${projectId.value}/chat`)
}
</script>

<template>
  <section class="agent-page">
    <div class="agent-head">
      <h2>Agent 详情 · {{ profile.name }}</h2>
      <div class="head-actions">
        <span>来源：{{ collabFlowState.irFileName }} · #task-02</span>
        <button class="btn ghost" type="button" @click="backToChat">返回群聊页</button>
      </div>
    </div>

    <div class="agent-grid">
      <article class="agent-summary">
        <h3>角色概览</h3>
        <p>职责：{{ profile.summary }}</p>
        <p>风格：{{ profile.style }}</p>
        <p class="status">运行状态：{{ profile.status }}</p>
      </article>

      <article class="asset-panel">
        <h3>核心资产文件</h3>
        <div class="asset-grid">
          <section class="asset-card">
            <div class="file-icon"></div>
            <div>
              <p class="asset-main">{{ profile.roleFile }}</p>
              <p class="asset-title">角色模板</p>
              <p class="asset-desc">{{ profile.roleDesc }}</p>
            </div>
          </section>
          <section class="asset-card">
            <div class="file-icon"></div>
            <div>
              <p class="asset-main">{{ profile.skillFile }}</p>
              <p class="asset-title">技能说明</p>
              <p class="asset-desc">{{ profile.skillDesc }}</p>
            </div>
          </section>
        </div>
      </article>

      <article class="tool-panel">
        <h3>工具接口（图标 + 文字）</h3>
        <div class="tool-list">
          <section v-for="tool in profile.tools" :key="tool.name" class="tool-row">
            <span class="tool-icon">{{ tool.icon }}</span>
            <strong>{{ tool.name }}</strong>
            <span>{{ tool.desc }}</span>
          </section>
        </div>
      </article>
    </div>
  </section>
</template>

<style scoped>
.agent-page {
  display: grid;
  gap: 14px;
}

.agent-head {
  border-radius: 14px;
  padding: 12px 14px;
  background: rgba(20, 18, 14, 0.75);
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.agent-head h2 {
  margin: 0;
  font-size: 32px;
}

.head-actions {
  display: inline-flex;
  gap: 10px;
  align-items: center;
  color: var(--muted);
}

.agent-grid {
  display: grid;
  grid-template-columns: 0.85fr 1.25fr;
  gap: 14px;
}

.agent-summary,
.asset-panel,
.tool-panel {
  border-radius: 16px;
  background: rgba(21, 18, 14, 0.8);
  padding: 16px;
}

.agent-summary {
  grid-row: 1 / span 2;
}

h3 {
  margin: 0 0 12px;
  color: #f2dc9f;
  font-size: 24px;
}

.agent-summary p {
  margin: 0 0 10px;
  font-size: 18px;
}

.agent-summary .status {
  color: #96d0ac;
}

.asset-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.asset-card {
  border-radius: 14px;
  background: rgba(27, 23, 18, 0.92);
  padding: 14px;
  display: grid;
  grid-template-columns: 76px 1fr;
  gap: 10px;
}

.file-icon {
  width: 62px;
  height: 74px;
  border: 1px solid rgba(212, 175, 55, 0.3);
  border-radius: 4px;
  background: linear-gradient(180deg, rgba(70, 60, 44, 0.5), rgba(40, 35, 26, 0.8));
  position: relative;
}

.file-icon::after {
  content: '';
  position: absolute;
  right: -1px;
  top: -1px;
  width: 18px;
  height: 18px;
  clip-path: polygon(0 0, 100% 100%, 100% 0);
  background: rgba(128, 109, 72, 0.8);
}

.asset-main {
  margin: 0;
  font-size: 38px;
  line-height: 1.06;
  font-weight: 700;
}

.asset-title {
  margin: 6px 0 0;
  color: #e4d4ac;
  font-size: 24px;
}

.asset-desc {
  margin: 4px 0 0;
  color: var(--muted);
  font-size: 16px;
}

.tool-list {
  display: grid;
  gap: 8px;
}

.tool-row {
  border-radius: 10px;
  background: rgba(23, 20, 16, 0.94);
  padding: 10px 12px;
  display: grid;
  grid-template-columns: 38px 180px 1fr;
  align-items: center;
  gap: 8px;
}

.tool-icon {
  width: 30px;
  height: 30px;
  border-radius: 999px;
  display: grid;
  place-items: center;
  background: rgba(85, 71, 48, 0.95);
  color: #f0d98f;
  font-size: 16px;
  font-weight: 700;
}

.tool-row strong {
  font-size: 17px;
}

.tool-row span {
  color: var(--muted);
  font-size: 15px;
}

@media (max-width: 1380px) {
  .agent-grid {
    grid-template-columns: 1fr;
  }

  .agent-summary {
    grid-row: auto;
  }

  .asset-grid {
    grid-template-columns: 1fr;
  }
}
</style>
