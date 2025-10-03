<template>
  <article
    class="relative flex flex-col gap-3 rounded-md bg-white px-5 py-5 shadow-sm transition duration-200 hover:-translate-y-0.5 hover:shadow-md"
  >
    <header class="flex items-center gap-3">
      <div class="flex h-12 w-12 shrink-0 items-center justify-center overflow-hidden rounded-full shadow-sm ring-inset ring-1 ring-gray-200 bg-amber-50">
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

      <BaseButton
        v-if="actor.id"
        class="shrink-0 rounded-md border-dashed shadow-none border-amber-500 bg-transparent px-3 py-1 text-sm font-medium text-amber-600 transition  focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-300 focus-visible:ring-offset-2 focus-visible:ring-offset-white"
        :disabled="isFollowDisabled"
        @click="handleFollowClick"
      >
        <span v-if="followSuccess">已发出关注</span>
        <span v-else-if="followLoading">关注中...</span>
        <span v-else>关注</span>
      </BaseButton>
    </header>

    <section v-if="sanitizedSummary" class="prose prose-amber max-w-none text-sm text-stone-600">
      <div v-html="sanitizedSummary" />
    </section>

    <footer class="flex flex-wrap items-center gap-0.5 text-xs text-stone-500">
      <span v-if="actor.type">类型: {{ actor.type }}</span>
      <span v-if="actor.inbox">收件箱: <span class="text-amber-600">{{ actor.inbox }}</span></span>
    </footer>
  </article>
</template>

<script setup lang="ts">
import { computed, withDefaults } from 'vue'
import BaseButton from '@/components/common/BaseButton.vue'

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

const props = withDefaults(
  defineProps<{ actor: ActivityPubActor; followLoading?: boolean; followSuccess?: boolean }>(),
  {
    followLoading: false,
    followSuccess: false,
  },
)

const emit = defineEmits<{
  (e: 'follow', actor: ActivityPubActor): void
}>()

const actor = computed(() => props.actor ?? {})
const followLoading = computed(() => props.followLoading ?? false)
const followSuccess = computed(() => props.followSuccess ?? false)
const isFollowDisabled = computed(() => followLoading.value || followSuccess.value)

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

const handleFollowClick = () => {
  if (isFollowDisabled.value) return
  emit('follow', actor.value)
}
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
