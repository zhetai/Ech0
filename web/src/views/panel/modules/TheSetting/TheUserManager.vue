<template>
  <PanelCard>
    <div class="flex flex-row items-center justify-between mb-3">
      <h1 class="text-gray-600 font-bold text-lg">用户管理</h1>
      <div class="flex flex-row items-center justify-end gap-2 w-14">
        <button @click="userEditMode = !userEditMode" title="编辑">
          <Edit v-if="!userEditMode" class="w-5 h-5 text-gray-400 hover:w-6 hover:h-6" />
          <Close v-else class="w-5 h-5 text-gray-400 hover:w-6 hover:h-6" />
        </button>
      </div>
    </div>

    <!-- 用户列表 -->
    <div v-if="allusers.length === 0">
      <h2 class="text-gray-500 font-semibold text-center">暂无其它用户</h2>
    </div>
    <div v-else>
      <div
        v-for="user in allusers"
        :key="user.id"
        class="flex flex-row items-center justify-start text-gray-500 gap-2 h-10"
      >
        <h2 class="font-semibold w-30">{{ user.username }}</h2>
        <BaseSwitch
          v-model="user.is_admin"
          :disabled="!userEditMode"
          class="w-14"
          @click="handleUpdateUserPermission(user.id)"
        />
        <BaseButton
          :icon="Deluser"
          class="rounded-md text-center w-auto text-align-center h-8"
          :disabled="!userEditMode"
          @click="handleDeleteUser(user.id)"
        />
      </div>
    </div>
  </PanelCard>
</template>

<script setup lang="ts">
import PanelCard from '@/layout/PanelCard.vue'
import Edit from '@/components/icons/edit.vue'
import Close from '@/components/icons/close.vue'
import { ref, onMounted } from 'vue'
import BaseSwitch from '@/components/common/BaseSwitch.vue'
import BaseButton from '@/components/common/BaseButton.vue'
import Deluser from '@/components/icons/deluser.vue'
import { theToast } from '@/utils/toast'
import { useBaseDialog } from '@/composables/useBaseDialog'
const { openConfirm } = useBaseDialog()

import { fetchGetAllUsers, fetchUpdateUserPermission, fetchDeleteUser } from '@/service/api'

const allusers = ref<App.Api.User.User[]>([])
const userEditMode = ref<boolean>(false)

const handleDeleteUser = async (userId: number) => {
  openConfirm({
    title: '确定要删除该用户吗？',
    description: '删除后将无法恢复，请谨慎操作',
    onConfirm: () => {
      fetchDeleteUser(userId).then((res) => {
        if (res.code === 1) {
          getAllUsers()
        }
      })
    },
  })
}

const handleUpdateUserPermission = async (userId: number) => {
  fetchUpdateUserPermission(userId)
    .then((res) => {
      if (res.code === 1) {
        theToast.success(res.msg)
      }
    })
    .finally(() => {
      // 重新获取设置
      getAllUsers()
    })
}

const getAllUsers = async () => {
  await fetchGetAllUsers().then((res) => {
    if (res.code === 1) {
      allusers.value = res.data
    }
  })
}

onMounted(() => {
  getAllUsers()
})
</script>

<style scoped></style>
