<template>
  <div class="rounded-md shadow-sm ring-1 ring-gray-200 ring-inset bg-white p-4 mb-3">
    <!-- 数据管理 -->
    <div>
      <div class="flex items-center justify-start mb-3">
        <h1 class="text-gray-600 font-bold text-lg">数据管理</h1>
      </div>

      <div class="flex flex-col gap-4">
        <!-- 备份数据 -->
        <div class="flex flex-start items-center gap-2">
          <p class="text-gray-400">创建快照:</p>
          <BaseButton
            :icon="CreateBackup"
            @click="handleBackup"
            class="rounded-lg !bg-gray-100 !text-gray-600 hover:!bg-gray-200"
            title="创建快照"
          />
        </div>
        <!-- 导出快照 -->
        <div class="flex flex-start items-center gap-2">
          <p class="text-gray-400">导出快照:</p>
          <BaseButton
            :icon="ExportBackup"
            @click="handleBackupExport"
            class="rounded-lg !bg-gray-100 !text-gray-600 hover:!bg-gray-200"
            title="导出快照"
          />
        </div>
        <!-- 恢复数据 -->
        <div class="flex flex-start items-center gap-2">
          <p class="text-gray-400">恢复快照:</p>
          <BaseButton
            :icon="RestoreBackup"
            @click="handleBackupRestore"
            class="rounded-lg !bg-gray-100 !text-gray-600 hover:!bg-gray-200"
            title="恢复快照"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import BaseButton from '@/components/common/BaseButton.vue'
import CreateBackup from '@/components/icons/createbackup.vue'
import ExportBackup from '@/components/icons/exportbackup.vue'
import RestoreBackup from '@/components/icons/restorebackup.vue'
import { fetchBackup, fetchImportBackup } from '@/service/api'
import { theToast } from '@/utils/toast'

const handleBackup = async () => {
  await theToast.promise(
    fetchBackup(),
    {
      loading: '备份中...',
      success: (res) => (res.code === 1 ? '备份成功' : '备份失败'),
      error: '备份失败',
    }
  )
}


const handleBackupExport = async () => {
  try {
    theToast.info('导出中...请稍等', {
      duration: 4000,
    })

    // 1. 获取 token
    const token = localStorage.getItem('token')
    const baseURL = import.meta.env.VITE_SERVICE_BASE_URL === "/" ? window.location.origin : import.meta.env.VITE_SERVICE_BASE_URL
    const downloadUrl = `${baseURL}/api/backup/export?token=${token}`

    // 创建隐藏的 a 标签触发下载
    const link = document.createElement('a')
    link.href = downloadUrl
    link.style.display = 'none'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)

    theToast.success('备份导出开始！')
  } catch (error) {
    theToast.error('导出失败')
    console.error('导出备份失败:', error)
  }
}

const handleBackupRestore = async () => {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = '.zip'
  input.onchange = async (event: Event) => {
    const target = event.target as HTMLInputElement
    if (target.files && target.files.length > 0) {
      const file = target.files[0]

      await theToast.promise(
        fetchImportBackup(file),
        {
          loading: '导入中,请不要关闭页面...',
          success: (res) => (res.code === 1 ? '快照恢复成功🎉' : `导入失败: ${res.msg}`),
          error: '导入失败',
        },
        {
          duration: 5000,
        }
      )
    }
  }
  input.click()
}
</script>
