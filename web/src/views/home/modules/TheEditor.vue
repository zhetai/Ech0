<template>
  <div
    class="bg-white rounded-lg ring-1 ring-gray-200 ring-inset mx-auto shadow-sm hover:shadow-md"
  >
    <div class="mx-auto w-full px-3 py-4">
      <!-- Title && Nav -->
      <TheTitleAndNav />

      <!-- Editor -->
      <div class="rounded-lg p-2 sm:p-3 mb-1">
        <TheMdEditor
          v-model="echoToAdd.content"
          class="rounded-lg"
          v-if="currentMode === Mode.ECH0"
        />

        <!-- ImageMode -->
        <div v-if="currentMode === Mode.Image">
          <h2 class="text-gray-500 font-bold">插入图片（支持本地上传、直链）</h2>
          <p class="text-gray-400 text-sm mb-2">注意：仅允许添加一张</p>
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-2">
              <span class="text-gray-500">选择添加方式：</span>
              <!-- 直链 -->
              <BaseButton
                :icon="Url"
                class="w-7 h-7 sm:w-7 sm:h-7 rounded-md"
                @click="imageToAdd.image_source = ImageSource.URL"
                title="插入图片链接"
              />
              <!-- 上传本地 -->
              <BaseButton
                :icon="Upload"
                class="w-7 h-7 sm:w-7 sm:h-7 rounded-md"
                @click="imageToAdd.image_source = ImageSource.LOCAL"
                title="上传本地图片"
              />
            </div>
            <div>
              <BaseButton
                v-if="imageToAdd.image_url != ''"
                :icon="Addmore"
                class="w-7 h-7 sm:w-7 sm:h-7 rounded-md"
                @click="handleAddMoreImage"
                title="添加更多图片"
              />
            </div>
          </div>
          <div class="my-1">
            <!-- 图片上传本地 -->
            <input
              id="file-input"
              class="hidden"
              type="file"
              accept="image/*"
              ref="fileInput"
              @change="handleUploadImage"
            />
            <BaseButton
              v-if="imageToAdd.image_source === ImageSource.LOCAL"
              @click="handleTriggerUpload"
              class="rounded-md"
              title="上传本地图片"
            >
              <span class="text-gray-400">点击上传</span>
            </BaseButton>
            <!-- 图片直链 -->
            <BaseInput
              v-if="imageToAdd.image_source === ImageSource.URL"
              v-model="imageToAdd.image_url"
              class="rounded-lg h-auto w-full"
              placeholder="请输入图片链接..."
            />
          </div>
        </div>

        <!-- todoMode -->
        <BaseTextArea
          v-if="currentMode === Mode.TODO"
          v-model="todoToAdd.content"
          class="rounded-lg h-auto sm:min-h-[6rem] md:min-h-[9rem]"
          placeholder="请输入待办事项..."
          :rows="3"
        />
        <!-- PlayMusic(上传音乐文件) -->
        <div v-if="currentMode === Mode.PlayMusic">
          <h2 class="text-gray-500 font-bold mb-1">欢迎使用音乐播放模式（仅PC）</h2>
          <div class="mb-1 flex items-center gap-2">
            <p class="text-gray-500">上传音乐：</p>
            <input
              id="file-input"
              class="hidden"
              type="file"
              accept="audio/*"
              ref="fileInput"
              @change="handleUploadMusic"
            />
            <BaseButton
              :icon="Audio"
              @click="handleTriggerUpload"
              class="w-7 h-7 sm:w-7 sm:h-7 rounded-md"
              title="上传音乐"
            />
          </div>
          <div class="flex items-center gap-2">
            <p class="text-gray-500">删除音乐：</p>
            <BaseButton
              :icon="Delete"
              @click="handleDeleteMusic"
              class="w-7 h-7 sm:w-7 sm:h-7 rounded-md"
              title="删除音乐"
            />
          </div>
        </div>

        <!-- Panel -->
        <TheModePanel
          v-if="currentMode === Mode.Panel"
          @switch-todo="handleSwitchTodoMode"
          @switch-extension="handleSwitchExtensionMode"
          @switch-play-music="handleSwitchPlayMusicMode"
        />
        <!-- Extension -->
        <div v-if="currentMode === Mode.EXTEN">
          <!-- 音乐分享 -->
          <div v-if="currentExtensionType === ExtensionType.MUSIC">
            <h2 class="text-gray-500 font-bold mb-1">音乐分享（支持网易云/QQ音乐）</h2>
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
              <span v-if="parseMusicURL(extensionToAdd.extension)" class="text-green-400"
                >成功</span
              >
              <span v-else class="text-red-300">失败</span>
            </div>
          </div>
          <!-- Bilibili视频分享 -->
          <div v-if="currentExtensionType === ExtensionType.VIDEO">
            <div class="text-gray-500 font-bold mb-1">Bilibili视频分享(粘贴自动提取BV号)</div>
            <BaseInput
              v-model="bilibiliURL"
              class="rounded-lg h-auto w-full my-2"
              placeholder="B站分享链接..."
            />
            <div class="text-gray-500 my-1">提取的BV号为：{{ extensionToAdd.extension }}</div>
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
      </div>

      <!-- Editor Buttons -->
      <TheEditorButtons
        :echo-to-add="echoToAdd"
        :current-mode="currentMode"
        @handle-addor-update="handleAddorUpdate"
        @handle-change-mode="handleChangeMode"
        @handle-add-image-mode="handleAddImageMode"
        @handle-exit-update-mode="handleExitUpdateMode"
        @handle-private="handlePrivate"
        />

      <!-- Editor Image -->
      <TheEditorImage
        :imagesToAdd="imagesToAdd"
        :current-mode="currentMode"
        @handleAddorUpdateEcho="handleAddorUpdateEcho"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import Upload from '@/components/icons/upload.vue'
