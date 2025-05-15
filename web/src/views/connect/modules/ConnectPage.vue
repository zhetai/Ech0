<template>
  <div class="px-4 pb-4 py-2 mt-4 mb-10 mx-auto flex justify-center items-center h-screen">
    <div class="h-3/5 sm:h-1/2 max-w-sm sm:max-w-md md:max-w-lg px-2 sm:px-4 md:px-6 my-4 sm:my-5 md:my-6">
      <h1 class="text-5xl italic font-bold text-center text-gray-300 mb-5">Ech0s Connect</h1>
      <div class="mb-5">
        <!-- 返回首页 -->
        <div class="mb-2">
          <BaseButton
            @click="$router.push('/')"
            class="text-gray-600 rounded-md !shadow-none !border-none !ring-0 !bg-transparent group"
            title="返回首页"
          >
            <Arrow
              class="w-9 h-9 rotate-180 transition-transform duration-200 group-hover:-translate-x-1"
            />
          </BaseButton>
        </div>

        <div class="w-full px-2">
          <!-- 列出所有连接（列出每个连接的头像） -->
          <div class="rounded-md shadow-sm ring-1 ring-gray-200 ring-inset bg-white p-4">
            <h2 class="text-gray-600 font-bold text-lg mb-2">我连接的Ech0s:</h2>
            <div>
              <div v-if="!connectsInfo.length" class="text-gray-500 text-sm mb-2">
                当前暂无连接
              </div>
              <div v-else class="flex flex-wrap gap-4">
                <div
                  v-for="(connect, index) in connectsInfo"
                  :key="index"
                  class="relative flex flex-col items-center justify-center w-8 h-8 border-2 border-gray-200 shadow-sm rounded-full hover:shadow-md transition duration-200 ease-in-out group"
                >
                  <a :href="connect.server_url" target="_blank">
                    <img
                      :src="connect.logo"
                      alt="avatar"
                      class="w-8 h-8 rounded-full object-cover"
                    />
                    <!-- 小绿点 -->
                    <span
                      class="absolute top-0 right-0 w-2.5 h-2.5 bg-green-500 border-2 border-white rounded-full"
                      style="transform: translate(35%, -35%)"
                    ></span>
                  </a>
                  <!-- Tooltip -->
                  <div
                    class="absolute z-10 left-1/2 -translate-x-1/2 top-10 min-w-max bg-gray-800 text-white text-xs rounded px-3 py-2 opacity-0 group-hover:opacity-100 pointer-events-none transition-opacity duration-200 shadow-lg"
                  >
                    <div class="font-bold mb-1">{{ connect.server_name }}</div>
                    <div v-if="connect.sys_username">管理员: {{ connect.sys_username }}</div>
                    <div v-if="connect.ech0s">共有: {{ connect.ech0s }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import BaseButton from '@/components/common/BaseButton.vue'
import Arrow from '@/components/icons/arrow.vue'
import { useConnectStore } from '@/stores/connect'
import { storeToRefs } from 'pinia'
import { onMounted } from 'vue'

const connectStore = useConnectStore()
const { getConnectInfo } = connectStore
const { connectsInfo } = storeToRefs(connectStore)

onMounted(async () => {
  await getConnectInfo()
})
</script>
