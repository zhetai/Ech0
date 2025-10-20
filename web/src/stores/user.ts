import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { fetchLogin, fetchSignup, fetchGetCurrentUser } from '@/service/api'
import { saveAuthToken } from '@/service/request/shared'
import { localStg } from '@/utils/storage'
import { theToast } from '@/utils/toast'
import router from '@/router'
import { useEchoStore } from './echo'

export const useUserStore = defineStore('userStore', () => {
  /**
   * state
   */
  const user = ref<App.Api.User.User | null>(null)
  const isLogin = computed(() => !!user.value)
  const initialized = ref<boolean>(false)

  /**
   * actions
   */
  // 登录
  async function login(userInfo: App.Api.Auth.LoginParams) {
    await fetchLogin(userInfo).then((res) => {
      const token = String(res.data)

      if (token && token.length > 0) {
        // 保存token到localStorage
        saveAuthToken(token)

        // 获取当前登录用户信息
        refreshCurrentUser()

        // 登录成功
        theToast.success('登录成功,欢迎回来！🎉')

        // 清除echo数据
        const echoStore = useEchoStore()
        echoStore.clearEchos()

        // 跳转到首页
        router.push({ name: 'home' })
      }
    })
  }

  // 使用token登录（自动登录或OAuth2登录后使用）
  async function loginWithToken(token: string) {
    if (token && token.length > 0) {
      // 保存token到localStorage
      saveAuthToken(token)

      // 获取当前登录用户信息
      await refreshCurrentUser()

      // 登录成功
      theToast.success('登录成功,欢迎回来！🎉')

      // 清除echo数据
      const echoStore = useEchoStore()
      echoStore.clearEchos()

      // 跳转到首页
      router.push({ name: 'home' })
    }
  }

  // 注册
  async function signup(userInfo: App.Api.Auth.SignupParams) {
    return await fetchSignup(userInfo).then((res) => {
      // 注册成功，前往登录
      if (res.code === 1) {
        theToast.success('注册成功,请登录！')
        return true
      }

      // 注册失败
      return false
    })
  }

  // 退出登录
  async function logout() {
    // 清除token
    user.value = null

    // 清除echo数据
    const echoStore = useEchoStore()
    echoStore.clearEchos()

    // 标记需要重定向到登录页
    localStg.setItem('needLoginRedirect', true)

    // 重新登录(⚠️：交给路由守卫处理)
    // router.push({ name: 'auth' })
  }

  // 自动登录
  async function autoLogin() {
    // 检查localStorage中是否有token
    const token = String(localStg.getItem('token'))
    if (token && token.length > 0 && token !== 'undefined' && token !== 'null') {
      // 如果有token，则获取用户信息
      await refreshCurrentUser()
    }
  }

  // 获取当前登录用户信息
  async function refreshCurrentUser() {
    const res = await fetchGetCurrentUser()
    if (res.code === 1) {
      console.log('获取用户信息成功,自动登录', res.data)
      user.value = res.data
    } else {
      // 获取用户信息失败，清除token
      console.log('获取用户信息失败，清除token，重新登录')
      await logout()
    }
  }

  // 初始化
  const init = async () => {
    await autoLogin()
    initialized.value = true
  }

  return {
    initialized,
    user,
    isLogin,
    login,
    loginWithToken,
    signup,
    logout,
    autoLogin,
    refreshCurrentUser,
    init,
  }
})
