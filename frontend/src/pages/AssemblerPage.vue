<script setup>
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { collabFlowState, loadProjectContext } from '../content/collabFlowMock'
import { assemblerMock } from '../content/workbenchMock'

const route = useRoute()

const tasks = ref(assemblerMock.tasks)
const agents = ref(assemblerMock.agents)

const nodes = computed(() => {
  if (collabFlowState.agentGraph?.nodes?.length) {
    return collabFlowState.agentGraph.nodes.map((item) => ({
      id: item.id,
      label: item.name,
      state: 'running',
    }))
  }
  return assemblerMock.graph.nodes
})

const edges = computed(() => {
  if (collabFlowState.agentGraph?.edges?.length) {
    return collabFlowState.agentGraph.edges
  }
  return assemblerMock.graph.edges
})

const selectedTaskId = ref(tasks.value[0].id)

const selectedAgent = computed(() => agents.value.find((a) => a.taskId === selectedTaskId.value) || agents.value[0])

function nodeClass(state) {
  if (state === 'ready') return 'ready'
  if (state === 'running') return 'running'
  return 'pending'
}

watch(
  () => route.params.projectId,
  (projectId) => {
    loadProjectContext(String(projectId || 'proj-001')).catch((error) => {
      console.error(error)
    })
  },
  { immediate: true },
)
</script>

<template>
  <section class="panel assembler">
    <div class="panel-head">
      <h2>组装器页面</h2>
      <span class="stage-pill running">图形化组装中</span>
    </div>

    <p class="panel-intro">根据 IR 任务块组装 Agent 实例图。每个 Agent 明确展示：角色模板、技能集、工具接口、记忆。</p>

    <div class="assembler-grid">
      <article class="comp-card">
        <h3>任务块</h3>
        <div class="task-stack">
          <button
            v-for="task in tasks"
            :key="task.id"
            type="button"
            class="task-chip"
            :class="{ active: task.id === selectedTaskId }"
            @click="selectedTaskId = task.id"
          >
            <span class="id">{{ task.id }}</span>
            <span class="name">{{ task.name }}</span>
            <span class="priority">{{ task.priority }}</span>
          </button>
        </div>
      </article>

      <article class="comp-card">
        <h3>Agent 四件套（{{ selectedAgent.agentId }}）</h3>
        <div class="assembler-four-grid">
          <section class="mini-box">
            <p class="mini-title">角色模板</p>
            <div class="row">
              <span class="dot role"></span>
              <span>{{ selectedAgent.roleTemplate }}</span>
            </div>
          </section>

          <section class="mini-box">
            <p class="mini-title">技能集</p>
            <div class="tag-group">
              <span v-for="skill in selectedAgent.skills" :key="skill" class="tag skill">{{ skill }}</span>
            </div>
          </section>

          <section class="mini-box">
            <p class="mini-title">工具接口</p>
            <div class="tag-group">
              <span v-for="tool in selectedAgent.tools" :key="tool" class="tag tool">{{ tool }}</span>
            </div>
          </section>

          <section class="mini-box">
            <p class="mini-title">记忆</p>
            <p class="memory-label">工作记忆</p>
            <ul class="memory-list">
              <li v-for="item in selectedAgent.memory.work" :key="item">{{ item }}</li>
            </ul>
            <p class="memory-label">长期记忆</p>
            <ul class="memory-list">
              <li v-for="item in selectedAgent.memory.long" :key="item">{{ item }}</li>
            </ul>
          </section>
        </div>
      </article>

      <article class="comp-card">
        <h3>AgentGraph</h3>
        <div class="graph-box">
          <div class="node-row">
            <div v-for="node in nodes" :key="node.id" class="node" :class="nodeClass(node.state)">
              <span>{{ node.id }}</span>
              <small>{{ node.label }}（{{ node.id }}）</small>
            </div>
          </div>
          <div class="edge-row">
            <span v-for="edge in edges" :key="`${edge.from}-${edge.to}`" class="edge">{{ edge.from }} → {{ edge.to }}</span>
          </div>
        </div>
      </article>
    </div>
  </section>
</template>
