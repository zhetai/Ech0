import type { User, Status, UserToLogin, UserToRegister, Response } from "~/types/models"


export const useUserStore = defineStore("userStore", () => {
    // 状态
    const user = ref<User | null>(null);
    const status = ref<Status | null>(null);
    const token = ref<string | null>(null);
    const isLogin = ref<boolean>(false);
    const toast = useToast()

    // 检查是否已经登录过了
    if (localStorage.getItem("token")) {
        token.value = localStorage.getItem("token");
        isLogin.value = true;
    } else {
        token.value = null;
        isLogin.value = false;
    }

    // 注册
    const register = async (userToRegister: UserToRegister) => {
        const response = await postRequest<any>("register", userToRegister);
        if (!response || response.code !== 1) {
            console.log("注册失败");
            toast.add({
                title: "注册失败",
                description: response?.msg,
                icon: "i-fluent-error-circle-16-filled",
                color: "red",
                timeout: 2000,
            });
            return false;
        }

        if (response && response.code === 1) {
            console.log("注册成功");
            return true;
        }

        return false;
    };

    // 登录
    const login = async (userToLogin: UserToLogin) => {
        const response = await postRequest<string>("login", userToLogin);
        if (!response || response.code !== 1) {
            console.log("登录失败");
            toast.add({
                title: "登录失败",
                description: response?.msg,
                icon: "i-fluent-error-circle-16-filled",
                color: "red",
                timeout: 2000,
            });
            return false
        }

        if (response && response.code === 1 && response.data) {
            token.value = response.data;
            localStorage.setItem("token", token.value);
            isLogin.value = true;

            // 获取用户信息
            getStatus();

            return true;
        }

        return false;
    }

    // 获取状态
    const getStatus = async () => {
        const response = await getRequest<Status>("status");
        if (!response || response.code !== 1) {
            console.log("获取系统信息失败");
            toast.add({
                title: "获取系统信息失败",
                description: response?.msg,
                icon: "i-fluent-error-circle-16-filled",
                color: "red",
                timeout: 2000,
            });
            return false;
        }

        if (response && response.code === 1 && response.data) {
            status.value = response.data;
            return true;
        }
    }

    // 获取当前登录用户信息
    const getUser = async () => {
        const response = await getRequest<User>("user");
        if (!response || response.code !== 1) {
            console.log("获取用户信息失败");
            toast.add({
                title: "获取用户信息失败",
                description: response?.msg,
                icon: "i-fluent-error-circle-16-filled",
                color: "red",
                timeout: 2000,
            });
            return false;
        }

        if (response && response.code === 1 && response.data) {
            user.value = response.data;
            return true;
        }
    }

    // 退出登录
    const logout = async () => {
        isLogin.value = false;
        token.value = null;
        localStorage.removeItem("token");
        status.value = null;
        

        return true;
    }

    return {
        user,
        status,
        token,
        isLogin,
        register,
        login,
        getStatus,
        logout,
        getUser,
    }
})