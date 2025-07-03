<template>
  <div class="w-full">
    <!-- 日期时间 && 操作按钮 -->
    <div class="flex justify-between items-center">
      <!-- 日期时间 -->
      <div class="flex justify-start items-center h-auto">
        <!-- 小点 -->
        <div class="w-2 h-2 rounded-full bg-orange-600 mr-2"></div>
        <!-- 具体日期时间 -->
        <div class="flex justify-start text-sm text-orange-500">
          {{ formatDate(props.echo.created_at) }}
        </div>
        <!-- 用户名称 -->
        <!-- <div class="flex justify-start items-center">
          <span class="text-md mx-1 text-orange-400"> @ </span>
          <span class="text-sm text-orange-400">{{ props.echo.username }}</span>
        </div> -->
      </div>
      <!-- 操作按钮 -->
      <div class="flex justify-end items-center gap-1 h-auto">
        <!-- 是否隐私 -->
        <span v-if="props.echo.private" title="私密状态"><Lock /></span>
        <!-- 删除 -->
        <button
          v-if="userStore.isLogin"
          @click="handleDeleteEcho(props.echo.id)"
          title="删除"
          class="transform transition-transform duration-200 active:scale-160 active:animate-pulse"
        >
          <Roll />
        </button>
        <button
          v-if="userStore.isLogin"
          @click="handleUpdateEcho()"
          title="更新"
          class="transform transition-transform duration-200 active:scale-160 active:animate-pulse"
        >
          <EditEcho />
        </button>

        <!-- 点赞 -->
        <div class="flex items-center justify-end" title="点赞">
          <div class="flex items-center gap-1">
            <!-- 点赞按钮   -->
            <button
              @click="handleLikeEcho(props.echo.id)"
              title="点赞"
              class="transform transition-transform duration-200 active:scale-160 active:animate-pulse"
            >
              <GrayLike class="w-4 h-4 transition-colors duration-200 hover:text-red-500" />
            </button>

            <!-- 点赞数量   -->
            <span class="text-sm text-gray-400">
              <!-- 如果点赞数不超过99，则显示数字，否则显示99+ -->
              {{ props.echo.fav_count > 99 ? '99+' : props.echo.fav_count }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- 图片 && 内容 -->
    <div class="border-l-2 border-gray-300 ml-1 mb-1">
      <div class="p-6">
        <!-- 图片 -->
        <div v-if="props.echo.images && props.echo.images.length > 0" class="w-5/6 mx-auto">
          <div class="shadow-lg rounded-lg overflow-hidden mb-2">
            <a :href="getImageUrl(props.echo.images[imageIndex])" data-fancybox>
              <img
                :src="getImageUrl(props.echo.images[imageIndex])"
                alt="Image"
                class="max-w-full object-cover"
                loading="lazy"
              />
            </a>
          </div>
          <!-- 图片切换 -->
          <div v-if="props.echo.images.length > 1" class="flex items-center justify-center">
            <button @click="imageIndex = Math.max(imageIndex - 1, 0)">
              <Prev class="w-6 h-6" />
            </button>
            <span class="text-gray-500 text-sm mx-2">
              {{ imageIndex + 1 }} / {{ props.echo.images.length }}
            </span>
            <button @click="imageIndex = Math.min(imageIndex + 1, props.echo.images.length - 1)">
              <Next class="w-6 h-6" />
            </button>
          </div>
        </div>

        <!-- 内容 -->
        <div>
          <MdPreview
            :id="previewOptions.proviewId"
            :modelValue="props.echo.content"
            :theme="previewOptions.theme"
            :show-code-row-number="previewOptions.showCodeRowNumber"
            :preview-theme="previewOptions.previewTheme"
            :code-theme="previewOptions.codeTheme"
            :code-style-reverse="previewOptions.codeStyleReverse"
            :no-img-zoom-in="previewOptions.noImgZoomIn"
            :code-foldable="previewOptions.codeFoldable"
            :auto-fold-threshold="previewOptions.autoFoldThreshold"
          />
        </div>

        <!-- 扩展内容 -->
        <div v-if="props.echo.extension" class="my-2">
          <div v-if="props.echo.extension_type === ExtensionType.MUSIC">
            <TheAPlayerCard :echo="props.echo" />
          </div>
          <div v-if="props.echo.extension_type === ExtensionType.VIDEO">
            <TheVideoCard :videoId="props.echo.extension" class="px-2 mx-auto hover:shadow-md" />
          </div>
          <TheGithubCard
            v-if="props.echo.extension_type === ExtensionType.GITHUBPROJ"
            :GithubURL="props.echo.extension"
            class="px-2 mx-auto hover:shadow-md"
          />
          <TheWebsiteCard
            v-if="props.echo.extension_type === ExtensionType.WEBSITE"
            :website="props.echo.extension"
            class="px-2 mx-auto hover:shadow-md"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Fancybox } from '@fancyapps/ui'
