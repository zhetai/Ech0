<template>
  <div class="px-4 pb-4 py-2 mt-4 mb-10 mx-auto">
    <TheSearch />
    <TheEditor />
    <TheEchos />
  </div>
</template>

<script setup lang="ts">
import TheSearch from './TheSearch.vue';
import TheEditor from './TheEditor.vue'
import TheEchos from './TheEchos.vue'
import { onMounted } from 'vue'
import { useSettingStore } from '@/stores/settting';

const { getSystemSetting } = useSettingStore()

onMounted(async () => {
  await getSystemSetting()
})
</script>
