<template>
  <!-- todoMode -->
  <BaseTextArea
    v-if="currentMode === Mode.TODO"
    v-model="content"
    class="rounded-lg h-auto sm:min-h-[6rem] md:min-h-[9rem]"
    placeholder="请输入待办事项..."
    :rows="3"
  />
</template>

<script setup lang="ts">
import { Mode } from '@/enums/enums'
import BaseTextArea from '@/components/common/BaseTextArea.vue'

defineProps<{
  currentMode: Mode
}>()

// 使用 defineModel 创建双向绑定
const content = defineModel<string>('content', {
  required: true,
  default: '',
})
</script>

<style scoped></style>
