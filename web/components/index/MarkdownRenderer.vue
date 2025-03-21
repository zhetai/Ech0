<template>
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
  background-color: #000000;
  padding: 2px 4px;
  border-radius: 8px;
}

#markdown-content pre {
  background-color: #000000;
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


</style>
