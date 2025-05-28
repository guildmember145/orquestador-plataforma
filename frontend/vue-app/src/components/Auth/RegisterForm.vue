<template>
  <div class="auth-form">
    <h2>Register</h2>
    <form @submit.prevent="handleRegister">
      <div class="input-group">
        <label for="username">Username:</label>
        <input type="text" id="username" v-model="username" required />
      </div>
      <div class="input-group">
        <label for="email">Email:</label>
        <input type="email" id="email" v-model="email" required />
      </div>
      <div class="input-group">
        <label for="password">Password:</label>
        <input type="password" id="password" v-model="password" required />
      </div>
      <button type="submit" :disabled="authStore.isLoading">
        {{ authStore.isLoading ? 'Registering...' : 'Register' }}
      </button>
      <p v-if="authStore.error" class="error">{{ authStore.error }}</p>
      <p v-if="message" class="success">{{ message }}</p>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useAuthStore } from '../../stores/auth';

const authStore = useAuthStore();
const username = ref('');
const email = ref('');
const password = ref('');
const message = ref('');

const handleRegister = async () => {
  message.value = '';
  try {
    await authStore.register({
      username: username.value,
      email: email.value,
      password: password.value,
    });
    message.value = 'Registration successful! Please login.';
    username.value = '';
    email.value = '';
    password.value = '';
  } catch (error) {
    console.error('Registration failed:', authStore.error);
  }
};
</script>

<style scoped>
.auth-form {
  max-width: 400px;
  margin: auto;
  padding: 30px;
  background: #22a568;
  color: #fff;
  border-radius: 12px;
  box-shadow: 0px 4px 10px rgba(0, 0, 0, 0.2);
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
  border: 2px solid #444;
  border-radius: 8px;
  background: #333;
  color: #fff;
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
  background: #555;
}

button:hover:not(:disabled) {
  background: #0056b3;
}

.error {
  color: red;
  font-size: 14px;
  margin-top: 10px;
}

.success {
  color: limegreen;
  font-size: 14px;
  margin-top: 10px;
}
</style>
