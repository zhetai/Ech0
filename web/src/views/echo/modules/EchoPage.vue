<template>
  <div class="px-4 pb-4 py-2 mt-4 mb-10 mx-auto flex justify-center items-center">
    <div class="h-auto max-w-sm sm:max-w-md md:max-w-lg">
      <div v-if="echo">
        <!-- TODO: 该页面暂时处于开发中，后续会展示Echo的详细信息 -->
        <TheEchoDetail :echo="echo" />
      </div>
      <div v-else class="text-gray-500">当前暂无Echo详情可展示</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { fetchGetEchoById } from '@/service/api'
import { ref } from 'vue'
import TheEchoDetail from '@/components/advanced/TheEchoDetail.vue'

const route = useRoute()
const echoId = route.params.echoId as string

const echo = ref<App.Api.Ech0.Echo | null>(null)

onMounted(async () => {
  // 在这里可以添加获取Echo详情的逻辑
  await fetchGetEchoById(echoId).then((res) => {
    if (res.code === 1) {
      echo.value = res.data
    }
  })
})
</script>
