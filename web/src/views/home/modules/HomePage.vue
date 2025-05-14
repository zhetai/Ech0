<template>
  <div class="px-4 pb-4 py-2 mt-4 mb-10 mx-auto">
    <TheTop />
    <TheEditor />
    <TheEchos v-if="!todoMode" />
    <TheTodos v-else />
  </div>
</template>

<script setup lang="ts">
import TheTop from './TheTop.vue';
import TheEditor from './TheEditor.vue'
import TheEchos from './TheEchos.vue'
import TheTodos from './TheTodos.vue';
import { onMounted } from 'vue'
import { useSettingStore } from '@/stores/settting';
import { useTodoStore } from '@/stores/todo';
import { storeToRefs } from 'pinia';

const { getSystemSetting } = useSettingStore()
const todoStore = useTodoStore()
const { todoMode } = storeToRefs(todoStore)

onMounted(async () => {
  await getSystemSetting()
})
</script>
