import { createApp } from 'vue'
import { createPinia } from 'pinia' // Importar Pinia
import App from './App.vue'
// Si usas router, impórtalo aquí
// import router from './router'
import './style.css' // O cualquier estilo global

const app = createApp(App)

app.use(createPinia()) // Usar Pinia
// app.use(router) // Si usas router

app.mount('#app')