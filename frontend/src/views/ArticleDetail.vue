<template>
  <div class="article-detail-container fade-in" v-if="article">
    <div class="article-layout">
      <div class="main-content">
        <div class="card">
          <h1 class="article-title">{{ article.title }}</h1>
          <div class="meta text-sm text-gray">
            <span>{{ $t('article.by') }} {{ article.author?.username }}</span>
            <span> • {{ formatDate(article.created_at) }}</span>
            <span v-if="article.category"> • {{ $t('article.in') }} <strong>{{ article.category.name }}</strong></span>
            <span> • {{ article.views }} {{ $t('article.views') }}</span>
          </div>
          
          <div class="tags" v-if="article.tags && article.tags.length">
            <span v-for="tag in article.tags" :key="tag.id" class="tag">#{{ tag.name }}</span>
          </div>

          <div class="content markdown-body" v-html="processedContent"></div>
          
          <div class="interactions">
            <button @click="likeArticle" :disabled="liked" class="btn" :class="liked ? 'btn-primary' : 'btn-ghost'">
              <span v-if="liked">❤️ {{ $t('article.liked') }}</span>
              <span v-else>🤍 {{ $t('article.like') }}</span>
              <span class="ml-2">({{ article.likes }})</span>
            </button>
          </div>

          <div class="actions" v-if="isAuthor">
            <router-link :to="'/articles/' + article.id + '/edit'" class="btn btn-ghost">{{ $t('common.edit') }}</router-link>
            <button @click="deleteArticle" class="btn btn-danger">{{ $t('common.delete') }}</button>
          </div>
        </div>

        <div class="comments-section card mt-4">
          <h3>{{ $t('comments.title') }}</h3>
          
          <form @submit.prevent="addComment" v-if="authStore.isAuthenticated" class="comment-form">
            <textarea v-model="newComment" :placeholder="$t('comments.placeholder')" required rows="3"></textarea>
            <div class="form-actions">
              <button type="submit" class="btn btn-primary">{{ $t('comments.post') }}</button>
            </div>
          </form>
          <div v-else class="login-prompt">
            <router-link to="/login">{{ $t('auth.login_title') }}</router-link> to comment.
          </div>

          <div class="comment-list">
            <div v-for="comment in comments" :key="comment.id" class="comment">
              <div class="comment-header">
                <strong>{{ comment.user?.username }}</strong>
                <span class="text-sm text-gray">{{ formatDate(comment.created_at) }}</span>
              </div>
              <p class="comment-body">{{ comment.content }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Table of Contents -->
      <div class="toc-sidebar" :class="{ 'collapsed': !showToc }" v-if="toc.length > 0">
        <div class="toc-card card">
          <div class="toc-header" :class="{ 'header-collapsed': !showToc }">
            <h3 v-show="showToc">TOC</h3>
            <div class="toc-controls">
              <button v-show="showToc" @click="expandAll" title="Expand All" class="icon-btn">+</button>
              <button v-show="showToc" @click="collapseAll" title="Collapse All" class="icon-btn">-</button>
              <button @click="toggleToc" title="Toggle TOC" class="icon-btn toggle-btn">
                {{ showToc ? '>>' : '<<' }}
              </button>
            </div>
          </div>
          <div class="toc-content" v-show="showToc">
            <ul class="toc-list">
              <li v-for="item in toc" :key="item.id" :class="'level-' + item.level">
                <div class="toc-item-row">
                  <span 
                    v-if="item.children.length" 
                    @click="toggleNode(item)" 
                    class="toggle-icon"
                  >
                    {{ item.expanded ? '▼' : '▶' }}
                  </span>
                  <span v-else class="spacer"></span>
                  <a 
                    :href="'#' + item.id" 
                    @click.prevent="scrollTo(item.id)"
                    :class="{ active: activeHeaderId === item.id }"
                  >
                    {{ item.text }}
                  </a>
                </div>
                <ul v-if="item.children.length && item.expanded" class="toc-sublist">
                  <li v-for="child in item.children" :key="child.id" :class="'level-' + child.level">
                    <div class="toc-item-row">
                      <span 
                        v-if="child.children.length" 
                        @click="toggleNode(child)" 
                        class="toggle-icon"
                      >
                        {{ child.expanded ? '▼' : '▶' }}
                      </span>
                      <span v-else class="spacer"></span>
                      <a 
                        :href="'#' + child.id" 
                        @click.prevent="scrollTo(child.id)"
                        :class="{ active: activeHeaderId === child.id }"
                      >
                        {{ child.text }}
                      </a>
                    </div>
                    <!-- Support 3 levels for now, can be recursive component -->
                     <ul v-if="child.children.length && child.expanded" class="toc-sublist">
                        <li v-for="subchild in child.children" :key="subchild.id" :class="'level-' + subchild.level">
                           <div class="toc-item-row">
                              <span class="spacer"></span>
                              <a 
                                :href="'#' + subchild.id" 
                                @click.prevent="scrollTo(subchild.id)"
                                :class="{ active: activeHeaderId === subchild.id }"
                              >
                                {{ subchild.text }}
                              </a>
                           </div>
                        </li>
                     </ul>
                  </li>
                </ul>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div v-else-if="loading" class="loading">{{ $t('common.loading') }}</div>
  <div v-else class="error">{{ error }}</div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed, nextTick } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import api from '../api/axios';
