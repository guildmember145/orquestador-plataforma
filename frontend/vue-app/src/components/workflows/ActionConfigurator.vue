<template>
  <div class="action-configurator">
    <div class="action-header">
      <span class="action-title">Acción: {{ editableAction.name || 'Nueva Acción' }}</span>
      <button @click="$emit('delete')" type="button" class="delete-action-btn">Eliminar</button>
    </div>

    <div class="action-body">
      <div class="form-group half">
        <label>Nombre de la Acción</label>
        <input type="text" placeholder="Ej: LlamarAPIExterna" v-model="editableAction.name" required />
      </div>
      <div class="form-group half">
        <label>Tipo de Acción</label>
        <select v-model="editableAction.type" @change="onTypeChange">
          <option value="log_message">Registrar Mensaje (Log)</option>
          <option value="http_endpoint">Llamar a Endpoint HTTP</option>
        </select>
      </div>

      <div v-if="editableAction.type === 'log_message'" class="config-section">
        <div class="form-group">
          <label>Mensaje a Registrar</label>
          <input type="text" placeholder="El workflow se ejecutó con éxito" v-model="editableAction.config.message" />
        </div>
      </div>

      <div v-if="editableAction.type === 'http_endpoint'" class="config-section">
        <div class="form-group">
          <label>URL</label>
          <input type="url" placeholder="https://api.example.com/data" v-model="editableAction.config.url" />
        </div>
        <div class="form-group half">
          <label>Método HTTP</label>
          <select v-model="editableAction.config.method">
            <option>GET</option>
            <option>POST</option>
            <option>PUT</option>
            <option>DELETE</option>
            <option>PATCH</option>
          </select>
        </div>
        <div class="form-group">
          <label>Cabeceras (Headers) - Formato JSON</label>
          <textarea placeholder='{ "Content-Type": "application/json" }' v-model="editableAction.config.headers_str"></textarea>
        </div>
        <div class="form-group">
          <label>Cuerpo (Body) - JSON o Texto Plano</label>
          <textarea placeholder='{ "key": "value" }' v-model="editableAction.config.body_str"></textarea>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';

// Definimos una interfaz para el objeto 'action' para mayor claridad
interface Action {
  name: string;
  type: string;
  config: any; // 'any' por flexibilidad
}

const props = defineProps({
  modelValue: {
    type: Object as () => Action,
    required: true,
  }
});

const emit = defineEmits(['update:modelValue', 'delete']);

// Creamos una copia local profunda para poder modificar los datos sin afectar al padre directamente
const editableAction = ref(JSON.parse(JSON.stringify(props.modelValue)));

// Función que se ejecuta cuando el usuario cambia el tipo de acción
const onTypeChange = () => {
    // Reiniciamos la configuración para evitar mantener datos de un tipo de acción anterior
    if (editableAction.value.type === 'log_message') {
        editableAction.value.config = { message: '' };
    } else if (editableAction.value.type === 'http_endpoint') {
        editableAction.value.config = { url: '', method: 'GET', headers_str: '{}', body_str: '' };
    }
};

// onMounted se ejecuta cuando el componente se crea por primera vez
onMounted(() => {
    // Nos aseguramos que el objeto 'config' exista
    if (!editableAction.value.config) {
      editableAction.value.config = {};
    }
    // Convertimos los objetos de headers/body a string para los <textarea>
    if (editableAction.value.config.headers) {
        editableAction.value.config.headers_str = JSON.stringify(editableAction.value.config.headers, null, 2);
    }
    if (editableAction.value.config.body) {
        if (typeof editableAction.value.config.body === 'object') {
            editableAction.value.config.body_str = JSON.stringify(editableAction.value.config.body, null, 2);
        } else {
            editableAction.value.config.body_str = editableAction.value.config.body;
        }
    }
});

// 'watch' se ejecuta cada vez que 'editableAction' cambia
watch(editableAction, (newValue) => {
    // Creamos una copia para no modificar el valor mientras procesamos
    const finalValue = JSON.parse(JSON.stringify(newValue));
    
    // Intentamos convertir los strings de headers/body de vuelta a objetos JSON
    try {
        if (finalValue.config.headers_str) {
            finalValue.config.headers = JSON.parse(finalValue.config.headers_str);
        }
    } catch (e) { /* Ignoramos JSON inválido mientras el usuario está escribiendo */ }
    
    try {
        if (finalValue.config.body_str) {
            finalValue.config.body = JSON.parse(finalValue.config.body_str);
        }
    } catch (e) {
        // Si no es un JSON válido, lo tratamos como un string plano
        finalValue.config.body = finalValue.config.body_str;
    }
    
    // No queremos guardar las versiones en string en nuestro estado principal
    delete finalValue.config.headers_str;
    delete finalValue.config.body_str;
    
    // Emitimos el objeto limpio al componente padre (WorkflowForm.vue)
    emit('update:modelValue', finalValue);
}, { deep: true });

</script>

<style scoped>
.action-configurator { background-color: var(--color-background); border: 1px solid var(--color-border); border-radius: 5px; margin-bottom: 15px; overflow: hidden; }
.action-header { display: flex; justify-content: space-between; align-items: center; background-color: rgba(255, 255, 255, 0.05); padding: 10px 15px; }
.action-title { font-weight: bold; }
.delete-action-btn { background-color: var(--color-error); color: white; border: none; border-radius: 4px; padding: 5px 10px; cursor: pointer; }
.action-body { padding: 15px; display: flex; flex-wrap: wrap; gap: 15px; }
.form-group { flex: 1 1 100%; }
.form-group.half { flex: 1 1 calc(50% - 8px); }
label { display: block; margin-bottom: 8px; font-weight: bold; font-size: 0.9em; color: var(--color-text-secondary); }
input, select, textarea { width: 100%; padding: 10px; background-color: var(--color-surface); border: 1px solid var(--color-border); color: var(--color-text-primary); border-radius: 5px; box-sizing: border-box; font-family: inherit; }
textarea { min-height: 100px; resize: vertical; font-family: 'Courier New', Courier, monospace; }
.config-section { width: 100%; display: flex; flex-wrap: wrap; gap: 15px; }
</style>