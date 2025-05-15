import { defineStore } from 'pinia'
import { ref } from 'vue'
import { fetchGetSettings, fetchGetStatus } from '@/service/api'
import { localStg } from '@/utils/storage'
import { theToast } from '@/utils/toast'
import { useRouter } from 'vue-router'

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
  })
  const router = useRouter()

  /**
   * Actions
   */
  const getSystemSetting = async () => {
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
      await fetchGetStatus().then((res) => {
        if (res.code === 666) {
          isSystemReady.value = false
          theToast.info(res.msg)
          // 跳转到登录页面
          router.push({ name: 'auth' })
        } else {
          isSystemReady.value = true
        }
      }).finally(() => {
        // 保存系统状态到localStorage
        localStg.setItem('systemStatus', isSystemReady.value)
      })
    }

    // 获取系统设置
    await fetchGetSettings().then((res) => {
      if (res.code === 1) {
        SystemSetting.value = res.data
      }
    })
  }

  const setSystemReady = (status :boolean) => {
    isSystemReady.value = status
  }

  return { isSystemReady, SystemSetting, getSystemSetting, setSystemReady }
})
