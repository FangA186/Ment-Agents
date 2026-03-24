import { computed, reactive } from 'vue'
import { apiListProjects } from '../services/api'

const fallbackProjects = [
  { id: 'proj-001', name: '营销活动自动化系统', status: '进行中', updatedAt: '2026-03-23 11:20' },
]

export const projectState = reactive({
  projects: [...fallbackProjects],
  loading: false,
  loaded: false,
  error: '',
})

export const projects = computed(() => projectState.projects)

export async function ensureProjects() {
  if (projectState.loaded || projectState.loading) return
  await refreshProjects()
}

export async function refreshProjects() {
  projectState.loading = true
  projectState.error = ''
  try {
    const list = await apiListProjects()
    projectState.projects = list.length ? list : [...fallbackProjects]
    projectState.loaded = true
  } catch (error) {
    projectState.error = error.message || '加载项目失败'
    projectState.projects = [...fallbackProjects]
  } finally {
    projectState.loading = false
  }
}

export function getProjectById(projectId) {
  return projectState.projects.find((item) => item.id === projectId) || projectState.projects[0]
}
