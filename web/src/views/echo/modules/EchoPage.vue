<template>
  <div class="flex justify-center items-center">
    <div v-if="echo">
      <!-- TODO: 该页面暂时处于开发中，后续会展示Echo的详细信息 -->
      <h1>Echo: {{ echoId }}</h1>
      <p>该页面暂时处于开发中，后续会展示Echo的详细信息</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { fetchGetEchoById } from '@/service/api'
import { ref } from 'vue'

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
