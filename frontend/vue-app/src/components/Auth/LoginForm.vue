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
      <button type="submit" :disabled="authStore.isLoading">
        {{ authStore.isLoading ? 'Logging in...' : 'Login' }}
      </button>
      <p v-if="authStore.error" class="error">{{ authStore.error }}</p>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useAuthStore } from '../../stores/auth';

const authStore = useAuthStore();
const email = ref('');
const password = ref('');

const handleLogin = async () => {
  try {
    await authStore.login({ email: email.value, password: password.value });
    console.log('Login successful, user:', authStore.currentUser);
  } catch (error) {
    console.error('Login failed:', authStore.error);
  }
};
</script>

<style scoped>
.auth-form {
  max-width: 400px;
  margin: auto;
  padding: 30px;
  background: #22a568;
  border-radius: 12px;
  box-shadow: 0px 4px 10px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.input-group {
  display: flex;
  flex-direction: column;
  margin-bottom: 15px;
}

label {
  font-weight: bold;
  margin-bottom: 5px;
}

input {
  padding: 10px;
  border: 2px solid #ddd;
  border-radius: 8px;
  transition: all 0.3s ease-in-out;
}

input:focus {
  border-color: #007bff;
  outline: none;
}

button {
  background: #007bff;
  color: white;
  padding: 10px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.3s;
}

button:disabled {
  background: #ccc;
}

button:hover:not(:disabled) {
  background: #0056b3;
}

.error {
  color: red;
  font-size: 14px;
  margin-top: 10px;
}
</style>
