<template>
  <!-- 瀑布流缩略图 -->
  <div
    v-if="images?.length"
    :class="[
      'w-5/6 mx-auto grid gap-2 mb-4',
      images.length === 1 ? 'grid-cols-1 justify-items-center' : 'grid-cols-2',
    ]"
  >
    <button
      v-for="(src, idx) in images"
      :key="idx"
      class="bg-transparent border-0 p-0 cursor-pointer w-fit"
      :class="getColSpan(idx, images.length)"
      @click="active = idx"
    >
      <img
        :src="getImageUrl(src)"
        alt="`预览图片${idx + 1}`"
        loading="lazy"
        class="block rounded-md max-w-full h-auto"
      />
    </button>
  </div>

  <!-- 灯箱层 -->
  <Teleport to="body">
    <transition name="fade">
      <div
        v-if="active !== null"
        class="fixed inset-0 w-screen h-screen bg-black/80 backdrop-blur-[12px] flex justify-center items-center z-[9999] overflow-hidden"
        @click.self="active = null"
      >
        <img
          :src="getImageUrl(images[active])"
          class="max-w-[80vw] max-h-[80vh] rounded-md cursor-pointer object-contain shadow-[0_4px_12px_rgba(0,0,0,0.4)] transition ease-in-out duration-200"
          @click="active = null"
        />
      </div>
    </transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { getImageUrl } from '@/utils/other'

defineProps<{
  images: App.Api.Ech0.Echo['images']
}>()

const active = ref<number | null>(null)
const getColSpan = (idx: number, total: number) => {
  // 单张图片占满
  if (total === 1) return 'col-span-1 justify-self-center'
  // 第一张在奇数张时跨两列
  if (idx === 0 && total % 2 !== 0) return 'col-span-2'
  // 其他图片默认占一列
  return ''
}

function close() {
  active.value = null
}

function onKeyDown(e: KeyboardEvent) {
  if (e.key === 'Escape' && active.value !== null) {
    close()
  }
}

onMounted(() => {
  window.addEventListener('keydown', onKeyDown)
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', onKeyDown)
})
</script>

<style scoped>
/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.05s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
