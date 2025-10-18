<template>
  <div class="w-full px-2">The Dashboard</div>
  <pre>{{ dashboardMetrics }}</pre>
</template>

<script setup lang="ts">
import { onMounted, watch, ref } from 'vue'
import { useOWebSocket } from '@/service/request/websocket'

const dashboardMetrics = ref<App.Api.Dashboard.Metrics | null>(null)

const { status, sendMessage, onMessage, open, ws } = useOWebSocket<
  App.Api.Response<App.Api.Dashboard.Metrics>
>({
  url: 'ws://localhost:6277/ws/dashboard/metrics',
  autoReconnect: true,
  heartbeat: true,
})

onMounted(() => {
  open() // 等组件挂载后再连接

  // 每次收到服务器推送的 metrics 更新 metrics
  onMessage((payload) => {
    if (payload.code === 1 && payload.data) {
      dashboardMetrics.value = payload.data
    }
  })
})
</script>

<style scoped></style>
