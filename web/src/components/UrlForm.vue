<template>
  <div class="url-form">
    <el-form @submit.prevent="handleSubmit">
      <el-form-item>
        <el-input
          v-model="url"
          :placeholder="$t('form.placeholder')"
          size="large"
          clearable
          @keyup.enter="handleSubmit"
        >
          <template #append>
            <el-button type="primary" @click="handleSubmit" :loading="loading">
              {{ $t('form.button') }}
            </el-button>
          </template>
        </el-input>
      </el-form-item>
    </el-form>

    <div v-if="result" class="result">
      <el-divider />
      <p class="result-label">{{ $t('result.shortUrl') }}</p>
      <el-input
        :model-value="result"
        size="large"
        readonly
      >
        <template #append>
          <el-button @click="copyUrl">
            {{ copied ? $t('form.copied') : $t('form.copy') }}
          </el-button>
        </template>
      </el-input>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'

const { t } = useI18n()

const url = ref('')
const result = ref('')
const loading = ref(false)
const copied = ref(false)

async function handleSubmit() {
  if (!url.value.trim()) {
    ElMessage.warning(t('error.emptyUrl'))
    return
  }

  try {
    new URL(url.value)
  } catch {
    ElMessage.warning(t('error.invalidUrl'))
    return
  }

  loading.value = true
  result.value = ''

  try {
    const resp = await fetch('/api/shorten', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ url: url.value }),
    })

    const data = await resp.json()
    if (data.code === 0) {
      result.value = data.data.short_url
    } else {
      ElMessage.error(data.msg || t('error.generateFailed'))
    }
  } catch {
    ElMessage.error(t('error.generateFailed'))
  } finally {
    loading.value = false
  }
}

async function copyUrl() {
  try {
    await navigator.clipboard.writeText(result.value)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  } catch {
    ElMessage.error('Copy failed')
  }
}
</script>

<style scoped>
.url-form {
  width: 100%;
  max-width: 600px;
}

.result {
  text-align: left;
}

.result-label {
  margin-bottom: 8px;
  font-weight: 500;
  color: #606266;
}
</style>
