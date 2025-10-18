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

    <div class="my-2 text-gray-400">访问令牌配置开发中……</div>
  </PanelCard>
</template>

<script setup lang="ts">
import PanelCard from '@/layout/PanelCard.vue'
// import BaseInput from '@/components/common/BaseInput.vue'
// import BaseButton from '@/components/common/BaseButton.vue'
import Edit from '@/components/icons/edit.vue'
// import Trashbin from '@/components/icons/trashbin.vue'
import Close from '@/components/icons/close.vue'
import { ref, onMounted } from 'vue'
import { useSettingStore } from '@/stores/setting'
// import { storeToRefs } from 'pinia'
// import { fetchDeleteWebhook, fetchCreateWebhook } from '@/service/api'
// import { useBaseDialog } from '@/composables/useBaseDialog'
// import { theToast } from '@/utils/toast'

const accessTokenEdit = ref<boolean>(false)
const useSetting = useSettingStore()
// const { AccessTokens } = storeToRefs(useSetting)

onMounted(async () => {
  await useSetting.getAllAccessTokens()
})
</script>

<style scoped></style>
