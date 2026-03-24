<script setup>
import { computed, onMounted } from 'vue'
import { RouterLink, RouterView, useRoute } from 'vue-router'
import { ensureProjects, projects } from '../content/projectStore'

const route = useRoute()

const activeProjectId = computed(() => String(route.params.projectId || ''))
const activeProject = computed(() => projects.value.find((project) => project.id === activeProjectId.value))
const totalProjects = computed(() => projects.value.length)

const moduleMenus = computed(() => {
  if (!activeProjectId.value) return []
  return [
    { key: 'chat', label: '群聊编排', to: `/studio/projects/${activeProjectId.value}/chat` },
    { key: 'orchestrator', label: '编排执行', to: `/studio/projects/${activeProjectId.value}/orchestrator` },
    { key: 'debug', label: '调试与自愈', to: `/studio/projects/${activeProjectId.value}/debug` },
    { key: 'memory', label: '上下文记忆', to: `/studio/projects/${activeProjectId.value}/memory` },
    { key: 'communication', label: '多 Agent 通信', to: `/studio/projects/${activeProjectId.value}/communication` },
    { key: 'assembler', label: '组装器', to: `/studio/projects/${activeProjectId.value}/assembler` },
  ]
})

onMounted(() => {
  ensureProjects().catch((error) => {
    console.error(error)
  })
})
</script>

<template>
  <div class="wb-shell">
    <aside class="wb-sidebar">
      <div class="wb-brand">
        <div class="logo-dot">M</div>
        <div>
          <p class="brand-title">Ment-Agents</p>
          <p class="brand-sub">项目库驱动工作台</p>
        </div>
      </div>

      <nav v-if="!activeProject" class="single-nav">
        <RouterLink class="single-nav-item" :class="{ active: route.path.startsWith('/studio/projects') }" to="/studio/projects">项目库</RouterLink>
      </nav>

      <nav v-if="activeProject" class="module-nav">
        <RouterLink
          v-for="menu in moduleMenus"
          :key="menu.key"
          class="module-nav-item"
          :class="{ active: route.path === menu.to }"
          :to="menu.to"
        >
          {{ menu.label }}
        </RouterLink>
      </nav>

      <div class="nav-note">
        <p>项目库：共 {{ totalProjects }} 个项目</p>
        <p>当前项目：{{ activeProject ? activeProject.name : '未选择' }}</p>
      </div>
    </aside>

    <div class="wb-main">
      <header class="wb-topbar">
        <h1>{{ activeProject ? `${activeProject.name} · 多智能体控制台` : '项目库' }}</h1>
        <div class="top-status">
          <span class="status-dot"></span>
          <span>{{ activeProject ? `${activeProject.id} · ${activeProject.status}` : '请选择项目进入智能体编排' }}</span>
        </div>
      </header>

      <main class="wb-content">
        <RouterView :key="route.fullPath" />
      </main>
    </div>
  </div>
</template>

<style scoped>
.single-nav {
  display: grid;
  align-content: start;
}

.single-nav-item {
  border: 1px solid transparent;
  border-radius: 12px;
  background: rgba(24, 20, 16, 0.6);
  color: var(--text);
  text-decoration: none;
  padding: 12px;
  font-size: 16px;
  font-weight: 700;
}

.single-nav-item.active {
  border-color: rgba(212, 175, 55, 0.24);
  background: linear-gradient(135deg, rgba(212, 175, 55, 0.2), rgba(212, 175, 55, 0.08));
}

.module-nav {
  display: grid;
  align-content: start;
  gap: 8px;
}

.module-nav-item {
  border: 1px solid transparent;
  border-radius: 12px;
  background: rgba(24, 20, 16, 0.55);
  color: var(--text);
  text-decoration: none;
  padding: 10px;
  font-size: 14px;
  font-weight: 600;
}

.module-nav-item.active {
  border-color: rgba(212, 175, 55, 0.22);
  background: linear-gradient(135deg, rgba(212, 175, 55, 0.18), rgba(212, 175, 55, 0.06));
}
</style>
