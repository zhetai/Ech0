import { useUserStore } from './user'
import { useSettingStore } from './setting'
import { useTodoStore } from './todo'

export async function initStores() {
  const userStore = useUserStore()
  const settingStore = useSettingStore()
  const todoStore = useTodoStore()

  await userStore.init()
  settingStore.init()
  todoStore.init()
}
