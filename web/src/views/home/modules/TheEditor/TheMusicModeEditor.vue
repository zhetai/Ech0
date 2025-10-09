<template>
  <div>
    <h2 class="text-gray-500 font-bold mb-1">欢迎使用音乐播放模式（仅PC）</h2>
    <div class="mb-1 flex items-center gap-2">
      <p class="text-gray-500">上传音乐：</p>
      <input
        id="file-input"
        class="hidden"
        type="file"
        accept="audio/*"
        ref="fileInput"
        @change="handleUploadMusic"
      />
      <BaseButton
        :icon="Audio"
        @click="handleTriggerUpload"
        class="w-7 h-7 sm:w-7 sm:h-7 rounded-md"
        title="上传音乐"
      />
    </div>
    <div class="flex items-center gap-2">
      <p class="text-gray-500">删除音乐：</p>
      <BaseButton
        :icon="Delete"
        @click="handleDeleteMusic"
        class="w-7 h-7 sm:w-7 sm:h-7 rounded-md"
        title="删除音乐"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import Audio from '@/components/icons/audio.vue'
import Delete from '@/components/icons/delete.vue'
import BaseButton from '@/components/common/BaseButton.vue'
import { ref } from 'vue'
import { fetchUploadMusic, fetchDeleteMusic } from '@/service/api'
import { theToast } from '@/utils/toast'
import { useBaseDialog } from '@/composables/useBaseDialog'

const { openConfirm } = useBaseDialog()

const emit = defineEmits(['refreshAudio'])

const fileInput = ref<HTMLInputElement | null>(null)
const handleTriggerUpload = () => {
  if (fileInput.value) {
    fileInput.value.click()
  }
}

const handleUploadMusic = async (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return

  try {
    const res = await theToast.promise(fetchUploadMusic(file), {
      loading: '音乐上传中...',
      success: '音乐上传成功！',
      error: '音乐上传失败，请稍后再试',
    })

    if (res.code === 1) {
      emit('refreshAudio')
    }
  } catch (err) {
    console.error('音乐上传异常:', err)
  } finally {
    target.value = ''
  }
}

const handleDeleteMusic = () => {
  openConfirm({
    title: '确定要删除音乐吗？',
    description: '删除后将无法恢复，请谨慎操作',
    onConfirm: () => {
      fetchDeleteMusic().then((res) => {
        if (res.code === 1) {
          theToast.success('音乐删除成功！')
          emit('refreshAudio')
        }
      })
    },
  })
}
</script>

<style scoped></style>
