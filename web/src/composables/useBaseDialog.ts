import { ref, shallowRef, nextTick } from 'vue'

const title = ref('')
const description = ref('')
const onConfirmCallback = shallowRef<(() => void) | null>(null)
import type { ComponentPublicInstance } from 'vue'
import type BaseDialog from '@/components/common/BaseDialog.vue'

type ConfirmDialogInstance = ComponentPublicInstance<typeof BaseDialog>

let confirmDialogRef: ConfirmDialogInstance | null = null

export function useBaseDialog() {
  // 注册全局 ConfirmDialog 的引用
  function register(refInstance: ConfirmDialogInstance) {
    confirmDialogRef = refInstance
  }

  function openConfirm(options: { title: string; description: string; onConfirm?: () => void }) {
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
