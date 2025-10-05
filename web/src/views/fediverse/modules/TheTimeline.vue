<template>
  <div class="text-left">
    <p v-if="loading" class="text-sm text-gray-400">正在同步时间线…</p>
    <p v-else-if="error" class="text-sm text-red-400">时间线获取失败</p>
    <div v-else-if="resolvedItems.length" class="space-y-6">
      <article
        v-for="item in resolvedItems"
        :key="item.id"
        class="rounded-xl border border-gray-200/80 bg-white/95 p-5 shadow-md transition-shadow duration-200 hover:shadow-lg dark:border-gray-700/60 dark:bg-gray-900/70"
      >
        <header class="mb-4 flex items-start justify-between gap-4">
          <div class="flex flex-col gap-1">
            <span class="text-base font-semibold text-slate-900 dark:text-gray-100">
              {{ item.displayName }}
            </span>
            <span class="text-xs text-slate-500 dark:text-gray-400">@ {{ item.actorHandle }}</span>
          </div>
          <time class="text-xs text-slate-400 dark:text-gray-500">
            {{ item.timeText }}
          </time>
        </header>

        <div v-if="item.galleryImages.length" class="mb-4">
          <TheImageGallery :images="item.galleryImages" />
        </div>

        <section v-if="item.contentHtml" class="prose prose-sm max-w-none text-slate-700 dark:prose-invert dark:text-gray-200" v-html="item.contentHtml"></section>
        <p v-else-if="item.contentText" class="whitespace-pre-line text-sm text-slate-700 dark:text-gray-200">
          {{ item.contentText }}
        </p>

        <footer class="mt-4 flex flex-wrap items-center justify-between gap-2 text-xs text-slate-400 dark:text-gray-500">
          <span class="truncate" :title="item.objectId">{{ item.objectId }}</span>
          <span>来自联邦宇宙</span>
        </footer>
      </article>
    </div>
    <p v-else class="text-sm text-gray-500">关注一些联邦好友，时间线才会热闹起来～</p>
  </div>
</template>

<script setup lang="ts">
import { computed, toRefs } from 'vue'
import TheImageGallery from '@/components/advanced/TheImageGallery.vue'
import { ImageSource } from '@/enums/enums'

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

const extractRecord = (value: unknown): Record<string, unknown> | null => {
  if (value && typeof value === 'object' && !Array.isArray(value)) {
    return value as Record<string, unknown>
  }
  return null
}

const sanitizeHtml = (input: string) =>
  input
    .replace(/<script[\s\S]*?>[\s\S]*?<\/script>/gi, '')
    .replace(/on\w+\s*=\s*"[^"]*"/gi, '')
    .replace(/on\w+\s*=\s*'[^']*'/gi, '')

const extractImageUrlsFromHtml = (html: string) => {
  const matches = html.match(/<img[^>]+src=["']([^"']+)["'][^>]*>/gi) || []
  return matches
    .map((tag) => {
      const srcMatch = tag.match(/src=["']([^"']+)["']/i)
      return srcMatch ? srcMatch[1] : ''
    })
    .filter((src) => !!src)
}

const extractAttachmentUrls = (raw: unknown, fallbackHtml: string) => {
  const urls: string[] = []
  const push = (candidate: unknown) => {
    if (typeof candidate === 'string') {
      if (candidate) urls.push(candidate)
      return
    }
    const record = extractRecord(candidate)
    if (!record) return
    const possible = typeof record.url === 'string' ? (record.url as string) : typeof record.href === 'string' ? (record.href as string) : ''
    if (possible) {
      urls.push(possible)
    }
  }

  const record = extractRecord(raw)
  if (record) {
    const imageField = record.image
    if (imageField) {
      if (Array.isArray(imageField)) {
        imageField.forEach((img) => push(img))
      } else {
        push(imageField)
      }
    }

    const attachmentField = record.attachment
    if (attachmentField) {
      const allowAttachment = (att: unknown) => {
        const attRecord = extractRecord(att)
        if (!attRecord) {
          return true
        }
        const type = typeof attRecord.type === 'string' ? (attRecord.type as string).toLowerCase() : ''
        const mediaType = typeof attRecord.mediaType === 'string' ? (attRecord.mediaType as string).toLowerCase() : ''
        if (type && type !== 'image' && !mediaType.startsWith('image/')) {
          return false
        }
        return true
      }

      if (Array.isArray(attachmentField)) {
        attachmentField.forEach((att) => {
          if (allowAttachment(att)) {
            push(att)
          }
        })
      } else if (allowAttachment(attachmentField)) {
        push(attachmentField)
      }
    }
  }

  if (!urls.length && fallbackHtml) {
    urls.push(...extractImageUrlsFromHtml(fallbackHtml))
  }

  return Array.from(new Set(urls))
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

const resolveActorHandle = (item: App.Api.Fediverse.TimelineItem) => {
  const actorUrl = item.actorId || ''
  try {
    const parsed = new URL(actorUrl)
    const pathSegments = parsed.pathname.split('/').filter(Boolean)
    const userSlug = item.actorPreferredUsername || pathSegments[pathSegments.length - 1] || parsed.hostname
    return `${userSlug}@${parsed.hostname}`
  } catch {
    const username = item.actorPreferredUsername || item.actorDisplayName
    if (!username) {
      return actorUrl
    }
    const domain = actorUrl.split('://')[1]?.split('/')[0]
    return domain ? `${username}@${domain}` : username
  }
}

const resolvedItems = computed(() =>
  items.value.map((item) => {
    const contentHtmlCandidate = (item.content && item.content.trim()) || ''
    const sanitizedHtml = contentHtmlCandidate ? sanitizeHtml(contentHtmlCandidate) : ''
    const contentText = getTimelineContent(item)
    const images = extractAttachmentUrls(item.rawObject, contentHtmlCandidate)
    const galleryImages: App.Api.Ech0.Image[] = images.map((url, index) => ({
      id: index,
      message_id: item.id ?? 0,
      image_url: url,
      image_source: ImageSource.URL,
    }))

    return {
      id: item.id,
      displayName: item.actorDisplayName || item.actorPreferredUsername || '联邦好友',
      actorHandle: resolveActorHandle(item),
      timeText: formatTimelineTime(item.publishedAt || item.createdAt),
      contentHtml: sanitizedHtml,
      contentText,
      images,
      galleryImages,
      objectId: item.objectId,
    }
  })
)

</script>

<style scoped>

</style>
