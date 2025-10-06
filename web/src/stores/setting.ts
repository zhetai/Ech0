import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
  fetchGetSettings,
  fetchGetStatus,
  fetchGetCommentSettings,
  fetchGetS3Settings,
  fetchGetOAuth2Settings
} from '@/service/api'
import { localStg } from '@/utils/storage'
import { theToast } from '@/utils/toast'
import { useRouter } from 'vue-router'
import { CommentProvider, S3Provider, OAuth2Provider } from '@/enums/enums'

export const useSettingStore = defineStore('settingStore', () => {
  /**
   * State
   */
  const isSystemReady = ref<boolean>(false)
  const SystemSetting = ref<App.Api.Setting.SystemSetting>({
    site_title: import.meta.env.VITE_APP_TITLE,
    server_name: import.meta.env.VITE_APP_NAME,
    server_url: '',
    allow_register: true,
    ICP_number: '',
    meting_api: '',
    custom_css: '',
    custom_js: '',
  })
  const CommentSetting = ref<App.Api.Setting.CommentSetting>({
    enable_comment: false,
    provider: CommentProvider.TWIKOO,
    comment_api: '',
  })
  const S3Setting = ref<App.Api.Setting.S3Setting>({
    enable: false,
    provider: S3Provider.AWS,
    endpoint: '',
    access_key: '',
    secret_key: '',
    bucket_name: '',
    region: '',
    use_ssl: false,
    cdn_url: '',
    path_prefix: '',
    public_read: true,
  })
  const OAuth2Setting = ref<App.Api.Setting.OAuth2Setting>({
    enable: false,
    provider: OAuth2Provider.GITHUB,
    client_id: '',
    client_secret: '',
    redirect_uri: '',
    scopes: [],
    auth_url: '',
    token_url: '',
    user_info_url: '',
  })
  const loading = ref<boolean>(true)
  const router = useRouter()

  /**
   * Actions
   */
  const getSystemReady = async () => {
    // 检查localStorage中是否有系统状态
    const systemStatus = localStg.getItem<boolean>('systemStatus')
    if (systemStatus !== null) {
      // 如果有，直接使用localStorage中的值
      isSystemReady.value = systemStatus
    } else {
      // 如果没有，默认设置为false
      isSystemReady.value = false
    }

    // 检查系统是否准备好
    if (!isSystemReady.value) {
      // 如果系统未准备好，调用接口获取系统状态
      await fetchGetStatus()
        .then((res) => {
          if (res.code === 666) {
            isSystemReady.value = false
            theToast.info(res.msg)
            // 跳转到登录页面
            router.push({ name: 'auth' })
          } else {
            isSystemReady.value = true
            console.log('系统已准备好')
          }
        })
        .finally(() => {
          // 保存系统状态到localStorage
          localStg.setItem('systemStatus', isSystemReady.value)
        })
    }
  }

  const getSystemSetting = async () => {
    await fetchGetSettings().then((res) => {
      if (res.code === 1) {
        SystemSetting.value = res.data
        loading.value = false
      }
    })
  }

  const getCommentSetting = async () => {
    fetchGetCommentSettings().then((res) => {
      if (res.code === 1) {
        CommentSetting.value = res.data
      }
    })
  }

  const getS3Setting = async () => {
    fetchGetS3Settings().then((res) => {
      if (res.code === 1) {
        S3Setting.value = res.data
      }
    })
  }

  const getOAuth2Setting = async () => {
    fetchGetOAuth2Settings().then((res) => {
      if (res.code === 1) {
        OAuth2Setting.value = res.data
      }
    })
  }

  const setSystemReady = (status: boolean) => {
    isSystemReady.value = status
  }

  const init = () => {
    if (!isSystemReady.value) {
      getSystemReady()
    }
    getSystemSetting()
    getCommentSetting()
    getS3Setting()
  }

  return {
    isSystemReady,
    SystemSetting,
    CommentSetting,
    S3Setting,
    OAuth2Setting,
    loading,
    getSystemReady,
    getSystemSetting,
    getCommentSetting,
    getS3Setting,
    getOAuth2Setting,
    setSystemReady,
    init,
  }
})
