<template>
  <div id="app-layout">
    <header class="main-header" v-if="authStore.isAuthenticated">
      <div class="logo">
        Plataforma de Tareas
      </div>
      <nav>
        <span>Bienvenido, {{ authStore.currentUser?.username }}</span>
        <button @click="handleLogout">Logout</button>
      </nav>
    </header>

    <main>
      <router-view /> </main>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useAuthStore } from './stores/auth';
import { useRouter } from 'vue-router';

const authStore = useAuthStore();
const router = useRouter();

// Cuando la app se carga, intenta inicializar el estado de autenticación
// para mantener al usuario logueado si ya tenía un token válido.
onMounted(() => {
  authStore.init();
});

const handleLogout = () => {
  authStore.logout();
  // Después de desloguear, redirigimos al usuario a la página de login.
  router.push('/login');
};
</script>

<style>
/* 1. Definimos nuestra paleta de colores como variables CSS */
:root {
  --color-background: #18181b; /* Un negro más suave */
  --color-surface: #27272a;  /* Superficie un poco más clara */
  --color-surface-light: #3f3f46; /* Para hover y elementos secundarios */
  --color-border: #3f3f46;
  --color-text-primary: #f4f4f5;
  --color-text-secondary: #a1a1aa;
  --color-accent: #f59e0b; /* Ámbar como acento cálido */
  --color-accent-hover: #d97706;
  --color-success: #10b981;
  --color-error: #ef4444;

  --font-sans: 'Inter', -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
}

/* Reseteo y estilos base */
*, *::before, *::after {
  box-sizing: border-box;
}

/* 2. Aplicamos estilos globales básicos */
body {
  margin: 0;
  font-family: var(--font-sans);
  background-color: var(--color-background);
  color: var(--color-text-primary);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

h1, h2, h3 {
  font-weight: 700;
  letter-spacing: -0.02em;
}

#app-layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

main {
  flex-grow: 1;
  padding: 20px 40px;
}

/* 3. Estilos generales para el Header (usando las variables) */
.main-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2.5rem;
  background-color: var(--color-surface);
  border-bottom: 1px solid var(--color-border);
}

.main-header .logo {
  font-weight: bold;
  font-size: 1.2em;
}

.main-header nav {
  display: flex;
  align-items: center;
  gap: 15px;
}

.main-header nav span {
  color: var(--color-text-secondary);
  font-size: 0.9em;
}

.main-header button {
  background-color: var(--color-error);
  color: white;
  border: none;
  padding: 8px 12px;
  border-radius: 5px;
  cursor: pointer;
  font-weight: bold;
  transition: background-color 0.2s;
}

.main-header button:hover {
  background-color: #c9302c; /* Un rojo más oscuro al pasar el mouse */
}

/* 4. Estilos generales para botones y formularios */
button.primary-action {
  background-color: var(--color-accent);
  color: #1a1a1d;
  font-weight: bold;
  border: none;
  padding: 10px 15px;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.2s;
}

button.primary-action:hover {
  background-color: var(--color-accent-hover);
}

.auth-form, .workflow-form {
  background-color: var(--color-surface);
  padding: 30px;
  border-radius: 8px;
  border: 1px solid var(--color-border);
}

.auth-form input, .workflow-form input[type="text"], .workflow-form textarea {
    background-color: var(--color-background);
    color: var(--color-text-primary);
    border: 1px solid var(--color-border);
}

.auth-form input:focus, .workflow-form input[type="text"]:focus, .workflow-form textarea:focus {
    border-color: var(--color-accent);
    outline: none;
}
</style>