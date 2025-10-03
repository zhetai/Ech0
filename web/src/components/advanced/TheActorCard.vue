<template>
  <article
    class="relative flex flex-col gap-3 rounded-xl border border-amber-100/80 bg-white/90 px-5 py-5 shadow-md transition duration-200 hover:-translate-y-0.5 hover:shadow-lg"
  >
    <header class="flex items-center gap-3">
      <div class="flex h-12 w-12 shrink-0 items-center justify-center overflow-hidden rounded-full border border-amber-100 bg-amber-50">
        <img
          v-if="avatarUrl"
          :src="avatarUrl"
          :alt="displayName"
          class="h-full w-full object-cover"
        />
        <span v-else class="text-lg font-semibold text-amber-600">{{ initials }}</span>
      </div>
      <div class="min-w-0 flex-1">
        <p class="truncate text-lg font-semibold text-stone-800">
          {{ displayName }}
        </p>
        <p v-if="username" class="truncate text-sm text-stone-500">@{{ username }}</p>
      </div>
      <a
        v-if="actor.id"
        :href="actor.id"
        class="shrink-0 text-sm text-amber-600 transition hover:text-amber-500"
        target="_blank"
        rel="noopener noreferrer"
      >
        查看
      </a>
    </header>

    <section v-if="sanitizedSummary" class="prose prose-amber max-w-none text-sm text-stone-600">
      <div v-html="sanitizedSummary" />
    </section>

    <footer class="flex flex-wrap items-center gap-3 text-xs text-stone-500">
      <span v-if="actor.type">类型: {{ actor.type }}</span>
      <span v-if="actor.inbox">收件箱: <span class="text-amber-600">{{ actor.inbox }}</span></span>
    </footer>
  </article>
</template>

<script setup lang="ts">
import { computed } from 'vue'

type ActivityPubMedia = {
  type?: string
  mediaType?: string
  url?: string
}

type ActivityPubActor = {
  id?: string
  type?: string
  name?: string
  preferredUsername?: string
  summary?: string
  icon?: ActivityPubMedia | ActivityPubMedia[]
  image?: ActivityPubMedia | ActivityPubMedia[]
  inbox?: string
}

const props = defineProps<{ actor: ActivityPubActor }>()

const actor = computed(() => props.actor ?? {})

const displayName = computed(() => actor.value.name || actor.value.preferredUsername || '未知用户')
const username = computed(() => actor.value.preferredUsername || actor.value.name || '')

const avatarUrl = computed(() => {
  const icon = actor.value.icon ?? actor.value.image
  if (!icon) return ''
  if (Array.isArray(icon)) {
    const candidate = icon.find((item) => item?.url)
    return candidate?.url ?? ''
  }
  return icon.url ?? ''
})

const initials = computed(() => {
  const name = displayName.value.trim()
  if (!name) return 'AP'
  const words = name.split(/\s+/)
  if (words.length === 1) return words[0].charAt(0).toUpperCase()
  return (words[0].charAt(0) + words[1].charAt(0)).toUpperCase()
})

const sanitizedSummary = computed(() => {
  const summary = actor.value.summary?.trim()
  if (!summary) return ''
  return summary
})
</script>

<style scoped>
.prose :deep(a) {
  color: rgb(217 119 6);
  text-decoration: underline;
}

.prose :deep(a:hover) {
  color: rgb(180 83 9);
}
</style>
