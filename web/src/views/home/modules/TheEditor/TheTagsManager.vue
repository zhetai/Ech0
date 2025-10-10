<template>
  <div>

  </div>
</template>

<script setup lang="ts">
import { useEchoStore } from '@/stores/echo';
import { fetchDeleteTagById } from '@/service/api';
import { onMounted } from 'vue';

const echoStore = useEchoStore();

const handleDeleteTag = (tagId: number) => {
  fetchDeleteTagById(tagId).then((res) => {
    if (res.code === 1) {
      echoStore.getTags();
    }
  })
}

</script>

<style scoped>

</style>
