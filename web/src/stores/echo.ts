import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import { fetchGetEchosByPage, fetchGetTags } from '@/service/api'

export const useEchoStore = defineStore('echoStore', () => {
  /**
   * state
   */
  const echoList = ref<App.Api.Ech0.Echo[]>([]) // 存储Echo列表
  const echoIndexMap = ref(new Map<number, number>()) // id -> index
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
  const tagList = ref<App.Api.Ech0.Tag[]>([]) // 标签列表
  const tagOptions = computed<string[]>(() => {
    return tagList.value.map((tag) => (tag.name))
  })

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

    isLoading.value = true

    await fetchGetEchosByPage({
      page: current.value,
      pageSize: pageSize.value,
      search: searchValue.value,
    })
      .then((res) => {
        if (res.code === 1) {
          total.value = res.data.total

          // 同步更新 echoMap
          res.data.items.forEach((item: App.Api.Ech0.Echo) => {
            const idx = echoIndexMap.value.get(item.id)
            if (idx !== undefined) {
              echoList.value[idx] = item // 更新已有数据
            } else {
              echoList.value.push(item) // 添加新数据
              echoIndexMap.value.set(item.id, echoList.value.length - 1)
            }
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
    echoIndexMap.value.clear()
    getEchosByPage()
  }

  const clearEchos = () => {
    current.value = 1
    page.value = 0
    echoList.value = []
    echoIndexMap.value.clear()
    total.value = 0
  }

  const refreshForSearch = () => {
    current.value = 1
    page.value = 0
    echoList.value = []
    echoIndexMap.value.clear()
  }

  const updateEcho = (echo: App.Api.Ech0.Echo) => {
    const idx = echoIndexMap.value.get(echo.id)
    if (idx !== undefined) {
      echoList.value[idx] = echo // 更新
    }
  }

  const updateLikeCount = (echoId: number, delta: number = 1) => {
    const idx = echoIndexMap.value.get(echoId)
    if (idx !== undefined) {
      const targetEcho = echoList.value[idx]
      if (targetEcho) {
        targetEcho.fav_count = (targetEcho.fav_count || 0) + delta
        echoList.value[idx] = { ...targetEcho } // 保证响应式触发
      }
    }
  }

  const getTags = async () => {
    const res = await fetchGetTags()
    if (res.code === 1) {
      // 清空原数组再插入新数据
      tagList.value.splice(0, tagList.value.length, ...res.data)
      console.log('标签列表更新', tagList.value)
      console.log('标签选项更新', tagOptions.value)
    }
  }

  const init = () => {
    getTags()
  }

  return {
    echoList,
    echoIndexMap,
    isLoading,
    total,
    pageSize,
    page,
    current,
    searchValue,
    searchingMode,
    hasMore,
    echoToUpdate,
    tagList,
    tagOptions,
    getEchosByPage,
    refreshEchos,
    clearEchos,
    refreshForSearch,
    updateEcho,
    updateLikeCount,
    getTags,
    init,
  }
})
