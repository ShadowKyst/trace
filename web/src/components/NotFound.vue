<script setup lang="ts">
const emit = defineEmits<{ (e: 'home'): void }>();
</script>

<template>
  <div class="not-found-container">
    <div class="content">
      <h1 class="glitch" data-text="DATA NOT FOUND">DATA NOT FOUND</h1>
      <p class="subtitle">The paste you are looking for does not exist in the void.</p>
      
      <button class="btn-home" @click="$emit('home')">
        INITIALIZE NEW SEQUENCE
      </button>
    </div>
    
    <!-- Декоративные элементы -->
    <div class="scanline"></div>
    <div class="noise"></div>
  </div>
</template>

<style scoped>
.not-found-container {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  width: 100%;
  background-color: var(--bg-deep);
  position: relative;
  overflow: hidden;
  font-family: var(--font-mono);
}

.content {
  text-align: center;
  z-index: 2;
  position: relative;
}

/* --- Glitch Effect --- */
.glitch {
  font-size: 4rem;
  font-weight: 900;
  color: #fff;
  position: relative;
  letter-spacing: 4px;
  margin-bottom: 1rem;
}

.glitch::before,
.glitch::after {
  content: attr(data-text);
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: var(--bg-deep);
}

/* Красный слой смещен влево */
.glitch::before {
  left: 2px;
  text-shadow: -1px 0 #ff00c1;
  clip-path: inset(44% 0 61% 0);
  animation: glitch-anim-1 2.5s infinite linear alternate-reverse;
}

/* Синий слой смещен вправо */
.glitch::after {
  left: -2px;
  text-shadow: -1px 0 #00fff9;
  clip-path: inset(54% 0 10% 0);
  animation: glitch-anim-2 3s infinite linear alternate-reverse;
}

.subtitle {
  color: var(--text-muted);
  font-size: 1.1rem;
  margin-bottom: 3rem;
  opacity: 0.8;
}

/* --- Button --- */
.btn-home {
  background: transparent;
  border: 1px solid var(--primary);
  color: var(--primary);
  padding: 1rem 2rem;
  font-family: var(--font-mono);
  font-weight: bold;
  letter-spacing: 1px;
  cursor: pointer;
  transition: all 0.2s;
  text-transform: uppercase;
  position: relative;
  overflow: hidden;
}

.btn-home:hover {
  background: var(--primary);
  color: #000;
  box-shadow: 0 0 20px var(--primary-glow);
}

/* --- Background Effects --- */
.scanline {
  position: absolute;
  top: 0; left: 0; width: 100%; height: 100%;
  background: linear-gradient(to bottom, transparent 50%, rgba(0, 0, 0, 0.3) 51%);
  background-size: 100% 4px;
  pointer-events: none;
  z-index: 10;
}

.noise {
  position: absolute;
  top: 0; left: 0; width: 100%; height: 100%;
  opacity: 0.03;
  pointer-events: none;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noiseFilter'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.65' numOctaves='3' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noiseFilter)'/%3E%3C/svg%3E");
}

/* --- Keyframes for randomness --- */
@keyframes glitch-anim-1 {
  0% { clip-path: inset(20% 0 80% 0); }
  20% { clip-path: inset(60% 0 10% 0); }
  40% { clip-path: inset(40% 0 50% 0); }
  60% { clip-path: inset(80% 0 5% 0); }
  80% { clip-path: inset(10% 0 70% 0); }
  100% { clip-path: inset(30% 0 20% 0); }
}

@keyframes glitch-anim-2 {
  0% { clip-path: inset(10% 0 60% 0); }
  20% { clip-path: inset(30% 0 20% 0); }
  40% { clip-path: inset(70% 0 20% 0); }
  60% { clip-path: inset(20% 0 50% 0); }
  80% { clip-path: inset(50% 0 30% 0); }
  100% { clip-path: inset(5% 0 80% 0); }
}

/* Mobile responsive */
@media (max-width: 600px) {
  .glitch { font-size: 2.5rem; }
}
</style>