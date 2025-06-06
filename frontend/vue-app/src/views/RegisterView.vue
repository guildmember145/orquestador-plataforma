<template>
  <div class="register-view">
    <h1>Crear una Cuenta Nueva</h1>
    <p>Completa tus datos para empezar a crear workflows.</p>
    <div class="form-wrapper">
      <RegisterForm />
    </div>
    <div class="login-link">
      <p>¿Ya tienes una cuenta? <router-link to="/login">Inicia sesión aquí</router-link>.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { RegisterForm } from '@/components/Auth';
</script>

<style scoped>
.register-view {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-top: 50px;
  max-width: 500px;
  margin: 0 auto;
}
.form-wrapper {
  width: 100%;
  margin-top: 30px;
}
.login-link {
  margin-top: 20px;
  font-size: 0.9em;
}
</style>