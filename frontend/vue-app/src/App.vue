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
  --color-background: #1a1a1d; /* Un negro/gris muy oscuro, ligeramente cálido */
  --color-surface: #2c2c34;  /* Una superficie un poco más clara para tarjetas y formularios */
  --color-border: #4a4a52;   /* Un borde sutil */
  --color-text-primary: #f2f2f2; /* Un blanco no tan puro, más suave a la vista */
  --color-text-secondary: #a9a9b2; /* Un gris más suave para texto secundario */
  
  --color-accent: #fca311;       /* Un acento cálido, tipo ámbar/naranja */
  --color-accent-hover: #e85d04;  /* Un naranja más intenso para hover */

  --color-success: #2a9d8f;      /* Un verde azulado cálido */
  --color-error: #e76f51;        /* Un rojo/coral cálido */
}

/* 2. Aplicamos estilos globales básicos */
body {
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji";
  background-color: var(--color-background);
  color: var(--color-text-primary);
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
  padding: 15px 40px;
  background-color: var(--color-surface);
  border-bottom: 1px solid var(--color-border);
  color: var(--color-text-primary);
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