import { useAuthStore } from '../stores/auth';
import { marked } from 'marked';
import DOMPurify from 'dompurify';
import { useToast } from '../composables/useToast';
import hljs from 'highlight.js';
import 'highlight.js/styles/atom-one-dark.css'; // Sci-Fi theme compatible

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();
const toast = useToast();

const article = ref(null);
const comments = ref([]);
const newComment = ref('');
const loading = ref(true);
const error = ref('');
const liked = ref(false);
const showToc = ref(true);

// Configure marked with highlight.js using modern API
marked.use({
  gfm: true,
  breaks: true,
  highlight: function (code, lang) {
    const language = hljs.getLanguage(lang) ? lang : 'plaintext';
    return hljs.highlight(code, { language }).value;
  },
  langPrefix: 'hljs language-',
});

// Helper to strip Markdown and HTML
const stripHtml = (html) => {
  const tmp = document.createElement('DIV');
  tmp.innerHTML = html;
  return tmp.textContent || tmp.innerText || '';
};

const generateId = (text) => {
  return text
    .toLowerCase()
    .replace(/\s+/g, '-')
    .replace(/[^\w\u4e00-\u9fa5-]/g, '')
    .replace(/^-+|-+$/g, '') || 'heading';
};

// Computed property for TOC
const toc = computed(() => {
  if (!article.value || !article.value.content) return [];
  
  try {
    const tokens = marked.lexer(article.value.content);
    const headings = tokens.filter(t => t.type === 'heading');
    
    const seenIds = new Map();
    const tocList = headings.map(h => {
      const text = h.text; 
      const cleanText = text.replace(/\*\*|__|\*|_/g, ''); 
      
      let id = generateId(cleanText);
      
      if (seenIds.has(id)) {
        const count = seenIds.get(id) + 1;
        seenIds.set(id, count);
        id = `${id}-${count}`;
      } else {
        seenIds.set(id, 1);
      }
      
      return {
        level: h.depth,
        text: cleanText,
        id: id,
        children: [],
        expanded: true
      };
    });
    
    // Build tree
    const tree = [];
    const stack = [];
    
    tocList.forEach(node => {
      const treeNode = { ...node };
      while (stack.length > 0 && stack[stack.length - 1].level >= treeNode.level) {
        stack.pop();
      }
      if (stack.length === 0) {
        tree.push(treeNode);
      } else {
        stack[stack.length - 1].children.push(treeNode);
      }
      stack.push(treeNode);
    });
    
    return tree;
  } catch (e) {
    console.error('Error generating TOC:', e);
    return [];
  }
});

