import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import { fetchGetEchosByPage } from '@/service/api'

export const useEchoStore = defineStore('echoStore', () => {
  /**
   * state
   */
  const echoList = ref<App.Api.Ech0.Echo[]>([]) // 存储Echo列表
  const echoMap = ref(new Map<number, App.Api.Ech0.Echo>()) // 存储Echo的Map，便于快速查找
  const isLoading = ref<boolean>(true) // 是否正在加载数据
  const total = ref<number>(0) // 总数据量
  const pageSize = ref<number>(5) // 每页显示的数量
  const page = ref<number>(0) // 当前页码，从0开始计数
  const current = ref<number>(1) // 当前页码，从1开始计数
  const searchValue = ref<string>('') // 搜索关键词
  const hasMore = computed(() => {
    return total.value > echoList.value.length
  }) // 是否还有更多数据可加载
  const searchingMode = computed(() => {
    return searchValue.value.length > 0
  }) // 是否处于搜索模式
  const echoToUpdate = ref<App.Api.Ech0.EchoToUpdate | null>(null) // 用于更新Echo的临时存储
  const isUpdateMode = ref<boolean>(false) // 是否处于更新模式

  // 监听 searchingMode 的变化
  watch(searchingMode, (newValue, oldValue) => {
    // 如果从搜索模式切换到非搜索模式，重置当前页码和数据列表
    if (newValue === false && oldValue === true) {
      refreshEchos()
    }
  })

  /**
   * actions
   */
  async function getEchosByPage() {
    if (current.value <= page.value) return

    await fetchGetEchosByPage({
      page: current.value,
      pageSize: pageSize.value,
      search: searchValue.value,
    })
      .then((res) => {
        if (res.code === 1) {
          total.value = res.data.total
          echoList.value = [...echoList.value, ...res.data.items]

          // 同步更新 echoMap
          res.data.items.forEach((item: App.Api.Ech0.Echo) => {
            echoMap.value.set(item.id, item)
          })

          page.value += 1
        }
      })
      .finally(() => {
        isLoading.value = false
      })
  }

  const refreshEchos = () => {
    current.value = 1
    page.value = 0
    echoList.value = []
    echoMap.value.clear()
    getEchosByPage()
  }

  const refreshForSearch = () => {
    current.value = 1
    page.value = 0
    echoList.value = []
    echoMap.value.clear()
  }

  return {
    echoList,
    echoMap,
    isLoading,
    total,
    pageSize,
    page,
    current,
    searchValue,
    searchingMode,
    hasMore,
    echoToUpdate,
    isUpdateMode,
    getEchosByPage,
    refreshEchos,
    refreshForSearch,
  }
})
