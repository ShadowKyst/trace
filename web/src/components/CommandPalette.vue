<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue';
import { supportedLanguages } from '../config/languages';

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

// Режимы палитры: 'commands' (главная) или 'languages' (выбор языка)
const mode = ref<'commands' | 'languages'>('commands');

// Список команд
const commands = [
  { id: 'lang', label: 'Change Language...', action: () => { mode.value = 'languages'; search.value = ''; } },
  { id: 'copy', label: 'Copy Link to Clipboard', action: () => emit('copy-link') },
  { id: 'new',  label: 'New Paste', action: () => emit('new-paste') },
];

// Фильтрация списков
const filteredItems = computed(() => {
  const query = search.value.toLowerCase();
  
  if (mode.value === 'languages') {
    return supportedLanguages
      .filter(l => l.name.toLowerCase().includes(query))
      .map(l => ({ id: l.name, label: l.name, action: () => emit('select-language', l.name) }));
  }
  
  return commands.filter(c => c.label.toLowerCase().includes(query));
});

// Навигация клавишами
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
  } else if (e.key === 'Escape') {
    if (mode.value === 'languages' && search.value === '') {
       mode.value = 'commands'; // Назад
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

// Фокус при открытии
import { watch } from 'vue';
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
    <div class="palette">
      <div class="input-wrapper">
        <span class="prompt">></span>
        <input 
          ref="searchInput"
          v-model="search" 
          type="text" 
          placeholder="Type a command..." 
          @input="selectedIndex = 0"
        >
        <span class="badge" v-if="mode === 'languages'">LANG</span>
      </div>
      
      <ul class="list">
        <li 
          v-for="(item, index) in filteredItems" 
          :key="item.id || item.label"
          :class="{ active: index === selectedIndex }"
          @click="execute(item)"
          @mouseenter="selectedIndex = index"
        >
          {{ item.label }}
        </li>
        <li v-if="filteredItems.length === 0" class="empty">No results found</li>
      </ul>
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
  z-index: 100;
}

.palette {
  width: 500px;
  max-width: 90%;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 8px;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.5), 0 0 0 1px var(--border);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  max-height: 400px;
}

.input-wrapper {
  display: flex;
  align-items: center;
  padding: 1rem;
  border-bottom: 1px solid var(--border);
}

.prompt {
  color: var(--accent);
  font-weight: bold;
  margin-right: 0.8rem;
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
  background: var(--border);
  padding: 2px 6px;
  border-radius: 4px;
  color: var(--text-muted);
}

.list {
  list-style: none;
  overflow-y: auto;
  padding: 0.5rem 0;
}

.list li {
  padding: 0.75rem 1rem;
  cursor: pointer;
  color: var(--text-muted);
  font-size: 0.95rem;
  display: flex;
  align-items: center;
}

.list li.active {
  background: rgba(139, 92, 246, 0.1); /* Прозрачный accent */
  color: var(--text-main);
  border-left: 2px solid var(--accent);
}

.empty {
  padding: 1rem;
  color: var(--text-muted);
  text-align: center;
  font-style: italic;
}
</style>