// Computed property for Content
// Computed property for Content
const processedContent = computed(() => {
  if (!article.value || !article.value.content) return '';
  
  try {
    const renderer = new marked.Renderer();
    const seenIds = new Map();
    
    renderer.heading = function(text, level, raw) {
      let idText = '';
      let contentHtml = '';
      let headingLevel = level;

      if (typeof text === 'object' && text !== null && !Array.isArray(text)) {
        // New signature: { tokens, depth }
        const { tokens, depth } = text;
        headingLevel = depth;
        // Parse inline tokens to get HTML content
        if (this.parser && this.parser.parseInline) {
           contentHtml = this.parser.parseInline(tokens);
        } else {
           contentHtml = tokens.map(t => t.text).join(''); 
        }
        idText = contentHtml.replace(/<[^>]*>/g, '');
      } else {
        // Old signature: text, level, raw
        contentHtml = text;
        const source = raw || text || '';
        idText = source.replace(/\*\*|__|\*|_/g, '').replace(/<[^>]*>/g, '');
      }

      let id = generateId(idText);
      
      if (seenIds.has(id)) {
        const count = seenIds.get(id) + 1;
        seenIds.set(id, count);
        id = `${id}-${count}`;
      } else {
        seenIds.set(id, 1);
      }
      
      return `<h${headingLevel} id="${id}">${contentHtml}</h${headingLevel}>`;
    };
    
    // Use marked.parse instead of marked()
    const html = marked.parse(article.value.content, { renderer });
    
    // Trigger side effects after render
    nextTick(() => {
      addCopyButtons();
      setupIntersectionObserver();
    });

    return DOMPurify.sanitize(html, { ADD_ATTR: ['id'] });
  } catch (e) {
    console.error('Error rendering content:', e);
    return `<div class="alert alert-danger">Error rendering content: ${e.message}</div>`;
  }
});

const addCopyButtons = () => {
  const blocks = document.querySelectorAll('pre code');
  blocks.forEach((block) => {
    if (block.parentElement.querySelector('.copy-btn')) return;
    
    const button = document.createElement('button');
    button.className = 'copy-btn';
    button.innerText = 'Copy';
    button.addEventListener('click', () => {
      navigator.clipboard.writeText(block.innerText);
      button.innerText = 'Copied!';
      setTimeout(() => { button.innerText = 'Copy'; }, 2000);
    });
    
    const pre = block.parentElement;
    pre.style.position = 'relative';
    pre.appendChild(button);
  });
};

const activeHeaderId = ref('');
let observer = null;

const scrollTo = (id) => {
  const element = document.getElementById(id);
  if (element) {
    const offset = 80; // Navbar height + buffer
    const bodyRect = document.body.getBoundingClientRect().top;
    const elementRect = element.getBoundingClientRect().top;
    const elementPosition = elementRect - bodyRect;
    const offsetPosition = elementPosition - offset;

    window.scrollTo({
      top: offsetPosition,
      behavior: 'smooth'
    });
    
    // Manually set active to avoid observer lag
    activeHeaderId.value = id;
  }
};

const setupIntersectionObserver = () => {
  if (observer) observer.disconnect();
  
  observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        activeHeaderId.value = entry.target.id;
      }
    });
  }, {
    rootMargin: '-100px 0px -60% 0px'
  });

  const headers = document.querySelectorAll('.markdown-body h1, .markdown-body h2, .markdown-body h3, .markdown-body h4');
  headers.forEach(header => {
    observer.observe(header);
  });
};

const toggleToc = () => {
  showToc.value = !showToc.value;
};

const toggleNode = (node) => {
  node.expanded = !node.expanded;
};

const expandAll = () => {
  const traverse = (nodes) => {
    nodes.forEach(node => {
      node.expanded = true;
      if (node.children) traverse(node.children);
    });
  };
  traverse(toc.value);
};

const collapseAll = () => {
  const traverse = (nodes) => {
    nodes.forEach(node => {
      node.expanded = false;
      if (node.children) traverse(node.children);
    });
  };
  traverse(toc.value);
};

const fetchComments = async () => {
  try {
    const response = await api.get(`/articles/${route.params.id}/comments`);
    comments.value = response.data;
  } catch (err) {
    console.error('Failed to load comments');
  }
};

