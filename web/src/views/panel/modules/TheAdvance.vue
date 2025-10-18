<template>
  <div class="w-full px-2">
    <!-- Webhook 设置 -->
    <TheWebhookSetting />

    <!-- 数据管理 -->
    <TheBackupSetting />
  </div>
</template>

<script setup lang="ts">
import TheBackupSetting from './TheSetting/TheBackupSetting.vue'
import TheWebhookSetting from './TheSetting/TheWebhookSetting.vue'
</script>
