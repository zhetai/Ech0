<template>
  <div class="w-full px-2">
    <!-- 用户设置 -->
    <TheUserSetting />

    <!-- 用户管理 -->
    <TheUserManager />
  </div>
</template>

<script setup lang="ts">
import TheUserSetting from './TheSetting/TheUserSetting.vue'
import TheUserManager from './TheSetting/TheUserManager.vue'
</script>