import Url from '@/components/icons/url.vue'
import Audio from '@/components/icons/audio.vue'
import Delete from '@/components/icons/delete.vue'
import Addmore from '@/components/icons/addmore.vue'

import BaseButton from '@/components/common/BaseButton.vue'
import BaseInput from '@/components/common/BaseInput.vue'
import BaseTextArea from '@/components/common/BaseTextArea.vue'

import TheMdEditor from '@/components/advanced/TheMdEditor.vue'
import TheModePanel from './TheModePanel.vue'
import TheTitleAndNav from './TheEditor/TheTitleAndNav.vue'
import TheEditorImage from './TheEditor/TheEditorImage.vue'
import TheEditorButtons from './TheEditor/TheEditorButtons.vue'

import { theToast } from '@/utils/toast'
import {  ref, watch } from 'vue'
import {
  fetchUploadImage,
  fetchAddEcho,
  fetchAddTodo,
  fetchUploadMusic,
  fetchDeleteMusic,
  fetchUpdateEcho,
} from '@/service/api'
import { useEchoStore } from '@/stores/echo'
import { useTodoStore } from '@/stores/todo'
import { storeToRefs } from 'pinia'
import { parseMusicURL } from '@/utils/other'
import { Mode, ExtensionType, ImageSource } from '@/enums/enums'

const emit = defineEmits(['refreshAudio'])

const echoStore = useEchoStore()
const todoStore = useTodoStore()

const { setTodoMode, getTodos } = todoStore

const { todoMode } = storeToRefs(todoStore)
const { echoToUpdate, isUpdateMode } = storeToRefs(echoStore)

const currentMode = ref<Mode>(Mode.ECH0)
const currentExtensionType = ref<ExtensionType>()

const handleChangeMode = () => {
  if (currentMode.value === Mode.ECH0) {
    currentMode.value = Mode.Panel
  } else if (
    currentMode.value === Mode.TODO ||
    currentMode.value === Mode.PlayMusic ||
    currentMode.value === Mode.EXTEN
  ) {
    currentMode.value = Mode.Panel
    setTodoMode(false)

    if (!echoToAdd.value.image_url || echoToAdd.value.image_url.length === 0) {
      echoToAdd.value.image_source = null
    }
  } else {
    currentMode.value = Mode.ECH0
  }
}
const handleSwitchExtensionMode = (extensiontype: ExtensionType) => {
  currentMode.value = Mode.EXTEN
  currentExtensionType.value = extensiontype
  extensionToAdd.value.extension_type = extensiontype
}
const handleSwitchTodoMode = () => {
  setTodoMode(true)
  currentMode.value = Mode.TODO
}
const handleSwitchPlayMusicMode = () => {
  currentMode.value = Mode.PlayMusic
}

const websiteToAdd = ref<{
  title: string
  site: string
}>({
  title: '',
  site: '',
}) // 临时网站链接变量
const bilibiliURL = ref<string>('') // 临时Bilibili链接变量
const extensionToAdd = ref({
  extension: '',
  extension_type: '',
}) // 临时扩展变量
const imageIndex = ref<number>(0) // 临时图片索引变量
const imageSourceMemory = ref<string>() // 临时图片来源变量
const imageToAdd = ref<App.Api.Ech0.ImageToAdd>({
  image_url: '',
  image_source: '',
}) // 临时图片添加变量
const imagesToAdd = ref<App.Api.Ech0.ImageToAdd[]>([])
const echoToAdd = ref<App.Api.Ech0.EchoToAdd>({
  content: '',
  image_url: null,
  image_source: null,
  images: [],
  private: false,
  extension: null,
  extension_type: null,
})

