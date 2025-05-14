import { ref } from 'vue'
import { defineStore } from 'pinia'
import { fetchGetConnects, fetchGetConnect } from '@/service/api'

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
    await fetchGetConnects()
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
    // 根据 connects 的 url 获取 connect 的详细信息
    const promises = connects.value.map((connect) =>
      fetchGetConnect(connect.connect_url)
        .then((res) => {
          if (res.code === 1) {
            return res.data
          }
          return null
        })
        .catch((err) => {
          console.error(err)
          return null
        }),
    )
    const results = await Promise.all(promises)
    connectsInfo.value = results.filter((item) => item !== null) as App.Api.Connect.Connect[]
  }

  return {
    connects,
    connectsInfo,
    getConnect,
    getConnectInfo,
  }
})
