<template>
  <meting-js
    class="!shadow-none"
    v-if="musicInfo"
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
import Music from '@/components/icons/music.vue'
type Echo = App.Api.Ech0.Echo
const enum ExtensionType {
  MUSIC = 'MUSIC',
  VIDEO = 'VIDEO',
  GITHUBPROJ = 'GITHUBPROJ',
}

const props = defineProps<{
  echo: Echo
}>()

const musicInfo = computed(() => {
  if (props.echo.extension_type !== ExtensionType.MUSIC || !props.echo.extension) return null
  return parseMusicURL(props.echo.extension)
})

// 解析音乐链接（网易云、QQ音乐）
const parseMusicURL = (url: string) => {
  url = url.trim()

  const neteaseMatch = url.match(/music\.163\.com\/#\/(song|playlist|album)\?id=(\d+)/)
  if (neteaseMatch) {
    return {
      server: 'netease',
      type: neteaseMatch[1], // song, playlist, album
      id: neteaseMatch[2],
    }
  }

  // QQ音乐 新格式支持，songDetail 路径，id一般是字母数字混合
  const qqNewSongMatch = url.match(/y\.qq\.com\/n\/ryqq\/songDetail\/([a-zA-Z0-9]+)/)
  if (qqNewSongMatch) {
    return {
      server: 'tencent',
      type: 'song',
      id: qqNewSongMatch[1],
    }
  }

  // 解析失败
  return null
}
</script>

<style scoped></style>
