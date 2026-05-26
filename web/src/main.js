import { createApp } from 'vue'
import 'element-plus/theme-chalk/index.css'
import App from './App.vue'
import i18n, { loadLocaleMessages } from './i18n'

const initialLocale = i18n.global.locale.value

;(async () => {
  await loadLocaleMessages('en')
  if (initialLocale !== 'en') {
    await loadLocaleMessages(initialLocale)
  }
  const app = createApp(App)
  app.use(i18n)
  app.mount('#app')
})()