const fetchArticle = async () => {
  try {
    const response = await api.get(`/articles/${route.params.id}`);
    article.value = response.data;
    // Increment view count
    api.post(`/articles/${route.params.id}/view`);
    fetchComments();
  } catch (err) {
    error.value = 'Failed to load article';
    console.error('Fetch article error:', err);
  } finally {
    loading.value = false;
  }
};

const likeArticle = async () => {
  if (!authStore.isAuthenticated) {
    toast.error('Please login to like');
    return;
  }
  try {
    const response = await api.post(`/articles/${route.params.id}/like`);
    article.value.likes = response.data.likes;
    liked.value = true;
    toast.success('Article liked!');
  } catch (err) {
    toast.error('Failed to like article');
  }
};

const addComment = async () => {
  try {
    await api.post(`/articles/${route.params.id}/comments`, { content: newComment.value });
    newComment.value = '';
    fetchComments();
    toast.success('Comment posted!');
  } catch (err) {
    toast.error('Failed to post comment');
  }
};

const deleteArticle = async () => {
  if (!confirm('Are you sure?')) return;
  try {
    await api.delete(`/articles/${route.params.id}`);
    toast.success('Article deleted');
    router.push('/');
  } catch (err) {
    toast.error('Failed to delete article');
  }
};

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString();
};

onMounted(() => {
  fetchArticle();
});

onUnmounted(() => {
  if (observer) observer.disconnect();
});
</script>

<style scoped>
.article-detail-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 var(--spacing-md);
}

.article-title {
  font-size: 3rem;
  margin-bottom: var(--spacing-sm);
  color: var(--primary-color);
  text-shadow: var(--neon-text-shadow);
}

.meta {
  margin-bottom: var(--spacing-lg);
  padding-bottom: var(--spacing-md);
  border-bottom: 1px solid var(--surface-border);
  font-family: var(--font-heading);
  letter-spacing: 1px;
}

.content {
  font-size: 1.1rem;
  line-height: 1.8;
  color: var(--text-primary);
  margin-bottom: var(--spacing-xl);
}

.tags {
  display: flex;
  gap: 0.5rem;
  margin-bottom: var(--spacing-md);
}

.tag {
  background-color: rgba(0, 243, 255, 0.1);
  padding: 0.25rem 0.75rem;
  border: 1px solid var(--primary-color);
  font-family: var(--font-heading);
  font-size: 0.8rem;
  color: var(--primary-color);
  text-transform: uppercase;
}

.interactions {
  display: flex;
  justify-content: flex-end;
  margin-top: var(--spacing-lg);
  padding-top: var(--spacing-md);
  border-top: 1px solid var(--surface-border);
}

.actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-md);
  margin-top: var(--spacing-md);
}

.comments-section {
  padding: var(--spacing-lg);
  border: 1px solid var(--surface-border);
}

.comment-form {
  margin-bottom: var(--spacing-xl);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: var(--spacing-sm);
}

.comment-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.comment {
  padding: var(--spacing-md);
  background-color: rgba(0, 0, 0, 0.3);
  border: 1px solid var(--surface-border);
  border-left: 3px solid var(--secondary-color);
}

.comment-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: var(--spacing-xs);
  font-family: var(--font-heading);
  color: var(--secondary-color);
}

.ml-2 { margin-left: 0.5rem; }

/* Layout */
.article-layout {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 300px;
  gap: var(--spacing-xl);
  position: relative;
  align-items: start;
}

.main-content {
  min-width: 0; /* Prevent grid blowout */
}

@media (max-width: 1024px) {
  .article-layout {
    grid-template-columns: 1fr;
  }
  .toc-sidebar {
    display: none; /* Hide TOC on mobile for now, or make it a drawer */
  }
}

/* TOC Sidebar */
.toc-sidebar {
  position: sticky;
  top: 100px;
  height: fit-content;
  max-height: calc(100vh - 120px);
  overflow-y: auto;
  transition: all 0.3s ease;
}

.toc-sidebar.collapsed {
  width: 60px; /* Slightly wider to fit button comfortably */
  overflow: hidden;
}

.toc-sidebar.collapsed .toc-card {
  padding: 10px 5px;
  display: flex;
  justify-content: center;
}

