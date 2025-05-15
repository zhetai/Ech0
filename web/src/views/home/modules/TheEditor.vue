<template>
  <div
    class="bg-white rounded-lg ring-1 ring-gray-200 ring-inset mx-auto shadow-sm hover:shadow-md"
  >
    <div class="mx-auto w-full px-3 py-4">
      <!-- Title && Nav -->
      <div class="flex justify-between items-center py-1 px-2">
        <div class="flex flex-row items-center gap-2 justify-between">
          <!-- <div class="text-xl">👾</div> -->
          <div>
            <img
              :src="logo"
              alt="logo"
              class="w-6 md:w-7 h-6 md:h-7 rounded-full ring-1 ring-gray-200 shadow-sm object-cover"
            />
          </div>
          <h2 class="text-slate-600 font-bold italic">{{ SystemSetting.server_name }}</h2>
        </div>
        <div class="flex flex-row items-center gap-2">
          <!-- Github -->
          <div>
            <a href="https://github.com/lin-snow/Ech0" target="_blank" title="Github">
              <Github class="w-6 h-6 text-gray-400" />
            </a>
          </div>
        </div>
      </div>

      <!-- Editor -->
      <div class="rounded-lg p-2 mb-1">
        <TheMdEditor v-model="echoToAdd.content" class="rounded-lg" v-if="!todoMode" />
        <BaseTextArea
          v-else
          v-model="todoToAdd.content"
          class="rounded-lg h-auto sm:min-h-[6rem] md:min-h-[9rem]"
          placeholder="请输入待办事项..."
          :rows="3"
        />
      </div>

      <!-- Buttons -->
      <div class="flex flex-row items-center justify-between px-2">
        <div class="flex flex-row items-center gap-2">
          <!-- Todo -->
          <div>
            <BaseButton
              :icon="Todo"
              @click="handleChangeMode"
              :class="
                [
                  'w-8 h-8 sm:w-9 sm:h-9 rounded-md',
                  todoMode
                    ? 'bg-orange-100 shadow-[0_0_12px_-4px_rgba(255,140,0,0.6)] !ring-0 !text-white'
                    : '',
                ].join(' ')
              "
              title="切换待办模式"
            />
          </div>
          <!-- Photo Upload -->
          <div v-if="!todoMode">
            <input
              id="file-input"
              class="hidden"
              type="file"
              accept="image/*"
              ref="fileInput"
              @change="handleUploadImage"
            />
            <BaseButton
              :icon="ImageUpload"
              @click="handTriggerUpload"
              class="w-8 h-8 sm:w-9 sm:h-9 rounded-md"
              title="上传图片"
            />
          </div>
          <!-- Privacy Set -->
          <div v-if="!todoMode">
            <BaseButton
              :icon="echoToAdd.private ? Private : Public"
              @click="handlePrivate"
              class="w-8 h-8 sm:w-9 sm:h-9 rounded-md"
              title="是否私密"
            />
          </div>
        </div>

        <div class="flex flex-row items-center gap-2">
          <!-- Clear -->
          <!-- <div>
            <BaseButton
              :icon="Clear"
              @click="handleClear"
              class="w-8 h-8 rounded-md"
              title="清空输入和图片"
            />
          </div> -->
          <!-- Publish -->
          <div>
            <BaseButton
              :icon="Publish"
              @click="handleAdd"
              class="w-8 h-8 sm:w-9 sm:h-9 rounded-md"
              title="发布"
            />
          </div>
        </div>
      </div>

      <!-- Preview Image -->
      <div v-if="echoToAdd.image_url" class="relative rounded-lg shadow-lg w-5/6 mx-auto my-7">
        <button
          @click="handleRemoveImage"
          class="absolute -top-3 -right-4 bg-red-100 hover:bg-red-300 text-gray-600 rounded-lg w-7 h-7 flex items-center justify-center shadow"
          title="移除图片"
        >
          <Close class="w-4 h-4" />
        </button>
        <div class="rounded-lg overflow-hidden">
          <a :href="`${apiUrl}${echoToAdd.image_url}`" data-fancybox>
            <img
              :src="`${apiUrl}${echoToAdd.image_url}`"
              alt="Image"
              class="max-w-full object-cover"
              loading="lazy"
            />
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Github from '@/components/icons/github.vue'
import Rss from '@/components/icons/rss.vue'
import Close from '@/components/icons/close.vue'
import Todo from '@/components/icons/todo.vue'
import Panel from '@/components/icons/panel.vue'
import ImageUpload from '@/components/icons/image.vue'
import Public from '@/components/icons/public.vue'
import Private from '@/components/icons/private.vue'
import Clear from '@/components/icons/clear.vue'
import Publish from '@/components/icons/publish.vue'
import BaseButton from '@/components/common/BaseButton.vue'
import TheMdEditor from '@/components/advanced/TheMdEditor.vue'
import { theToast } from '@/utils/toast'
import { Fancybox } from '@fancyapps/ui'
import { onMounted, ref } from 'vue'
import { fetchUploadImage, fetchAddEcho, fetchGetStatus, fetchAddTodo, fetchDeleteImage } from '@/service/api'
import { getApiUrl } from '@/service/request/shared'
import { useEchoStore } from '@/stores/echo'
import { useSettingStore } from '@/stores/settting'
import { useTodoStore } from '@/stores/todo'
import '@fancyapps/ui/dist/fancybox/fancybox.css'
import { storeToRefs } from 'pinia'
import BaseTextArea from '@/components/common/BaseTextArea.vue'