const todoToAdd = ref<App.Api.Todo.TodoToAdd>({
  content: '',
})

const handleAddMoreImage = () => {
  imagesToAdd.value.push({
    image_url: String(imageToAdd.value.image_url),
    image_source: String(imageToAdd.value.image_source),
  })

  imageToAdd.value.image_url = ''
  imageToAdd.value.image_source = ''
}

const handleAddImageMode = () => {
  currentMode.value = Mode.Image
}
const fileInput = ref<HTMLInputElement | null>(null)
const handleTriggerUpload = () => {
  imageSourceMemory.value = imageToAdd.value.image_source
  if (fileInput.value) {
    fileInput.value.click()
  }
}
const handleUploadImage = async (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    const file = target.files[0]
    fetchUploadImage(file)
      .then((res) => {
        if (res.code === 1) {
          // 改成新版的图片数组
          // echoToAdd.value.image_url = res.data
          imageToAdd.value.image_url = String(res.data)
          imageToAdd.value.image_source = ImageSource.LOCAL
          theToast.success('图片上传成功！')

          handleAddMoreImage()

          // 如果当前处于Echo更新模式，则需要立马执行更新（图片上传操作不可逆，需要立马更新确保后端数据同步）
          if (isUpdateMode.value && echoToUpdate.value) {
            handleAddorUpdateEcho(true)
          }

          // if (currentMode.value === Mode.Image) {
          //   currentMode.value = Mode.ECH0
          // }
        }
      })
      .finally(() => {
        // 重置文件输入
        if (fileInput.value) {
          fileInput.value.value = ''
        }
      })
  }
}

const handleUploadMusic = async (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    const file = target.files[0]
    fetchUploadMusic(file).then((res) => {
      if (res.code === 1) {
        theToast.success('音乐上传成功！')
        emit('refreshAudio')
      }
    })
  }
}
const handleDeleteMusic = () => {
  if (confirm('确定要删除音乐吗？')) {
    fetchDeleteMusic().then((res) => {
      if (res.code === 1) {
        theToast.success('音乐删除成功！')
        emit('refreshAudio')
      }
    })
  }
}

const handlePrivate = () => {
  echoToAdd.value.private = !echoToAdd.value.private
}

const handleClear = () => {
  echoToAdd.value.content = ''
  echoToAdd.value.image_url = null
  echoToAdd.value.image_source = null
  echoToAdd.value.images = []
  echoToAdd.value.private = false
  echoToAdd.value.extension = null
  echoToAdd.value.extension_type = null
  extensionToAdd.value.extension = ''
  extensionToAdd.value.extension_type = ''
  bilibiliURL.value = ''
  imagesToAdd.value = []
  imageToAdd.value.image_url = ''
  imageToAdd.value.image_source = ''
  imageIndex.value = 0
}

const handleAddorUpdateEcho = (justSyncImages: boolean) => {
  echoToAdd.value.images = imagesToAdd.value // 将图片数组添加到Echo中

  // 检查是否有外部链接分享
  if (extensionToAdd.value.extension_type === ExtensionType.WEBSITE) {
    // 检查是否存在网站链接
    if (websiteToAdd.value.title.length > 0 && websiteToAdd.value.site.length === 0) {
      theToast.error('网站链接不能为空！')
      return
    }

    // 检查是否存在网站标题
    if (websiteToAdd.value.title.length === 0 && websiteToAdd.value.site.length > 0) {
      websiteToAdd.value.title = '外部链接'
    }

    // 检查网站标题和链接是否都存在
    if (websiteToAdd.value.title.length > 0 && websiteToAdd.value.site.length > 0) {
      // 将网站标题和链接添加到扩展中 (序列化为json)
      extensionToAdd.value.extension = JSON.stringify({
        title: websiteToAdd.value.title,
        site: websiteToAdd.value.site,
      })
    } else {
      extensionToAdd.value.extension = ''
      extensionToAdd.value.extension_type = ''
    }
  }

  // 检查最终的Extension模块是否有内容
  if (extensionToAdd.value.extension.length > 0 && extensionToAdd.value.extension_type.length > 0) {
    echoToAdd.value.extension = extensionToAdd.value.extension
    echoToAdd.value.extension_type = extensionToAdd.value.extension_type
  } else {
    echoToAdd.value.extension = null
    echoToAdd.value.extension_type = null
  }

  // 检查Echo是否为空
  if (
    !echoToAdd.value.content &&
    (!echoToAdd.value.images || echoToAdd.value.images.length === 0) &&
    !echoToAdd.value.extension &&
    !echoToAdd.value.extension_type
  ) {
    if (isUpdateMode.value) {
      theToast.error('待更新的Echo不能为空！')
      return
    } else {
      theToast.error('待添加的Echo不能为空！')
      return
    }
  }

  // if (!echoToAdd.value.image_url || echoToAdd.value.image_url.length === 0) {
  //   echoToAdd.value.image_source = null
  // }

  // 检查是否处于更新模式
  if (isUpdateMode.value) {
    // 处于更新模式，执行更新操作
    if (!echoToUpdate.value) {
      theToast.error('没有待更新的Echo！')
      return
    }

    // 回填 echoToUpdate
    echoToUpdate.value.content = echoToAdd.value.content
    echoToUpdate.value.private = echoToAdd.value.private
    echoToUpdate.value.images = echoToAdd.value.images
    echoToUpdate.value.extension = echoToAdd.value.extension
    echoToUpdate.value.extension_type = echoToAdd.value.extension_type

    // 更新Echo
    fetchUpdateEcho(echoToUpdate.value).then((res) => {
      if (res.code === 1 && !justSyncImages) {
        theToast.success('更新成功！')
        handleClear()
        echoStore.refreshEchos()
        isUpdateMode.value = false
        echoToUpdate.value = null
        currentMode.value = Mode.ECH0
      } else if (res.code === 1 && justSyncImages) {
        theToast.success('发现图片更改，已自动更新同步Echo！')
      } else {
        theToast.error('更新失败，请稍后再试！')
      }
    })
    return
  }

  // 不是Echo更新模式，执行添加操作
  fetchAddEcho(echoToAdd.value).then((res) => {
    if (res.code === 1) {
      theToast.success('发布成功！')
      handleClear()
      echoStore.refreshEchos()
      currentMode.value = Mode.ECH0
    }
  })
}

