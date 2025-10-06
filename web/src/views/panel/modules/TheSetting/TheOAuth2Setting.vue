<template>
  <div class="rounded-md shadow-sm ring-1 ring-gray-200 ring-inset bg-white p-4 mb-3">
    <!-- OAuth2 设置 -->
    <div class="w-full">
      <div class="flex flex-row items-center justify-between mb-3">
        <h1 class="text-gray-600 font-bold text-lg">OAuth2设置</h1>
        <div class="flex flex-row items-center justify-end gap-2 w-14">
          <button v-if="oauth2EditMode" @click="handleUpdateOAuth2Setting" title="编辑">
            <Saveupdate class="w-5 h-5 text-gray-400 hover:w-6 hover:h-6" />
          </button>
          <button @click="oauth2EditMode = !oauth2EditMode" title="编辑">
            <Edit v-if="!oauth2EditMode" class="w-5 h-5 text-gray-400 hover:w-6 hover:h-6" />
            <Close v-else class="w-5 h-5 text-gray-400 hover:w-6 hover:h-6" />
          </button>
        </div>
      </div>

      <!-- 开启OAuth2 -->
      <div class="flex flex-row items-center justify-start text-gray-500 h-10">
        <h2 class="font-semibold w-30 flex-shrink-0">启用OAuth2:</h2>
        <BaseSwitch v-model="OAuth2Setting.enable" :disabled="!oauth2EditMode" />
      </div>

      <!-- OAuth2 Provider -->
      <div class="flex flex-row items-center justify-start text-gray-500 gap-2 h-10">
        <h2 class="font-semibold w-30 flex-shrink-0">OAuth2 模板:</h2>
        <BaseSelect
          v-model="OAuth2Setting.provider"
          :options="OAuth2ProviderOptions"
          :disabled="!oauth2EditMode"
          class="w-34 h-8"
        />
      </div>

      <!-- Client ID -->
      <div class="flex flex-row items-center justify-start text-gray-500 gap-2 h-10">
        <h2 class="font-semibold w-30 flex-shrink-0">Client ID:</h2>
        <span
          v-if="!oauth2EditMode"
          class="truncate max-w-40 inline-block align-middle"
          :title="OAuth2Setting.client_id"
          style="vertical-align: middle"
        >
          {{ OAuth2Setting.client_id.length === 0 ? '暂无' : OAuth2Setting.client_id }}
        </span>
        <BaseInput
          v-else
          v-model="OAuth2Setting.client_id"
          type="text"
          placeholder="请输入Client ID"
          class="w-full !py-1"
        />
      </div>

      <!-- Client Secret -->
      <div class="flex flex-row items-center justify-start text-gray-500 gap-2 h-10">
        <h2 class="font-semibold w-30 flex-shrink-0">Client Secret:</h2>
        <span
          v-if="!oauth2EditMode"
          class="truncate max-w-40 inline-block align-middle"
          :title="OAuth2Setting.client_secret"
          style="vertical-align: middle"
        >
          {{ OAuth2Setting.client_secret.length === 0 ? '暂无' : OAuth2Setting.client_secret }}
        </span>
        <BaseInput
          v-else
          v-model="OAuth2Setting.client_secret"
          type="text"
          placeholder="请输入Client Secret"
          class="w-full !py-1"
        />
      </div>

      <!-- Callback URL -->
      <div class="flex flex-row items-center justify-start text-gray-500 gap-2 h-10">
        <h2 class="font-semibold w-30 flex-shrink-0">Callback URL:</h2>
        <span
          v-if="!oauth2EditMode"
          class="truncate max-w-40 inline-block align-middle"
          :title="OAuth2Setting.redirect_uri"
          style="vertical-align: middle"
        >
          {{ redirect_uri.length === 0 ? '暂无' : redirect_uri }}
        </span>
        <BaseInput
          v-else
          v-model="redirect_uri"
          type="text"
          placeholder="请输入回调地址"
          class="w-full !py-1"
        />
      </div>

      <!-- Auth URL -->
      <div class="flex flex-row items-center justify-start text-gray-500 gap-2 h-10">
        <h2 class="font-semibold w-30 flex-shrink-0">Auth URL:</h2>
        <span
          v-if="!oauth2EditMode"
          class="truncate max-w-40 inline-block align-middle"
          :title="OAuth2Setting.auth_url"
          style="vertical-align: middle"
        >
          {{ OAuth2Setting.auth_url.length === 0 ? '暂无' : OAuth2Setting.auth_url }}
        </span>
        <BaseInput
          v-else
          v-model="OAuth2Setting.auth_url"
          type="text"
          placeholder="请输入授权地址"
          class="w-full !py-1"
        />
      </div>

      <!-- Token URL -->
      <div class="flex flex-row items-center justify-start text-gray-500 gap-2 h-10">
        <h2 class="font-semibold w-30 flex-shrink-0">Token URL:</h2>
        <span
          v-if="!oauth2EditMode"
          class="truncate max-w-40 inline-block align-middle"
          :title="OAuth2Setting.token_url"
          style="vertical-align: middle"
        >
          {{ OAuth2Setting.token_url.length === 0 ? '暂无' : OAuth2Setting.token_url }}
        </span>
        <BaseInput
          v-else
          v-model="OAuth2Setting.token_url"
          type="text"
          placeholder="请输入Token地址"
          class="w-full !py-1"
        />
      </div>

      <!-- User Info URL -->
      <div class="flex flex-row items-center justify-start text-gray-500 gap-2 h-10">
        <h2 class="font-semibold w-30 flex-shrink-0">User Info URL:</h2>
        <span
          v-if="!oauth2EditMode"
          class="truncate max-w-40 inline-block align-middle"
          :title="OAuth2Setting.user_info_url"
          style="vertical-align: middle"
        >
          {{ OAuth2Setting.user_info_url.length === 0 ? '暂无' : OAuth2Setting.user_info_url }}
        </span>
        <BaseInput
          v-else
          v-model="OAuth2Setting.user_info_url"
          type="text"
          placeholder="请输入用户信息地址"
          class="w-full !py-1"
        />
      </div>

      <!-- Scopes -->
      <div class="flex flex-row items-center justify-start text-gray-500 gap-2 h-10">
        <h2 class="font-semibold w-30 flex-shrink-0">Scopes:</h2>
        <span
          v-if="!oauth2EditMode"
          class="truncate max-w-40 inline-block align-middle"
          :title="OAuth2Setting.scopes.join(', ')"
          style="vertical-align: middle"
        >
          {{ OAuth2Setting.scopes.length === 0 ? '暂无' : OAuth2Setting.scopes.join(', ') }}
        </span>
        <BaseInput
          v-else
          v-model="scopeString"
          type="text"
          placeholder="请输入Scopes，多个用逗号分隔"
          class="w-full !py-1"
          @blur="OAuth2Setting.scopes = scopeString.split(',').map(s => s.trim())"
        />
      </div>


    </div>
  </div>
