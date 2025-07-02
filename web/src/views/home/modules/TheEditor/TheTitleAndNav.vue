<template>
  <div class="flex justify-between items-center py-1 px-3">
        <div class="flex flex-row items-center gap-2 justify-between">
          <!-- <div class="text-xl">👾</div> -->
          <div>
            <img
              :src="logo"
              alt="logo"
              class="w-6 sm:w-7 h-6 sm:h-7 rounded-full ring-1 ring-gray-200 shadow-sm object-cover"
            />
          </div>
          <h1 class="text-slate-600 font-bold italic sm:text-xl">
            {{ SystemSetting.server_name }}
          </h1>
        </div>
        <div class="flex flex-row items-center gap-2">
          <!-- Github -->
          <div>
            <a href="https://github.com/lin-snow/Ech0" target="_blank" title="Github">
              <Github class="w-6 sm:w-7 h-6 sm:h-7 text-gray-400" />
            </a>
          </div>
        </div>
      </div>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import Github from '@/components/icons/github.vue'
import { onMounted, ref } from 'vue';
import { fetchGetStatus } from '@/service/api'
import { useSettingStore } from '@/stores/settting';
import { getApiUrl } from '@/service/request/shared';

const settingStore = useSettingStore()

const { SystemSetting } = storeToRefs(settingStore)

const apiUrl = getApiUrl()
const logo = ref<string>('/favicon.svg')

onMounted(() => {
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

<style scoped>

</style>
