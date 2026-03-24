import { createRouter, createWebHistory } from 'vue-router'
import WorkbenchShell from '../layouts/WorkbenchShell.vue'
import LandingPage from '../pages/LandingPage.vue'
import AssemblerPage from '../pages/AssemblerPage.vue'
import HomeEntryPage from '../pages/HomeEntryPage.vue'
import GroupChatPage from '../pages/GroupChatPage.vue'
import AgentDetailPage from '../pages/AgentDetailPage.vue'
import CommunicationPage from '../pages/CommunicationPage.vue'
import DebugSelfHealPage from '../pages/DebugSelfHealPage.vue'
import OrchestratorPage from '../pages/OrchestratorPage.vue'
import MemoryPage from '../pages/MemoryPage.vue'
import ProjectPortalPage from '../pages/ProjectPortalPage.vue'

const defaultProjectId = 'proj-001'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'home', component: HomeEntryPage },
    { path: '/read', name: 'landing', component: LandingPage },
    { path: '/workbench', redirect: '/studio/projects' },
    { path: '/chat', redirect: `/studio/projects/${defaultProjectId}/chat` },
    { path: '/compiler', redirect: `/studio/projects/${defaultProjectId}/chat` },
    { path: '/assembler', redirect: `/studio/projects/${defaultProjectId}/assembler` },
    { path: '/communication', redirect: `/studio/projects/${defaultProjectId}/communication` },
    { path: '/debug', redirect: `/studio/projects/${defaultProjectId}/debug` },
    { path: '/orchestrator', redirect: `/studio/projects/${defaultProjectId}/orchestrator` },
    { path: '/memory', redirect: `/studio/projects/${defaultProjectId}/memory` },
    {
      path: '/studio',
      component: WorkbenchShell,
      children: [
        { path: '', redirect: '/studio/projects' },
        { path: 'projects', name: 'projects', component: ProjectPortalPage },
        { path: 'projects/:projectId', redirect: (to) => `/studio/projects/${to.params.projectId}/chat` },
        { path: 'projects/:projectId/chat', name: 'chat', component: GroupChatPage },
        { path: 'projects/:projectId/agent/:id', name: 'agent-detail', component: AgentDetailPage },
        { path: 'projects/:projectId/assembler', name: 'assembler', component: AssemblerPage },
        { path: 'projects/:projectId/communication', name: 'communication', component: CommunicationPage },
        { path: 'projects/:projectId/debug', name: 'debug', component: DebugSelfHealPage },
        { path: 'projects/:projectId/orchestrator', name: 'orchestrator', component: OrchestratorPage },
        { path: 'projects/:projectId/memory', name: 'memory', component: MemoryPage },
        { path: 'chat', redirect: `/studio/projects/${defaultProjectId}/chat` },
        { path: 'assembler', redirect: `/studio/projects/${defaultProjectId}/assembler` },
        { path: 'communication', redirect: `/studio/projects/${defaultProjectId}/communication` },
        { path: 'debug', redirect: `/studio/projects/${defaultProjectId}/debug` },
        { path: 'orchestrator', redirect: `/studio/projects/${defaultProjectId}/orchestrator` },
        { path: 'memory', redirect: `/studio/projects/${defaultProjectId}/memory` },
      ],
    },
  ],
})

export default router
