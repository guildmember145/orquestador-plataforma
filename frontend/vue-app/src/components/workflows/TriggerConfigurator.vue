<template>
  <div class="trigger-configurator">
    <div class="form-group">
      <label for="trigger-type">Tipo de Disparador</label>
      <select id="trigger-type" v-model="editableTrigger.type" @change="onTypeChange">
        <option value="schedule">Programado (Schedule)</option>
        <option value="webhook">Webhook</option>
      </select>
    </div>

    <div v-if="editableTrigger.type === 'schedule'" class="trigger-details">
      <label for="cron-expression">Expresión Cron</label>
      <input 
        id="cron-expression"
        type="text" 
        v-model="editableTrigger.config.cron" 
        placeholder="Ej: */5 * * * * (cada 5 minutos)"
      />
      <p class="help-text">Define la frecuencia de ejecución usando una expresión cron.</p>
    </div>

    <div v-if="editableTrigger.type === 'webhook'" class="trigger-details">
      <p class="help-text">
        Al guardar, se generará una URL única para este webhook.
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';

const props = defineProps({
  modelValue: {
    type: Object as () => ({ type: string; config: any }),
    required: true,
  }
});

const emit = defineEmits(['update:modelValue']);

// Usamos una copia local para poder modificar los datos internamente
const editableTrigger = ref({ ...props.modelValue });

// Cuando el tipo de trigger cambia, reiniciamos su configuración
const onTypeChange = () => {
  if (editableTrigger.value.type === 'schedule') {
    editableTrigger.value.config = { cron: '' };
  } else if (editableTrigger.value.type === 'webhook') {
    editableTrigger.value.config = {};
  }
};

// Observamos cambios en nuestra copia local y los emitimos al componente padre
watch(editableTrigger, (newValue) => {
  emit('update:modelValue', newValue);
}, { deep: true });

</script>

<style scoped>
.trigger-configurator {
  padding: 20px;
  background-color: var(--color-background);
  border: 1px solid var(--color-border);
  border-radius: 5px;
}
.form-group {
  margin-bottom: 15px;
}
label {
  display: block;
  margin-bottom: 8px;
  font-weight: bold;
  font-size: 0.9em;
  color: var(--color-text-secondary);
}
select, input {
  width: 100%;
  padding: 12px;
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  color: var(--color-text-primary);
  border-radius: 5px;
  box-sizing: border-box;
  font-size: 1em;
}
select:focus, input:focus {
  border-color: var(--color-accent);
  outline: none;
}
.trigger-details {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px dashed var(--color-border);
}
.help-text {
  font-size: 0.9em;
  color: var(--color-text-secondary);
  margin-top: 8px;
}
</style>