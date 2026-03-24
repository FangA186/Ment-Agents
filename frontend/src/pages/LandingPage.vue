<script setup>
import { computed, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { siteContent } from '../content/siteContent'

const locale = ref('zh')
const page = computed(() => siteContent[locale.value])

function setLocale(nextLocale) {
  locale.value = nextLocale
}
</script>

<template>
  <div class="landing-shell">
    <header class="landing-topbar">
      <a class="landing-brand" href="#hero">
        <span class="brand-mark">M</span>
        <span>{{ page.brand }}</span>
      </a>

      <nav class="landing-nav" aria-label="Section navigation">
        <a v-for="item in page.nav" :key="item.id" :href="`#${item.id}`">{{ item.label }}</a>
      </nav>

      <div class="locale-toggle" role="group" aria-label="Language switcher">
        <button :class="{ active: locale === 'zh' }" @click="setLocale('zh')" type="button">中文</button>
        <button :class="{ active: locale === 'en' }" @click="setLocale('en')" type="button">EN</button>
      </div>
    </header>

    <main class="landing-content">
      <section id="hero" class="landing-section hero">
        <p class="kicker">{{ page.hero.kicker }}</p>
        <h1>{{ page.hero.title }}</h1>
        <p class="summary">{{ page.hero.summary }}</p>
        <ul class="pill-list">
          <li v-for="point in page.hero.highlights" :key="point">{{ point }}</li>
        </ul>
        <div class="hero-actions">
          <RouterLink class="btn ghost" to="/">返回入口首页</RouterLink>
          <RouterLink class="btn primary" to="/studio/projects">进入群聊编排</RouterLink>
          <RouterLink class="btn ghost" to="/studio/projects">进入组装器页面</RouterLink>
        </div>
      </section>

      <section id="architecture" class="landing-section card">
        <h2>{{ page.architecture.title }}</h2>
        <p class="section-intro">{{ page.architecture.intro }}</p>
        <div class="grid-2">
          <article class="soft-block">
            <h3>{{ page.architecture.flowTitle }}</h3>
            <ol>
              <li v-for="node in page.architecture.flow" :key="node">{{ node }}</li>
            </ol>
          </article>
          <article class="soft-block">
            <h3>{{ page.architecture.bundleTitle }}</h3>
            <ul>
              <li v-for="item in page.architecture.bundle" :key="item">{{ item }}</li>
            </ul>
          </article>
        </div>
      </section>

      <section id="capabilities" class="landing-section card">
        <h2>{{ page.capabilities.title }}</h2>
        <p class="section-intro">{{ page.capabilities.intro }}</p>
        <div class="grid-2">
          <article v-for="block in page.capabilities.blocks" :key="block.title" class="soft-block">
            <h3>{{ block.title }}</h3>
            <ul>
              <li v-for="item in block.items" :key="item">{{ item }}</li>
            </ul>
          </article>
        </div>
      </section>

      <section id="api" class="landing-section card">
        <h2>{{ page.api.title }}</h2>
        <p class="section-intro">{{ page.api.intro }}</p>
        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                <th>{{ page.api.table.method }}</th>
                <th>{{ page.api.table.endpoint }}</th>
                <th>{{ page.api.table.request }}</th>
                <th>{{ page.api.table.response }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="row in page.api.rows" :key="row.endpoint">
                <td><span class="chip">{{ row.method }}</span></td>
                <td><code>{{ row.endpoint }}</code></td>
                <td><code>{{ row.request }}</code></td>
                <td><code>{{ row.response }}</code></td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <section id="roadmap" class="landing-section card">
        <h2>{{ page.roadmap.title }}</h2>
        <p class="section-intro">{{ page.roadmap.intro }}</p>
        <div class="timeline">
          <article v-for="phase in page.roadmap.phases" :key="phase.name" class="timeline-item">
            <p class="phase-tag">{{ phase.stage }}</p>
            <h3>{{ phase.name }}</h3>
            <p>{{ phase.goal }}</p>
          </article>
        </div>
      </section>

      <section id="risk" class="landing-section card">
        <h2>{{ page.risk.title }}</h2>
        <p class="section-intro">{{ page.risk.intro }}</p>
        <div class="grid-2">
          <article v-for="risk in page.risk.items" :key="risk.title" class="soft-block">
            <h3>{{ risk.title }}</h3>
            <p>{{ risk.detail }}</p>
          </article>
        </div>
      </section>
    </main>

    <footer class="landing-footer">
      <p>{{ page.footer.reference }}</p>
      <p>{{ page.footer.notice }}</p>
    </footer>
  </div>
</template>

<style scoped>
.landing-shell {
  max-width: 1180px;
  margin: 0 auto;
  padding: 24px 20px 40px;
}

.landing-topbar {
  position: sticky;
  top: 10px;
  z-index: 10;
  display: grid;
  grid-template-columns: auto 1fr auto;
  align-items: center;
  gap: 16px;
  padding: 10px 14px;
  border: 1px solid transparent;
  border-radius: 14px;
  background: rgba(16, 15, 13, 0.88);
  backdrop-filter: blur(8px);
  box-shadow: none;
}

.landing-brand {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  color: var(--text);
  text-decoration: none;
}

.brand-mark {
  width: 28px;
  height: 28px;
  border-radius: 999px;
  display: grid;
  place-items: center;
  font-weight: 700;
  color: #1b160d;
  background: linear-gradient(145deg, #8a6b24, var(--accent-2));
}

.landing-nav {
  display: flex;
  gap: 12px;
  justify-content: center;
  flex-wrap: wrap;
}

.landing-nav a {
  color: var(--muted);
  text-decoration: none;
  font-size: 14px;
}

.landing-nav a:hover {
  color: var(--accent-2);
}

.landing-content {
  margin-top: 16px;
  display: grid;
  gap: 14px;
}

.landing-section {
  border: 1px solid transparent;
  border-radius: 18px;
  background: linear-gradient(180deg, rgba(23, 20, 16, 0.88), rgba(16, 14, 12, 0.92));
  box-shadow: 0 10px 24px rgba(0, 0, 0, 0.36);
  padding: 20px;
}

.hero-actions {
  margin-top: 16px;
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.btn {
  border: 1px solid transparent;
  padding: 8px 12px;
  border-radius: 10px;
  text-decoration: none;
  font-weight: 700;
}

.btn.primary {
  background: linear-gradient(145deg, #8a6b24, var(--accent-2));
  color: #18130b;
}

.btn.ghost {
  border-color: rgba(212, 175, 55, 0.26);
  color: var(--text);
}

.card .grid-2,
.grid-2,
.timeline {
  margin-top: 14px;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.soft-block,
.timeline-item {
  border: 1px solid transparent;
  border-radius: 12px;
  padding: 12px;
  background: rgba(21, 18, 14, 0.84);
  box-shadow: none;
}

.section-intro,
.summary {
  color: var(--muted);
}

.pill-list {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  padding: 0;
  margin: 14px 0 0;
  list-style: none;
}

.pill-list li,
.chip {
  border: 1px solid rgba(212, 175, 55, 0.2);
  border-radius: 999px;
  padding: 4px 10px;
  color: #f2e2bd;
  background: rgba(212, 175, 55, 0.14);
  font-size: 12px;
}

.table-wrap {
  margin-top: 12px;
  overflow: auto;
  border: 1px solid transparent;
  border-radius: 12px;
  box-shadow: none;
}

table {
  border-collapse: collapse;
  width: 100%;
  min-width: 840px;
}

th,
td {
  border-bottom: 1px solid rgba(212, 175, 55, 0.08);
  text-align: left;
  vertical-align: top;
  padding: 10px;
  font-size: 13px;
}

code {
  font-family: Consolas, 'Courier New', monospace;
  background: rgba(10, 9, 8, 0.88);
  color: #e8dcc5;
  padding: 2px 6px;
  border-radius: 6px;
}

.landing-section:hover,
.soft-block:hover,
.timeline-item:hover {
  border-color: rgba(212, 175, 55, 0.2);
}

.landing-footer {
  margin-top: 14px;
  color: var(--muted);
  font-size: 13px;
  display: grid;
  gap: 4px;
}

@media (max-width: 980px) {
  .landing-topbar {
    grid-template-columns: 1fr;
  }

  .grid-2,
  .timeline {
    grid-template-columns: 1fr;
  }
}
</style>



