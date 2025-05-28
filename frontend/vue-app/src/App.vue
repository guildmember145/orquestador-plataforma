<template>
  <div class="container">
    <header>
      <h1>Orquestador de Tareas</h1>
      <nav v-if="authStore.isAuthenticated">
        <span>Bienvenido, {{ authStore.currentUser?.username }}!</span>
        <button @click="handleLogout">Logout</button>
      </nav>
    </header>
    <main>
      <div v-if="!authStore.isAuthenticated" class="auth-section">
        <LoginForm />
        <hr>
        <RegisterForm />
      </div>
      <div v-else class="dashboard">
        <p>¡Estás logueado!</p>
        <p>Aquí iría el contenido de tu Dashboard o la vista principal.</p>
        <button @click="fetchMyData" v-if="!authStore.currentUser?.id">Ver mis datos</button>
        <pre v-if="authStore.currentUser">{{ authStore.currentUser }}</pre>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import LoginForm from './components/Auth/LoginForm.vue';
import RegisterForm from './components/Auth/RegisterForm.vue';
import { useAuthStore } from './stores/auth';

const authStore = useAuthStore();

onMounted(async () => {
  await authStore.init();
});

const handleLogout = () => {
  authStore.logout();
};

const fetchMyData = async () => {
  try {
    await authStore.fetchUser();
  } catch (error) {
    console.error("Failed to fetch user data on button click");
  }
};
</script>

<style scoped>
.container {
  background: #1e1e1e;
  color: #fff;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
}

header {
  background: #2a2a2a;
  padding: 15px;
  width: 100%;
  border-radius: 8px;
  text-align: center;
  box-shadow: 0px 4px 10px rgba(0, 0, 0, 0.2);
}

nav {
  margin-top: 10px;
}

nav span {
  font-weight: bold;
  margin-right: 10px;
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

button:hover {
  background: #0056b3;
}

.auth-section {
  background: #292929;
  padding: 15px;
  border-radius: 10px;
  box-shadow: 0px 4px 10px rgba(0, 0, 0, 0.3);
  text-align: center;
  width: 400px;
}

.dashboard {
  text-align: center;
  margin-top: 20px;
}

hr {
  border: none;
  height: 2px;
  background: #444;
  margin: 20px 0;
}

pre {
  background: #333;
  padding: 15px;
  border-radius: 8px;
  overflow-x: auto;
}
</style>
