<template>
  <div class="max-w-sm flex justify-around items-center bg-white rounded-lg shadow-sm p-2 gap-2">
    <img v-if="CardData?.owner?.avatar_url" :src="CardData?.owner?.avatar_url" alt="头像" class="w-10 h-10 rounded-full shadow" />
    <Githubproj v-else class="w-10 h-10 shadow" />
    <div class="px-3">
      <a :href="props.GithubURL" target="_blank">
        <span class="text-md font-bold text-gray-500">{{ CardData?.name || repo }}</span>
        <div v-if="CardData" class="">
          <p
            class="text-sm text-gray-400"
            :title="CardData?.description"
          >
            {{ CardData?.description }}
          </p>
          <div class="flex justify-start items-center h-auto text-gray-500">
            <!-- star -->
            <Star class="w-4 h-4 mr-1" /> <span> {{ CardData?.stargazers_count }} </span>
            <!-- fork -->
            <Fork class="w-4 h-4 mx-1" /> <span> {{ CardData?.forks_count }} </span>
          </div>
        </div>
      </a>
    </div>
  </div>
</template>

<script setup lang="ts">
import Githubproj from '../icons/githubproj.vue'
import Star from '../icons/star.vue';
import Fork from '../icons/fork.vue';
import { fetchGetGithubRepo } from '@/service/api'
import { onMounted, ref } from 'vue'
const props = defineProps<{
  GithubURL: string
}>()

// 处理GithubURL(提取owner和repo)
const [owner, repo] = props.GithubURL.split('/').slice(-2)
const CardData = ref<App.Api.Ech0.GithubCardData>()

onMounted(async () => {
  await fetchGetGithubRepo({ owner, repo }).then((res) => {
    if (res) {
      CardData.value = res
      console.log('CardData', CardData.value)
    }
  })
})
</script>

<style scoped></style>
