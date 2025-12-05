<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { Codemirror } from 'vue-codemirror';
import { javascript } from '@codemirror/lang-javascript';
import { oneDark } from '@codemirror/theme-one-dark';
import api from '../services/api';

// --- State ---
const content = ref('');
const language = ref('javascript'); // Пока хардкод, потом сделаем автодетект
const isReadOnly = ref(false);
const isLoading = ref(false);
const isSaving = ref(false);
const statusMessage = ref(''); // Для отображения ошибок или успеха

const route = useRoute();
const router = useRouter();

// Конфигурация редактора
const extensions = [javascript(), oneDark];

// --- Logic ---

// Функция сохранения
const savePaste = async () => {
  if (!content.value.trim()) return;
  
  isSaving.value = true;
  statusMessage.value = 'Saving...';
  
  try {
    const paste = await api.createPaste(content.value, language.value);
    // После сохранения редиректим на страницу просмотра
    await router.push({ name: 'view', params: { id: paste.id } });
    statusMessage.value = 'Saved!';
  } catch (e) {
    console.error(e);
    statusMessage.value = 'Error saving paste';
  } finally {
    isSaving.value = false;
  }
};

// Функция загрузки (если открыли ссылку)
const loadPaste = async (id: string) => {
  isLoading.value = true;
  isReadOnly.value = true; // Блокируем редактирование при просмотре
  statusMessage.value = 'Loading...';
  
  try {
    const paste = await api.getPaste(id);
    content.value = paste.content;
    language.value = paste.language; // Если бэк вернет язык
    statusMessage.value = '';
  } catch (e) {
    statusMessage.value = 'Paste not found or expired';
    content.value = '// Error 404: Data not found in the void.';
  } finally {
    isLoading.value = false;
  }
};

// Создать новую пасту (кнопка New)
const newPaste = () => {
  content.value = '';
  isReadOnly.value = false;
  statusMessage.value = '';
  router.push({ name: 'home' });
};

// --- Lifecycle ---

// Следим за ID в URL
watch(
  () => route.params.id,
  (newId) => {
    if (newId) {
      loadPaste(newId as string);
    } else {
      // Если вернулись на главную
      isReadOnly.value = false;
      content.value = '';
    }
  },
  { immediate: true } // Запустить сразу при загрузке
);
</script>

<template>
  <div class="layout">
    <!-- Header -->
    <header class="navbar">
      <div class="brand" @click="newPaste">
        <span class="logo-text">TRACE</span>
        <span class="status" v-if="statusMessage">:: {{ statusMessage }}</span>
      </div>
      
      <div class="actions">
        <!-- Показываем кнопку SAVE только если мы не в режиме просмотра -->
        <button v-if="!isReadOnly" class="btn save-btn" @click="savePaste" :disabled="isSaving">
          {{ isSaving ? 'SAVING...' : 'SAVE' }}
        </button>
        
        <button v-else class="btn new-btn" @click="newPaste">
          NEW
        </button>
      </div>
    </header>

    <!-- Editor Area -->
    <main class="editor-container">
      <codemirror
        v-model="content"
        placeholder="// Paste your code here..."
        :style="{ height: '100%', fontSize: '14px' }"
        :autofocus="true"
        :indent-with-tab="true"
        :tab-size="2"
        :extensions="extensions"
        :disabled="isReadOnly"
      />
    </main>
  </div>
</template>

<style scoped>
.layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: var(--bg-primary);
}

.navbar {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 1.5rem;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border);
}

.brand {
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 1rem;
}

.logo-text {
  font-family: var(--font-mono);
  font-weight: 700;
  font-size: 1.2rem;
  letter-spacing: -1px;
  color: var(--text-main);
}

.status {
  font-family: var(--font-mono);
  font-size: 0.8rem;
  color: var(--text-muted);
  animation: pulse 2s infinite;
}

.editor-container {
  flex: 1; /* Занимает всё оставшееся место */
  overflow: hidden;
}

/* Переопределяем стили CodeMirror, чтобы убрать белые рамки */
:deep(.cm-editor) {
  height: 100%;
  background-color: var(--bg-primary);
}
:deep(.cm-gutters) {
  background-color: var(--bg-primary);
  border-right: 1px solid var(--border);
}

@keyframes pulse {
  0% { opacity: 0.6; }
  50% { opacity: 1; }
  100% { opacity: 0.6; }
}
</style>