<template>
  <div class="w-full px-2">
    <!-- 系统设置 -->
    <TheSystemSetting class="mb-3" />

    <!-- 评论设置 -->
    <TheCommentSetting class="mb-3" />

    <!-- Connect设置 -->
    <TheConnectSetting class="mb-3" />
  </div>
</template>

<script setup lang="ts">
import TheSystemSetting from './TheSetting/TheSystemSetting.vue'
import TheCommentSetting from './TheSetting/TheCommentSetting.vue'
import TheConnectSetting from './TheSetting/TheConnectSetting.vue'
</script>
