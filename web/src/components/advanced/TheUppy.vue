<template>
  <!-- Uppy Dashboard 容器 -->
  <div id="uppy-dashboard">
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { getAuthToken } from '@/service/request/shared'
import { useUserStore } from '@/stores/user';
import { theToast } from '@/utils/toast';
import { storeToRefs } from 'pinia';
/* --------------- 与Uppy相关 ---------------- */
import Uppy from '@uppy/core';
import Dashboard from '@uppy/dashboard';
import XHRUpload from '@uppy/xhr-upload';
import StatusBar from '@uppy/status-bar';
import '@uppy/core/css/style.min.css';
import '@uppy/dashboard/css/style.min.css';
import zh_CN from '@uppy/locales/lib/zh_CN'


let uppy: Uppy | null = null

const emit = defineEmits(["uppyUploaded"])

const files = ref<string[]>([]);

const userStore = useUserStore();
const { isLogin } = storeToRefs(userStore);

onMounted(() => {
  uppy = new Uppy({
    restrictions: {
      maxNumberOfFiles: 6,
      allowedFileTypes: ['image/*'],
    },
    autoProceed: true,
  })

  uppy.use(Dashboard, {
    inline: true,
    target: '#uppy-dashboard',
    hideProgressDetails: false,
    hideUploadButton: false,
    hideCancelButton: false,
    hideRetryButton: false,
    hidePauseResumeButton: false,
    proudlyDisplayPoweredByUppy: false,
    height: 200,
    locale: zh_CN,
    note: '',
  })

  uppy.use(XHRUpload, {
    endpoint: 'http://localhost:6277/api/images/upload', // 换成你的后端上传接口
    fieldName: 'file',
    formData: true,
    headers: {
      "Authorization": `${getAuthToken()}`
    }
  })

  // uppy.on("file-added", file => {})
  uppy.on("upload", () => {
    if (!isLogin.value) {
      theToast.error("请先登录再上传图片 😢")
      uppy?.cancelAll()
      return
    }
    theToast.info("正在上传图片，请稍等... ⏳")
  })
  uppy.on("upload-error", (file, error, response) => {
  type ResponseBody = {
    code: number;
    msg: string;
    data: any;
  };

  let errorMsg = "上传图片时发生错误 😢";

  const resp = response as any; // 忽略 TS 类型限制

  if (resp?.response) {
    let resObj: ResponseBody;

    if (typeof resp.response === "string") {
      resObj = JSON.parse(resp.response) as ResponseBody;
    } else {
      resObj = resp.response as ResponseBody;
    }

    if (resObj?.msg) {
      errorMsg = resObj.msg;
    }
  }

  theToast.error(errorMsg);
});

  uppy.on("upload-success", (file, response) => {
    theToast.success(`好耶,上传成功！🎉`)
    const fileUrl = String(response.body?.data);
    files.value.push(fileUrl);
  });
  uppy.on("complete", () => {
    emit("uppyUploaded", files.value); // 发射事件到父组件
  })
})

onBeforeUnmount(() => {
  uppy?.destroy()
})
</script>

<style scoped>

</style>
