import { createApp } from 'vue'
import { createPinia } from 'pinia'
import Toast, { type PluginOptions, POSITION } from 'vue-toastification' // <-- IMPORTAR
import 'vue-toastification/dist/index.css' // <-- IMPORTAR ESTILOS

import App from './App.vue'
import router from './router'
import './style.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)

// Configuración opcional para las notificaciones
const options: PluginOptions = {
    timeout: 3000, // 3 segundos
    closeOnClick: true,
    pauseOnFocusLoss: true,
    pauseOnHover: true,
    draggable: true,
    position: POSITION.TOP_RIGHT,
};

app.use(Toast, options) // <-- USAR EL PLUGIN CON LAS OPCIONES

app.mount('#app')