<template>
  <div>
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

    <!-- 返回首页 / 登录 / 注册 -->
    <div class="flex justify-between items-center">
      <!-- 返回首页 -->
      <div>
        <BaseButton
          @click="$router.push('/')"
          class="text-gray-600 mt-6 rounded-md !shadow-none !border-none !ring-0 !bg-transparent group"
          title="返回首页"
        >
          <Arrow
            class="w-7 h-7 rotate-180 transition-transform duration-200 group-hover:-translate-x-1"
          />
        </BaseButton>
      </div>
      <!-- 登录 / 注册 -->
      <div>
        <BaseButton
          :icon="Auth"
          @click="$router.push('/auth')"
          class="text-gray-600 mt-6 rounded-md w-8 h-8"
          title="登录 / 注册"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import BaseButton from '@/components/common/BaseButton.vue'
import Arrow from '@/components/icons/arrow.vue'
import Auth from '@/components/icons/auth.vue'
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
