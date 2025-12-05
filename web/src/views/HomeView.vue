<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, computed } from 'vue'; // UPDATED: добавил onUnmounted, computed
import { useRoute, useRouter } from 'vue-router';
import { Codemirror } from 'vue-codemirror';
import { oneDark } from '@codemirror/theme-one-dark';
import api from '../services/api';

// UPDATED: Импорты
import { getExtensionByName } from '../config/languages'; 
import CommandPalette from '../components/CommandPalette.vue';
import SettingsModal from '../components/SettingsModal.vue';
import PasswordPrompt from '../components/PasswordPrompt.vue';

// --- State ---
const content = ref('');
const languageName = ref('Plain Text'); // UPDATED: Храним имя языка
const isReadOnly = ref(false);
const isLoading = ref(false);
const isSaving = ref(false);
const statusMessage = ref('');
const isPaletteOpen = ref(false); // UPDATED: Состояние палитры

const route = useRoute();
const router = useRouter();

const isSettingsOpen = ref(false);
const showPasswordPrompt = ref(false);
const settings = ref({
  ttl: 1209600, // 2 weeks
  burn: false,
  password: ''
});

// UPDATED: Динамические расширения для редактора
const extensions = computed(() => {
  const langExt = getExtensionByName(languageName.value);
  return [oneDark, langExt];
});

// --- Logic ---

const savePaste = async () => {
  if (!content.value.trim()) return;
  isSaving.value = true;
  statusMessage.value = 'Saving...';
  
  try {
    // UPDATED: Передаем объект настроек
    const paste = await api.createPaste({
      content: content.value,
      language: languageName.value,
      ttl: settings.value.ttl,
      password: settings.value.password,
      burn: settings.value.burn
    });
    await router.push({ name: 'view', params: { id: paste.id } });
    statusMessage.value = 'Saved!';
  } catch (e) {
    console.error(e);
    statusMessage.value = 'Error saving';
  } finally {
    isSaving.value = false;
  }
};

const loadPaste = async (id: string, pwd?: string) => {
  isLoading.value = true;
  isReadOnly.value = true;
  showPasswordPrompt.value = false; // Сброс
  statusMessage.value = 'Loading...';
  
  try {
    const paste = await api.getPaste(id, pwd);
    content.value = paste.content;
    languageName.value = paste.language;
    statusMessage.value = '';
    
    if (paste.burn_after_reading) {
        statusMessage.value = '🔥 Burned after reading';
    }
  } catch (e: any) {
    // UPDATED: Обработка защиты паролем
    if (e.status === 403 && e.isProtected) {
      showPasswordPrompt.value = true;
      statusMessage.value = 'Password required';
      content.value = ''; // Скрываем контент
    } else {
      statusMessage.value = 'Not found';
      content.value = '// Error 404: Data not found in the void.';
    }
  } finally {
    isLoading.value = false;
  }
};

const newPaste = () => {
  content.value = '';
  isReadOnly.value = false;
  languageName.value = 'Plain Text'; // Reset
  statusMessage.value = '';
  router.push({ name: 'home' });
};

// UPDATED: Actions for Command Palette
const openPalette = () => isPaletteOpen.value = true;
const closePalette = () => isPaletteOpen.value = false;

const setLanguage = (name: string) => {
  languageName.value = name;
  statusMessage.value = `Language: ${name}`;
  // Если мы в режиме редактирования, ничего не делаем, если просмотра - не даем менять (или даем?)
  // Лучше пока просто менять подсветку.
};

const copyLink = async () => {
  await navigator.clipboard.writeText(window.location.href);
  statusMessage.value = 'Link copied!';
  setTimeout(() => statusMessage.value = '', 2000);
};

const handlePasswordSubmit = (pwd: string) => {
  // Пробуем загрузить снова с паролем
  loadPaste(route.params.id as string, pwd);
};

// Обновление настроек из модалки
const applySettings = (newSettings: any) => {
  settings.value = newSettings;
  statusMessage.value = 'Settings applied';
};

// --- Lifecycle ---

// UPDATED: Global Hotkeys
const handleGlobalKeydown = (e: KeyboardEvent) => {
  // Ctrl+K or Cmd+K
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault();
    isPaletteOpen.value = !isPaletteOpen.value;
  }
  // Ctrl+S or Cmd+S to save
  if ((e.ctrlKey || e.metaKey) && e.key === 's') {
    e.preventDefault();
    if (!isReadOnly.value) savePaste();
  }
};

watch(
  () => route.params.id,
  (newId) => {
    if (newId) loadPaste(newId as string);
    else {
      isReadOnly.value = false;
      content.value = '';
      languageName.value = 'Plain Text';
    }
  },
  { immediate: true }
);

onMounted(() => window.addEventListener('keydown', handleGlobalKeydown));
onUnmounted(() => window.removeEventListener('keydown', handleGlobalKeydown));
</script>

<template>
  <div class="layout">
    <header class="navbar">
      <!-- ... Brand ... -->
      <div class="brand" @click="newPaste">
        <span class="logo-text">TRACE</span>
        <span class="status" v-if="statusMessage">:: {{ statusMessage }}</span>
      </div>

      <div class="actions">
        <!-- Кнопка языка -->
        <button class="btn-text" @click="openPalette" title="Cmd+K">
          {{ languageName }}
        </button>

        <!-- Кнопка настроек (только в режиме редактирования) -->
        <button v-if="!isReadOnly" class="btn-icon" @click="isSettingsOpen = true" title="Settings">
          ⚙️
        </button>
      
        <button v-if="!isReadOnly" class="btn save-btn" @click="savePaste" :disabled="isSaving">
          {{ isSaving ? 'SAVING...' : 'SAVE' }}
        </button>
        <button v-else class="btn new-btn" @click="newPaste">NEW</button>
      </div>
    </header>

    <main class="editor-container">
      <!-- Показываем промпт пароля ВМЕСТО редактора, если нужно -->
      <PasswordPrompt v-if="showPasswordPrompt" @submit="handlePasswordSubmit" />
      
      <codemirror
        v-else
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

    <!-- UPDATED: Вставляем палитру -->
    <CommandPalette 
      :is-open="isPaletteOpen" 
      @close="closePalette"
      @select-language="setLanguage"
      @copy-link="copyLink"
      @new-paste="newPaste"
    />
    <SettingsModal 
      :is-open="isSettingsOpen" 
      :defaults="settings"
      @close="isSettingsOpen = false"
      @apply="applySettings"
    />
  </div>
</template>

<style scoped>
.btn-icon {
  background: transparent;
  border: none;
  font-size: 1.2rem;
  cursor: pointer;
  margin-right: 1rem;
  padding: 0.2rem;
  border-radius: 4px;
  transition: background 0.2s;
}
.btn-icon:hover { background: rgba(255,255,255,0.1); }

.btn-text {
  background: transparent;
  border: 1px solid var(--border);
  color: var(--text-muted);
  padding: 0.4rem 0.8rem;
  border-radius: 4px;
  margin-right: 1rem;
  cursor: pointer;
  font-family: var(--font-mono);
  font-size: 0.8rem;
}
.btn-text:hover {
  color: var(--text-main);
  border-color: var(--text-muted);
}

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