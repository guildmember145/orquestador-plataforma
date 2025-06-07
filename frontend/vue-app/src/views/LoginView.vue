<template>
  <div class="login-view">
    <h1>Bienvenido de Nuevo</h1>
    <p>Inicia sesión para acceder a tu dashboard.</p>
    <div class="form-wrapper">
      <LoginForm />
    </div>
    <div class="register-link">
      <p>¿No tienes una cuenta? <router-link to="/register"> Regístrate aquí </router-link>.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import LoginForm from '@/components/Auth/LoginForm.vue';
</script>

<style scoped>
.login-view {
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
.register-link {
  margin-top: 20px;
  font-size: 0.9em;
}

@media screen and (max-width: 768px) {
  .forms-container {
    flex-direction: column;
    gap: 20px;
  }
}







</style>