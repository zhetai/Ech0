import { defineStore } from "pinia";
import { ref } from "vue";
import { fetchGetSettings } from "@/service/api";

export const useSettingStore = defineStore("settingStore", () => {
  /**
   * State
   */
  const SystemSetting = ref<App.Api.Setting.SystemSetting>({
    site_title: import.meta.env.VITE_APP_TITLE,
    server_name: import.meta.env.VITE_APP_NAME,
    server_url: "",
    allow_register: true,
    ICP_number: ""
  })

  /**
   * Actions
   */
  const getSystemSetting = async () => {
    await fetchGetSettings()
      .then((res) => {
        if (res.code === 1) {
          SystemSetting.value = res.data
        }
      })
  }

  return { SystemSetting, getSystemSetting }
})
