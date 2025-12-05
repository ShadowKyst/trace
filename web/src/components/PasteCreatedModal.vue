<script setup lang="ts">
import { ref } from 'vue';
import { Copy, Check, ExternalLink } from 'lucide-vue-next';

const props = defineProps<{
  isOpen: boolean;
  url: string;
  isBurn: boolean;
}>();

const emit = defineEmits<{ (e: 'close'): void }>();

const copied = ref(false);

const copyToClipboard = async () => {
  await navigator.clipboard.writeText(props.url);
  copied.value = true;
  setTimeout(() => copied.value = false, 2000);
};
</script>

<template>
  <div v-if="isOpen" class="overlay">
    <div class="modal glass-panel">
      <div class="header">
        <h3>Paste Secured</h3>
        <div class="badge">Success</div>
      </div>
      
      <p class="desc">
        Your paste is ready. Share this link to allow access.
      </p>

      <div class="link-box">
        <input type="text" :value="url" readonly>
        <button class="btn-copy" @click="copyToClipboard">
          <Check v-if="copied" :size="18" color="#4ade80" />
          <Copy v-else :size="18" />
        </button>
      </div>

      <div v-if="isBurn" class="warning-box">
        <strong>⚠️ Burn After Reading Active</strong>
        <p>Do not open this link yourself! The paste will be destroyed immediately after the first view.</p>
      </div>

      <div class="actions">
        <button class="btn-primary" @click="$emit('close')">Done</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.overlay {
  position: fixed; top: 0; left: 0; width: 100%; height: 100%;
  background: rgba(0,0,0,0.8);
  backdrop-filter: blur(5px);
  display: flex; justify-content: center; align-items: center; z-index: 100;
}

.modal {
  width: 450px;
  max-width: 90%;
  padding: 2rem;
  border-radius: 16px;
  background: var(--bg-surface);
  border: 1px solid var(--border);
}

.header { display: flex; align-items: center; gap: 1rem; margin-bottom: 0.5rem; }
h3 { font-size: 1.5rem; color: var(--text-main); }
.badge { background: rgba(74, 222, 128, 0.1); color: #4ade80; padding: 2px 8px; border-radius: 4px; font-size: 0.8rem; border: 1px solid rgba(74, 222, 128, 0.2); }

.desc { color: var(--text-muted); margin-bottom: 1.5rem; }

.link-box {
  display: flex; gap: 0.5rem; margin-bottom: 1.5rem;
}
input {
  flex: 1; background: var(--bg-deep); border: 1px solid var(--border);
  color: var(--accent); padding: 0.8rem; border-radius: 8px; font-family: var(--font-mono); outline: none;
}
.btn-copy {
  background: var(--bg-deep); border: 1px solid var(--border); color: var(--text-main);
  width: 48px; display: flex; align-items: center; justify-content: center;
  border-radius: 8px; cursor: pointer; transition: all 0.2s;
}
.btn-copy:hover { border-color: var(--accent); }

.warning-box {
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.2);
  padding: 1rem; border-radius: 8px; margin-bottom: 1.5rem;
  color: #fca5a5; font-size: 0.9rem;
}
.warning-box strong { display: block; margin-bottom: 0.3rem; color: #f87171; }

.btn-primary {
  width: 100%; background: var(--primary); color: white; border: none; padding: 0.8rem;
  border-radius: 8px; font-weight: 600; cursor: pointer;
}
.btn-primary:hover { background: #565add; }
</style>