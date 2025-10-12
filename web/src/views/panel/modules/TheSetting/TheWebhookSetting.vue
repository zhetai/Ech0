<template>
  <div class="rounded-md shadow-sm ring-1 ring-gray-200 ring-inset bg-white p-4 mb-3">
    <!-- Webhook 设置 -->
    <div class="w-full">
      <div class="flex flex-row items-center justify-between mb-3">
        <h1 class="text-gray-600 font-bold text-lg">Webhook</h1>
        <div class="flex flex-row items-center justify-end gap-2 w-14">
          <button @click="webhookEdit = !webhookEdit" title="编辑">
            <Edit v-if="!webhookEdit" class="w-5 h-5 text-gray-400 hover:w-6 hover:h-6" />
            <Close v-else class="w-5 h-5 text-gray-400 hover:w-6 hover:h-6" />
          </button>
        </div>
      </div>

      <!-- 添加 Webhook -->
      <div v-if="webhookEdit" class="mb-2 border border-gray-300 border-dashed rounded-md"></div>

      <!-- Webhook 列表 -->
      <div
        v-if="Webhooks.length === 0 && !webhookEdit"
        class="flex flex-col items-center justify-center mt-2"
      >
        <span class="text-gray-400">暂无 Webhook...</span>
      </div>
      <div v-else class="mt-2">
        <div
          v-for="(webhook, index) in Webhooks"
          :key="index"
          class="flex flex-row items-center justify-between text-gray-500 gap-3 h-10 border-b border-gray-200 last:border-0"
        >
          <div class="flex items-center gap-2 flex-1 min-w-0">
            <h2 class="font-semibold w-auto flex-shrink-0 w-30 truncate">{{ webhook.name }}</h2>
            <span class="truncate max-w-full" :title="webhook.url" style="display: inline-block">
              {{ webhook.url }}
            </span>
          </div>
          <BaseButton
            :icon="Trashbin"
            :disabled="!webhookEdit"
            @click="handleDeleteWebhook(webhook.id)"
            class="w-7 h-7 rounded-md flex-shrink-0"
            title="删除 Webhook"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import BaseInput from '@/components/common/BaseInput.vue'
import BaseButton from '@/components/common/BaseButton.vue'
import Edit from '@/components/icons/edit.vue'
import Trashbin from '@/components/icons/trashbin.vue'
import Close from '@/components/icons/close.vue'
import Publish from '@/components/icons/publish.vue'
import { ref, onMounted } from 'vue'
import { useSettingStore } from '@/stores/setting'
import { storeToRefs } from 'pinia'
import { fetchDeleteWebhook } from '@/service/api'
import { useBaseDialog } from '@/composables/useBaseDialog'
import { theToast } from '@/utils/toast'

const webhookEdit = ref<boolean>(false)

const settingStore = useSettingStore()
const { Webhooks } = storeToRefs(settingStore)
const { openConfirm } = useBaseDialog()

const webhookToAdd = ref<App.Api.Setting.WebhookDto>()

const handleDeleteWebhook = (id: number) => {
  openConfirm({
    title: '确认删除此 Webhook 吗？',
    description: '删除后将无法恢复，请谨慎操作！',
    onConfirm: () => {
      fetchDeleteWebhook(id).then((res) => {
        if (res.code === 1) {
          theToast.success('删除 Webhook 成功')
        }
      })
    },
  })
}

onMounted(async () => {
  await settingStore.getAllWebhooks()
})
</script>

<style scoped></style>
