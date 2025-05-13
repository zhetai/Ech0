import { ref } from "vue";
import { defineStore } from "pinia";
import { fetchGetConnects, fetchGetConnect } from "@/service/api";

export const useConnectStore = defineStore('connectStore', () => {
  /**
   * State
   */
  const connects = ref<App.Api.Connect.Connected[]>([])
  const connectsInfo = ref<App.Api.Connect.Connect[]>([])

  /**
   * Actions
   */
  async function getConnect() {
    fetchGetConnects()
      .then((res) => {
        if (res.code === 1) {
          connects.value = res.data
        }
      })
      .catch((err) => {
        console.error(err)
      })
  }

  async function getConnectInfo() {
    await getConnect()
    // 根据connects的url获取connect的详细信息
    const connecteds = ref<App.Api.Connect.Connect[]>([])
    connects.value.forEach((connect) => {
      fetchGetConnect(connect.connect_url)
        .then((res) => {
          if (res.code === 1) {
            connecteds.value.push(res.data)
          }
        })
        .catch((err) => {
          console.error(err)
        })
    })
    connectsInfo.value = connecteds.value
  }

  return {
    connects,
    connectsInfo,
    getConnect,
    getConnectInfo
  }
})
