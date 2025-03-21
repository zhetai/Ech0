// composables/useMarkdown.ts
import { marked } from 'marked';
import hljs from 'highlight.js';
import 'highlight.js/styles/atom-one-light.min.css'; // 引入 highlight.js 的样式
import javascript from 'highlight.js/lib/languages/javascript'; // 引入 JavaScript 语言支持

hljs.registerLanguage('javascript', javascript);

export const useMarkdown = () => {
    // 配置 Markdown 渲染选项
    const rendererOptions = {
        gfm: true,        // 启用 GitHub 风格的 Markdown 渲染
        breaks: true,     // 启用换行支持
        pedantic: false,  // 不严格遵守原始 Markdown 语法规范
        silent: false,    // 显示渲染错误和警告
        async: false,     // 渲染过程中不使用异步
    };
    
  // 创建一个自定义的 marked 渲染器
  const renderer = new marked.Renderer();

  // 自定义处理图片渲染，使其包含 data-fancybox 属性
  renderer.image = ({ href, title, text }: { href: string, title: string | null, text: string }) => {
    return `<div class="rounded-lg overflow-hidden mb-2 w-5/6 mx-auto shadow-lg"><a href="${href}" data-fancybox="gallery" data-caption="${text}"><img src="${href}" alt="${text}" class="max-w-full w-full object-cover" loading="lazy" /></a></div>`;
  };

  // 自定义处理代码块渲染，并加上 highlight.js 的高亮
  renderer.code = ({ text, lang }: { text: string, lang?: string }) => {
    if (lang && hljs.getLanguage(lang)) {
        try {
            const highlightedCode = hljs.highlight(text, { language: lang }).value;
            return `<pre class="hljs"><code>${highlightedCode}</code></pre>`;
        } catch (error) {
            console.error(error);
            return `<pre class="hljs"><code>${text}</code></pre>`;
        }
    }
    // 如果没有语言支持或失败，使用默认代码块渲染
    return `<pre class="hljs"><code>${text}</code></pre>`;
};
  
  // 渲染 Markdown 内容的函数
  const renderMarkdown = (markdownText: string) => {
    return marked(markdownText, {
      ...rendererOptions,
      renderer, // 使用自定义的渲染器
    });
  };

  return { renderMarkdown };
};