</template>

<script setup lang="ts">
import BaseInput from '@/components/common/BaseInput.vue'
import BaseSwitch from '@/components/common/BaseSwitch.vue'
import BaseSelect from '@/components/common/BaseSelect.vue'
import Edit from '@/components/icons/edit.vue'
import Close from '@/components/icons/close.vue'
import Saveupdate from '@/components/icons/saveupdate.vue'
import { ref, onMounted } from 'vue'
import { useSettingStore } from '@/stores/setting'
import { theToast } from '@/utils/toast'
import { OAuth2Provider } from '@/enums/enums'
import { fetchUpdateOAuth2Settings } from '@/service/api'

const settingStore = useSettingStore()
const { getOAuth2Setting } = settingStore
const { OAuth2Setting } = settingStore

const oauth2EditMode = ref(false)

const OAuth2ProviderOptions = [
  { label: 'GitHub', value: OAuth2Provider.GITHUB },
  // { label: 'Google', value: OAuth2Provider.GOOGLE },
]

const redirect_uri = ref(`${window.location.origin}/oauth/github/callback`)
const scopeString = ref('read:user')

const handleUpdateOAuth2Setting = async () => {
  // 修改Scopes
  OAuth2Setting.scopes = scopeString.value.split(',').map((s) => s.trim())
  // 修改回调地址为当前域名加上固定路径
  OAuth2Setting.redirect_uri = redirect_uri.value || `${window.location.origin}/oauth/github/callback`

  // 提交更新
  await fetchUpdateOAuth2Settings(OAuth2Setting)
    .then((res) => {
      if (res.code === 1) {
        theToast.success(res.msg)
      }
    })
    .finally(() => {
      oauth2EditMode.value = false
      // 重新获取OAuth2设置
      getOAuth2Setting()
    })
}

onMounted(() => {
  getOAuth2Setting()
})
</script>

<style scoped></style>
