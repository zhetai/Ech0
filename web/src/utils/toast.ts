// 处理全局通知工具函数 (基于 vue-sonner)
import { toast } from 'vue-sonner'

// 定义自定义通知选项接口
interface customToastOptions {
  duration?: number
}

// 默认通知选项
const defaultToastOptions: customToastOptions = {
  duration: 1500, // 默认持续时间为1500毫秒
}

function show(
  type: 'success' | 'error' | 'info' | 'warning',
  content: string,
  options?: customToastOptions,
) {
  const config = {
    duration: options?.duration ?? defaultToastOptions.duration,
  }

  toast[type](content, config)
}

export const theToast = {
  success: (content: string, options?: customToastOptions) => show('success', content, options),
  error: (content: string, options?: customToastOptions) => show('error', content, options),
  info: (content: string, options?: customToastOptions) => show('info', content, options),
  warning: (content: string, options?: customToastOptions) => show('warning', content, options),
}
