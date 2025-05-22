<template>
  <div
    class="max-w-sm flex justify-center items-center bg-white rounded-lg shadow-sm ring-1 ring-inset ring-gray-100 p-2 gap-2"
  >
    <a
      :href="websiteURL"
      target="_blank"
      rel="noopener noreferrer"
      class="flex items-center gap-2 hover:underline hover:decoration-orange-500 hover:decoration-wavy"
    >
      <Website class="w-6 h-6" />
      <span class="text-gray-400"> 外 部 链 接 </span>
      <Website class="w-6 h-6" />
    </a>
  </div>
</template>

<script setup lang="ts">
import Website from '../icons/website.vue'

defineProps<{
  websiteURL: string
}>()
</script>
