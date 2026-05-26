<template>
  <div class="language-selector">
    <el-select
      v-model="currentLocale"
      @change="handleChange"
      size="small"
      :loading="loading"
      filterable
    >
      <el-option
        v-for="(name, code) in supportedLocales"
        :key="code"
        :label="name"
        :value="code"
      />
    </el-select>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { SUPPORTED_LOCALES, setLocale } from '../i18n'

const { locale } = useI18n()
const loading = ref(false)
const supportedLocales = SUPPORTED_LOCALES

const currentLocale = computed({
  get: () => locale.value,
  set: () => {},
})

async function handleChange(code) {
  loading.value = true
  try {
    await setLocale(code)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.language-selector {
  position: absolute;
  top: 16px;
  right: 16px;
  z-index: 100;
}

.language-selector :deep(.el-select) {
  width: 160px;
}
</style>
