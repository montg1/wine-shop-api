<template>
  <div id="app">
    <header class="navbar">
      <router-link to="/" class="logo">
        <img src="/logo.png" alt="Wine Shop Logo" class="logo-img" />
        Wine Shop
      </router-link>
      <nav class="nav-links">
        <router-link to="/products">{{ $t('nav.wines') }}</router-link>
        <router-link v-if="authStore.isLoggedIn" to="/cart" class="cart-link">
          🛒 <span v-if="cartStore.totalItems" class="badge">{{ cartStore.totalItems }}</span>
        </router-link>
        <router-link v-if="authStore.isLoggedIn" to="/orders">{{ $t('nav.orders') }}</router-link>
        <!-- Admin link only visible to admin users -->
        <router-link v-if="authStore.isAdmin" to="/admin" class="admin-link">⚙️ {{ $t('nav.admin') }}</router-link>
        
        <LanguageSwitcher />
        
        <template v-if="authStore.isLoggedIn">
          <button @click="handleLogout" class="btn-link">{{ $t('nav.logout') }}</button>
        </template>
        <template v-else>
          <router-link to="/login">{{ $t('nav.login') }}</router-link>
        </template>
      </nav>
    </header>
    
    <main>
      <router-view />
    </main>
    
    <footer>
      <p>© 2026 Wine Shop. All rights reserved.</p>
    </footer>

    <!-- Wine Chatbot Widget -->
    <ChatbotWidget />
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from './stores/auth'
import { useCartStore } from './stores/cart'
import ChatbotWidget from './components/ChatbotWidget.vue'
import LanguageSwitcher from './components/LanguageSwitcher.vue'

const router = useRouter()
const authStore = useAuthStore()
const cartStore = useCartStore()

onMounted(async () => {
  // Initialize auth state (fetch user info if token exists)
  await authStore.init()
  if (authStore.isLoggedIn) {
    cartStore.fetchCart()
  }
})

const handleLogout = () => {
  authStore.logout()
  router.push('/')
}
</script>

<style>
/* Styles are in style.css */
</style>
