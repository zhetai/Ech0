import { useUserStore } from './user'
import { useSettingStore } from './setting'
import { useTodoStore } from './todo'
import { useEditorStore } from './editor'

export async function initStores() {
  const userStore = useUserStore()
  const settingStore = useSettingStore()
  const todoStore = useTodoStore()
  const editorStore = useEditorStore()

  await userStore.init()
  settingStore.init()
  todoStore.init()
  editorStore.init()
}
