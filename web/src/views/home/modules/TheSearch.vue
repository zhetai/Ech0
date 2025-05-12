<template>
  <div class="sm:max-w-sm mx-auto flex justify-between items-center mb-2">
    <div class="w-full">
        <BaseInput
          title="搜索"
          type="text"
          v-model="searchContent"
          placeholder="搜索..."
          class="w-1/3 h-9"
          @keyup.enter="handleSearch"
        />
    </div>
  </div>
</template>

<script setup lang="ts">
import BaseInput from '@/components/common/BaseInput.vue';
import { useEchoStore } from '@/stores/echo';
import { ref } from 'vue';
const echoStore = useEchoStore()
const { refreshForSearch, getEchosByPage } = echoStore

const searchContent = ref<string>('')

const handleSearch = () => {
  // 初始化搜索
  refreshForSearch();
  // 设置搜索内容
  echoStore.searchValue = searchContent.value
  // 开始搜索
  getEchosByPage();
}
</script>
