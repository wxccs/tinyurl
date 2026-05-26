import { createI18n } from 'vue-i18n'

export const SUPPORTED_LOCALES = {
  en: 'English',
  zh: '简体中文',
  'zh-Hant': '繁體中文',
  hi: 'हिन्दी',
  es: 'Español',
  fr: 'Français',
  ar: 'العربية',
  pt: 'Português',
  ru: 'Русский',
  ja: '日本語',
  de: 'Deutsch',
  ko: '한국어',
  vi: 'Tiếng Việt',
  tr: 'Türkçe',
  it: 'Italiano',
  th: 'ไทย',
  pl: 'Polski',
  uk: 'Українська',
  ro: 'Română',
  nl: 'Nederlands',
  hu: 'Magyar',
  id: 'Bahasa Indonesia',
  ms: 'Bahasa Melayu',
  sv: 'Svenska',
  cs: 'Čeština',
  bn: 'বাংলা',
  ta: 'தமிழ்',
  fa: 'فارسی',
  ur: 'اردو',
  el: 'Ελληνικά',
  he: 'עברית',
}

const LOCALE_CODES = new Set(Object.keys(SUPPORTED_LOCALES))
const RTL_LOCALE_CODES = new Set(['ar', 'fa', 'ur', 'he'])
const HANT_REGIONS = new Set(['TW', 'HK', 'MO'])
const loadedLocales = new Set()

const localeModules = import.meta.glob('./locales/*.json')

function resolveLocale(lang) {
  if (LOCALE_CODES.has(lang)) {
    return lang
  }
  const [primary, region] = lang.split('-')
  if (primary === 'zh' && HANT_REGIONS.has(region)) {
    return 'zh-Hant'
  }
  if (LOCALE_CODES.has(primary)) {
    return primary
  }
  return null
}

function detectLocale() {
  const stored = localStorage.getItem('locale')
  if (stored && LOCALE_CODES.has(stored)) {
    return stored
  }

  const browserLangs = [
    navigator.language,
    ...(navigator.languages || []),
  ]

  for (const lang of browserLangs) {
    const resolved = resolveLocale(lang)
    if (resolved) {
      return resolved
    }
  }

  return 'en'
}

const i18n = createI18n({
  legacy: false,
  locale: detectLocale(),
  fallbackLocale: 'en',
  messages: {},
})

export async function loadLocaleMessages(locale) {
  if (loadedLocales.has(locale)) {
    return
  }

  const key = `./locales/${locale}.json`
  const load = localeModules[key]
  if (!load) {
    console.warn(`Locale file not found: ${key}`)
    return
  }

  const module = await load()
  i18n.global.setLocaleMessage(locale, module.default || module)
  loadedLocales.add(locale)
}

export async function setLocale(locale) {
  if (!LOCALE_CODES.has(locale)) {
    console.warn(`Unsupported locale: ${locale}`)
    return
  }

  await loadLocaleMessages(locale)
  i18n.global.locale.value = locale
  localStorage.setItem('locale', locale)
  document.documentElement.setAttribute('lang', locale)
  document.documentElement.setAttribute('dir', RTL_LOCALE_CODES.has(locale) ? 'rtl' : 'ltr')
}

export default i18n
