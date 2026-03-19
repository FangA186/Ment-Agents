<script setup>
import { computed, ref } from 'vue'
import { siteContent } from './content/siteContent'

const locale = ref('zh')
const page = computed(() => siteContent[locale.value])

function setLocale(nextLocale) {
  locale.value = nextLocale
}
</script>

<template>
  <div class="page-shell">
    <header class="topbar">
      <a class="brand" href="#hero">
        <span class="brand-mark">M</span>
        <span>{{ page.brand }}</span>
      </a>

      <nav class="nav" aria-label="Section navigation">
        <a v-for="item in page.nav" :key="item.id" :href="`#${item.id}`">{{ item.label }}</a>
      </nav>

      <div class="locale-toggle" role="group" aria-label="Language switcher">
        <button :class="{ active: locale === 'zh' }" @click="setLocale('zh')" type="button">中文</button>
        <button :class="{ active: locale === 'en' }" @click="setLocale('en')" type="button">EN</button>
      </div>
    </header>

    <main class="content">
      <section id="hero" class="section hero">
        <p class="kicker">{{ page.hero.kicker }}</p>
        <h1>{{ page.hero.title }}</h1>
        <p class="summary">{{ page.hero.summary }}</p>
        <ul class="pill-list">
          <li v-for="point in page.hero.highlights" :key="point">{{ point }}</li>
        </ul>
      </section>

      <section id="architecture" class="section card">
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

      <section id="capabilities" class="section card">
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

      <section id="api" class="section card">
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
                <td><span class="mono chip">{{ row.method }}</span></td>
                <td><code>{{ row.endpoint }}</code></td>
                <td><code>{{ row.request }}</code></td>
                <td><code>{{ row.response }}</code></td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <section id="roadmap" class="section card">
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

      <section id="risk" class="section card">
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

    <footer class="footer">
      <p>{{ page.footer.reference }}</p>
      <p>{{ page.footer.notice }}</p>
    </footer>
  </div>
</template>
