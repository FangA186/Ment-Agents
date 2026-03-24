<script setup>
import { onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { ensureProjects, projects, projectState } from '../content/projectStore'

onMounted(() => {
  ensureProjects().catch((error) => {
    console.error(error)
  })
})
</script>

<template>
  <section class="panel project-library-page">
    <div class="panel-head">
      <h2>项目库</h2>
      <span class="stage-pill running">共 {{ projects.length }} 个项目</span>
    </div>

    <p class="panel-intro">这里只展示项目文件夹。点击任意项目后，直接进入该项目的智能体编排页面。</p>

    <p v-if="projectState.error" class="error-tip">{{ projectState.error }}</p>

    <div class="folder-grid">
      <RouterLink
        v-for="project in projects"
        :key="project.id"
        :to="`/studio/projects/${project.id}/chat`"
        class="folder-card"
      >
        <div class="folder-icon" aria-hidden="true"></div>
        <div>
          <h3>{{ project.name }}</h3>
          <p>{{ project.status }} · {{ project.updatedAt }}</p>
        </div>
      </RouterLink>
    </div>
  </section>
</template>

<style scoped>
.project-library-page {
  width: 100%;
  display: grid;
  gap: 12px;
}

.error-tip {
  margin: 0;
  color: #e6b6b6;
  font-size: 13px;
}

.folder-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.folder-card {
  border-radius: 14px;
  border: 1px solid rgba(212, 175, 55, 0.12);
  background: rgba(16, 14, 11, 0.86);
  padding: 14px;
  text-decoration: none;
  color: inherit;
  display: grid;
  grid-template-columns: 42px 1fr;
  align-items: center;
  gap: 10px;
}

.folder-icon {
  width: 36px;
  height: 30px;
  border-radius: 8px;
  background: linear-gradient(140deg, rgba(212, 175, 55, 0.55), rgba(97, 76, 33, 0.7));
  position: relative;
}

.folder-icon::before {
  content: '';
  position: absolute;
  left: 3px;
  top: -6px;
  width: 16px;
  height: 8px;
  border-radius: 4px 4px 0 0;
  background: rgba(212, 175, 55, 0.5);
}

.folder-card h3 {
  margin: 0;
  color: #f2dc9f;
  font-size: 22px;
}

.folder-card p {
  margin: 6px 0 0;
  color: var(--muted);
  font-size: 13px;
}

@media (max-width: 1024px) {
  .folder-grid {
    grid-template-columns: 1fr;
  }
}
</style>
