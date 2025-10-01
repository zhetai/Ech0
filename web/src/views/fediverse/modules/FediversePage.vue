<template>
  <div class="px-3 pb-4 py-2 mt-4 sm:mt-6 mb-10 mx-auto flex justify-center items-center">
    <div class="w-full sm:max-w-lg mx-auto">
      <!-- 返回上一页 -->
      <!-- <div class="mx-auto max-w-sm">
        <BaseButton @click="goBack"
          class="text-gray-600 rounded-md !shadow-none !border-none !ring-0 !bg-transparent group" title="返回首页">
          <Arrow class="w-9 h-9 rotate-180 transition-transform duration-200 group-hover:-translate-x-1" />
        </BaseButton>
      </div> -->

      <h1 class="text-4xl text-center font-bold text-gray-400">欢迎来到联邦宇宙</h1>

    </div>
  </div>
</template>

<script setup lang="ts">
// import { useRouter, useRoute } from 'vue-router'
// import BaseButton from '@/components/common/BaseButton.vue'
// import Arrow from '@/components/icons/arrow.vue'

// const router = useRouter()

// const goBack = () => {
//   if (window.history.length > 2) {
//     window.history.back()
//   } else {
//     router.push({ name: 'home' }) // 没有历史记录则跳首页
//   }
// }

</script>
