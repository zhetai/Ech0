<template>
  <meting-js
    v-if="musicInfo"
    :api="`${SystemSetting.meting_api ? SystemSetting.meting_api : 'https://meting.soopy.cn/api'}?server=:server&type=:type&id=:id&auth=:auth&r=:r`"
    :server="musicInfo.server"
    :type="musicInfo.type"
    :id="musicInfo.id"
    :auto="`{{ props.echo.extension }}`"
  >
  </meting-js>
  <div
    v-else
    class="max-w-sm flex justify-center items-center bg-white rounded-lg shadow-sm p-2 gap-2 text-gray-400"
  >
    <Music />非常抱歉，该音乐播放源已失效😭
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import Music from '@/components/icons/music.vue'
import { useSettingStore } from '@/stores/settting'
import { parseMusicURL } from '@/utils/other'
type Echo = App.Api.Ech0.Echo
const enum ExtensionType {
  MUSIC = 'MUSIC',
  VIDEO = 'VIDEO',
  GITHUBPROJ = 'GITHUBPROJ',
}
const settingStore = useSettingStore()
const { SystemSetting } = storeToRefs(settingStore)

const props = defineProps<{
  echo: Echo
}>()

const musicInfo = computed(() => {
  if (props.echo.extension_type !== ExtensionType.MUSIC || !props.echo.extension) return null
  return parseMusicURL(props.echo.extension)
})
</script>

<style scoped></style>
