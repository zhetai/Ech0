import { ref, shallowRef, nextTick } from 'vue'

const title = ref('')
const description = ref('')
const onConfirmCallback = shallowRef<(() => void) | null>(null)
let confirmDialogRef: any = null

export function useBaseDialog() {

  // 注册全局 ConfirmDialog 的引用
  function register(refInstance: any) {
    confirmDialogRef = refInstance
  }

  function openConfirm(options: {
    title: string
    description: string
    onConfirm?: () => void
  }) {
    title.value = options.title
    description.value = options.description
    onConfirmCallback.value = options.onConfirm || null

    nextTick(() => {
      confirmDialogRef?.open()
    })
  }

  function handleConfirm() {
    onConfirmCallback.value?.()
  }

  return {
    register,
    openConfirm,
    title,
    description,
    handleConfirm,
  }
}
