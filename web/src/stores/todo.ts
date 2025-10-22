import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import { fetchGetTodos } from '@/service/api'
import { useUserStore } from './user'

export const useTodoStore = defineStore('todoStore', () => {
  /**
   * State
   */
  const todos = ref<App.Api.Todo.Todo[]>([])
  const todoMode = ref<boolean>(false)
  const loading = ref<boolean>(true)

  const userStore = useUserStore()

  /**
   * Actions
   */
  function getTodos() {
    fetchGetTodos()
      .then((res) => {
        if (res.code === 1) {
          todos.value = res.data
        } else {
          todos.value = []
          loading.value = false
        }
      })
      .finally(() => {
        loading.value = false
      })
  }

  function setTodoMode(mode: boolean) {
    todoMode.value = mode
  }

  function init() {
    watch(
      () => userStore.isLogin,
      (isLogin) => {
        if (isLogin) {
          getTodos()
        } else {
          todos.value = []
          loading.value = false
        }
      },
      { immediate: true },
    )
  }

  return {
    todos,
    todoMode,
    loading,
    setTodoMode,
    getTodos,
    init,
  }
})
