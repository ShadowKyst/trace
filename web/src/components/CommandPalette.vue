<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue';
import { supportedLanguages } from '../config/languages';
// Импортируем иконки
import { 
  Search, 
  Globe, 
  Copy, 
  FilePlus, 
  ChevronRight, 
  Code2 
} from 'lucide-vue-next';

const props = defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'select-language', langName: string): void;
  (e: 'copy-link'): void;
  (e: 'new-paste'): void;
}>();

const search = ref('');
const searchInput = ref<HTMLInputElement | null>(null);
const selectedIndex = ref(0);
const mode = ref<'commands' | 'languages'>('commands');

// Определяем команды с иконками
const commands = [
  { 
    id: 'lang', 
    label: 'Change Language Mode...', 
    icon: Globe,
    action: () => { mode.value = 'languages'; search.value = ''; } 
  },
  { 
    id: 'copy', 
    label: 'Copy Link to Clipboard', 
    icon: Copy,
    action: () => emit('copy-link') 
  },
  { 
    id: 'new',  
    label: 'Create New Paste', 
    icon: FilePlus,
    action: () => emit('new-paste') 
  },
];

// Фильтрация
const filteredItems = computed(() => {
  const query = search.value.toLowerCase();
  
  if (mode.value === 'languages') {
    return supportedLanguages
      .filter(l => l.name.toLowerCase().includes(query))
      .map(l => ({ 
        id: l.name, 
        label: l.name, 
        icon: Code2, // Иконка для языков
        action: () => emit('select-language', l.name) 
      }));
  }
  
  return commands.filter(c => c.label.toLowerCase().includes(query));
});

// Сброс индекса при поиске
watch(search, () => selectedIndex.value = 0);

// Навигация
const handleKeydown = (e: KeyboardEvent) => {
  if (!props.isOpen) return;

  if (e.key === 'ArrowDown') {
    selectedIndex.value = (selectedIndex.value + 1) % filteredItems.value.length;
    e.preventDefault();
  } else if (e.key === 'ArrowUp') {
    selectedIndex.value = (selectedIndex.value - 1 + filteredItems.value.length) % filteredItems.value.length;
    e.preventDefault();
  } else if (e.key === 'Enter') {
    execute(filteredItems.value[selectedIndex.value]);
    e.preventDefault();
  } else if (e.key === 'Escape') {
    if (mode.value === 'languages' && search.value === '') {
       mode.value = 'commands';
    } else {
       close();
    }
  }
};

const execute = (item: any) => {
  if (item) {
    item.action();
    if (item.id !== 'lang') {
      close();
    }
  }
};

const close = () => {
  search.value = '';
  mode.value = 'commands';
  selectedIndex.value = 0;
  emit('close');
};

watch(() => props.isOpen, (newVal) => {
  if (newVal) {
    nextTick(() => searchInput.value?.focus());
  }
});

onMounted(() => window.addEventListener('keydown', handleKeydown));
onUnmounted(() => window.removeEventListener('keydown', handleKeydown));
</script>

<template>
  <div v-if="isOpen" class="overlay" @click.self="close">
    <div class="palette glass-panel">
      <!-- Input Area -->
      <div class="input-wrapper">
        <Search :size="18" class="search-icon" />
        <input 
          ref="searchInput"
          v-model="search" 
          type="text" 
          :placeholder="mode === 'commands' ? 'Type a command...' : 'Select language...'"
        >
        <div class="badge" v-if="mode === 'languages'">LANG</div>
      </div>
      
      <!-- List Area -->
      <ul class="list">
        <li 
          v-for="(item, index) in filteredItems" 
          :key="item.id || item.label"
          :class="{ active: index === selectedIndex }"
          @click="execute(item)"
          @mouseenter="selectedIndex = index"
        >
          <component :is="item.icon" :size="16" class="item-icon" />
          <span class="label">{{ item.label }}</span>
          
          <ChevronRight v-if="item.id === 'lang'" :size="14" class="arrow-right" />
        </li>
        
        <li v-if="filteredItems.length === 0" class="empty">
          No results found
        </li>
      </ul>
      
      <!-- Footer hints -->
      <div class="footer">
        <span><kbd>↑↓</kbd> to navigate</span>
        <span><kbd>↵</kbd> to select</span>
        <span><kbd>esc</kbd> to close</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.overlay {
  position: fixed;
  top: 0; left: 0; width: 100%; height: 100%;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  justify-content: center;
  padding-top: 15vh;
  z-index: 200;
}

.palette {
  width: 600px;
  max-width: 90%;
  /* Используем глобальный класс glass-panel, но переопределяем цвет для большей непрозрачности */
  background: rgba(10, 10, 12, 0.85); 
  border-radius: 12px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  max-height: 400px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.5);
}

.input-wrapper {
  display: flex;
  align-items: center;
  padding: 1rem 1.2rem;
  border-bottom: 1px solid var(--glass-border);
  gap: 0.8rem;
}

.search-icon {
  color: var(--text-muted);
}

input {
  background: transparent;
  border: none;
  color: var(--text-main);
  font-size: 1.1rem;
  flex: 1;
  outline: none;
  font-family: var(--font-sans);
}

.badge {
  font-size: 0.7rem;
  background: var(--primary-glow);
  border: 1px solid var(--primary);
  color: var(--primary);
  padding: 2px 6px;
  border-radius: 4px;
  font-weight: bold;
}

.list {
  list-style: none;
  overflow-y: auto;
  padding: 0.5rem;
  flex: 1;
}

.list li {
  padding: 0.8rem 1rem;
  cursor: pointer;
  color: var(--text-muted);
  font-size: 0.95rem;
  display: flex;
  align-items: center;
  gap: 0.8rem;
  border-radius: 8px;
  margin-bottom: 2px;
  transition: all 0.1s;
}

.item-icon {
  opacity: 0.7;
}

.label {
  flex: 1;
}

.list li.active {
  background: rgba(255, 255, 255, 0.08);
  color: var(--text-main);
}
.list li.active .item-icon {
  color: var(--primary);
  opacity: 1;
}

.arrow-right {
  opacity: 0.5;
}

.empty {
  padding: 2rem;
  color: var(--text-muted);
  text-align: center;
  font-style: italic;
}

.footer {
  padding: 0.6rem 1.2rem;
  border-top: 1px solid var(--glass-border);
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  background: rgba(0,0,0,0.2);
}

.footer span {
  font-size: 0.75rem;
  color: var(--text-muted);
  display: flex;
  align-items: center;
  gap: 4px;
}

kbd {
  background: rgba(255,255,255,0.1);
  padding: 2px 4px;
  border-radius: 3px;
  font-family: var(--font-mono);
  font-size: 0.7rem;
}
</style>