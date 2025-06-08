<template>
  <div class="workflow-edit-view">
    <header>
      <h1>Editar Workflow</h1>
      <router-link to="/dashboard/workflows" class="back-link">Volver al Dashboard</router-link>
    </header>

    <div v-if="isLoading" class="loading">Cargando datos del workflow...</div>
    <div v-else-if="!workflowToEdit" class="error-message">
      Workflow no encontrado o no tienes permiso para editarlo.
    </div>
    <!-- :key es importante para forzar la reinicialización del formulario si navegamos entre diferentes páginas de edición -->
    <WorkflowForm v-else :key="workflowId" :initial-data="workflowToEdit" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import { useWorkflowStore } from '@/stores';
import type { Workflow } from '@/types'; // Adjust the path if Workflow is defined elsewhere
import { WorkflowForm } from '@/components/Workflows';

const route = useRoute();
const workflowStore = useWorkflowStore();

const workflowId = ref(route.params.id as string);
const workflowToEdit = ref<Workflow | null>(null);
const isLoading = ref(true);

// Función para cargar los datos del workflow
const loadWorkflow = async (id: string) => {
  isLoading.value = true;
  // Asegurarnos de tener la lista de workflows en el store
  if (workflowStore.allWorkflows.length === 0) {
    await workflowStore.fetchWorkflows();
  }
  // Encontrar el workflow específico
  const foundWorkflow = workflowStore.allWorkflows.find(wf => wf.id === id);
  workflowToEdit.value = foundWorkflow ? JSON.parse(JSON.stringify(foundWorkflow)) : null;
  isLoading.value = false;
};

// Cargar los datos cuando el componente se monta
onMounted(() => {
  loadWorkflow(workflowId.value);
});

// Si el usuario navega de una página de edición a otra, recargamos los datos
watch(() => route.params.id, (newId) => {
    if (newId && newId !== workflowId.value) {
        workflowId.value = newId as string;
        loadWorkflow(workflowId.value);
    }
});
</script>

<style scoped>
/* Tus estilos existentes aquí... */
.workflow-edit-view header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; border-bottom: 1px solid var(--color-border); padding-bottom: 10px; }
.back-link { text-decoration: none; color: var(--color-accent); }
.loading, .error-message { text-align: center; margin-top: 50px; padding: 20px; background-color: var(--color-surface); border-radius: 8px; }
.error-message { color: var(--color-error); }
</style>