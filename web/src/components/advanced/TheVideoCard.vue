<template>
  <div class="max-w-sm flex justify-center items-center bg-white rounded-lg shadow-sm p-3 gap-2">
    <iframe
      class="w-full aspect-video rounded"
      :src="`https://player.bilibili.com/player.html?bvid=${props.bvnumber}`"
      loading="lazy"
      allowfullscreen
    ></iframe>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  bvnumber: string
}>()
</script>

<style scoped></style>
