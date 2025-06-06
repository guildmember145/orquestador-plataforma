<template>
  <form @submit.prevent="handleSubmit" class="workflow-form">
    <div class="form-section">
      <h2>Detalles Básicos</h2>
      <div class="form-group">
        <label for="name">Nombre del Workflow</label>
        <input id="name" type="text" v-model="formData.name" placeholder="Ej: Notificarme del Clima cada mañana" required />
      </div>
      <div class="form-group">
        <label for="description">Descripción</label>
        <textarea id="description" v-model="formData.description" placeholder="Describe brevemente qué hace este workflow"></textarea>
      </div>
      <div class="form-group-inline">
        <label for="is_enabled">Habilitado:</label>
        <input id="is_enabled" type="checkbox" v-model="formData.is_enabled" />
      </div>
    </div>

    <div class="form-section">
      <h2>Disparador (Trigger)</h2>
      <TriggerConfigurator v-model="formData.trigger" />
    </div>

    <div class="form-section">
      <header class="actions-section-header">
        <h2>Acciones (Actions)</h2>
        <button @click="addAction" type="button" class="add-action-btn">+ Añadir Acción</button>
      </header>
      <div class="actions-list">
        <ActionConfigurator 
          v-for="(action, index) in formData.actions"
          :key="index"
          :modelValue="action"
          @update:modelValue="updateAction(index, $event)"
          @delete="removeAction(index)"
        />
        <p v-if="formData.actions.length === 0" class="placeholder">
          Añade al menos una acción para tu workflow.
        </p>
      </div>
    </div>
    
    <div class="form-actions">
      <button type="submit" class="submit-btn primary-action">Guardar Workflow</button>
    </div>
  </form>
</template>

// En src/components/Workflows/WorkflowForm.vue

<script setup lang="ts">
import { reactive } from 'vue';
import { useWorkflowStore } from '@/stores'; // Asegúrate de importar desde el barril
import { useRouter } from 'vue-router';
import TriggerConfigurator from './TriggerConfigurator.vue';
import ActionConfigurator from './ActionConfigurator.vue';

const workflowStore = useWorkflowStore();
const router = useRouter();

const formData = reactive({
  name: '',
  description: '',
  is_enabled: true,
  trigger: {
    type: 'schedule',
    config: { cron: '' },
  },
  actions: [] as any[],
});

const addAction = () => {
  formData.actions.push({
    name: `Paso ${formData.actions.length + 1}`,
    type: 'log_message',
    config: { message: '' }
  });
};

const removeAction = (index: number) => {
  formData.actions.splice(index, 1);
};

const updateAction = (index: number, updatedAction: any) => {
  formData.actions[index] = updatedAction;
};

// --- INICIO DE LA MODIFICACIÓN ---
const handleSubmit = async () => {
  // Mostramos el objeto final que se enviará, es útil para depurar
  console.log('Enviando workflow para crear:', JSON.parse(JSON.stringify(formData)));
  
  try {
    // Llamamos a la acción 'createWorkflow' de nuestro store Pinia
    await workflowStore.createWorkflow(formData);
    
    // Si la creación es exitosa, redirigimos al usuario al dashboard
    alert('¡Workflow creado exitosamente!'); // Opcional: una mejor notificación podría ir aquí
    router.push('/dashboard/workflows');

  } catch (error) {
    // Si la acción 'createWorkflow' lanza un error, lo atrapamos aquí
    console.error("Error desde el componente de formulario:", error);
    // Mostramos una alerta al usuario. El error ya se logueó en la consola desde el store.
    alert('Hubo un error al crear el workflow. Revisa la consola para más detalles.');
  }
};
// --- FIN DE LA MODIFICACIÓN ---
</script>

<style scoped>
.workflow-form { max-width: 800px; margin: 20px auto; padding: 30px; background-color: var(--color-surface); border: 1px solid var(--color-border); border-radius: 8px; }
.form-section { margin-bottom: 30px; }
.form-section h2 { color: var(--color-text-primary); border-bottom: 1px solid var(--color-border); padding-bottom: 10px; margin-bottom: 20px; font-size: 1.4em; }
.form-group { margin-bottom: 20px; }
.form-group label, .form-group-inline label { display: block; margin-bottom: 8px; font-weight: bold; color: var(--color-text-secondary); }
.form-group input[type="text"], .form-group textarea { width: 100%; padding: 12px; background-color: var(--color-background); border: 1px solid var(--color-border); color: var(--color-text-primary); border-radius: 5px; box-sizing: border-box; font-size: 1em; transition: border-color 0.3s; }
.form-group-inline { display: flex; align-items: center; gap: 10px; }
.form-group-inline input[type="checkbox"] { width: 20px; height: 20px; }
.form-actions { text-align: right; margin-top: 40px; }
.submit-btn { font-size: 16px; padding: 12px 24px; }
.placeholder { color: var(--color-text-secondary); font-style: italic; padding: 20px; text-align: center; border: 1px dashed var(--color-border); border-radius: 5px; }
.actions-section-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.add-action-btn { background-color: var(--color-surface); color: var(--color-accent); border: 1px solid var(--color-accent); padding: 8px 12px; border-radius: 5px; cursor: pointer; font-weight: bold; }
</style>