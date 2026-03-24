import { nextTick } from 'vue'
import { mount } from '@vue/test-utils'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

vi.mock('../services/api', () => {
  const project = {
    id: 'proj-001',
    name: '测试项目 Alpha',
    status: '进行中',
    updatedAt: '2026-03-24 10:00',
  }

  const agentGraph = {
    nodes: [
      { id: 'planner', name: '规划Agent' },
      { id: 'executor', name: '执行Agent' },
    ],
    edges: [{ from: 'planner', to: 'executor' }],
  }

  return {
    apiListProjects: vi.fn(async () => [project]),
    apiGetProject: vi.fn(async () => project),
    apiGetChat: vi.fn(async () => [
      { id: 'm1', role: 'user', content: '需求: 搭建测试流程', createdAt: '2026-03-24 10:01:00' },
      { id: 'm2', role: 'assistant', content: '已记录目标与约束，等待进一步确认。', createdAt: '2026-03-24 10:01:10' },
    ]),
    apiSendChat: vi.fn(async (_projectId, message) => [
      { id: 'm1', role: 'user', content: '需求: 搭建测试流程', createdAt: '2026-03-24 10:01:00' },
      { id: 'm2', role: 'assistant', content: '已记录目标与约束，等待进一步确认。', createdAt: '2026-03-24 10:01:10' },
      { id: 'm3', role: 'user', content: message, createdAt: '2026-03-24 10:02:00' },
      { id: 'm4', role: 'assistant', content: `已收到：${message}`, createdAt: '2026-03-24 10:02:08' },
    ]),
    apiCompile: vi.fn(async () => ({
      version: 'v2',
      generatedAt: '2026-03-24 10:03:00',
    })),
    apiAssemble: vi.fn(async () => agentGraph),
    apiGetArtifacts: vi.fn(async () => ({
      ir: {
        version: 'v2',
        generatedAt: '2026-03-24 10:03:00',
      },
      agentGraph,
    })),
  }
})

async function flushUi() {
  await Promise.resolve()
  await nextTick()
  await Promise.resolve()
}

async function mountAt(path) {
  vi.resetModules()

  const [{ default: App }, { default: router }] = await Promise.all([import('../App.vue'), import('../router')])

  await router.push(path)
  await router.isReady()

  const wrapper = mount(App, {
    attachTo: document.body,
    global: {
      plugins: [router],
    },
  })

  await flushUi()

  return { wrapper, router }
}

function getButtonByText(wrapper, label) {
  return wrapper.findAll('button').find((button) => button.text() === label)
}

describe('工作台多页面冒烟检查', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  afterEach(() => {
    vi.useRealTimers()
    document.body.innerHTML = ''
  })

  it('项目库页面可渲染并进入群聊编排页', async () => {
    const { wrapper, router } = await mountAt('/studio/projects')

    expect(wrapper.text()).toContain('项目库')
    expect(wrapper.text()).toContain('测试项目 Alpha')

    const folderLink = wrapper.find('.folder-card')
    expect(folderLink.attributes('href')).toBe('/studio/projects/proj-001/chat')

    await router.push(folderLink.attributes('href'))
    await flushUi()

    expect(router.currentRoute.value.fullPath).toBe('/studio/projects/proj-001/chat')
    expect(wrapper.text()).toContain('群聊对话（你 + 分析Agent）')

    wrapper.unmount()
  })

  it('群聊页可发送消息并显示回包', async () => {
    const { wrapper } = await mountAt('/studio/projects/proj-001/chat')

    const input = wrapper.find('.chat-input-row input')
    await input.setValue('补充一个验收条件')
    await wrapper.find('.chat-input-row button').trigger('click')
    await flushUi()

    expect(wrapper.text()).toContain('已收到：补充一个验收条件')

    wrapper.unmount()
  })

  it('群聊生成流程可推进并跳转到组装器页', async () => {
    const { wrapper, router } = await mountAt('/studio/projects/proj-001/chat')

    vi.useFakeTimers()
    await wrapper.find('.chat-head .btn').trigger('click')
    await vi.advanceTimersByTimeAsync(5000)
    await flushUi()

    expect(router.currentRoute.value.fullPath).toBe('/studio/projects/proj-001/assembler')
    expect(wrapper.text()).toContain('组装器页面')

    wrapper.unmount()
  })

  it('组装器页可展示 AgentGraph 和四件套信息', async () => {
    const { wrapper } = await mountAt('/studio/projects/proj-001/assembler')

    expect(wrapper.text()).toContain('组装器页面')
    expect(wrapper.text()).toContain('AgentGraph')
    expect(wrapper.text()).toContain('Agent 四件套')

    wrapper.unmount()
  })

  it('编排页可执行一步推进并刷新日志', async () => {
    const { wrapper } = await mountAt('/studio/projects/proj-001/orchestrator')

    const stepButton = getButtonByText(wrapper, '推进一步')
    expect(stepButton).toBeTruthy()

    await stepButton.trigger('click')
    await flushUi()

    expect(wrapper.find('.log-panel p').text()).toContain('推进一步')

    wrapper.unmount()
  })

  it('调试与自愈页可执行重试并更新结果提示', async () => {
    const { wrapper } = await mountAt('/studio/projects/proj-001/debug')

    expect(wrapper.text()).toContain('调试与自愈页面')

    const retryButton = getButtonByText(wrapper, '触发重试')
    expect(retryButton).toBeTruthy()

    await retryButton.trigger('click')
    await flushUi()

    expect(wrapper.text()).toContain('最近动作：已对')
    expect(wrapper.text()).toContain('触发第 2 次重试')

    wrapper.unmount()
  })

  it('记忆页可切换上下文层并更新详情展示', async () => {
    const { wrapper } = await mountAt('/studio/projects/proj-001/memory')

    expect(wrapper.find('.detail-box').text()).toContain('Global Memory')

    const layerButtons = wrapper.findAll('.layer-item')
    await layerButtons[1].trigger('click')
    await flushUi()

    expect(wrapper.find('.detail-box').text()).toContain('Task Memory')
    expect(wrapper.text()).toContain('RAG 检索命中')

    wrapper.unmount()
  })
})
