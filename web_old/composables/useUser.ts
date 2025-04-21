import type { UserToRegister, UserToLogin } from '~/types/models';
import { useUserStore } from '~/store/user';

export const useUser = () => {
    const userStore = useUserStore();
    const toast = useToast()
    const router = useRouter()

    const register = async (userToRegister: UserToRegister) => {
        const response = await userStore.register(userToRegister);
        if (response) {
            toast.add({
                title: '注册成功',
                description: '欢迎使用Ech0s！请登录',
                icon: 'i-fluent-checkmark-starburst-16-filled',
                color: 'green',
                timeout: 1000,
            });
        }
    }

    const login = async (userToLogin: UserToLogin) => {
        const response = await userStore.login(userToLogin);
        if (response) {
            toast.add({
                title: '登录成功',
                description: '欢迎回来！',
                icon: 'i-fluent-checkmark-starburst-16-filled',
                color: 'green',
                timeout: 1000,
            });
            // 跳转到首页
            router.push('/')
        }
    }

    const logout = () => {
        userStore.logout()
        toast.add({
            title: '注销成功',
            description: '欢迎再次使用！',
            icon: 'i-fluent-checkmark-starburst-16-filled',
            color: 'green',
            timeout: 1000,
        });
        router.push('/')
    }

    const getStatus = () => {
        return userStore.getStatus()
    }

    return {
        register,
        login,
        logout,
        getStatus,
    }
}