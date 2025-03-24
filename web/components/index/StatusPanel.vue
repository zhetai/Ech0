<template>
    <div class="h-3/5 max-w-sm">
        <h1 class="text-5xl italic font-bold text-center text-gray-300 mb-14">Ech0s Status</h1>

        <!-- 登录后显示的内容 -->
        <div v-if="showStatus">
            <!-- 当前系统状态 -->
            <div class="flex justify-start items-center">
                <span class="text-gray-500 text-md">当前系统管理员:</span>
                <span class="font-bold text-md text-gray-600 ml-2">{{ userStore?.status?.username }}</span>
            </div>
            <div class="flex justify-start items-center mt-2">
                <span class="text-gray-500 text-md">当前Ech0共有:</span>
                <span class="font-bold text-md text-gray-600 ml-2">{{ userStore?.status?.total_messages }}</span>
                <span class="text-gray-500 text-md ml-2">条</span>
            </div>
            <!-- 退出登录按钮 -->
            <div v-if="isLogin">
                <UButton size="sm" icon="i-mdi-logout" @click="logout" color="gray" variant="solid"
                    class="text-gray-500 mt-6" />
            </div>
            <div v-else>
                <UButton size="sm" icon="i-mdi-account-key-outline" @click="showStatusChange" color="gray"
                    variant="solid" class="text-gray-500 mt-6" />
            </div>
        </div>
        <div v-else>
            <!-- 登录或注册 -->
            <div class="w-full h-full">
                <div v-if="authmode" class="p-3">
                    <div class="flex items-center justify-between mb-2">
                        <span class="text-gray-500 text-md ml-1">登录</span>
                        <UButton size="sm" icon="i-mdi-account-plus-outline" @click="authmode = false" color="gray"
                            variant="solid" class="text-gray-400" />
                    </div>

                    <UForm :state="authForm" class="flex flex-col gap-2 mt-1">
                        <UInput v-model="authForm.username" placeholder="用户名" />
                        <UInput v-model="authForm.password" type="password" placeholder="密码" />
                        <UButton size="md" @click="login(authForm)" color="gray" variant="solid"
                            class="text-gray-500 mt-6 w-1/3 mx-auto">
                            <span class="mx-auto">
                                登录
                            </span>
                        </UButton>
                    </UForm>
                </div>
                <div v-else class="p-3">
                    <div class="flex items-center justify-between mb-2">
                        <span class="text-gray-500 text-md ml-1">注册</span>
                        <UButton size="sm" icon="i-mdi-account-key-outline" @click="authmode = true" color="gray"
                            variant="solid" class="text-gray-400" />
                    </div>
                    <UForm :state="authForm" class="flex flex-col gap-2 mt-1">
                        <UInput v-model="authForm.username" placeholder="用户名" />
                        <UInput v-model="authForm.password" type="password" placeholder="密码" />
                        <UButton size="md" @click="register(authForm)" color="gray" variant="solid"
                            class="text-gray-500 mt-6 w-1/3 mx-auto">
                            <span class="mx-auto">
                                注册
                            </span>
                        </UButton>
                    </UForm>
                </div>
            </div>
        </div>

    </div>
</template>

<script setup lang="ts">

import { ref, reactive } from 'vue'
import type { UserToLogin, UserToRegister } from '~/types/models';
import { useUser } from '~/composables/useUser';
import { useUserStore } from '~/store/user';

const isLogin = computed(() => userStore?.isLogin ?? false);
const authmode = ref<boolean>(true); // authmode : true 登录， false 注册
const showStatus = ref<boolean>(true); // 是否显示状态

const userStore = useUserStore()
const { login, register, logout } = useUser()

const showStatusChange = () => {
    showStatus.value = !showStatus.value
}

const authForm = reactive<UserToLogin | UserToRegister>({
    username: '',
    password: ''
})

onMounted(() => {
    showStatus.value = true
    userStore.getStatus()
})

</script>