import { defineStore } from 'pinia'
import { ref } from 'vue'
import { fetchGetTodos } from '@/service/api'
import { useUserStore } from './user'

export const useTodoStore = defineStore('todoStore', () => {
  /**
   * State
   */
  const todos = ref<App.Api.Todo.Todo[]>([])
  const todoMode = ref<boolean>(false)
  const loading = ref<boolean>(true)

  /**
   * Actions
   */
  function getTodos() {
    fetchGetTodos()
      .then((res) => {
        if (res.code === 1) {
          todos.value = res.data
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
    if (useUserStore().isLogin) {
      getTodos()
    }
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
