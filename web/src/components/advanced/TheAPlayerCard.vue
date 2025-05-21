<template>
  <meting-js
    v-if="musicInfo && metingAPI.length > 0"
    :api="metingAPI"
    :server="musicInfo.server"
    :type="musicInfo.type"
    :id="musicInfo.id"
    :auto="props.echo.extension"
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
import { computed, ref, onMounted } from 'vue'
import Music from '@/components/icons/music.vue'
import { parseMusicURL } from '@/utils/other'
import { fetchGetSettings } from '@/service/api'
type Echo = App.Api.Ech0.Echo
const enum ExtensionType {
  MUSIC = 'MUSIC',
  VIDEO = 'VIDEO',
  GITHUBPROJ = 'GITHUBPROJ',
}

const props = defineProps<{
  echo: Echo
}>()

const metingAPI = ref<string>('')
const musicInfo = computed(() => {
  if (props.echo.extension_type !== ExtensionType.MUSIC || !props.echo.extension) return null
  return parseMusicURL(props.echo.extension)
})

onMounted(async () => {
  await fetchGetSettings().then((res) => {
    if (res.code === 1 && res.data.meting_api.length > 0) {
      metingAPI.value = res.data.meting_api + '?server=:server&type=:type&id=:id&auth=:auth&r=:r'
      console.log(metingAPI)
    } else {
      metingAPI.value =
        'https://meting.soopy.cn/api?server=:server&type=:type&id=:id&auth=:auth&r=:r'
    }
  })
})
</script>

<style scoped></style>
