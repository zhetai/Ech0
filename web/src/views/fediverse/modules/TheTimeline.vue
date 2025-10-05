<template>
  <div class="text-left">
    <p v-if="loading" class="text-sm text-gray-400">正在同步时间线…</p>
    <p v-else-if="error" class="text-sm text-red-400">时间线获取失败</p>
    <div v-else-if="items.length" class="space-y-4">
      <article
        v-for="item in items"
        :key="item.id"
        class="rounded-lg border border-gray-700/50 bg-gray-900/40 p-4 shadow-sm"
      >
        <div class="mb-2 flex items-center justify-between text-xs text-gray-400">
          <span class="font-medium text-gray-300">
            {{ item.actorDisplayName || item.actorPreferredUsername || item.actorId }}
          </span>
          <span>{{ formatTimelineTime(item.publishedAt || item.createdAt) }}</span>
        </div>
        <p class="whitespace-pre-line text-sm text-gray-200">
          {{ getTimelineContent(item) }}
        </p>
      </article>
    </div>
    <p v-else class="text-sm text-gray-500">关注一些联邦好友，时间线才会热闹起来～</p>
  </div>
</template>

<script setup lang="ts">
import { toRefs } from 'vue'

const props = defineProps<{
  loading: boolean
  error: string
  items: App.Api.Fediverse.TimelineItem[]
}>()

const { loading, error, items } = toRefs(props)

const stripHtml = (value: string) =>
  value
    ? value
        .replace(/<br\s*\/?>(?=\s|$)/gi, '\n')
        .replace(/<\/p>/gi, '\n')
        .replace(/<[^>]*>/g, '')
        .replace(/\n{3,}/g, '\n\n')
        .trim()
    : ''

const getTimelineContent = (item: App.Api.Fediverse.TimelineItem) => {
  if (item.summary && item.summary.trim()) {
    return item.summary.trim()
  }
  if (item.content && item.content.trim()) {
    const sanitized = stripHtml(item.content)
    return sanitized || item.content.trim()
  }
  return ''
}

const formatTimelineTime = (input: string | undefined) => {
  if (!input) {
    return ''
  }

  const date = new Date(input)
  if (Number.isNaN(date.getTime())) {
    return ''
  }

  return date.toLocaleString()
}

</script>

<style scoped>

</style>
