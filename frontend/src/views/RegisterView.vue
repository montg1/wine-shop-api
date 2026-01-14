<template>
  <div class="auth-page">
    <div class="auth-card">
      <h1>Create Account</h1>
      <p class="subtitle">Join our wine community</p>
      <form @submit.prevent="handleRegister">
        <div class="form-group">
          <label>Email</label>
          <input v-model="email" type="email" placeholder="you@example.com" required />
        </div>
        <div class="form-group">
          <label>Password</label>
          <input v-model="password" type="password" placeholder="••••••••" required />
        </div>
        <div class="form-group">
          <label>Confirm Password</label>
          <input v-model="confirmPassword" type="password" placeholder="••••••••" required />
        </div>
        <button type="submit" class="btn btn-primary btn-block" :disabled="loading">
          {{ loading ? 'Creating...' : 'Create Account' }}
        </button>
        <p v-if="error" class="error">{{ error }}</p>
        <p v-if="success" class="success">Account created! Redirecting to login...</p>
      </form>
      <p class="switch-link">
        Already have an account? <router-link to="/login">Sign in</router-link>
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
const confirmPassword = ref('')
const loading = ref(false)
const error = ref('')
const success = ref(false)

const handleRegister = async () => {
  if (password.value !== confirmPassword.value) {
    error.value = 'Passwords do not match'
    return
  }
  
  loading.value = true
  error.value = ''
  try {
    await authStore.register(email.value, password.value)
    success.value = true
    setTimeout(() => router.push('/login'), 1500)
  } catch (err) {
    error.value = err.response?.data?.error || 'Registration failed'
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
  padding: 40px 20px;
  background: var(--bg-warm);
}

.auth-card {
  background: var(--card-bg);
  padding: 50px;
  border-radius: 8px;
  width: 100%;
  max-width: 420px;
  box-shadow: var(--shadow-md);
  border: 1px solid var(--border);
}

.auth-card h1 {
  text-align: center;
  font-size: 2rem;
  color: var(--primary);
  margin-bottom: 5px;
}

.subtitle {
  text-align: center;
  color: var(--text-muted);
  margin-bottom: 35px;
}

.form-group {
  margin-bottom: 22px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: var(--text);
  font-size: 0.9rem;
  font-weight: 500;
}

.form-group input {
  width: 100%;
  padding: 14px 16px;
  border-radius: 6px;
  border: 1px solid var(--border);
  background: var(--bg);
  color: var(--text);
  font-size: 1rem;
  transition: border-color 0.3s;
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
  color: #c44;
  text-align: center;
  margin-top: 15px;
  font-size: 0.9rem;
}

.success {
  color: #4a4;
  text-align: center;
  margin-top: 15px;
  font-size: 0.9rem;
}

.switch-link {
  text-align: center;
  margin-top: 30px;
  color: var(--text-muted);
  font-size: 0.95rem;
}

.switch-link a {
  color: var(--primary);
}
</style>
