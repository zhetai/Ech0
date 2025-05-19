<template>
  <div class="audio-card">
    <audio :src="url" controls />
  </div>
</template>

<script setup lang="ts">
import { defineProps } from 'vue'

const props = defineProps<{
  url: string
}>()
</script>

<style scoped>
.audio-card {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 12px;
}
audio {
  width: 100%;
  max-width: 400px;
}
</style>
