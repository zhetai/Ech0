<template>
  <div
    class="px-4 pb-4 py-2 mx-auto flex flex-col max-w-screen-lg border-3 border-stone-300 rounded-md mt-4"
  >
    <h1 class="text-4xl md:text-6xl italic font-bold font-serif text-center text-gray-300 mb-8">
      Ech0 Panel
    </h1>

    <!-- 移动端选择器 -->
    <div class="md:hidden mb-6 px-2">
      <div class="w-1/2">
        <BaseSelect
          v-model="selectedRoute"
          :options="routeOptions"
          placeholder="选择页面"
          @change="handleRouteChange"
        />
      </div>
      <div class="flex gap-2 mt-3">
        <button
          @click="router.push('/')"
          class="flex-1 px-4 py-2 rounded-md transition-all duration-300 border-none !shadow-none !ring-0 text-gray-600 hover:opacity-75 bg-transparent"
        >
          返回首页
        </button>
        <button
          v-if="userStore.isLogin"
          @click="handleLogout"
          class="flex-1 px-4 py-2 rounded-md transition-all duration-300 border-none !shadow-none !ring-0 text-gray-600 hover:opacity-75 bg-transparent"
        >
          退出登录
        </button>
        <button
          v-else
          @click="router.push('/auth')"
          class="flex-1 px-4 py-2 rounded-md transition-all duration-300 border-none !shadow-none !ring-0 text-gray-600 hover:opacity-75 bg-transparent"
        >
          登录 / 注册
        </button>
      </div>
    </div>

    <!-- 主内容区 -->
    <div class="mx-auto flex my-4 w-full max-w-screen-lg">
      <!-- 桌面端侧边栏 -->
      <div class="hidden md:flex flex-col gap-2 w-48 pr-8 shrink-0">
        <!-- 返回首页 -->
        <BaseButton @click="router.push('/')" :class="getButtonClasses('', true)" title="返回首页">
          <Arrow
            class="w-9 h-9 rotate-180 transition-transform duration-200 group-hover:-translate-x-1"
          />
        </BaseButton>

        <div class="h-px bg-gray-300 mx-2" />

        <!-- 仪表盘 -->
        <BaseButton
          :icon="Dashboard"
          @click="router.push('/panel/dashboard')"
          :class="getButtonClasses('panel-dashboard')"
          title="仪表盘"
        >
          仪表盘
        </BaseButton>

        <!-- 状态 -->
        <!-- <BaseButton
          :icon="Status"
          @click="router.push('/panel/status')"
          :class="getButtonClasses('panel-status')"
          title="状态"
        >
          状态
        </BaseButton> -->

        <!-- 设置 -->
        <BaseButton
          :icon="Setting"
          @click="router.push('/panel/setting')"
          :class="getButtonClasses('panel-setting')"
          title="系统"
        >
          系统
        </BaseButton>

        <!-- 个人中心 -->
        <BaseButton
          :icon="User"
          @click="router.push('/panel/user')"
          :class="getButtonClasses('panel-user')"
          title="成员"
        >
          成员
        </BaseButton>

        <!-- 存储 -->
        <BaseButton
          :icon="Storage"
          @click="router.push('/panel/storage')"
          :class="getButtonClasses('panel-storage')"
          title="存储"
        >
          存储
        </BaseButton>

        <!-- 单点登录 -->
        <BaseButton
          :icon="Sso"
          @click="router.push('/panel/sso')"
          :class="getButtonClasses('panel-sso')"
          title="单点登录"
        >
          单点登录
        </BaseButton>

        <!-- 高级 -->
        <BaseButton
          :icon="Others"
          @click="router.push('/panel/advance')"
          :class="getButtonClasses('panel-advance')"
          title="高级"
        >
          高级
        </BaseButton>

        <div class="h-px bg-gray-300 mx-2" />

        <!-- 退出登录 -->
        <BaseButton
          :icon="Logout"
          @click="handleLogout"
          :class="getBottomButtonClasses()"
          title="退出登录"
        >
          退出登入
        </BaseButton>

        <!-- 登录 / 注册 -->
        <BaseButton
          :icon="Auth"
          @click="router.push('/auth')"
          :class="getBottomButtonClasses()"
          title="登录 / 注册"
        >
          登录 / 注册
        </BaseButton>
      </div>

      <!-- 路由内容 -->
      <div class="flex-1 min-w-0">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import BaseButton from '@/components/common/BaseButton.vue'
import BaseSelect from '@/components/common/BaseSelect.vue'
import Arrow from '@/components/icons/arrow.vue'
import User from '@/components/icons/user.vue'
import Auth from '@/components/icons/auth.vue'
import Status from '@/components/icons/status.vue'
import Dashboard from '@/components/icons/dashboard.vue'
import Others from '@/components/icons/theothers.vue'
import Setting from '@/components/icons/setting.vue'
import Storage from '@/components/icons/storage.vue'
import Sso from '@/components/icons/sso.vue'
import Logout from '@/components/icons/logout.vue'
import { computed, ref, watch } from 'vue'
import { useUserStore } from '@/stores/user'
import { useRouter, useRoute } from 'vue-router'
import { theToast } from '@/utils/toast'
import { useBaseDialog } from '@/composables/useBaseDialog'

const { openConfirm } = useBaseDialog()

const userStore = useUserStore()
const router = useRouter()
const route = useRoute()

const currentRoute = computed(() => route.name as string)
const selectedRoute = ref(route.path)

// 统一的按钮样式计算函数
const getButtonClasses = (routeName: string, isBackButton = false) => {
  const baseClasses = isBackButton
    ? 'text-stone-600 rounded-md transition-all duration-300 border-none !shadow-none !ring-0 hover:opacity-75 p-2 group bg-transparent'
    : 'flex items-center gap-2 pl-3 py-1 rounded-md transition-all duration-300 border-none !shadow-none !ring-0 justify-start bg-transparent'

  const activeClasses =
    currentRoute.value === routeName
      ? 'text-stone-800 bg-orange-200'
      : 'text-stone-600 hover:opacity-75'

  return `${baseClasses} ${activeClasses}`
}

// 底部按钮样式
const getBottomButtonClasses = () => {
  return 'flex items-center gap-2 pl-3 py-1 rounded-md transition-all duration-300 border-none !shadow-none !ring-0 text-gray-600 hover:opacity-75 justify-start bg-transparent'
}

// 路由选项
const routeOptions = [
  { label: '仪表盘', value: '/panel/dashboard' },
  // { label: '状态', value: '/panel/status' },
  { label: '设置', value: '/panel/setting' },
  { label: '个人中心', value: '/panel/user' },
  { label: '存储', value: '/panel/storage' },
  { label: '单点登录', value: '/panel/sso' },
  { label: '高级', value: '/panel/advance' },
]

// 监听路由变化，更新选择器
watch(
  () => route.path,
  (newPath) => {
    selectedRoute.value = newPath
  },
)

// 处理选择器变化
const handleRouteChange = () => {
  router.push(selectedRoute.value)
}

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
