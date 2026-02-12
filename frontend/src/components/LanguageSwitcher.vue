<template>
  <div class="language-switcher">
    <button @click="toggleDropdown" class="lang-btn">
      <span class="flag">{{ currentFlag }}</span>
      <span class="code">{{ currentLocale.toUpperCase() }}</span>
      <span class="arrow">▼</span>
    </button>
    
    <div v-if="isOpen" class="lang-dropdown">
      <button 
        v-for="locale in availableLocales" 
        :key="locale.code" 
        @click="switchLanguage(locale.code)"
        class="lang-option"
        :class="{ active: currentLocale === locale.code }"
      >
        <span class="flag">{{ locale.flag }}</span>
        {{ locale.name }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'

const { locale } = useI18n()
const isOpen = ref(false)

const availableLocales = [
  { code: 'en', name: 'English', flag: '🇺🇸' },
  { code: 'fr', name: 'Français', flag: '🇫🇷' },
  { code: 'es', name: 'Español', flag: '🇪🇸' }
]

const currentLocale = computed(() => locale.value)
const currentFlag = computed(() => availableLocales.find(l => l.code === locale.value)?.flag || '🌐')

function toggleDropdown() {
  isOpen.value = !isOpen.value
}

function switchLanguage(code) {
  locale.value = code
  localStorage.setItem('language', code)
  isOpen.value = false
}

// Close dropdown when clicking outside
function handleClickOutside(event) {
  if (!event.target.closest('.language-switcher')) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.language-switcher {
  position: relative;
  font-family: 'Inter', sans-serif;
}

.lang-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  background: none;
  border: 1px solid var(--border);
  padding: 6px 12px;
  border-radius: 20px;
  cursor: pointer;
  font-size: 0.9rem;
  color: var(--text);
  transition: all 0.2s;
}

.lang-btn:hover {
  border-color: var(--primary);
  background: var(--bg);
}

.flag {
  font-size: 1.1rem;
}

.arrow {
  font-size: 0.7rem;
  opacity: 0.6;
}

.lang-dropdown {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 8px;
  background: white;
  border: 1px solid var(--border);
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  overflow: hidden;
  min-width: 140px;
  z-index: 1000;
  display: flex;
  flex-direction: column;
}

.lang-option {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 16px;
  background: none;
  border: none;
  text-align: left;
  cursor: pointer;
  font-size: 0.9rem;
  color: var(--text);
  transition: background 0.2s;
}

.lang-option:hover {
  background: var(--bg-warm);
}

.lang-option.active {
  background: var(--primary);
  color: white;
}
</style>
