<template>
  <div class="max-w-sm sm:max-w-full px-2 pb-4 py-2 mt-4 sm:mt-6 mb-10 mx-auto flex flex-col sm:flex-row justify-center items-start sm:gap-8">
    <div class="sm:max-w-sm w-full ">
      <TheTop class="sm:hidden" />
      <TheEditor />
    </div>
    <div class="sm:max-w-lg w-full sm:mt-1">
      <TheTop class="hidden sm:block sm:px-4" />
      <TheEchos v-if="!todoMode" />
      <TheTodos v-else />
    </div>
    <div class="hidden xl:block sm:max-w-sm w-full px-6">
      <TheHeatMap class="mb-2" />
      <TheConnects />
    </div>
  </div>
</template>

<script setup lang="ts">
import TheTop from './TheTop.vue'
import TheEditor from './TheEditor.vue'
import TheEchos from './TheEchos.vue'
import TheTodos from './TheTodos.vue'
import TheConnects from '@/views/connect/modules/TheConnects.vue'
import TheHeatMap from '@/components/advanced/TheHeatMap.vue'
import { onMounted } from 'vue'
import { useSettingStore } from '@/stores/settting'
import { useTodoStore } from '@/stores/todo'
import { storeToRefs } from 'pinia'

const { getSystemSetting } = useSettingStore()
const todoStore = useTodoStore()
const { todoMode } = storeToRefs(todoStore)

onMounted(async () => {

  await getSystemSetting()
})
</script>
