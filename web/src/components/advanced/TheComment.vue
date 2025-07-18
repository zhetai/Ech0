<template>
  <div id="comments" class="max-w-sm h-auto p-4 my-4 mx-auto">
    <div v-if="enabled" id="comment"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useSettingStore } from '@/stores/settting'
import { storeToRefs } from 'pinia'

declare global {
  // @ts-nocheck
  /* eslint-disable */
  interface Window {
    twikoo: any
  }
}

const { SystemSetting, loading } = storeToRefs(useSettingStore())

const enabled = ref<boolean>(false)
const commentapi = ref<string>('')

// 加载 twikoo 脚本
function loadScript(src: string) {
  return new Promise((resolve, reject) => {
    const script = document.createElement('script')
    script.src = src
    script.onload = resolve
    script.onerror = reject
    document.head.appendChild(script)
  })
}

watch(
  () => loading.value,
  async (newVal) => {
    if (!newVal && SystemSetting.value.comment_api) {
      commentapi.value = SystemSetting.value.comment_api
      enabled.value = true

      await loadScript('/others/scripts/twikoo.all.min.js')
      window.twikoo.init({
        envId: commentapi.value,
        el: '#comment',
      })
    }
  },
  { immediate: true },
)
</script>

<style>
.twikoo {
  width: 100% !important;
  max-width: 500px !important;
}
.tk-meta-input {
  flex-direction: column !important;
  padding: 1px 3px !important;
}

.tk-meta-input .el-input {
  width: 100% !important;
  flex: 1 !important;
  margin: 4px 0px !important;
}

.tk-meta-input .el-input + .el-input {
  margin-left: 0 !important;
}

.el-textarea {
  position: relative;
  display: inline-block;
  width: 100%;
  vertical-align: bottom;
  font-size: 14px;
  padding: 0px 2px !important;
}

.tk-comments-container {
  width: 100% !important;
  max-width: 500px !important;
  padding: 0 0.4rem !important;
  min-height: 10rem;
  display: flex;
  flex-direction: column;
}
</style>
