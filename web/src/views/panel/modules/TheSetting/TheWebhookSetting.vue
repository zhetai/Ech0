<template>
  <PanelCard>
    <!-- Webhook 设置 -->
    <div class="w-full">
      <div class="flex flex-row items-center justify-between mb-4">
        <h1 class="text-gray-600 font-bold text-lg">Webhook</h1>
        <div class="flex flex-row items-center justify-end gap-2 w-14">
          <button @click="webhookEdit = !webhookEdit" title="编辑">
            <Edit v-if="!webhookEdit" class="w-5 h-5 text-gray-400 hover:w-6 hover:h-6" />
            <Close v-else class="w-5 h-5 text-gray-400 hover:w-6 hover:h-6" />
          </button>
        </div>
      </div>

      <!-- 添加 Webhook -->
      <div
        v-if="webhookEdit"
        class="mb-2 border border-gray-300 border-dashed rounded-md flex flex-col gap-2 p-2 text-gray-400"
      >
        <div>
          <span>Webhook 名称：</span>
          <BaseInput class="w-full" v-model="webhookToAdd.name" placeholder="Webhook 名称" />
        </div>

        <div>
          <span>Webhook 地址：</span>
          <BaseInput
            class="w-full"
            v-model="webhookToAdd.url"
            placeholder="Webhook 地址（带https/http）"
          />
        </div>

        <div class="flex items-center justify-center my-2">
          <BaseButton
            :disabled="isSubmitting"
            @click="handleCancelAddWebhook"
            class="w-1/3 h-8 rounded-md flex justify-center mr-2"
            title="取消添加"
          >
            <span>取消</span>
          </BaseButton>

          <BaseButton
            :loading="isSubmitting"
            @click="handleAddWebhook"
            class="w-1/3 h-8 rounded-md flex justify-center"
            title="添加 Webhook"
          >
            <span class="text-gray-600">添加</span>
          </BaseButton>
        </div>
      </div>

      <!-- Webhook 列表 -->
      <div
        v-if="Webhooks.length === 0 && !webhookEdit"
        class="flex flex-col items-center justify-center mt-2"
      >
        <span class="text-gray-400">暂无 Webhook...</span>
      </div>
      <div
        v-else-if="Webhooks.length !== 0 && !webhookEdit"
        class="mt-2 border border-dashed border-gray-300 rounded-md p-2 flex flex-col gap-2 max-h-60 overflow-y-auto"
      >
        <div
          v-for="(webhook, index) in Webhooks"
          :key="index"
          class="flex w-full flex-row items-center justify-between text-gray-500 gap-3 h-10 border-b border-gray-200 last:border-0 flex-shrink-0"
        >
          <div class="w-60 md:w-full flex-nowrap flex items-start gap-2">
            <span class="w-26 md:w-32 font-bold text-nowrap overflow-x-auto">{{
              webhook.name
            }}</span>
            <span class="text-gray-700 font-mono font-bold"> | </span>
            <span class="w-32 md:w-72 overflow-x-auto" :title="webhook.url">
              {{ webhook.url }}
            </span>
          </div>
          <BaseButton
            :icon="Trashbin"
            @click="handleDeleteWebhook(webhook.id)"
            class="w-7 h-7 rounded-md"
            title="删除 Webhook"
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
import Trashbin from '@/components/icons/trashbin.vue'
import Close from '@/components/icons/close.vue'
import { ref, onMounted } from 'vue'
import { useSettingStore } from '@/stores/setting'
import { storeToRefs } from 'pinia'
import { fetchDeleteWebhook, fetchCreateWebhook } from '@/service/api'
import { useBaseDialog } from '@/composables/useBaseDialog'
import { theToast } from '@/utils/toast'

const webhookEdit = ref<boolean>(false)

const settingStore = useSettingStore()
const { Webhooks } = storeToRefs(settingStore)
const { openConfirm } = useBaseDialog()

const webhookToAdd = ref<App.Api.Setting.WebhookDto>({
  name: '',
  url: '',
  is_active: true,
})
const isSubmitting = ref<boolean>(false)

const handleAddWebhook = () => {
  if (isSubmitting.value) return
  isSubmitting.value = true

  if (!webhookToAdd.value?.name || !webhookToAdd.value?.url) {
    theToast.error('请填写完整的 Webhook 信息')
    isSubmitting.value = false
    return
  }

  fetchCreateWebhook(webhookToAdd.value)
    .then((res) => {
      if (res.code === 1) {
        theToast.success('添加 Webhook 成功')
        webhookToAdd.value = { name: '', url: '', is_active: true }
        settingStore.getAllWebhooks()
      }

      isSubmitting.value = false

      handleCancelAddWebhook()
    })
    .finally(() => {
      isSubmitting.value = false
    })
}

const handleCancelAddWebhook = () => {
  webhookEdit.value = false
  webhookToAdd.value = { name: '', url: '', is_active: true }
}

const handleDeleteWebhook = (id: number) => {
  openConfirm({
    title: '确认删除此 Webhook 吗？',
    description: '删除后将无法恢复，请谨慎操作！',
    onConfirm: () => {
      fetchDeleteWebhook(id).then((res) => {
        if (res.code === 1) {
          theToast.success('删除 Webhook 成功')
          settingStore.getAllWebhooks()
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
