const API_BASE = (import.meta.env.VITE_API_BASE_URL || 'http://127.0.0.1:8080/api').replace(/\/$/, '')

async function request(path, options = {}) {
  const response = await fetch(`${API_BASE}${path}`, {
    headers: {
      'Content-Type': 'application/json',
      ...(options.headers || {}),
    },
    ...options,
  })

  const text = await response.text()
  const data = text ? JSON.parse(text) : {}

  if (!response.ok) {
    throw new Error(data.error || `请求失败: ${response.status}`)
  }

  return data
}

export async function apiListProjects() {
  const data = await request('/projects')
  return data.projects || []
}

export async function apiCreateProject(name) {
  const data = await request('/projects', {
    method: 'POST',
    body: JSON.stringify({ name }),
  })
  return data.project
}

export async function apiGetProject(projectId) {
  const data = await request(`/projects/${projectId}`)
  return data.project
}

export async function apiGetChat(projectId) {
  const data = await request(`/projects/${projectId}/chat`)
  return data.messages || []
}

export async function apiSendChat(projectId, message) {
  const data = await request(`/projects/${projectId}/chat`, {
    method: 'POST',
    body: JSON.stringify({ message }),
  })
  return data.messages || []
}

export async function apiCompile(projectId) {
  const data = await request(`/projects/${projectId}/compile`, { method: 'POST' })
  return data.ir
}

export async function apiAssemble(projectId) {
  const data = await request(`/projects/${projectId}/assemble`, { method: 'POST' })
  return data.agentGraph
}

export async function apiGetArtifacts(projectId) {
  const data = await request(`/projects/${projectId}/artifacts`)
  return data.artifacts || {}
}
