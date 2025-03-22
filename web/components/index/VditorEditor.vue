<template>
  <div ref="editorContainer" class="vditor-container"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from "vue";
import Vditor from "vditor";
import "vditor/dist/index.css";

const props = defineProps({
  modelValue: {
    type: String,
    default: "",
  },
});

const emit = defineEmits(["update:modelValue"]);

// 编辑器实例引用
const editorContainer = ref<HTMLElement>();
let vditorInstance: Vditor | null = null;

// 编辑器配置
const editorOptions: IOptions = {
  mode: "ir",
  height: 140,
  icon: "ant",
  lang: "zh_CN" as keyof II18n,
  theme: "classic",
  toolbar: ["preview", "fullscreen"],
  toolbarConfig: {
    pin: true,
  },
  counter: {
    enable: false,
  },
  cache: {
    enable: true,
    id: "vue-vditor",
  },
  input: (content: string) => {
    emit("update:modelValue", content);
  },
  preview: {
    hljs: {
      style: "native",
    },
    actions: [],
  },
  placeholder: "一吐为快~"
};

// 初始化编辑器
onMounted(async () => {
  if (!editorContainer.value) return;

  vditorInstance = new Vditor(editorContainer.value, {
    ...editorOptions,
    after: () => {
      vditorInstance?.setValue(props.modelValue);
    },
  });
});

// 清理资源
onBeforeUnmount(() => {
  vditorInstance?.destroy();
});

// 暴露 clear 方法
defineExpose({
  clear: () => {
    if (vditorInstance) {
      vditorInstance.setValue(''); // 通过 setValue 清空内容
    }
  },
});
</script>

<style>
.vditor-container {
  border-radius: 8px;
  overflow: hidden;
  /* box-shadow: 0px 0px 4px rgba(0, 0, 0, 0.1); */
  margin-bottom: 12px;
}

.vditor-toolbar--pin {
  padding-left: 6px !important;
}

.vditor-ir pre.vditor-reset {
    padding: 5px 10px !important;
}

@media screen and (max-width: 520px) {
    .vditor-toolbar__item {
        padding: 0px 4px;
    }
}
</style>