import { MdPreview } from 'md-editor-v3'
import { getImageUrl } from '@/utils/other'
import { onMounted, ref } from 'vue'
import { fetchDeleteEcho, fetchLikeEcho } from '@/service/api'
import { theToast } from '@/utils/toast'
import { useUserStore } from '@/stores/user'
import TheGithubCard from './TheGithubCard.vue'
import TheVideoCard from './TheVideoCard.vue'
import '@fancyapps/ui/dist/fancybox/fancybox.css'
import 'md-editor-v3/lib/preview.css'
import Roll from '../icons/roll.vue'
import Lock from '../icons/lock.vue'
import Prev from '../icons/prev.vue'
import Next from '../icons/next.vue'
import GrayLike from '../icons/graylike.vue'
import EditEcho from '../icons/editecho.vue'
import TheAPlayerCard from './TheAPlayerCard.vue'
import TheWebsiteCard from './TheWebsiteCard.vue'
import { useEchoStore } from '@/stores/echo'
import { localStg } from '@/utils/storage'

const emit = defineEmits(['refresh', 'updateLikeCount'])

type Echo = App.Api.Ech0.Echo
const enum ExtensionType {
  MUSIC = 'MUSIC',
  VIDEO = 'VIDEO',
  GITHUBPROJ = 'GITHUBPROJ',
  WEBSITE = 'WEBSITE',
}

const props = defineProps<{
  echo: Echo
}>()
const imageIndex = ref<number>(0)
const userStore = useUserStore()
const previewOptions = {
  proviewId: 'preview-only',
  theme: 'light' as 'light' | 'dark',
  showCodeRowNumber: false,
  previewTheme: 'github',
  codeTheme: 'atom',
  codeStyleReverse: true,
  noImgZoomIn: false,
  codeFoldable: true,
  autoFoldThreshold: 15,
}

const echoStore = useEchoStore()

const handleDeleteEcho = (echoId: number) => {
  // 浏览器alert弹窗确认删除
  if (confirm('确定要删除吗？')) {
    fetchDeleteEcho(echoId).then(() => {
      theToast.success('删除成功！')
      // 触发父组件的刷新事件emit
      emit('refresh')
    })
  }
}

const handleUpdateEcho = () => {
  if (echoStore.isUpdateMode) {
    // 如果已经在更新模式，返回顶部并提示用户先退出更新模式
    window.scrollTo({ top: 0, behavior: 'smooth' })
    theToast.warning('请先退出更新模式！')
    return
  }

  echoStore.echoToUpdate = props.echo
  echoStore.isUpdateMode = true
}

const LIKE_LIST_KEY = 'likedEchoIds'
const likedEchoIds: number[] = localStg.getItem(LIKE_LIST_KEY) || []
const hasLikedEcho = (echoId: number): boolean => {
  return likedEchoIds.includes(echoId)
}
const handleLikeEcho = (echoId: number) => {
  // 检查LocalStorage中是否已经点赞过
  if (hasLikedEcho(echoId)) {
    theToast.success('你已经点赞过了,感谢你的喜欢！')
    return
  }

  fetchLikeEcho(echoId).then((res) => {
    if (res.code === 1) {
      likedEchoIds.push(echoId)
      localStg.setItem(LIKE_LIST_KEY, likedEchoIds)
      // 发送更新事件
      emit('updateLikeCount', echoId)
      theToast.success('点赞成功！')
    }
  })
}

const formatDate = (dateString: string) => {
  // 当天则显示（时：分）
  // 非当天但是三内天则显示几天前
  // 超过三天则显示（时：分 年月日）
  const date = new Date(dateString)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const diffInDays = Math.floor(diff / (1000 * 60 * 60 * 24))
  const diffInHours = Math.floor(diff / (1000 * 60 * 60))
  const diffInMinutes = Math.floor(diff / (1000 * 60))

  const diffInSeconds = Math.floor(diff / 1000)
  if (diffInSeconds < 60) {
    return '刚刚'
  } else if (diffInMinutes < 60) {
    return `${diffInMinutes}分钟前`
  } else if (diffInHours < 24) {
    return `${diffInHours}小时前`
  } else if (diffInDays < 3) {
    return `${diffInDays}天前`
  } else {
    return date.toLocaleString() // 返回完整的日期和时间
  }
}

onMounted(() => {
  Fancybox.bind('[data-fancybox]', {
    // Your custom options
  })
})
</script>

<style scoped lang="css">
#preview-only {
  background-color: inherit;
}

.md-editor {
  font-family: 'LXGW WenKai Screen';
}

:deep(ul li) {
  list-style-type: disc;
}
:deep(ul li li) {
  list-style-type: circle;
}
:deep(ul li li li) {
  list-style-type: square;
}
:deep(ol li) {
  list-style-type: decimal;
}
</style>
