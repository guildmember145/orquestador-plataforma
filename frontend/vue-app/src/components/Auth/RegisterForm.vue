<template>
  <div class="auth-form">
    <h2>Register</h2>
    <form @submit.prevent="handleRegister">
      <div class="input-group">
        <label for="register-username">Username:</label>
        <input id="register-username" type="text" v-model="username" required />
      </div>
      <div class="input-group">
        <label for="register-email">Email:</label>
        <input id="register-email" type="email" v-model="email" required />
      </div>
      <div class="input-group">
        <label for="register-password">Password:</label>
        <input id="register-password" type="password" v-model="password" required />
      </div>
      <button type="submit" :disabled="authStore.isLoading" class="primary-action">
        {{ authStore.isLoading ? 'Registering...' : 'Register' }}
      </button>

      <p v-if="error" class="error">{{ error }}</p>
      <p v-if="successMessage" class="success">{{ successMessage }}</p>

    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useAuthStore } from '@/stores'; // Usamos el barril de stores

const authStore = useAuthStore();
const username = ref('');
const email = ref('');
const password = ref('');

// Estados locales para los mensajes
const error = ref<string | null>(null);
const successMessage = ref<string | null>(null);

const handleRegister = async () => {
  // Limpiamos mensajes anteriores
  error.value = null;
  successMessage.value = null;

  try {
    // --- ESTA ES LA LÍNEA CORREGIDA ---
    // Llamamos a la acción 'register' del store
    await authStore.register({
      username: username.value,
      email: email.value,
      password: password.value,
    });

    // Si el registro es exitoso
    successMessage.value = '¡Registro exitoso! Ahora puedes iniciar sesión.';

    // Limpiar el formulario
    username.value = '';
    email.value = '';
    password.value = '';

  } catch (err: any) {
    // Si la acción 'register' falla, atrapamos el error
    // y lo mostramos en la interfaz.
    error.value = err.message || 'Ocurrió un error durante el registro.';
    console.error('Registration failed:', err);
  }
};
</script>

<style scoped>
/* Usaremos los mismos estilos que ya definiste para LoginForm.vue */
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
}
.error {
  color: var(--color-error);
  font-size: 0.9em;
  margin-top: 15px;
  text-align: center;
}
.success {
  color: var(--color-success);
  font-size: 0.9em;
  margin-top: 15px;
  text-align: center;
}
</style>