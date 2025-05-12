<template>
  <div class="sm:max-w-sm mx-auto flex justify-between items-center mb-2">
    <div class="w-full flex justify-between items-center">
        <BaseInput
          title="搜索"
          type="text"
          v-model="searchContent"
          placeholder="搜索..."
          class="w-1/3 h-9"
          @blur="handleSearch"
        />
        <BaseButton
          :icon="Other"
          @click="notDeveloper"
          class="w-10 h-9 rounded-md !shadow-none !border-none !ring-0 !bg-transparent"
          title="外界"
        />
    </div>
  </div>
</template>

<script setup lang="ts">
import BaseButton from '@/components/common/BaseButton.vue';
import BaseInput from '@/components/common/BaseInput.vue';
import Other from '@/components/icons/other.vue';
import router from '@/router';
import { useEchoStore } from '@/stores/echo';
import { storeToRefs } from 'pinia';
import { ref } from 'vue';
import { theToast } from '@/utils/toast';
const echoStore = useEchoStore()
const { refreshForSearch, getEchosByPage } = echoStore
const { searchingMode } = storeToRefs(echoStore)

const searchContent = ref<string>('')

const handleSearch = () => {
  // 设置搜索内容
  echoStore.searchValue = searchContent.value

  // 判断是否是搜索模式
  if (searchingMode.value) {
    // 初始化搜索
    refreshForSearch();
    // 开始搜索
    getEchosByPage();
  }
}

const notDeveloper = () => {
  theToast.info('该功能正在开发中，请耐心等待！')
}
</script>
