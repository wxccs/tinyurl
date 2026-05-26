<template>
  <div class="home">
    <LanguageSelector />
    <div class="content">
      <h1 class="title">{{ pageTitle }}</h1>
      <p class="subtitle">{{ $t('app.subtitle') }}</p>
      <UrlForm />
    </div>
    <AppFooter />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import UrlForm from '../components/UrlForm.vue'
import AppFooter from '../components/AppFooter.vue'
import LanguageSelector from '../components/LanguageSelector.vue'

const { t } = useI18n()
const pageTitle = ref(t('app.title'))

onMounted(async () => {
  try {
    const resp = await fetch('/api/config')
    const data = await resp.json()
    if (data.code === 0 && data.data.title) {
      pageTitle.value = data.data.title
      document.title = data.data.title
    }
  } catch {
    // ignore
  }
})
</script>

<style scoped>
.home {
  position: relative;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.content {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
}

.title {
  font-size: 28px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 8px;
}

.subtitle {
  font-size: 16px;
  color: #909399;
  margin-bottom: 32px;
}
</style>
