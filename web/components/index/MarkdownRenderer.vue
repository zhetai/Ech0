<template>
  <div ref="previewElement" class="markdown-preview"></div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue';
import Vditor from 'vditor';  // 如果你使用的是导入 Vditor 库
// 或者直接通过方法.min.js 加载

const previewElement = ref<HTMLDivElement | null>(null);

const props = defineProps({
  content: {
    type: String,
    required: true,
  },
});

const renderMarkdown = async (markdown: string) => {
  if (!previewElement.value) return;

  try {
    // 检查是否已经加载 Vditor
    if (typeof Vditor === 'undefined') {
      console.error('Vditor is not loaded.');
      return;
    }

    // 调用 Vditor 的 preview 方法渲染 markdown 内容
    Vditor.preview(previewElement.value, markdown, {
      mode: 'light', // 设置渲染模式（light 或 dark）
      lang: 'zh_CN', // 设置语言
      theme: {
        current: 'light'
      }, // 设置主题
      hljs: {
        style: 'catppuccin-macchiato',
      },
      markdown: {
        gfmAutoLink: false,
        footnotes: false,
        toc: false,
        sanitize: false,
      },
      after: () => {
        console.log('Rendering complete.');
      }
    });
  } catch (error) {
    console.error("Error rendering markdown:", error);
    previewElement.value.innerHTML = ''; // 如果发生错误，清空内容
  }
};

// Watcher: 当 content 改变时，重新渲染 Markdown
watch(
  () => props.content,
  async (newContent) => {
    await renderMarkdown(newContent);
  },
  { immediate: true } // 初始化时就渲染一次
);

onMounted(() => {
  // 如果初始化时需要渲染，可以放在这里
  renderMarkdown(props.content);
});
</script>

<style>
.markdown-preview {
  font-family: "LXGW WenKai Screen";
}

.markdown-preview h1,
.markdown-preview h2,
.markdown-preview h3,
.markdown-preview h4,
.markdown-preview h5,
.markdown-preview h6 {
  color: #3d3b3b;
}

.markdown-preview p {
  color: #272727;
}

.markdown-preview table thead tr {
  background-color: #fafbfc85 !important;
}

.markdown-preview table tbody tr {
  background-color: #f0f8ff63 !important;
}
</style> 