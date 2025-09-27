<template>
  <div class="w-full px-2">
    <!-- 系统设置 -->
    <TheSystemSetting />

    <!-- 评论设置 -->
    <TheCommentSetting />

    <!-- Connect设置 -->
    <TheConnectSetting />

    <!-- 存储设置 -->
    <TheStorageSetting />
  </div>
</template>

<script setup lang="ts">
import TheSystemSetting from './TheSetting/TheSystemSetting.vue'
import TheCommentSetting from './TheSetting/TheCommentSetting.vue'
import TheConnectSetting from './TheSetting/TheConnectSetting.vue'
import TheStorageSetting from './TheSetting/TheStorageSetting.vue'
</script>
