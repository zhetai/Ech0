<template>
  <div class="w-full px-2">
    <!-- 用户设置 -->
    <TheUserSetting />

    <!-- OAuth2 设置 -->
    <TheOAuth2Setting />

    <!-- 用户管理 -->
    <TheUserManager />
  </div>
</template>

<script setup lang="ts">
import TheUserSetting from './TheSetting/TheUserSetting.vue'
import TheOAuth2Setting from './TheSetting/TheOAuth2Setting.vue'
import TheUserManager from './TheSetting/TheUserManager.vue'
</script>
