<template>
  <div class="w-full px-2">
    <!-- 存储设置 -->
    <TheStorageSetting />

    <!-- 数据管理 -->
    <TheBackupSetting />

    <!-- Webhook 设置 -->
    <TheWebhookSetting />
  </div>
</template>

<script setup lang="ts">
import TheBackupSetting from './TheSetting/TheBackupSetting.vue'
import TheStorageSetting from './TheSetting/TheStorageSetting.vue'
import TheWebhookSetting from './TheSetting/TheWebhookSetting.vue'
</script>
