<script setup lang="ts">
import PanelPage from './modules/PanelPage.vue'
import { onMounted } from 'vue'
import { theToast } from '@/utils/toast'

onMounted(() => {
  const url = new URL(window.location.href)
  const bind = url.searchParams.get('bind')
  if (bind === 'success') {
    theToast.success('OAuth2账号绑定成功', { duration: 5000 })
  } else if (bind === 'failed') {
    theToast.error('OAuth2账号绑定失败，请重试')
  }
})
</script>

<template>
  <div class="w-full">
    <PanelPage />
  </div>
</template>

<style></style>
