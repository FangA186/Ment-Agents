<script setup>
import { RouterLink } from 'vue-router'

const quickStats = [
  { label: 'Mock 任务流', value: '12 条' },
  { label: 'Agent 角色模板', value: '9 个' },
  { label: '当前演示模式', value: '纯前端' },
]

const pathways = [
  {
    title: '阅读产品介绍',
    desc: '先快速理解平台定位、架构主线、API 示例与风险控制。',
    to: '/read',
    action: '进入阅读页',
    mode: 'read',
  },
  {
    title: '进入群聊编排',
    desc: '直接在群聊页输入需求，并完成 IR 产出与 Agent 详情查看。',
    to: '/studio/projects',
    action: '进入控制台',
    mode: 'build',
  },
]
</script>

<template>
  <main class="entry-home">
    <section class="entry-hero">
      <p class="entry-kicker">Ment-Agents</p>
      <h1>从这里开始你的多智能体编排</h1>
      <p class="entry-summary">
        首页作为统一入口。你可以先阅读，再进入实操；也可以直接进入群聊编排进行需求确认与产出推进。
      </p>
    </section>

    <section class="entry-path-grid">
      <RouterLink v-for="item in pathways" :key="item.to" :to="item.to" class="entry-path" :class="item.mode">
        <p class="path-mode">{{ item.mode === 'read' ? '阅读路径' : '实操路径' }}</p>
        <h2>{{ item.title }}</h2>
        <p class="path-desc">{{ item.desc }}</p>
        <span class="path-action">{{ item.action }} →</span>
      </RouterLink>
    </section>

    <section class="entry-signal">
      <article v-for="item in quickStats" :key="item.label" class="signal-card">
        <p>{{ item.label }}</p>
        <strong>{{ item.value }}</strong>
      </article>
    </section>
  </main>
</template>

<style scoped>
.entry-home {
  min-height: 100dvh;
  width: min(1180px, 100% - 40px);
  margin: 0 auto;
  display: grid;
  align-content: center;
  gap: 18px;
  padding: 28px 0 36px;
}

.entry-hero {
  padding: 24px 22px;
  border-radius: 20px;
  background: linear-gradient(140deg, rgba(18, 16, 13, 0.85), rgba(11, 10, 8, 0.92));
  box-shadow: 0 10px 24px rgba(0, 0, 0, 0.36);
}

.entry-kicker {
  margin: 0;
  color: var(--accent-2);
  letter-spacing: 0.08em;
  text-transform: uppercase;
  font-size: 12px;
}

.entry-hero h1 {
  margin: 10px 0 8px;
  font-size: clamp(30px, 4.2vw, 54px);
  line-height: 1.08;
  font-family: 'Noto Serif SC', 'Songti SC', serif;
}

.entry-summary {
  margin: 0;
  max-width: 760px;
  color: var(--muted);
  font-size: 15px;
}

.entry-path-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
}

.entry-path {
  border-radius: 18px;
  text-decoration: none;
  color: inherit;
  min-height: 230px;
  padding: 22px;
  display: grid;
  align-content: space-between;
  transition: transform 0.24s ease, background 0.24s ease;
}

.entry-path.read {
  background:
    linear-gradient(180deg, rgba(26, 22, 16, 0.78), rgba(16, 14, 11, 0.92)),
    radial-gradient(circle at 20% 24%, rgba(212, 175, 55, 0.18), transparent 50%);
}

.entry-path.build {
  background:
    linear-gradient(180deg, rgba(20, 18, 14, 0.78), rgba(13, 12, 10, 0.92)),
    radial-gradient(circle at 82% 18%, rgba(121, 168, 198, 0.2), transparent 46%);
}

.entry-path:hover {
  transform: translateY(-3px);
}

.path-mode {
  margin: 0;
  font-size: 12px;
  color: var(--muted);
}

.entry-path h2 {
  margin: 8px 0;
  font-size: clamp(23px, 3.2vw, 34px);
}

.path-desc {
  margin: 0;
  color: #ddd2be;
  font-size: 14px;
  max-width: 38ch;
}

.path-action {
  justify-self: start;
  margin-top: 18px;
  border-radius: 999px;
  border: 1px solid rgba(212, 175, 55, 0.26);
  padding: 7px 13px;
  font-size: 13px;
  color: #f2e3be;
}

.entry-signal {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
}

.signal-card {
  margin: 0;
  border-radius: 14px;
  padding: 12px 14px;
  background: rgba(18, 16, 13, 0.7);
}

.signal-card p {
  margin: 0;
  color: var(--muted);
  font-size: 13px;
}

.signal-card strong {
  display: block;
  margin-top: 5px;
  font-size: 24px;
  color: #f2e3be;
}

@media (max-width: 920px) {
  .entry-path-grid,
  .entry-signal {
    grid-template-columns: 1fr;
  }
}
</style>