const apiUrl = getApiUrl()
const echoStore = useEchoStore()
const todoStore = useTodoStore()
const settingStore = useSettingStore()

const { setTodoMode, getTodos } = todoStore
const { SystemSetting } = storeToRefs(settingStore)
const { todoMode } = storeToRefs(todoStore)

const logo = ref<string>('/favicon.svg')

const handleChangeMode = () => {
  setTodoMode(!todoMode.value)
}

const echoToAdd = ref<App.Api.Ech0.EchoToAdd>({
  content: '',
  image_url: null,
  private: false,
})

const todoToAdd = ref<App.Api.Todo.TodoToAdd>({
  content: '',
})

const fileInput = ref<HTMLInputElement | null>(null)
const handTriggerUpload = () => {
  if (fileInput.value) {
    fileInput.value.click()
  }
}
const handleUploadImage = async (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    const file = target.files[0]
    fetchUploadImage(file).then((res) => {
      if (res.code === 1) {
        echoToAdd.value.image_url = res.data
        theToast.success('图片上传成功！')
      }
    })
  }
}

const handleRemoveImage = () => {
  if (confirm('确定要移除图片吗？')) {
    fetchDeleteImage({
      url: echoToAdd.value.image_url ?? '',
    }).then((res) => {
      if (res.code === 1) {
        // theToast.success('图片已移除')
      }
    })
    .finally(() => {
      echoToAdd.value.image_url = null
    })
  }
}

const handlePrivate = () => {
  echoToAdd.value.private = !echoToAdd.value.private
}

const handleClear = () => {
  echoToAdd.value.content = ''
  echoToAdd.value.image_url = null
  echoToAdd.value.private = false
}

const handleAddEcho = () => {
  if (!echoToAdd.value.content && !echoToAdd.value.image_url) {
    theToast.error('内容不能为空！')
    return
  }

  fetchAddEcho(echoToAdd.value).then((res) => {
    if (res.code === 1) {
      theToast.success('发布成功！')
      handleClear()
      echoStore.refreshEchos()
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

const handleAdd = () => {
  if (todoMode.value) {
    handleAddTodo()
  } else {
    handleAddEcho()
  }
}

onMounted(() => {
  Fancybox.bind('[data-fancybox]', {})
  fetchGetStatus().then((res) => {
    if (res.code === 1) {
      const theLogo = res.data.logo
      if (theLogo && theLogo !== '') {
        logo.value = `${apiUrl}${theLogo}`
      }
    }
  })
})
</script>
