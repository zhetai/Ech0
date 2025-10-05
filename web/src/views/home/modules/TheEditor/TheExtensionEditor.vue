<template>
  <div v-if="currentMode === Mode.EXTEN">
    <!-- 音乐分享 -->
    <div v-if="currentExtensionType === ExtensionType.MUSIC">
      <h2 class="text-gray-500 font-bold mb-1">音乐分享（支持网易云/QQ音乐/Apple Music）</h2>
      <p class="text-gray-400 text-sm mb-1">注意：不支持VIP歌曲，建议使用自建API</p>
      <BaseInput
        v-model="extensionToAdd.extension"
        class="rounded-lg h-auto w-full"
        placeholder="音乐链接..."
      />
      <div
        v-if="
          extensionToAdd.extension.length > 0 &&
          extensionToAdd.extension_type === ExtensionType.MUSIC
        "
        class="mt-1 text-gray-400 text-md"
      >
        解析结果：
        <span v-if="parseMusicURL(extensionToAdd.extension)" class="text-green-400">成功</span>
        <span v-else class="text-red-300">失败</span>
      </div>
    </div>
    <!-- Bilibili/YouTube视频分享 -->
    <div v-if="currentExtensionType === ExtensionType.VIDEO">
      <div class="text-gray-500 font-bold mb-1">视频分享（支持Bilibili、YouTube）</div>
      <div class="text-gray-400 mb-1">粘贴自动提取ID</div>
      <BaseInput
        v-model="videoURL"
        class="rounded-lg h-auto w-full my-2"
        placeholder="B站/YouTube链接..."
      />
      <div class="text-gray-500 my-1">Video ID：{{ extensionToAdd.extension }}</div>
    </div>
    <!-- Github项目分享 -->
    <div v-if="currentExtensionType === ExtensionType.GITHUBPROJ">
      <div class="text-gray-500 font-bold mb-1">Github项目分享</div>
      <BaseInput
        v-model="extensionToAdd.extension"
        class="rounded-lg h-auto w-full"
        placeholder="https://github.com/username/repo"
      />
    </div>
    <!-- 网站链接分享 -->
    <div v-if="currentExtensionType === ExtensionType.WEBSITE">
      <div class="text-gray-500 font-bold mb-1">网站链接分享</div>
      <!-- 网站标题 -->
      <BaseInput
        v-model="websiteToAdd.title"
        class="rounded-lg h-auto w-full mb-2"
        placeholder="网站标题..."
      />
      <BaseInput
        v-model="websiteToAdd.site"
        class="rounded-lg h-auto w-full"
        placeholder="https://example.com"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import BaseInput from '@/components/common/BaseInput.vue'
import { ExtensionType, Mode } from '@/enums/enums'
import { parseMusicURL } from '@/utils/other'

const currentMode = defineModel<Mode>('currentMode')
const currentExtensionType = defineModel<ExtensionType>('currentExtensionType')
const extensionToAdd = defineModel<{
  extension_type: string
  extension: string
}>('extensionToAdd', { required: true })
const videoURL = defineModel<string>('videoURL')
const websiteToAdd = defineModel<{
  title: string
  site: string
}>('websiteToAdd', {
  required: true,
})
</script>

<style scoped></style>
