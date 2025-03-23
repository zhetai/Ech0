<template>
  <div class="mx-auto my-4 sm:max-w-sm px-2">
    <div
      v-for="msg in message.messages"
      :key="msg.id"
      class="w-full h-auto overflow-hidden flex flex-col justify-between"
    >
      <div class="flex justify-between items-center">
        <!-- 时间 -->
        <div class="flex justify-start items-center h-auto">
          <div class="w-2 h-2 rounded-full bg-orange-600 mr-2"></div>
          <!-- 保留年月日 -->
          <div class="flex justify-start text-sm text-orange-500">
            {{ formatDate(msg.created_at) }}
          </div>
        </div>
        <!-- 编辑按钮 -->
        <div v-if="isLogin" class="w-5 h-5" @click="deleteMsg(msg.id)">
          <UIcon name="i-mdi-paper-roll-outline" class="text-gray-400" />
        </div>
      </div>

      <div class="border-l-2 border-gray-300 p-6 ml-1">
        <!-- 留言内容 -->
        <div class="rounded-lg overflow-hidden mb-2 w-5/6 mx-auto shadow-lg">
          <a :href="`${BASE_API}${msg.image_url}`" data-fancybox>
            <img
              v-if="msg.image_url"
              :src="`${BASE_API}${msg.image_url}`"
              alt="Image"
              class="max-w-full object-cover"
              loading="lazy"
            />
          </a>
        </div>
        <!-- 留言 -->
        <div class="p-2 mb-2 rounded-lg h-3/5 overflow-y-auto">
          <!-- <div v-html="renderMarkdown(msg.content)" class="text-gray-900"></div> -->
          <MarkdownRenderer :content="msg.content" />
        </div>
        <!-- 昵称 -->
        <!-- <div class="p-1 h-1/3 flex items-center justify-end">
          <span class="text-xs text-gray-500 italic underline max-w-20">
            {{ msg.username ? msg.username : "匿名" }}
          </span>
        </div> -->
      </div>
    </div>
  </div>

  <!-- 加载更多 -->
  <div v-if="message.hasMore" class="mx-auto -mt-3 max-w-sm">
    <UButton
      color="gray"
      variant="outline"
      size="sm"
      class="rounded-full border-gray-200"
      @click="message.getMessages({ page: message.page + 1, pageSize: 10 })"
    >
      继续装填！
    </UButton>
  </div>
  <!-- 没有啦~ -->
  <div v-else class="text-center text-gray-500 mt-4">
    <UIcon name="i-fluent-emoji-flat-confetti-ball" size="lg" />
    没有啦~
  </div>
</template>

<script setup lang="ts">
import { Fancybox } from "@fancyapps/ui";
import "@fancyapps/ui/dist/fancybox/fancybox.css";
import { useMessageStore } from "~/store/message";
import { useUserStore } from "~/store/user";
import MarkdownRenderer from "~/components/index/MarkdownRenderer.vue"; // 导入组件

const BASE_API = useRuntimeConfig().public.baseApi;
const { deleteMessage } = useMessage();
const message = useMessageStore();
const isLogin = ref<boolean>(false);

const deleteMsg = (id: number) => {
  // 显示确认框
  const confirmDelete = confirm("确定要删除这条消息吗？");
  if (confirmDelete) {
    deleteMessage(id);
  }
};

const formatDate = (dateString: string) => {
  // 当天则显示（时：分）
  // 非当天但是三内天则显示几天前
  // 超过三天则显示（时：分 年月日）
  const date = new Date(dateString);
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  const diffInDays = Math.floor(diff / (1000 * 60 * 60 * 24));
  const diffInHours = Math.floor(diff / (1000 * 60 * 60));
  const diffInMinutes = Math.floor(diff / (1000 * 60));

  const diffInSeconds = Math.floor(diff / 1000);
  if (diffInSeconds < 60) {
    return "刚刚";
  } else if (diffInMinutes < 60) {
    return `${diffInMinutes}分钟前`;
  } else if (diffInHours < 24) {
    return `${diffInHours}小时前`;
  } else if (diffInDays < 3) {
    return `${diffInDays}天前`;
  } else {
    return date.toLocaleString(); // 返回完整的日期和时间
  }
};

onMounted(() => {
  isLogin.value = useUserStore()?.isLogin;

  message.getMessages({
    page: 1,
    pageSize: 10,
  });
  Fancybox.bind("[data-fancybox]", {});
});

onBeforeUnmount(() => {
  Fancybox.destroy();
});
</script>
