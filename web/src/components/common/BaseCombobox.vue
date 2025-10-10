<template>
  <div class="base-combobox">
    <!-- Label -->
    <label v-if="label" :for="id" class="block text-sm font-medium text-gray-700 mb-1">
      {{ label }}
    </label>

    <!-- Combobox Wrapper -->
    <Combobox
      v-model="internalValue"
      :by="by"
      :multiple="multiple"
      :nullable="!multiple"
      @update:model-value="onSelect"
    >
      <div class="relative">
        <!-- Input -->
        <div
          class="flex items-center px-0.5 py-0.5 rounded-md bg-white transition duration-150 ease-in-out"
          @focusout="onBlurOutside"
          @focusin="onFocusInput"
          @mousedown="onFocusInput"
          >
          <ComboboxInput
            :displayValue="displayValue"
            :placeholder="placeholder"
            @input="onInputChange"
            :class="['outline-none text-md', inputClass]"
          />

          <!-- 可选的 suffix slot -->
          <slot name="suffix">
            <ComboboxButton class="ml-1 text-gray-400">
              <!-- <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24">
                <path
                  fill="#888888"
                  d="m12 19.15l3.875-3.875q.3-.3.7-.3t.7.3t.3.713t-.3.712l-3.85 3.875q-.575.575-1.425.575t-1.425-.575L6.7 16.7q-.3-.3-.288-.712t.313-.713t.713-.3t.712.3zm0-14.3L8.15 8.7q-.3.3-.7.288t-.7-.288q-.3-.3-.312-.712t.287-.713l3.85-3.85Q11.15 2.85 12 2.85t1.425.575l3.85 3.85q.3.3.288.713t-.313.712q-.3.275-.7.288t-.7-.288z"
                />
              </svg> -->
            </ComboboxButton>
          </slot>
        </div>

        <!-- Dropdown -->
        <Transition
          enter="transition ease-out duration-100"
          enter-from="opacity-0 translate-y-1"
          enter-to="opacity-100 translate-y-0"
          leave="transition ease-in duration-75"
          leave-from="opacity-100 translate-y-0"
          leave-to="opacity-0 translate-y-1"
        >
          <ComboboxOptions
            v-if="dropdownOpen && (filteredOptions.length > 0 || allowCreate)"
            class="absolute z-10 mt-2 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-sm shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none"
          >
            <!-- Existing Options -->
            <ComboboxOption
              v-for="item in filteredOptions"
              :key="item.id || item.value || item"
              :value="item"
              class="text-gray-500 hover:text-gray-700 text-lg cursor-pointer select-none px-2 py-1"
            >
              <slot name="option" :option="item">
                {{ item[labelField] || item }}
              </slot>
            </ComboboxOption>

            <!-- Create new option -->
            <ComboboxOption
              v-if="
                allowCreate &&
                normalizedQuery &&
                !filteredOptions.some((o) => getOptionLabel(o).toLowerCase() === normalizedQuery)
              "
              :value="{ [labelField]: query, isNew: true }"
              class="cursor-pointer select-none px-3 py-2 text-orange-600 hover:bg-orange-50"
            >
              创建新标签 "{{ query }}"
            </ComboboxOption>
          </ComboboxOptions>
        </Transition>
      </div>
    </Combobox>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import {
  Combobox,
  ComboboxInput,
  ComboboxOptions,
  ComboboxOption,
  ComboboxButton,
} from '@headlessui/vue'
import { Transition } from 'vue'

type ClassValue = string | string[] | Record<string, boolean | number | string>

const props = defineProps<{
  /** 绑定到外部的值，支持单选或多选 */
  modelValue: any
  /** 可供选择的选项列表 */
  options: any[]
  /** 输入框上方显示的标签文本 */
  label?: string
  /** 关联 label 与输入框的 id */
  id?: string
  /** 输入框提示文本 */
  placeholder?: string
  /** 自定义选项比对逻辑或字段名 */
  by?: string | ((a: any, b: any) => boolean)
  /** 显示选项时使用的字段名 */
  labelField?: string
  /** 是否允许创建新选项 */
  allowCreate?: boolean
  /** 是否启用多选模式 */
  multiple?: boolean
  /** 输入框额外的样式类 */
  inputClass?: ClassValue
}>()

const emit = defineEmits(['update:modelValue', 'create'])

const query = ref('')
const dropdownOpen = ref(false)
const internalValue = ref(props.modelValue)
const labelField = props.labelField || 'name'
const allowCreate = props.allowCreate ?? false
const multiple = props.multiple ?? false

watch(
  () => props.modelValue,
  (val) => {
    internalValue.value = val
  },
)
watch(internalValue, (val) => {
  emit('update:modelValue', val)
})

const onSelect = (val: any) => {
  internalValue.value = val
  query.value = getOptionLabel(val) // 更新显示
  dropdownOpen.value = multiple
}

const onInputChange = (e: Event) => {
  const value = (e.target as HTMLInputElement).value.trim()
  query.value = value
  dropdownOpen.value = true

  // 输入框被清空时 -> 清空绑定值
  if (value === '') {
    internalValue.value = multiple ? [] : null
    emit('update:modelValue', internalValue.value)
    return
  }

  // 如果输入内容刚好匹配某个现有选项 -> 自动选择该项
  const matched = props.options.find(
    (option) => getOptionLabel(option).toLowerCase() === value.toLowerCase()
  )
  if (matched) {
    internalValue.value = matched
    emit('update:modelValue', matched)
    dropdownOpen.value = multiple
  } else {
    // 否则表示用户正在输入新的标签
    internalValue.value = { [labelField]: value, isNew: true }
    emit('create', value) // 可选：通知外部准备创建
    emit('update:modelValue', internalValue.value)
  }
}

const onFocusInput = () => {
  dropdownOpen.value = true
}

const onBlurOutside = (e: FocusEvent) => {
  // 确保焦点确实离开整个 Combobox（不是内部选项）
  const currentTarget = e.currentTarget as HTMLElement
  if (!currentTarget.contains(e.relatedTarget as Node)) {
    dropdownOpen.value = false
    if (query.value.trim() === '') {
      internalValue.value = multiple ? [] : null
      emit('update:modelValue', internalValue.value)
    }
  }
}



const getOptionLabel = (option: any): string => {
  if (option == null) return ''
  if (typeof option === 'object' && !Array.isArray(option)) {
    const record = option as Record<string, unknown>
    const candidate = record[labelField]
    if (candidate != null) return String(candidate)
  }
  return String(option ?? '')
}

const normalizedQuery = computed(() => query.value.trim().toLowerCase())

const filteredOptions = computed(() => {
  if (!normalizedQuery.value) return props.options
  const lowerQuery = normalizedQuery.value
  return props.options.filter((option) => getOptionLabel(option).toLowerCase().includes(lowerQuery))
})

const displayValue = (item: any) => {
  if (Array.isArray(item)) return item.map((i) => getOptionLabel(i)).join(', ')
  return getOptionLabel(item)
}
</script>

<style scoped>
.base-combobox {
  display: flex;
  flex-direction: column;
}
</style>
