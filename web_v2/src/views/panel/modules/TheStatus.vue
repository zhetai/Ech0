<template>
  <div class="w-full px-2">
    <!-- 系统管理员 -->
    <div class="flex justify-start items-center">
      <h1 class="text-gray-500 text-md">当前系统管理员：</h1>
      <span class="font-bold text-md text-gray-600 ml-2">{{ status?.username }}</span>
    </div>
    <!-- 当前登录用户 -->
    <div class="flex justify-start items-center">
      <h1 class="text-gray-500 text-md">当前登录用户：</h1>
      <span class="font-bold text-md text-gray-600 ml-2">
        {{ userStore?.user?.username ? userStore.user.username : '当前未登录' }}
      </span>
    </div>
    <!-- 当前共有Ech0 -->
    <div class="flex justify-start items-center">
      <h1 class="text-gray-500 text-md">当前Ech0共有：</h1>
      <span class="font-bold text-md text-gray-600 ml-2">{{ status?.total_messages }}</span>
      <span class="text-gray-500 text-md ml-2">条</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { fetchGetStatus } from '@/service/api'
import { onMounted, ref } from 'vue'
import { useUserStore } from '@/stores/user'

const status = ref<App.Api.Ech0.Status>()
const userStore = useUserStore()

onMounted(async () => {
  await fetchGetStatus().then((res) => {
    status.value = res.data
  })
})
</script>
