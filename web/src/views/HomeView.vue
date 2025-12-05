<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, computed } from 'vue'; // UPDATED: добавил onUnmounted, computed
import { useRoute, useRouter } from 'vue-router';
import { Codemirror } from 'vue-codemirror';
import { oneDark } from '@codemirror/theme-one-dark';
import { 
  Settings, 
  Save, 
  Plus, 
  ChevronDown, 
  Copy, 
  Lock 
} from 'lucide-vue-next';

import api from '../services/api';

// UPDATED: Импорты
import { getExtensionByName } from '../config/languages'; 
import CommandPalette from '../components/CommandPalette.vue';
import SettingsModal from '../components/SettingsModal.vue';
import PasswordPrompt from '../components/PasswordPrompt.vue';
import PasteCreatedModal from '../components/PasteCreatedModal.vue';
import NotFound from '../components/NotFound.vue';

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
const isNotFound = ref(false); // Новый стейт
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

const showCreatedModal = ref(false);
const createdPasteUrl = ref('');
const createdPasteBurn = ref(false);

// --- Logic ---

const savePaste = async () => {
  if (!content.value.trim()) return;
  isSaving.value = true;
  statusMessage.value = 'Securing data...'; // Более "хакерский" текст
  
  try {
    const paste = await api.createPaste({
      content: content.value,
      language: languageName.value,
      ttl: settings.value.ttl,
      password: settings.value.password,
      burn: settings.value.burn
    });
    
    // НОВАЯ ЛОГИКА:
    createdPasteUrl.value = `${window.location.origin}/${paste.id}`;
    createdPasteBurn.value = settings.value.burn ?? false;
    
    // Очищаем редактор, чтобы не сохранили дважды то же самое
    newPaste(); 
    
    // Показываем окно успеха
    showCreatedModal.value = true;
    statusMessage.value = 'Data secured';
    
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
  showPasswordPrompt.value = false;
  isNotFound.value = false; // Сбрасываем перед загрузкой
  statusMessage.value = 'Connecting to uplink...';
  
  try {
    const paste = await api.getPaste(id, pwd);
    content.value = paste.content;
    languageName.value = paste.language;
    statusMessage.value = ''; // Успех, убираем статус
    
    if (paste.burn_after_reading) {
        statusMessage.value = '🔥 Burned after reading';
    }
  } catch (e: any) {
    if (e.status === 403 && e.isProtected) {
      showPasswordPrompt.value = true;
      statusMessage.value = 'Encryption detected';
    } else {
      // ВОТ ТУТ МЕНЯЕМ ЛОГИКУ:
      isNotFound.value = true; 
      statusMessage.value = 'Signal lost';
    }
  } finally {
    isLoading.value = false;
  }
};

const newPaste = () => {
  isNotFound.value = false; // Скрываем 404
  showPasswordPrompt.value = false;
  content.value = '';
  isReadOnly.value = false;
  languageName.value = 'Plain Text';
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
    
    <!-- Floating Header -->
    <header class="navbar-container">
      <div class="glass-panel navbar">
        
        <!-- Logo Area -->
        <div class="brand" @click="newPaste">
          <div class="logo-box">TRACE</div> <!-- Лого как на скрине -->
        </div>

        <!-- Central Info (Filename / Status) -->
        <div class="status-bar">
           <!-- Можно выводить имя файла или статус -->
           <span v-if="statusMessage" class="status-text">{{ statusMessage }}</span>
           <span v-else class="file-name">untitled.txt</span>
        </div>
        
        <!-- Actions -->
        <div class="actions">
          <!-- Language Selector -->
          <button class="btn-ghost" @click="openPalette" title="Change Language">
            {{ languageName }}
            <ChevronDown :size="14" class="icon-right" />
          </button>

          <!-- Settings -->
          <button v-if="!isReadOnly" class="btn-icon" @click="isSettingsOpen = true">
            <Settings :size="18" />
          </button>

          <!-- Save / New Button -->
          <div class="divider"></div>
          
          <button v-if="!isReadOnly" class="btn-primary" @click="savePaste" :disabled="isSaving">
            <Save :size="16" class="icon-left" v-if="!isSaving"/>
            <span>{{ isSaving ? 'Saving...' : 'Save' }}</span>
          </button>
          
          <button v-else class="btn-primary" @click="newPaste">
            <Plus :size="16" class="icon-left"/>
            <span>New</span>
          </button>
        </div>
      </div>
    </header>

    <main class="editor-wrap">
      <NotFound v-if="isNotFound" @home="newPaste" />
      <!-- Показываем промпт пароля ВМЕСТО редактора, если нужно -->
      <PasswordPrompt v-if="showPasswordPrompt" @submit="handlePasswordSubmit" />
      
      <codemirror
        v-else
        v-model="content"
        placeholder="// Type something..."
        :style="{ height: '100%', fontSize: '14px', backgroundColor: 'transparent' }"
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
    <PasteCreatedModal
      :is-open="showCreatedModal"
      :url="createdPasteUrl"
      :is-burn="createdPasteBurn"
      @close="showCreatedModal = false"
    />
  </div>
</template>

<style scoped>
.layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  position: relative;
  /* Добавляем отступы, чтобы редактор не прилипал к краям, как на макете */
  padding: 80px 20px 20px 20px; 
}

/* Floating Navbar Container */
.navbar-container {
  position: absolute;
  top: 20px;
  left: 0;
  width: 100%;
  display: flex;
  justify-content: center;
  z-index: 10;
  padding: 0 1rem;
}

.navbar {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  padding: 0.5rem 0.5rem 0.5rem 1rem;
  border-radius: 12px;
  min-width: 600px; /* Ширина островка */
  max-width: 100%;
  justify-content: space-between;
}

/* Logo Box */
.logo-box {
  background: linear-gradient(135deg, var(--primary), var(--secondary));
  color: white;
  font-family: var(--font-mono);
  font-weight: bold;
  font-size: 0.9rem;
  padding: 4px 8px;
  border-radius: 6px;
  letter-spacing: 1px;
}

.status-bar {
  flex: 1;
  text-align: center;
  font-family: var(--font-mono);
  font-size: 0.85rem;
  color: var(--text-muted);
}

.actions {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.divider {
  width: 1px;
  height: 20px;
  background: var(--glass-border);
  margin: 0 0.5rem;
}

/* Buttons */
.btn-ghost {
  background: transparent;
  border: none;
  color: var(--text-muted);
  font-size: 0.85rem;
  padding: 6px 10px;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 6px;
  transition: all 0.2s;
}
.btn-ghost:hover {
  background: rgba(255, 255, 255, 0.05);
  color: var(--text-main);
}

.btn-icon {
  background: transparent;
  border: none;
  color: var(--text-muted);
  padding: 6px;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  transition: color 0.2s;
}
.btn-icon:hover { color: var(--text-main); }

.btn-primary {
  background: var(--primary);
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 0.85rem;
  font-weight: 500;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 6px;
  box-shadow: 0 4px 15px var(--primary-glow);
  transition: transform 0.2s, box-shadow 0.2s;
}
.btn-primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 20px var(--primary-glow);
}
.btn-primary:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.editor-wrap {
  flex: 1;
  border-radius: 12px;
  overflow: hidden;
  /* Тонкая рамка вокруг редактора */
  border: 1px solid var(--glass-border); 
  background: rgba(0,0,0,0.2); 
}

/* Переопределяем фон CodeMirror, чтобы он был прозрачным */
:deep(.cm-editor) { background-color: transparent !important; }
:deep(.cm-gutters) { background-color: transparent !important; border-right: 1px solid rgba(255,255,255,0.05); }
</style>