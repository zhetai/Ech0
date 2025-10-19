<template>
  <PanelCard>
    <!-- Webhook 设置 -->
    <div class="w-full">
      <div class="flex flex-row items-center justify-between mb-4">
        <h1 class="text-gray-600 font-bold text-lg">访问令牌</h1>
        <div class="flex flex-row items-center justify-end gap-2 w-14">
          <button @click="accessTokenEdit = !accessTokenEdit" title="编辑">
            <Edit v-if="!accessTokenEdit" class="w-5 h-5 text-gray-400 hover:w-6 hover:h-6" />
            <Close v-else class="w-5 h-5 text-gray-400 hover:w-6 hover:h-6" />
          </button>
        </div>
      </div>
    </div>

    <div v-if="!accessTokenEdit">
      <div class="mt-2 overflow-x-auto border border-stone-300 rounded-lg">
        <table class="min-w-full divide-y divide-gray-200">
          <thead>
            <tr class="bg-stone-50 opacity-70">
              <th class="px-3 py-2 text-left text-sm font-semibold text-stone-600">Token</th>
              <th class="px-3 py-2 text-left text-sm font-semibold text-stone-600">名称</th>
              <th class="px-3 py-2 text-left text-sm font-semibold text-stone-600">创建时间</th>
              <th class="px-3 py-2 text-left text-sm font-semibold text-stone-600">过期时间</th>
              <th class="px-3 py-2 text-right text-sm font-semibold text-stone-600">操作</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100 text-nowrap">
            <tr v-for="t in AccessTokens" :key="t.id">
              <td class="px-3 py-2 flex items-center gap-x-1 font-mono text-sm text-stone-700">
                {{ formatToken(t.token) }}
                <button class="p-1 hover:bg-gray-100 rounded" @click="handleCopyToken(t.token)">
                  <Clipboard class="w-5 h-5 text-gray-500" />
                </button>
              </td>
              <td class="px-3 py-2 text-sm text-stone-700">
                <span :title="t.name" class="truncate block max-w-xs">{{ t.name }}</span>
              </td>
              <td class="px-3 py-2 text-sm text-stone-500">
                {{ new Date(t.created_at).toLocaleString() }}
              </td>
              <td class="px-3 py-2 text-sm text-stone-500">
                {{ t.expiry ? new Date(t.expiry).toLocaleString() : '永不过期' }}
              </td>
              <td class="px-3 py-2 text-right">
                <button
                  class="p-1 hover:bg-gray-100 rounded"
                  @click="handleDeleteAccessToken(t.id)"
                  title="删除 Token"
                >
                  <Trashbin class="w-5 h-5 text-red-500" />
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
    <div v-else class="text-stone-500">
      <!-- 添加 AccessToken -->

      <div>
        <span>Token 名称：</span>
        <BaseInput class="w-full" v-model="accessTokenToAdd.name" placeholder="Token 名称" />
      </div>

      <div>
        <span>过期时间</span>
        <BaseSelect
          v-model="accessTokenToAdd.expiry"
          :options="ExpirationOptions"
          class="w-34 h-8"
        />
      </div>

      <div class="flex items-center justify-center my-2">
        <BaseButton
          :disabled="isSubmitting"
          @click="handleCancelAddAccessToken"
          class="w-1/3 h-8 rounded-md flex justify-center mr-2"
          title="取消添加"
        >
          <span>取消</span>
        </BaseButton>

        <BaseButton
          :loading="isSubmitting"
          @click="handleAddAccessToken"
          class="w-1/3 h-8 rounded-md flex justify-center"
          title="添加 Access Token"
        >
          <span class="text-gray-600">添加</span>
        </BaseButton>
      </div>
    </div>
  </PanelCard>
</template>

<script setup lang="ts">
import PanelCard from '@/layout/PanelCard.vue'
import BaseInput from '@/components/common/BaseInput.vue'
import BaseButton from '@/components/common/BaseButton.vue'
import BaseSelect from '@/components/common/BaseSelect.vue'
import Edit from '@/components/icons/edit.vue'
import Trashbin from '@/components/icons/trashbin.vue'
import Close from '@/components/icons/close.vue'
import { ref, onMounted } from 'vue'
import { useSettingStore } from '@/stores/setting'
import { storeToRefs } from 'pinia'
import { fetchCreateAccessToken, fetchDeleteAccessToken } from '@/service/api'
import { useBaseDialog } from '@/composables/useBaseDialog'
import { theToast } from '@/utils/toast'
import Clipboard from '@/components/icons/clipboard.vue'
import { AccessTokenExpiration } from '@/enums/enums'

const { openConfirm } = useBaseDialog()

const accessTokenEdit = ref<boolean>(false)
const useSetting = useSettingStore()
const { AccessTokens } = storeToRefs(useSetting)

const deleteTarget = ref<App.Api.Setting.AccessToken | null>(null)
const accessTokenToAdd = ref<App.Api.Setting.AccessTokenDto>({
  name: '',
  expiry: AccessTokenExpiration.EIGHT_HOUR_EXPIRY,
})
const ExpirationOptions = [
  { label: '8 Hours', value: AccessTokenExpiration.EIGHT_HOUR_EXPIRY },
  { label: '1 Month', value: AccessTokenExpiration.ONE_MONTH_EXPIRY },
  { label: 'Never', value: AccessTokenExpiration.NEVER_EXPIRY },
]

const isSubmitting = ref<boolean>(false)
const handleAddAccessToken = async () => {
  if (!accessTokenToAdd.value?.name) {
    theToast.error('请填写 Token 名称')
    return
  }

  isSubmitting.value = true

  const res = await fetchCreateAccessToken({
    name: accessTokenToAdd.value.name,
    expiry: accessTokenToAdd.value.expiry || AccessTokenExpiration.NEVER_EXPIRY,
  })
  if (res.code === 1) {
    theToast.success('Access Token 创建成功')
    accessTokenToAdd.value = {
      name: '',
      expiry: AccessTokenExpiration.EIGHT_HOUR_EXPIRY,
    }
    await useSetting.getAllAccessTokens()
    accessTokenEdit.value = false
  }
  isSubmitting.value = false
}

const handleCancelAddAccessToken = () => {
  accessTokenToAdd.value = { name: '', expiry: AccessTokenExpiration.EIGHT_HOUR_EXPIRY }
  accessTokenEdit.value = false
}

// 删除 Access Token
const handleDeleteAccessToken = async (id: number) => {
  if (!deleteTarget.value) return

  openConfirm({
    title: '确认删除此 Access Token 吗？',
    description: `Token 名称：${deleteTarget.value.name}`,
    onConfirm: async () => {
      const res = await fetchDeleteAccessToken(id)
      if (res.code === 1) {
        theToast.success('Access Token 删除成功')
        await useSetting.getAllAccessTokens()
        deleteTarget.value = null
      }
    },
  })
}

onMounted(async () => {
  await useSetting.getAllAccessTokens()
})

// 复制 Token
const handleCopyToken = (token: string) => {
  navigator.clipboard.writeText(token)
  theToast.success('Access Token 已复制到剪贴板')
}

// 显示格式化的 Token
const formatToken = (token: string) => `${token.slice(0, 4)}****${token.slice(-4)}`
</script>

<style scoped></style>
