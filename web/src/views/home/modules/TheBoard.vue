<template>
  <div
    class="bg-white rounded-lg ring-1 ring-gray-200 ring-inset mx-auto shadow-sm hover:shadow-md"
  >
    <div class="mx-auto w-full px-3 py-4">
      <!-- The Title && Nav -->
      <TheTitleAndNav />

      <div class="mt-4 px-3 gap-2 font-serif flex justify-between items-baseline">
        <p class="text-[#5b7083] text-md">
          <RouterLink to="/auth" class="font-bold">登录</RouterLink> 后一吐为快
        </p>
        <span class="text-gray-400 text-xs"
          >Powered by <a href="https://github.com/lin-snow/Ech0">Ech0</a></span
        >
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import TheTitleAndNav from './TheEditor/TheTitleAndNav.vue'
</script>
