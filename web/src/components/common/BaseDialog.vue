<!-- ConfirmDialog.vue -->
<template>
  <TransitionRoot :show="isOpen" as="template">
    <Dialog @close="close" class="relative z-50">
      <!-- 背景遮罩 -->
      <TransitionChild
        enter="duration-300 ease-out"
        enter-from="opacity-0"
        enter-to="opacity-100"
        leave="duration-200 ease-in"
        leave-from="opacity-100"
        leave-to="opacity-0"
      >
        <div class="fixed inset-0 bg-black/30" aria-hidden="true" />
      </TransitionChild>

      <!-- 对话框面板 -->
      <div class="fixed inset-0 flex items-center justify-center p-4">
        <TransitionChild
          enter="duration-300 ease-out"
          enter-from="opacity-0 scale-95"
          enter-to="opacity-100 scale-100"
          leave="duration-200 ease-in"
          leave-from="opacity-100 scale-100"
          leave-to="opacity-0 scale-95"
        >
          <DialogPanel
            class="w-full max-w-sm rounded-xl bg-white p-6 shadow-sm ring-1 ring-inset ring-gray-200"
          >
            <DialogTitle class="text-base font-semibold text-gray-600">
              {{ title }}
            </DialogTitle>
            <DialogDescription class="mt-2 text-sm text-gray-400 leading-relaxed">
              {{ description }}
            </DialogDescription>

            <div class="mt-6 flex justify-end gap-3">
              <button
                @click="cancel"
                class="cursor-pointer px-3 py-2 rounded-lg bg-white shadow-xs ring-1 ring-inset ring-gray-300 text-gray-600 hover:text-orange-400"
              >
                取消
              </button>
              <button
                @click="confirm"
                class="cursor-pointer px-3 py-2 rounded-lg bg-orange-400 text-white shadow-xs hover:bg-orange-500"
              >
                确认
              </button>
            </div>
          </DialogPanel>
        </TransitionChild>
      </div>
    </Dialog>
  </TransitionRoot>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import {
  Dialog,
  DialogPanel,
  DialogTitle,
  DialogDescription,
  TransitionChild,
  TransitionRoot,
} from '@headlessui/vue'

const props = defineProps({
  title: String,
  description: String,
})

const emit = defineEmits(['confirm', 'cancel'])

const isOpen = ref(false)

function open() {
  isOpen.value = true
}

function close() {
  isOpen.value = false
}

function confirm() {
  emit('confirm')
  close()
}

function cancel() {
  emit('cancel')
  close()
}

defineExpose({ open })
</script>
