<template>
  <div class="w-full px-2">
    <!-- 连接设置 -->
    <TheConnectSetting class="mb-3" />
    <!-- 评论设置 -->
    <TheCommentSetting />
  </div>
</template>

<script setup lang="ts">
import TheConnectSetting from './TheSetting/TheConnectSetting.vue'
import TheCommentSetting from './TheSetting/TheCommentSetting.vue'
</script>

<style scoped></style>
