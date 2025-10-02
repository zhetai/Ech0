<template>
  <div class="px-9 md:px-11">
    <div class="rounded-md shadow-sm hover:shadow-md ring-1 ring-gray-200 ring-inset bg-white p-4">
      <h2 class="text-gray-600 font-bold text-lg mb-1 flex items-center">
        <StatusInfo class="mr-2" />系统状态:
      </h2>

      <div class="text-gray-500 text-md">
        <!-- 系统管理员 -->
        <!-- <div>
          <h1>
            当前系统管理员：
            <span class="ml-2">{{ status?.username }}</span>
          </h1>
        </div> -->
        <!-- 当前登录用户 -->
        <div>
          <h1>
            <span class="font-bold">当前登录用户:</span>
            <span class="ml-2">
              {{ userStore?.user?.username ? userStore.user.username : '当前未登录' }}
            </span>
          </h1>
        </div>
        <!-- 当前共有Ech0 -->
        <div>
          <h1>
            <span class="font-bold">已发布Echos:</span>
            <span class="ml-2">{{ status?.total_echos }}</span>
            <span class="font-bold ml-2">条</span>
          </h1>
        </div>
      </div>
    </div>
  </div>

</template>

<script setup lang="ts">
import { fetchGetStatus } from '@/service/api'
import { onMounted, ref } from 'vue'
import { useUserStore } from '@/stores/user'
import StatusInfo from '../icons/statusinfo.vue'

const status = ref<App.Api.Ech0.Status>()
const userStore = useUserStore()

onMounted(() => {
  fetchGetStatus().then((res) => {
    status.value = res.data
  })
})
</script>

<style scoped></style>
