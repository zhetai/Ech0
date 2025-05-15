<script setup lang="ts">
import { RouterView } from 'vue-router'
import { onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { watch } from 'vue'
import { useSettingStore } from '@/stores/settting'
import { storeToRefs } from 'pinia'
import { fetchGetStatus } from './service/api'
import { useRouter } from 'vue-router'

const userStore = useUserStore()
const settingStore = useSettingStore()
const { setSystemReady } = settingStore
const { SystemSetting } = storeToRefs(settingStore)
const router = useRouter()

watch(
  () => SystemSetting.value.site_title,
  (title) => {
    if (title) document.title = title
  },
  { immediate: true },
)

onMounted(async () => {
  // 自动登录
  await userStore.autoLogin()
})
</script>

<template>
  <RouterView />
</template>

<style scoped></style>
