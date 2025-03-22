<!-- <template>
  <div id="markdown-content" v-html="renderedMarkdown" class="text-gray-900"></div>
</template>

<script setup lang="ts">
import hljs from "highlight.js";
import "highlight.js/styles/github-dark.css"; // 引入高亮样式
import { ref, watch } from "vue";
import { useMarkdown } from "~/composables/useMarkdown";
const { renderMarkdown } = useMarkdown();

// 接收传入的 Markdown 内容
const props = defineProps({
  content: {
    type: String,
    required: true,
  },
});

// 渲染后的 HTML
const renderedMarkdown = ref("");



// 监听 `content` 的变化，并重新渲染
watch(
    () => props.content,
    (newContent) => {
        const result = renderMarkdown(newContent);
        if (result instanceof Promise) {
            result.then(res => {
                renderedMarkdown.value = res;
            });
        } else {
            renderedMarkdown.value = result;
        }
    },
    { immediate: true } // 初始渲染时也执行一次
);
</script>
<style>
/* 这里通过 id 选择器为整个 Markdown 渲染的内容添加样式 */
#markdown-content {
  color: #333;
}

#markdown-content h1,
#markdown-content h2,
#markdown-content h3,
#markdown-content h4,
#markdown-content h5,
#markdown-content h6 {
  margin-bottom: 1rem;
  font-weight: bold;
}

#markdown-content p {
  margin-bottom: 1rem;
}
#markdown-content ul {
  list-style: disc;
  padding-inline-start: 30px;
}

#markdown-content ol {
  list-style: decimal;
  padding-inline-start: 30px;
}

#markdown-content ul input[type="checkbox"] {
  margin: 0;
  margin-right: 5px;
}

#markdown-content ul:has(input) {
  list-style: none;
  padding-inline-start: 8px;
}

#markdown-content li {
  margin-bottom: .5rem;
  line-height: 1.5;
}

#markdown-content ol,
#markdown-content ul,
#markdown-content blockquote,
#markdown-content .highlight {
  margin: 1rem 0;
}

#markdown-content code {
  background-color: #ffffff;
  padding: 2px 4px;
  border-radius: 8px;
}

#markdown-content pre {
  background-color: #ffffff;
  padding: 1rem;
  border-radius: 16px;
  white-space: pre-wrap;
  /* word-wrap: break-word; */
  overflow: auto;
  box-shadow: 0 0 8px #1d11044f;
}

#markdown-content a {
  color: #1b5b368f;
  background-color: #c9e7d421   ;
  /* text-decoration: none; */
}

#markdown-content a:hover {
  text-decoration: underline;
}


</style> -->

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
  margin-bottom: 1rem;
  font-weight: bold;
  font-size: 1.5rem; /* 你可以根据需要调整字体大小 */
  line-height: 1.5;  /* 行高 */
}

.markdown-preview p {
  margin-bottom: 1rem;
}

</style> 