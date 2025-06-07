<template>
  <div class="trigger-configurator">
    <div class="form-group">
      <label for="frequency-type">Ejecutar</label>
      <select id="frequency-type" v-model="schedule.type">
        <option value="minutely">Cada N Minutos</option>
        <option value="hourly">Cada N Horas</option>
        <option value="daily">Diariamente a una Hora</option>
        <option value="weekly">Semanalmente</option>
        <option value="custom">Personalizado (Cron)</option>
      </select>
    </div>

    <div v-if="schedule.type === 'minutely'" class="trigger-details">
      <label>Cada cuántos minutos</label>
      <input type="number" v-model.number="schedule.minuteInterval" min="1" max="59" />
    </div>
    
    <div v-if="schedule.type === 'hourly'" class="trigger-details">
      <label>Cada cuántas horas</label>
      <input type="number" v-model.number="schedule.hourInterval" min="1" max="23" />
    </div>

    <div v-if="schedule.type === 'daily'" class="trigger-details">
      <label>A la hora</label>
      <input type="time" v-model="schedule.time" />
    </div>
    
    <div v-if="schedule.type === 'weekly'" class="trigger-details">
       <label>El día</label>
       <select v-model.number="schedule.dayOfWeek">
         <option value="1">Lunes</option>
         <option value="2">Martes</option>
         <option value="3">Miércoles</option>
         <option value="4">Jueves</option>
         <option value="5">Viernes</option>
         <option value="6">Sábado</option>
         <option value="0">Domingo</option>
       </select>
       <label>A la hora</label>
       <input type="time" v-model="schedule.time" />
    </div>

    <div class="trigger-details">
      <label for="cron-expression">Expresión Cron Resultante</label>
      <input 
        id="cron-expression"
        type="text" 
        v-model="generatedCron" 
        :disabled="schedule.type !== 'custom'"
      />
      <p class="help-text">Esta es la expresión cron que se guardará.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';

// No necesitamos props porque este componente ahora es el que genera el valor.
// En su lugar, emitirá los cambios al formulario padre.
const emit = defineEmits(['update:trigger']);

// Estado local para manejar las selecciones del usuario
const schedule = ref({
  type: 'minutely',       // Tipo de frecuencia por defecto
  minuteInterval: 5,    // Cada 5 minutos
  hourInterval: 1,      // Cada 1 hora
  time: '09:00',          // A las 09:00 AM
  dayOfWeek: 1,         // Lunes
});

// Estado para la expresión cron generada
const generatedCron = ref('');

// El "traductor": esta función se ejecuta cada vez que cambia una selección
// del usuario y genera la cadena de texto cron correspondiente.
watch(schedule, (newSchedule) => {
  let cron = '* * * * *'; // Default: cada minuto
  const [hour, minute] = newSchedule.time.split(':');

  switch (newSchedule.type) {
    case 'minutely':
      cron = `*/${newSchedule.minuteInterval} * * * *`;
      break;
    case 'hourly':
      cron = `0 */${newSchedule.hourInterval} * * *`;
      break;
    case 'daily':
      cron = `${minute} ${hour} * * *`;
      break;
    case 'weekly':
      cron = `${minute} ${hour} * * ${newSchedule.dayOfWeek}`;
      break;
    case 'custom':
      // En modo custom, no hacemos nada, el usuario escribe directamente.
      // Pero mantenemos el valor actual en el input.
      return; 
  }
  generatedCron.value = cron;
}, { deep: true, immediate: true }); // 'immediate: true' lo ejecuta al inicio

// Este watcher se asegura de que cualquier cambio en la expresión cron
// (ya sea generada o escrita manualmente en modo custom) se emita al formulario padre.
watch(generatedCron, (newCron) => {
    emit('update:trigger', {
        type: 'schedule', // El tipo de trigger para el backend siempre es 'schedule'
        config: {
            cron: newCron
        }
    });
});
</script>

<style scoped>
/* Tus estilos existentes para TriggerConfigurator funcionan bien, 
   puedes mantenerlos o ajustarlos si lo deseas. */
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
  display: flex;
  gap: 15px;
  align-items: center;
  flex-wrap: wrap;
}
.trigger-details label {
  margin-bottom: 0;
}
.help-text {
  width: 100%;
  font-size: 0.9em;
  color: var(--color-text-secondary);
  margin-top: 8px;
}

@media screen and (max-width: 768px) {
  .form-group.half {
    flex-basis: 100%; /* El campo ahora ocupa todo el ancho */
  }
}



</style>