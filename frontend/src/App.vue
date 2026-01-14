<template>
  <div id="app">
    <header class="navbar">
      <router-link to="/" class="logo">🍷 Wine Shop</router-link>
      <nav class="nav-links">
        <router-link to="/products">Wines</router-link>
        <router-link v-if="authStore.isLoggedIn" to="/cart" class="cart-link">
          🛒 <span v-if="cartStore.totalItems" class="badge">{{ cartStore.totalItems }}</span>
        </router-link>
        <router-link v-if="authStore.isLoggedIn" to="/orders">Orders</router-link>
        <router-link v-if="authStore.isLoggedIn" to="/admin" class="admin-link">⚙️ Admin</router-link>
        <template v-if="authStore.isLoggedIn">
          <button @click="handleLogout" class="btn-link">Logout</button>
        </template>
        <template v-else>
          <router-link to="/login">Login</router-link>
        </template>
      </nav>
    </header>
    
    <main>
      <router-view />
    </main>
    
    <footer>
      <p>© 2026 Wine Shop. All rights reserved.</p>
    </footer>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from './stores/auth'
import { useCartStore } from './stores/cart'

const router = useRouter()
const authStore = useAuthStore()
const cartStore = useCartStore()

onMounted(() => {
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
