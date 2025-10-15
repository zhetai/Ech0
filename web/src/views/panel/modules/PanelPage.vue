<template>
  <div
    class="px-4 pb-4 py-2 mb-10 mx-auto flex flex-col min-h-screen max-w-screen-lg border border-gray-300 rounded-md mt-4"
  >
    <h1 class="text-6xl italic font-bold text-center text-gray-300 mb-8">Ech0 Panel</h1>

    <!-- 主内容区 -->
    <div class="mx-auto flex px-2 my-4 w-full">
      <div class="flex flex-col gap-3 w-1/5 pr-8">
        <!-- 返回首页 -->
        <BaseButton
          @click="router.push('/')"
          class="text-gray-600 rounded-md !shadow-none !border-none !ring-0 !bg-transparent group"
          title="返回首页"
        >
          <Arrow
            class="w-9 h-9 rotate-180 transition-transform duration-200 group-hover:-translate-x-1"
          />
        </BaseButton>

        <div class="h-px bg-gray-300 mx-2" />

        <!-- 状态 -->
        <BaseButton
          :icon="Status"
          @click="router.push('/panel/status')"
          class="flex items-center justify-center gap-2"
          :class="
            currentRoute === 'panel-status'
              ? 'text-gray-800 rounded-md transition-all !bg-gray-200'
              : 'text-gray-600 rounded-md transition-all'
          "
          title="状态"
        >
          状态
        </BaseButton>

        <!-- 设置 -->
        <BaseButton
          :icon="Setting"
          @click="router.push('/panel/setting')"
          class="flex items-center justify-center gap-2"
          :class="
            currentRoute === 'panel-setting'
              ? 'text-gray-800 rounded-md transition-all !bg-gray-200'
              : 'text-gray-600 rounded-md transition-all'
          "
          title="设置"
        >
          设置
        </BaseButton>

        <!-- 个人中心 -->
        <BaseButton
          :icon="User"
          @click="router.push('/panel/user')"
          class="flex items-center justify-center gap-2"
          :class="
            currentRoute === 'panel-user'
              ? 'text-gray-800 rounded-md transition-all !bg-gray-200'
              : 'text-gray-600 rounded-md transition-all'
          "
          title="个人中心"
        >
          个人中心
        </BaseButton>

        <!-- 高级 -->
        <BaseButton
          :icon="Others"
          @click="router.push('/panel/advance')"
          class="flex items-center justify-center gap-2"
          :class="
            currentRoute === 'panel-advance'
              ? 'text-gray-800 rounded-md transition-all !bg-gray-200'
              : 'text-gray-600 rounded-md transition-all'
          "
          title="高级"
        >
          高级
        </BaseButton>

        <div class="h-px bg-gray-300 mx-2" />

        <!-- 退出登录 -->
        <BaseButton
          :icon="Logout"
          @click="handleLogout"
          class="flex items-center justify-center gap-2 rounded-md"
          title="退出登录"
        >
          退出登入
        </BaseButton>

        <!-- 登录 / 注册 -->
        <BaseButton
          :icon="Auth"
          @click="router.push('/auth')"
          class="flex items-center justify-center gap-2 rounded-md"
          title="登录 / 注册"
        >
          登录 / 注册
        </BaseButton>
      </div>

      <!-- 路由内容 -->
      <router-view />
    </div>
  </div>
</template>

<script setup lang="ts">
import BaseButton from '@/components/common/BaseButton.vue'
import Arrow from '@/components/icons/arrow.vue'
import User from '@/components/icons/user.vue'
import Auth from '@/components/icons/auth.vue'
import Status from '@/components/icons/status.vue'
import Others from '@/components/icons/theothers.vue'
import Setting from '@/components/icons/setting.vue'
import Logout from '@/components/icons/logout.vue'
import { computed } from 'vue'
import { useUserStore } from '@/stores/user'
import { useRouter, useRoute } from 'vue-router'
import { theToast } from '@/utils/toast'
import { useBaseDialog } from '@/composables/useBaseDialog'

const { openConfirm } = useBaseDialog()

const userStore = useUserStore()
const router = useRouter()
const route = useRoute()

const currentRoute = computed(() => route.name as string)

const handleLogout = () => {
  // 检查是否登录
  if (!userStore.isLogin) {
    theToast.info('当前未登录')
    return
  }

  // 弹出浏览器确认框
  openConfirm({
    title: '确定要退出登录吗？',
    description: '',
    onConfirm: () => {
      // 清除用户信息
      userStore.logout()
    },
  })
}
</script>
