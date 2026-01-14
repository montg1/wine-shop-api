<template>
  <div class="auth-page">
    <div class="auth-card">
      <h1>Login</h1>
      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label>Email</label>
          <input v-model="email" type="email" placeholder="you@example.com" required />
        </div>
        <div class="form-group">
          <label>Password</label>
          <input v-model="password" type="password" placeholder="••••••••" required />
        </div>
        <button type="submit" class="btn btn-primary btn-block" :disabled="loading">
          {{ loading ? 'Logging in...' : 'Login' }}
        </button>
        <p v-if="error" class="error">{{ error }}</p>
      </form>
      <p class="switch-link">
        Don't have an account? <router-link to="/register">Register</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
  loading.value = true
  error.value = ''
  try {
    await authStore.login(email.value, password.value)
    router.push('/products')
  } catch (err) {
    error.value = 'Invalid email or password'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-page {
  min-height: calc(100vh - 80px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.auth-card {
  background: var(--card-bg);
  padding: 40px;
  border-radius: 20px;
  width: 100%;
  max-width: 400px;
}

.auth-card h1 {
  text-align: center;
  margin-bottom: 30px;
  color: var(--primary);
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #ccc;
}

.form-group input {
  width: 100%;
  padding: 14px;
  border-radius: 10px;
  border: 1px solid #444;
  background: #1a1a2e;
  color: #fff;
  font-size: 1rem;
}

.form-group input:focus {
  outline: none;
  border-color: var(--primary);
}

.btn-block {
  width: 100%;
  margin-top: 10px;
}

.error {
  color: #e55;
  text-align: center;
  margin-top: 15px;
}

.switch-link {
  text-align: center;
  margin-top: 25px;
  color: #888;
}

.switch-link a {
  color: var(--primary);
  text-decoration: none;
}

.switch-link a:hover {
  text-decoration: underline;
}
</style>
