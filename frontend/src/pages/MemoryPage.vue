<script setup>
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { ensureProjects, getProjectById, projects } from '../content/projectStore'
import { getProjectMemory } from '../content/memoryMock'

const route = useRoute()

const activeProjectId = computed(() => String(route.params.projectId || projects.value[0]?.id || 'proj-001'))
const activeProject = computed(() => getProjectById(activeProjectId.value))
const memoryData = computed(() => getProjectMemory(activeProjectId.value))

const activeLayerId = ref('')
const activeSlotId = ref('')

const activeLayer = computed(() => memoryData.value.layers.find((item) => item.id === activeLayerId.value) || memoryData.value.layers[0])
const slots = computed(() => activeLayer.value?.slots || [])
const activeSlot = computed(() => slots.value.find((item) => item.id === activeSlotId.value) || slots.value[0])

function switchLayer(layerId) {
  activeLayerId.value = layerId
  const firstSlot = memoryData.value.layers.find((item) => item.id === layerId)?.slots[0]
  activeSlotId.value = firstSlot?.id || ''
}

function resetSelection() {
  activeLayerId.value = memoryData.value.layers[0]?.id || ''
  activeSlotId.value = memoryData.value.layers[0]?.slots[0]?.id || ''
}

watch(activeProjectId, resetSelection, { immediate: true })

ensureProjects().catch((error) => {
  console.error(error)
})
</script>

<template>
  <section class="panel memory-page">
    <div class="panel-head">
      <h2>上下文与记忆</h2>
      <div class="panel-actions">
        <span class="stage-pill running">{{ activeProject.name }} · {{ activeProject.id }}</span>
      </div>
    </div>

    <p class="panel-intro">展示当前项目的 Global / Task / Agent 分层上下文与 RAG 检索命中。</p>

    <div class="memory-grid">
      <article class="memory-card">
        <h3>分层上下文</h3>
        <div class="layer-list">
          <button
            v-for="layer in memoryData.layers"
            :key="layer.id"
            type="button"
            class="layer-item"
            :class="{ active: layer.id === activeLayerId }"
            @click="switchLayer(layer.id)"
          >
            <strong>{{ layer.name }}</strong>
            <span>{{ layer.desc }}</span>
          </button>
        </div>
      </article>

      <article class="memory-card">
        <h3>上下文槽位</h3>
        <div class="slot-list">
          <button
            v-for="slot in slots"
            :key="slot.id"
            type="button"
            class="slot-item"
            :class="{ active: slot.id === activeSlotId }"
            @click="activeSlotId = slot.id"
          >
            <p class="slot-key">{{ slot.key }}</p>
            <p class="slot-value">{{ slot.value }}</p>
            <p class="slot-meta">{{ slot.source }} · {{ slot.updatedAt }}</p>
          </button>
        </div>
      </article>

      <article class="memory-card">
        <h3>槽位详情</h3>
        <div class="detail-box">
          <p><strong>层级：</strong>{{ activeLayer.name }}</p>
          <p><strong>键：</strong>{{ activeSlot.key }}</p>
          <p><strong>值：</strong>{{ activeSlot.value }}</p>
          <p><strong>来源：</strong>{{ activeSlot.source }}</p>
          <p><strong>更新时间：</strong>{{ activeSlot.updatedAt }}</p>
        </div>
      </article>
    </div>

    <article class="memory-card rag-card">
      <h3>RAG 检索命中</h3>
      <div class="rag-list">
        <section v-for="item in memoryData.rag" :key="item.id" class="rag-item">
          <p class="rag-query">查询：{{ item.query }}</p>
          <p class="rag-hit">命中槽位：{{ item.hit }}</p>
          <p class="rag-snippet">{{ item.snippet }}</p>
          <span class="rag-score">相似度 {{ item.score }}</span>
        </section>
      </div>
    </article>
  </section>
</template>

<style scoped>
.memory-page {
  width: 100%;
  display: grid;
  gap: 12px;
  min-height: 0;
}

.memory-grid {
  min-height: 0;
  display: grid;
  grid-template-columns: 0.9fr 1fr 0.9fr;
  gap: 12px;
}

.memory-card {
  border-radius: 14px;
  background: rgba(10, 9, 8, 0.9);
  padding: 12px;
  min-height: 0;
  display: grid;
  grid-template-rows: auto 1fr;
  gap: 8px;
}

.memory-card h3 {
  margin: 0;
  color: #f2dc9f;
  font-size: 20px;
}

.layer-list,
.slot-list,
.rag-list {
  min-height: 0;
  overflow-y: auto;
  display: grid;
  align-content: start;
  gap: 8px;
}

.layer-item,
.slot-item {
  border: none;
  border-radius: 12px;
  background: rgba(24, 20, 16, 0.9);
  color: var(--text);
  text-align: left;
  padding: 10px;
  display: grid;
  gap: 4px;
  cursor: pointer;
}

.layer-item.active,
.slot-item.active {
  box-shadow: inset 0 0 0 1px rgba(212, 175, 55, 0.42);
  background: rgba(32, 27, 20, 0.9);
}

.layer-item strong {
  font-size: 14px;
}

.layer-item span,
.slot-meta {
  color: var(--muted);
  font-size: 12px;
}

.slot-key {
  margin: 0;
  font-family: Consolas, 'Courier New', monospace;
  color: #e8cf8f;
  font-size: 12px;
}

.slot-value {
  margin: 0;
  font-size: 14px;
}

.detail-box {
  border-radius: 12px;
  background: rgba(24, 20, 16, 0.9);
  padding: 10px;
  display: grid;
  gap: 8px;
}

.detail-box p {
  margin: 0;
  font-size: 14px;
}

.rag-card {
  min-height: 220px;
}

.rag-item {
  border-radius: 12px;
  background: rgba(24, 20, 16, 0.9);
  padding: 10px;
  display: grid;
  gap: 4px;
}

.rag-query,
.rag-hit,
.rag-snippet {
  margin: 0;
  font-size: 13px;
}

.rag-hit {
  color: #e8cf8f;
}

.rag-snippet {
  color: var(--muted);
}

.rag-score {
  justify-self: end;
  border-radius: 999px;
  background: rgba(121, 168, 198, 0.2);
  padding: 2px 8px;
  font-size: 12px;
}

@media (max-width: 1460px) {
  .memory-grid {
    grid-template-columns: 1fr;
  }
}
</style>