.toc-header.header-collapsed {
  border-bottom: none;
  margin-bottom: 0;
  padding-bottom: 0;
  justify-content: center;
  width: 100%;
}

.toc-header.header-collapsed .toc-controls {
  width: 100%;
  justify-content: center;
}

.toc-card {
  padding: var(--spacing-md);
  border: 1px solid var(--primary-color);
  background: rgba(5, 5, 16, 0.9);
}

.toc-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-sm);
  border-bottom: 1px solid var(--surface-border);
  padding-bottom: var(--spacing-xs);
}

.toc-header h3 {
  font-size: 1rem;
  margin: 0;
  color: var(--primary-color);
}

.toc-controls {
  display: flex;
  gap: 2px;
}

.icon-btn {
  background: none;
  border: 1px solid var(--surface-border);
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 0.8rem;
  padding: 2px 5px;
  border-radius: 2px;
}

.icon-btn:hover {
  color: var(--primary-color);
  border-color: var(--primary-color);
}

.toc-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.toc-sublist {
  list-style: none;
  padding-left: 10px;
  margin: 0;
  border-left: 1px solid var(--surface-border);
}

.toc-item-row {
  display: flex;
  align-items: center;
  padding: 2px 0;
}

.toggle-icon {
  cursor: pointer;
  font-size: 0.7rem;
  width: 15px;
  text-align: center;
  color: var(--secondary-color);
}

.spacer {
  width: 15px;
}

.toc-list a {
  color: var(--text-secondary);
  font-size: 0.9rem;
  text-decoration: none;
  display: block;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  transition: color 0.2s;
}

.toc-list a:hover,
.toc-list a.active {
  color: var(--primary-color);
  text-shadow: 0 0 5px var(--primary-color);
}

.level-1 > .toc-item-row > a { font-weight: bold; }
.level-2 > .toc-item-row > a { padding-left: 5px; }
.level-3 > .toc-item-row > a { padding-left: 10px; }

/* Markdown Styles Override */
:deep(.markdown-body) {
  color: var(--text-primary);
  font-family: var(--font-body);
}

:deep(.markdown-body h1),
:deep(.markdown-body h2),
:deep(.markdown-body h3) {
  border-bottom: 1px solid var(--surface-border);
  padding-bottom: 0.3em;
  color: var(--primary-color);
  margin-top: 1.5em;
}

:deep(.markdown-body code) {
  background-color: rgba(0, 243, 255, 0.1);
  color: var(--primary-color);
  padding: 0.2em 0.4em;
  border-radius: 3px;
  font-family: monospace;
}

:deep(.markdown-body pre) {
  background-color: #1e1e1e; /* Atom One Dark bg */
  padding: 1rem;
  border-radius: 4px;
  overflow-x: auto;
  border: 1px solid var(--surface-border);
  position: relative;
}

:deep(.markdown-body pre code) {
  background-color: transparent;
  padding: 0;
  color: inherit;
}

/* Copy Button */
:deep(.copy-btn) {
  position: absolute;
  top: 5px;
  right: 5px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid var(--surface-border);
  color: var(--text-secondary);
  padding: 2px 8px;
  font-size: 0.7rem;
  border-radius: 3px;
  cursor: pointer;
  opacity: 0;
  transition: opacity 0.2s;
}

:deep(.markdown-body pre:hover .copy-btn) {
  opacity: 1;
}

:deep(.copy-btn:hover) {
  background: var(--primary-color);
  color: #000;
}

/* Tables */
:deep(.markdown-body table) {
  width: 100%;
  border-collapse: collapse;
  margin: 1rem 0;
  border: 1px solid var(--surface-border);
}

:deep(.markdown-body th),
:deep(.markdown-body td) {
  border: 1px solid var(--surface-border);
  padding: 0.75rem;
  text-align: left;
}

:deep(.markdown-body th) {
  background-color: rgba(0, 243, 255, 0.05);
  color: var(--primary-color);
  font-weight: bold;
}

:deep(.markdown-body tr:nth-child(even)) {
  background-color: rgba(255, 255, 255, 0.02);
}

:deep(.markdown-body tr:hover) {
  background-color: rgba(0, 243, 255, 0.05);
}
</style>