const handleAddTodo = () => {
  if (todoToAdd.value.content === '') {
    theToast.error('内容不能为空！')
    return
  }

  fetchAddTodo(todoToAdd.value).then((res) => {
    if (res.code === 1) {
      theToast.success('添加成功！')
      todoToAdd.value.content = ''
      getTodos()
    }
  })
}

const handleAddorUpdate = () => {
  if (todoMode.value) {
    handleAddTodo()
  } else {
    handleAddorUpdateEcho(false)
  }
}

const handleExitUpdateMode = () => {
  isUpdateMode.value = false
  echoToUpdate.value = null
  handleClear()
  theToast.info('已退出更新模式！')
}

// 监听用户输入
watch(
  () => bilibiliURL.value,
  (newVal) => {
    if (newVal.length > 0) {
      const bvRegex = /(BV[0-9A-Za-z]{10})/
      const match = newVal.match(bvRegex)
      if (match) {
        extensionToAdd.value.extension = match[0]
      } else {
        theToast.error('请输入正确的B站分享链接！')
      }
    }
  },
)

// 监听是否处于更新模式
watch(
  () => isUpdateMode.value,
  (newVal) => {
    if (newVal) {
      // 处于更新模式（将待更新的数据填充到编辑器）
      // 0. 清空编辑器
      handleClear()

      // 1. 填充本文
      echoToAdd.value.content = echoToUpdate.value?.content || ''
      echoToAdd.value.private = echoToUpdate.value?.private || false

      // 2. 填充图片
      if (echoToUpdate.value?.images && echoToUpdate.value.images.length > 0) {
        imagesToAdd.value = echoToUpdate.value.images.map((img) => ({
          image_url: img.image_url || '',
          image_source: img.image_source || '',
        }))
      } else {
        imagesToAdd.value = []
      }

      // 3. 填充扩展
      if (echoToUpdate.value?.extension && echoToUpdate.value.extension_type) {
        currentExtensionType.value = echoToUpdate.value.extension_type as ExtensionType
        extensionToAdd.value.extension = echoToUpdate.value.extension
        extensionToAdd.value.extension_type = echoToUpdate.value.extension_type
        // 根据扩展类型填充
        switch (echoToUpdate.value.extension_type) {
          case ExtensionType.MUSIC:
            break

          case ExtensionType.VIDEO:
            bilibiliURL.value = echoToUpdate.value.extension // 直接使用extension填充B站链接
            break

          case ExtensionType.GITHUBPROJ:
            break

          case ExtensionType.WEBSITE:
            // 反序列化网站链接
            const websiteData = JSON.parse(echoToUpdate.value.extension) as {
              title?: string
              site?: string
            }
            websiteToAdd.value.title = websiteData.title || ''
            websiteToAdd.value.site = websiteData.site || ''
            break
        }
      }

      // 4. 回到页面顶部（触发BackToTop）
      window.scrollTo({ top: 0, behavior: 'smooth' })

      // 5. 弹出通知，提示可以编辑了
      theToast.info('已进入更新模式，请编辑内容后点击更新按钮！')
    } else {
      // 退出更新模式
      echoToUpdate.value = null
    }
  },
)
</script>
