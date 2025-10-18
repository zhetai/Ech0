<template>
  <div class="w-full px-2">
    <!-- 系统设置 -->
    <TheSystemSetting class="mb-3" />

    <!-- 访问令牌设置 -->
    <TheAccessTokenSetting class="mb-3" />
  </div>
</template>

<script setup lang="ts">
import TheSystemSetting from './TheSetting/TheSystemSetting.vue'
import TheAccessTokenSetting from './TheSetting/TheAccessTokenSetting.vue'
</script>
