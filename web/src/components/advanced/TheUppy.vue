<template>
  <!-- Uppy Dashboard 容器 -->
  <div id="uppy-dashboard" class="rounded-md overflow-hidden shadow-inner ring-inset ring-1 ring-gray-200">
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onBeforeUnmount } from 'vue'
import { getAuthToken } from '@/service/request/shared'
import { useUserStore } from '@/stores/user';
import { theToast } from '@/utils/toast';
import { storeToRefs } from 'pinia';
import { ImageSource } from '@/enums/enums';
import { fetchGetPresignedUrl } from '@/service/api';
/* --------------- 与Uppy相关 ---------------- */
import Uppy from '@uppy/core';
import Dashboard from '@uppy/dashboard';
import XHRUpload from '@uppy/xhr-upload';
import AwsS3, { type AwsBody } from '@uppy/aws-s3';
import '@uppy/core/css/style.min.css';
import '@uppy/dashboard/css/style.min.css';
import zh_CN from '@uppy/locales/lib/zh_CN'

let uppy: Uppy | null = null

const props = defineProps<{
  TheImageSource: string
}>()
const emit = defineEmits(["uppyUploaded"])

const isUploading = ref<boolean>(false); // 是否正在上传
const files = ref<App.Api.Ech0.ImageToAdd[]>([]); // 已上传的文件列表
const tempFiles = ref<Map<string, string>>(new Map()); // 用于S3临时存储文件回显地址的 Map(key: fileName, value: url)

const userStore = useUserStore();
const { isLogin } = storeToRefs(userStore);
const envURL = import.meta.env.VITE_SERVICE_BASE_URL as string
const backendURL = envURL.endsWith('/') ? envURL.slice(0, -1) : envURL

// ✨ 监听粘贴事件
const handlePaste = (e: ClipboardEvent) => {
  if (!e.clipboardData) return

  for (const item of e.clipboardData.items) {
    if (item.type.startsWith("image/")) {
      const file = item.getAsFile()
      if (file) {
        uppy?.addFile({
          name: `pasted-${Date.now()}.png`,
          type: file.type,
          data: file,
          source: "PastedImage",
        })
        uppy?.upload()
      }
    }
  }
}

// 初始化 Uppy 实例
const initUppy = () => {
  // 创建 Uppy 实例
  uppy = new Uppy({
    restrictions: {
      maxNumberOfFiles: 6,
      allowedFileTypes: ['image/*'],
    },
    autoProceed: true,
  })

  // 使用 Dashboard 插件
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
    note: '支持粘贴或选择图片上传哦！',
  })

  // 根据 props.TheImageSource 动态切换上传插件
  if (props.TheImageSource === ImageSource.LOCAL) {
    uppy.use(XHRUpload, {
      endpoint: `${backendURL}/api/images/upload`, // 本地上传接口
      fieldName: 'file',
      formData: true,
      headers: {
        "Authorization": `${getAuthToken()}`
      }
    });
  } else if (props.TheImageSource === ImageSource.S3) {
    uppy.use(AwsS3, {
      endpoint: '', // 走自定义的签名接口
      shouldUseMultipart: false, // 禁用分块上传
      // 每来一个文件都调用一次该函数，获取签名参数
      async getUploadParameters(file) {
        // console.log("Uploading to S3:", file)
        const fileName = file.name ? file.name : ''
        const contentType = file.type ? file.type : ''
        // console.log("fileName, contentType", fileName, contentType)

        const res = await fetchGetPresignedUrl(fileName, contentType)
        if (res.code !== 1) {
          throw new Error(res.msg || '获取预签名 URL 失败')
        }
        const data = res.data as App.Api.Ech0.PresignResult
        tempFiles.value.set(data.file_name, data.file_url)

        return {
          method: 'PUT',
          url: data.presign_url, // 预签名 URL
          headers: {
            // 必须跟签名时的 Content-Type 完全一致
            'Content-Type': file.type
          },
          // PUT 上传没有 fields
          fields: {}
        }
      }
    });
  }

  // 监听粘贴事件
  document.addEventListener("paste", handlePaste)

  // 上传开始前，检查是否登录
  uppy.on("upload", (uploadID, files) => {
    if (!isLogin.value) {
      theToast.error("请先登录再上传图片 😢")
      uppy?.cancelAll()
      return
    }
    theToast.info("正在上传图片，请稍等... ⏳", { duration: 1000})
    isUploading.value = true;
  })
  // 单个文件上传失败后，显示错误信息
  uppy.on("upload-error", (file, error, response) => {
    if (props.TheImageSource === ImageSource.LOCAL) {
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
    }
    isUploading.value = false;
  });
  // 单个文件上传成功后，保存文件 URL 到 files 列表
  uppy.on("upload-success", (file, response) => {
    theToast.success(`好耶,上传成功！🎉`)
    // console.log("Upload success", file, response);
    // 分两种情况: Local 或者 S3
    if (props.TheImageSource === ImageSource.LOCAL) {
      const fileUrl = String(response.body?.data);
      const item = {
        image_url: fileUrl,
        image_source: ImageSource.LOCAL
      }
      files.value.push(item);
    } else if (props.TheImageSource === ImageSource.S3) {
      const fileUrl = tempFiles.value.get(file?.name || '') || '';
      if (fileUrl) {
        const item = {
          image_url: fileUrl,
          image_source: ImageSource.S3
        }
        files.value.push(item);
      }
    }
  });
  // 全部文件上传完成后，发射事件到父组件
  uppy.on("complete", () => {
    isUploading.value = false;
    emit("uppyUploaded", files.value); // 发射事件到父组件
  })
}

// 监听 props.TheImageSource 变化
watch(
  () => props.TheImageSource,
  (newSource, oldSource) => {
    if ((newSource !== oldSource) && (isUploading.value === false)) {
      // 销毁旧的 Uppy 实例
      uppy?.destroy()
      uppy = null
      files.value = [] // 清空已上传文件列表
      // 初始化新的 Uppy 实例
      initUppy();
    } else if ((newSource !== oldSource) && (isUploading.value === true)) {
      theToast.warning("图片正在上传中，请稍后再切换上传方式 😢")
    }
  }
);

onMounted(() => {
  initUppy();
})

onBeforeUnmount(() => {
  document.removeEventListener("paste", handlePaste)
  uppy?.destroy()
})
</script>

<style scoped>
:deep(.uppy-Root) {
  border: transparent;
}

:deep(.uppy-Dashboard-innerWrap) {
  background-color: #f4f1ec;
}

:deep(.uppy-Dashboard-AddFiles) {
  /* background-color: #fff; */
  /* 内阴影 */
  box-shadow: inset 0px 0px 2px rgba(80, 80, 80, 0.12), inset 0px 0px 2px rgba(80, 80, 80, 0.12);
}

:deep(.uppy-Dashboard-AddFiles-title) {
  color: #6f5427;
}

:deep(.uppy-Dashboard-browse) {
  color: #e5a437;
}

:deep(.uppy-DashboardContent-bar) {
  background-color: #fff;
}

:deep(.uppy-DashboardContent-back) {
  color: #cf8e12;
}

:deep(.uppy-DashboardContent-addMore) {
  color: #cf8e12;
}
</style>
