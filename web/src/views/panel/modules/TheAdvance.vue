<template>
  <div class="w-full px-2">
    <!-- Webhook 设置 -->
    <TheWebhookSetting class="mb-3" />

    <!-- 联邦设置 -->
    <TheFediverseSetting />
  </div>
</template>

<script setup lang="ts">
import TheWebhookSetting from './TheSetting/TheWebhookSetting.vue'
import TheFediverseSetting from './TheSetting/TheFediverseSetting.vue'
</script>
