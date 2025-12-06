<script setup lang="ts">
import { ref, watch } from 'vue';

const props = defineProps<{
  isOpen: boolean;
  defaults: { ttl: number; burn: boolean; password: string; wordWrap: boolean; };
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'apply', settings: { ttl: number; burn: boolean; password: string; wordWrap: boolean }): void;
}>();

// Варианты TTL в секундах
const ttlOptions = [
  { label: '1 Hour', value: 3600 },
  { label: '1 Day', value: 86400 },
  { label: '1 Week', value: 604800 },
  { label: '2 Weeks', value: 1209600 },
];

const form = ref({
  ttl: 1209600, // Default 2 weeks
  burn: false,
  password: '',
  wordWrap: false
});

// Синхронизация при открытии
watch(() => props.isOpen, (val) => {
  if (val) {
    form.value = { ...props.defaults };
  }
});

const save = () => {
  emit('apply', form.value);
  emit('close');
};
</script>

<template>
  <div v-if="isOpen" class="overlay" @click.self="$emit('close')">
    <div class="modal">
      <h3 class="title">Paste Settings</h3>
      
      <div class="field">
        <label>Expiration (TTL)</label>
        <select v-model="form.ttl">
          <option v-for="opt in ttlOptions" :key="opt.value" :value="opt.value">
            {{ opt.label }}
          </option>
        </select>
      </div>

      <div class="field">
        <label>Password Protection</label>
        <input 
          type="password" 
          v-model="form.password" 
          placeholder="Optional..."
          autocomplete="new-password"
        >
      </div>

      <div class="field checkbox-field">
        <label>
        <input type="checkbox" v-model="form.wordWrap">
        <span>Word Wrap</span>
        </label>
        <p class="hint">Break long lines to fit the screen.</p>
    </div>

      <div class="field checkbox-field">
        <label>
          <input type="checkbox" v-model="form.burn">
          <span>Burn after reading</span>
        </label>
        <p class="hint">Paste will be deleted immediately after first access.</p>
      </div>

      <div class="actions">
        <button class="btn secondary" @click="$emit('close')">Cancel</button>
        <button class="btn primary" @click="save">Apply</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.overlay {
  position: fixed;
  top: 0; left: 0; width: 100%; height: 100%;
  background: rgba(0, 0, 0, 0.7);
  backdrop-filter: blur(2px);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 200;
}

.modal {
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  padding: 1.5rem;
  border-radius: 8px;
  width: 400px;
  max-width: 90%;
  box-shadow: 0 10px 40px rgba(0,0,0,0.5);
}

.title {
  color: var(--text-main);
  margin-bottom: 1.5rem;
  font-family: var(--font-mono);
}

.field {
  margin-bottom: 1.2rem;
  display: flex;
  flex-direction: column;
}

label {
  color: var(--text-muted);
  font-size: 0.9rem;
  margin-bottom: 0.4rem;
}

input[type="password"], select {
  background: var(--bg-primary);
  border: 1px solid var(--border);
  color: var(--text-main);
  padding: 0.6rem;
  border-radius: 4px;
  font-family: var(--font-sans);
  outline: none;
}

input[type="password"]:focus, select:focus {
  border-color: var(--accent);
}

.checkbox-field label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: var(--text-main);
  cursor: pointer;
}

.hint {
  font-size: 0.75rem;
  color: var(--text-muted);
  margin-top: 0.2rem;
  margin-left: 1.5rem;
}

.actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 2rem;
}

.btn {
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
  border: none;
  font-weight: 500;
}

.btn.primary { background: var(--accent); color: white; }
.btn.secondary { background: transparent; color: var(--text-muted); border: 1px solid var(--border); }
.btn.secondary:hover { color: var(--text-main); border-color: var(--text-muted); }

select {
  appearance: none; /* Убираем системный стиль */
  cursor: pointer;
  background-image: url("data:image/svg+xml;charset=UTF-8,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='white' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3e%3cpolyline points='6 9 12 15 18 9'%3e%3c/polyline%3e%3c/svg%3e");
  background-repeat: no-repeat;
  background-position: right 1rem center;
  background-size: 1em;
  padding-right: 2.5rem;
  background-color: var(--bg-primary); /* Явно задаем фон самого селекта */
  color: var(--text-main);
}

select option {
  background-color: #18181b !important; /* var(--bg-secondary) хардкодом для надежности */
  color: #ffffff !important;
  padding: 10px;
}

@media (prefers-color-scheme: dark) {
  select {
    background-color: var(--bg-primary);
    color: white;
  }
}
</style>