<template>
  <div class="w-full px-2">
    <div>
      <div class="md:flex items-center gap-2">
        <!-- CPU -->
        <MetricCard title="CPU 使用率" class="md:w-1/2">
          <VChart class="h-60" :option="cpuOption" autoresize />
        </MetricCard>

        <!-- Memory -->
        <MetricCard title="内存使用率" class="md:w-1/2">
          <VChart class="h-60" :option="memoryOption" autoresize />
        </MetricCard>
      </div>

      <div class="md:flex items-center mt-4 gap-2">
        <!-- Disk -->
        <MetricCard title="磁盘使用情况" class="md:w-1/2">
          <VChart class="h-60" :option="diskOption" autoresize />
        </MetricCard>

        <!-- System -->
        <MetricCard title="系统信息" class="md:w-1/2">
          <div class="text-md text-stone-500 font-bold h-60 p-2">
            <p>
              主机名：<span class="text-sm font-normal">{{ metrics.System?.Hostname }}</span>
            </p>
            <p>
              操作系统：<span class="text-sm font-normal">{{ metrics.System?.OsName }}</span>
            </p>
            <p>
              内核版本：<span class="text-sm font-normal">{{ metrics.System?.KernelVersion }}</span>
            </p>
            <p>
              运行时长：<span class="text-sm font-normal">{{ metrics.System?.Uptime }} s</span>
            </p>
            <p>
              当前时间：<span class="text-sm font-normal">{{ metrics.System?.Time }}</span>
            </p>
          </div>
        </MetricCard>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useOWebSocket } from '@/service/request/websocket'
import MetricCard from '@/layout/MetricCard.vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent,
} from 'echarts/components'
import { GaugeChart, LineChart, BarChart, PieChart } from 'echarts/charts'
import { CanvasRenderer } from 'echarts/renderers'

// 注册按需组件
use([
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent,
  GaugeChart,
  LineChart,
  BarChart,
  PieChart,
  CanvasRenderer,
])

const metrics = ref<App.Api.Dashboard.Metrics>({
  CPU: {
    UsagePercent: 0,
    Cores: 0,
    FrequencyMHz: 0,
  },
  Memory: {
    Total: 0,
    Used: 0,
    Available: 0,
    Percentage: 0,
  },
  Disk: {
    Total: 0,
    Used: 0,
    Available: 0,
    Percentage: 0,
  },
  Network: {
    BytesSentPerSecond: 0,
    BytesReceivedPerSecond: 0,
    TotalBytesReceived: 0,
    TotalBytesSent: 0,
  },
  System: {
    Uptime: 0,
    Hostname: '',
    OsName: '',
    KernelVersion: '',
    Time: '',
    KernelArch: '',
  },
})

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
      metrics.value = payload.data
      updateCharts() // 更新图表
    }
  })
})

// ---- 图表配置 ----
const cpuOption = ref({})
const memoryOption = ref({})
const diskOption = ref({})
const networkOption = ref({})

// 实时更新函数
function updateCharts() {
  cpuOption.value = {
    series: [
      {
        type: 'gauge',
        progress: { show: true, width: 10 },
        detail: false,
        data: [{ value: metrics.value.CPU.UsagePercent, name: 'CPU' }],
      },
    ],
  }

  memoryOption.value = {
    tooltip: { trigger: 'item' },
    legend: { bottom: '0%' },
    series: [
      {
        name: '内存',
        type: 'pie',
        radius: ['40%', '70%'],
        data: [
          { value: metrics.value.Memory.Used, name: '已用' },
          { value: metrics.value.Memory.Available, name: '可用' },
        ],
      },
    ],
  }

  diskOption.value = {
    grid: { left: '10%', right: '10%', bottom: '10%', top: '15%' },
    xAxis: { type: 'category', data: ['磁盘使用率'] },
    yAxis: { type: 'value', max: 100 },
    series: [
      {
        type: 'bar',
        data: [metrics.value.Disk.Percentage],
        label: { show: true, position: 'top', formatter: '{c}%' },
      },
    ],
  }

  networkOption.value = {
    tooltip: { trigger: 'axis' },
    legend: { data: ['上传(B/s)', '下载(B/s)'] },
    xAxis: { type: 'category', data: [metrics.value.System.Time] },
    yAxis: { type: 'value' },
    series: [
      {
        name: '上传(B/s)',
        type: 'line',
        data: [metrics.value.Network.BytesSentPerSecond],
      },
      {
        name: '下载(B/s)',
        type: 'line',
        data: [metrics.value.Network.BytesReceivedPerSecond],
      },
    ],
  }
}
</script>

<style scoped></style>
