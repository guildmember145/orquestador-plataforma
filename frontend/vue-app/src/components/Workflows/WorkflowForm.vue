<template>
  <form @submit.prevent="handleSubmit" class="workflow-form">
    <div class="form-section">
      <h2>Detalles Básicos</h2>
      <div class="form-group">
        <label for="name">Nombre del Workflow</label>
        <input id="name" type="text" v-model="formData.name" placeholder="Ej: Notificarme del Clima cada mañana"
          required />
      </div>
      <div class="form-group">
        <label for="description">Descripción</label>
        <textarea id="description" v-model="formData.description"
          placeholder="Describe brevemente qué hace este workflow"></textarea>
      </div>
      <div class="form-group-inline">
        <label for="is_enabled">Habilitado:</label>
        <input id="is_enabled" type="checkbox" v-model="formData.is_enabled" />
      </div>
    </div>

    <div class="form-section">
      <h2>Disparador (Trigger)</h2>
      <TriggerConfigurator @update:trigger="updateTriggerData" />
    </div>

    <div class="form-section">
      <header class="actions-section-header">
        <h2>Acciones (Actions)</h2>
        <button @click="addAction" type="button" class="add-action-btn">+ Añadir Acción</button>
      </header>
      <div class="actions-list">
        <ActionConfigurator v-for="(action, index) in formData.actions" :key="index" :modelValue="action"
          @update:modelValue="updateAction(index, $event)" @delete="removeAction(index)" />
        <p v-if="formData.actions.length === 0" class="placeholder">
          Añade al menos una acción para tu workflow.
        </p>
      </div>
    </div>

    <div class="form-actions">
      <button type="submit" class="submit-btn primary-action">
        {{ isEditMode ? 'Actualizar Workflow' : 'Guardar Workflow' }}
      </button>
    </div>
  </form>
</template>

// En src/components/Workflows/WorkflowForm.vue

<script setup lang="ts">
import { reactive, computed, watch } from 'vue';
import { useWorkflowStore, type Workflow } from '@/stores';
import { useRouter } from 'vue-router';
import { useToast } from 'vue-toastification';
import TriggerConfigurator from './TriggerConfigurator.vue';
import ActionConfigurator from './ActionConfigurator.vue';


const props = defineProps({
  initialData: {
    type: Object as () => Workflow | null,
    default: null,
  }
});

const workflowStore = useWorkflowStore();
const router = useRouter();
const toast = useToast();

const isEditMode = computed(() => !!props.initialData);




const formData = reactive({
  name: '',
  description: '',
  is_enabled: true,
  trigger: { type: 'schedule', config: { cron: '*/5 * * * *' } },
  actions: [] as any[],
});

// Nuevo método para recibir los datos del TriggerConfigurator
const updateTriggerData = (newTriggerData: any) => {
  formData.trigger = newTriggerData;
};

watch(() => props.initialData, (newData) => {
  if (newData) {
    // Hacemos una copia profunda para no mutar las props directamente
    const dataCopy = JSON.parse(JSON.stringify(newData));
    Object.assign(formData, dataCopy);
  }
}, { immediate: true });

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


const handleSubmit = async () => {
  try {
    if (isEditMode.value) {
      // MODO EDICIÓN: Llamamos a 'updateWorkflow' con el ID y los datos del formulario.
      await workflowStore.updateWorkflow(props.initialData!.id, formData);
      toast.success('¡Workflow actualizado exitosamente!');
    } else {
      // MODO CREACIÓN: Llamamos a 'createWorkflow'.
      await workflowStore.createWorkflow(formData);
      toast.success('¡Workflow creado exitosamente!');
    }
    // Si todo va bien, redirigimos al dashboard.
    router.push('/dashboard/workflows');
  } catch (error) {
    toast.error('Hubo un error al guardar el workflow.');
  }
};

</script>

<style scoped>
.workflow-form {
  max-width: 800px;
  margin: 20px auto;
  padding: 30px;
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 8px;
}

.form-section {
  margin-bottom: 30px;
}

.form-section h2 {
  color: var(--color-text-primary);
  border-bottom: 1px solid var(--color-border);
  padding-bottom: 10px;
  margin-bottom: 20px;
  font-size: 1.4em;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label,
.form-group-inline label {
  display: block;
  margin-bottom: 8px;
  font-weight: bold;
  color: var(--color-text-secondary);
}

.form-group input[type="text"],
.form-group textarea {
  width: 100%;
  padding: 12px;
  background-color: var(--color-background);
  border: 1px solid var(--color-border);
  color: var(--color-text-primary);
  border-radius: 5px;
  box-sizing: border-box;
  font-size: 1em;
  transition: border-color 0.3s;
}

.form-group-inline {
  display: flex;
  align-items: center;
  gap: 10px;
}

.form-group-inline input[type="checkbox"] {
  width: 20px;
  height: 20px;
}

.form-actions {
  text-align: right;
  margin-top: 40px;
}

.submit-btn {
  font-size: 16px;
  padding: 12px 24px;
}

.placeholder {
  color: var(--color-text-secondary);
  font-style: italic;
  padding: 20px;
  text-align: center;
  border: 1px dashed var(--color-border);
  border-radius: 5px;
}

.actions-section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.add-action-btn {
  background-color: var(--color-surface);
  color: var(--color-accent);
  border: 1px solid var(--color-accent);
  padding: 8px 12px;
  border-radius: 5px;
  cursor: pointer;
  font-weight: bold;
}
</style>