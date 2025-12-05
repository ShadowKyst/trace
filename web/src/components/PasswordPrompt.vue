<script setup lang="ts">
import { ref } from 'vue';

const password = ref('');
const emit = defineEmits<{ (e: 'submit', pwd: string): void }>();

const submit = () => {
  if (password.value) emit('submit', password.value);
};
</script>

<template>
  <div class="prompt-container">
    <div class="lock-icon">🔒</div>
    <h2>Protected Paste</h2>
    <p>This content is password protected.</p>
    
    <form @submit.prevent="submit" class="form">
      <input 
        type="password" 
        v-model="password" 
        placeholder="Enter password..." 
        autofocus
      >
      <button type="submit" class="btn">Unlock</button>
    </form>
  </div>
</template>

<style scoped>
.prompt-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--text-main);
}

.lock-icon { font-size: 3rem; margin-bottom: 1rem; }
h2 { margin-bottom: 0.5rem; }
p { color: var(--text-muted); margin-bottom: 2rem; }

.form { display: flex; gap: 0.5rem; }
input {
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  color: var(--text-main);
  padding: 0.6rem 1rem;
  border-radius: 4px;
  outline: none;
  width: 200px;
}
input:focus { border-color: var(--accent); }

.btn {
  background: var(--accent);
  color: white;
  border: none;
  padding: 0 1rem;
  border-radius: 4px;
  cursor: pointer;
}
</style>