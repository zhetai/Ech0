import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useEditorStore = defineStore('editorStore', () => {
  /**
   * state
   */
  const ImageUploading = ref<boolean>(false)

  return {
    ImageUploading,
  }
})
