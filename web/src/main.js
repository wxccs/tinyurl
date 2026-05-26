import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'
import i18n, { loadLocaleMessages } from './i18n'

const initialLocale = i18n.global.locale.value

;(async () => {
  if (initialLocale !== 'en') {
    await loadLocaleMessages(initialLocale)
  }
  const app = createApp(App)
  app.use(ElementPlus)
  app.use(i18n)
  app.mount('#app')
})()
