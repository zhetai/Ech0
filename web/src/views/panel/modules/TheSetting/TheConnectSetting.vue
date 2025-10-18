<template>
  <PanelCard>
    <!-- Ech0 Connect设置 -->
    <div class="w-full">
      <div class="flex flex-row items-center justify-between mb-3">
        <h1 class="text-gray-600 font-bold text-lg">Ech0 Connect</h1>
        <div class="flex flex-row items-center justify-end gap-2 w-14">
          <button @click="connectsEdit = !connectsEdit" title="编辑">
            <Edit v-if="!connectsEdit" class="w-5 h-5 text-gray-400 hover:w-6 hover:h-6" />
            <Close v-else class="w-5 h-5 text-gray-400 hover:w-6 hover:h-6" />
          </button>
        </div>
      </div>
      <!-- 添加 Connect -->
      <div v-if="connectsEdit" class="flex flex-row items-center justify-between h-10">
        <BaseInput
          v-model="connectUrl"
          type="text"
          placeholder="请输入Connect地址（带https/http）"
          class="w-full h-7"
        />
        <BaseButton
          :icon="Publish"
          @click="handleAddConnect"
          class="w-7 h-7 ml-2 rounded-md"
          title="连接"
        />
      </div>
      <!-- Connect列表 -->
      <div
        v-if="connects.length === 0 && !connectsEdit"
        class="flex flex-col items-center justify-center mt-2 border border-dashed border-gray-300 rounded-md p-1"
      >
        <span class="text-gray-400">暂无连接...</span>
      </div>
      <div
        v-else
        class="mt-2 border border-dashed border-gray-300 rounded-md p-2 flex flex-col overflow-hidden gap-2 max-h-60 overflow-y-auto"
      >
        <div
          v-for="(connect, index) in connects"
          :key="index"
          class="flex flex-row items-center justify-between text-gray-500 gap-3 h-10 border-b border-gray-200 last:border-0 flex-shrink-0"
        >
          <div class="flex items-center gap-2 flex-1 min-w-0">
            <h2 class="font-semibold w-auto flex-shrink-0">Connect {{ index + 1 }}:</h2>
            <span
              class="overflow-x-auto max-w-full"
              :title="connect.connect_url"
              style="display: inline-block"
              >{{ connect.connect_url }}</span
            >
          </div>
          <BaseButton
            :icon="Disconnect"
            :disabled="!connectsEdit"
            @click="handleDisconnect(connect.id)"
            class="w-7 h-7 rounded-md flex-shrink-0"
            title="断开连接"
          />
        </div>
      </div>
    </div>
  </PanelCard>
</template>

<script setup lang="ts">
import PanelCard from '@/layout/PanelCard.vue'
import BaseInput from '@/components/common/BaseInput.vue'
import BaseButton from '@/components/common/BaseButton.vue'
import Edit from '@/components/icons/edit.vue'
import Disconnect from '@/components/icons/disconnect.vue'
import Close from '@/components/icons/close.vue'
import Publish from '@/components/icons/publish.vue'
import { ref, onMounted } from 'vue'
import { fetchAddConnect, fetchDeleteConnect } from '@/service/api'
import { theToast } from '@/utils/toast'

import { useConnectStore } from '@/stores/connect'
import { storeToRefs } from 'pinia'

import { useBaseDialog } from '@/composables/useBaseDialog'
const { openConfirm } = useBaseDialog()

const connectStore = useConnectStore()
const { getConnect } = connectStore
const { connects } = storeToRefs(connectStore)
const connectsEdit = ref<boolean>(false)
const connectUrl = ref<string>('')

const handleAddConnect = async () => {
  if (connectUrl.value.length === 0) {
    theToast.error('请输入Connect地址')
    return
  }
  await fetchAddConnect(connectUrl.value).then((res) => {
    if (res.code === 1) {
      theToast.success(res.msg)
      connectUrl.value = ''
      getConnect()
    }
  })
}

const handleDisconnect = async (connect_id: number) => {
  // 弹出确认框
  openConfirm({
    title: '确定要断开连接吗？',
    description: '',
    onConfirm: async () => {
      await fetchDeleteConnect(connect_id).then((res) => {
        if (res.code === 1) {
          theToast.success(res.msg)
          getConnect()
        }
      })
    },
  })
}

onMounted(() => {
  getConnect()
})
</script>

<style scoped></style>
