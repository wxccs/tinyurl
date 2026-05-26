<template>
  <div class="app-footer" v-if="beian.miit || beian.mps">
    <a
      v-if="beian.miit"
      href="https://beian.miit.gov.cn"
      target="_blank"
      class="footer-link"
    >
      <img
        class="beian-logo"
        src="/miit_logo.svg"
        alt=""
      />
      {{ beian.miit }}
    </a>
    <span v-if="beian.miit && beian.mps" class="footer-separator">|</span>
    <a
      v-if="beian.mps"
      :href="mpsUrl"
      target="_blank"
      class="footer-link"
    >
      <img
        class="beian-logo"
        src="/beian_logo.png"
        alt=""
      />
      {{ beian.mps }}
    </a>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

const beian = ref({ miit: '', mps: '' })

const mpsUrl = computed(() => {
  if (!beian.value.mps) return ''
  const digits = beian.value.mps.match(/\d+/g)
  const recordcode = digits ? digits.join('') : ''
  return `http://www.beian.gov.cn/portal/registerSystemInfo?recordcode=${recordcode}`
})

onMounted(async () => {
  try {
    const resp = await fetch('/api/config')
    const data = await resp.json()
    if (data.code === 0) {
      beian.value = data.data.beian || {}
    }
  } catch {
    // ignore
  }
})
</script>

<style scoped>
.app-footer {
  padding: 16px 0;
  text-align: center;
  color: #909399;
  font-size: 13px;
}

.footer-link {
  color: #909399;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.footer-link:hover {
  color: #409eff;
}

.footer-separator {
  display: inline-flex;
  align-items: center;
  margin: 0 8px;
}

.beian-logo {
  width: 20px;
  height: 22px;
}
</style>
