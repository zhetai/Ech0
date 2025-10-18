<template>
  <div class="w-full px-2">
    <!-- 存储设置 -->
    <TheStorageSetting class="mb-3" />

    <!-- 数据管理 -->
    <TheBackupSetting class="mb-3" />
  </div>
</template>

<script setup lang="ts">
import TheStorageSetting from './TheSetting/TheStorageSetting.vue'
import TheBackupSetting from './TheSetting/TheBackupSetting.vue'
</script>

<style scoped></style>
