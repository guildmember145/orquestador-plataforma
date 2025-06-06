<template>
  <div class="auth-form">
    <h2>Login</h2>
    <form @submit.prevent="handleLogin">
      <div class="input-group">
        <label for="login-email">Email:</label>
        <input type="email" id="login-email" v-model="email" required />
      </div>
      <div class="input-group">
        <label for="login-password">Password:</label>
        <input type="password" id="login-password" v-model="password" required />
      </div>
      <button type="submit" :disabled="authStore.isLoading" class="primary-action">
        {{ authStore.isLoading ? 'Logging in...' : 'Login' }}
      </button>
      <p v-if="authStore.error" class="error">{{ authStore.error }}</p>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '../../stores/auth'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const router = useRouter()
const email = ref('')
const password = ref('')

const handleLogin = async () => {
  try {
    await authStore.login({ email: email.value, password: password.value })
    console.log('Login successful, user:', authStore.currentUser)
    // Redirigir al dashboard después del login exitoso
    router.push('/dashboard')
  } catch (error) {
    console.error('Login failed:', authStore.error)
  }
}
</script>

<style scoped>
.auth-form {
  width: 100%;
  max-width: 400px;
  background-color: var(--color-surface);
  padding: 30px 40px;
  border-radius: 8px;
  border: 1px solid var(--color-border);
}

.auth-form h2 {
  text-align: center;
  margin-bottom: 25px;
  color: var(--color-text-primary);
}

.input-group {
  margin-bottom: 20px;
}

.input-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: bold;
  font-size: 0.9em;
  color: var(--color-text-secondary);
}

.input-group input {
  width: 100%;
  padding: 12px;
  background-color: var(--color-background);
  border: 1px solid var(--color-border);
  color: var(--color-text-primary);
  border-radius: 5px;
  box-sizing: border-box;
  transition: border-color 0.3s;
}

.input-group input:focus {
  border-color: var(--color-accent);
  outline: none;
}

button.primary-action {
  width: 100%;
  padding: 12px;
  font-size: 1em;
  font-weight: bold;
  /* El estilo de .primary-action ya está definido globalmente en App.vue */
}

.error {
  color: var(--color-error);
  font-size: 0.9em;
  margin-top: 15px;
  text-align: center;
}
</style>