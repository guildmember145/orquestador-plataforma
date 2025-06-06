<template>
  <div class="dashboard">
    <header class="dashboard-header">
      <h1>Mis Workflows</h1>
      <button @click="goToCreatePage" class="primary-action">Crear Nuevo Workflow</button>
    </header>

    <div v-if="workflowStore.isWorkflowsLoading && workflowStore.allWorkflows.length === 0" class="loading">
      Cargando workflows...
    </div>
    <div v-else-if="workflowStore.getWorkflowError" class="error-message">
      <p>Error al cargar los workflows: {{ workflowStore.getWorkflowError }}</p>
    </div>
    <div v-else-if="workflowStore.allWorkflows.length === 0" class="no-workflows">
      <p>No has creado ningún workflow todavía.</p>
    </div>

    <div v-else class="workflows-list">
      <table>
        <thead>
          <tr>
            <th>Nombre</th>
            <th>Descripción</th>
            <th>Estado</th>
            <th>Acciones</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="workflow in workflowStore.allWorkflows" :key="workflow.id">
            <td>{{ workflow.name }}</td>
            <td>{{ workflow.description }}</td>
            <td>
              <span :class="['status-badge', workflow.is_enabled ? 'enabled' : 'disabled']">
                {{ workflow.is_enabled ? 'Habilitado' : 'Deshabilitado' }}
              </span>
            </td>
            <td class="actions-cell">
              <button @click="handleEdit(workflow.id)" class="action-btn edit">Editar</button>
              <button @click="handleDelete(workflow.id, workflow.name)" class="action-btn delete">Eliminar</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useWorkflowStore } from '@/stores';
import { useRouter } from 'vue-router';

const workflowStore = useWorkflowStore();
const router = useRouter();

onMounted(() => {
  workflowStore.fetchWorkflows();
});

const goToCreatePage = () => {
  router.push('/dashboard/workflows/new');
};

// Función para navegar a la página de edición
const handleEdit = (workflowId: string) => {
  router.push(`/dashboard/workflows/edit/${workflowId}`);
};

// Función para manejar la eliminación
const handleDelete = async (workflowId: string, workflowName: string) => {
  const confirmed = confirm(`¿Estás seguro de que quieres eliminar el workflow "${workflowName}"?`);
  if (confirmed) {
    try {
      await workflowStore.deleteWorkflow(workflowId);
      alert(`Workflow "${workflowName}" eliminado.`);
    } catch (error) {
      alert(`Error al eliminar el workflow: ${workflowStore.getWorkflowError}`);
    }
  }
};
</script>

<style scoped>
  /* Tus estilos existentes aquí... */
  .dashboard { padding: 20px; }
  .dashboard-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
  .loading, .error-message, .no-workflows { text-align: center; margin-top: 50px; padding: 20px; background-color: var(--color-surface); border-radius: 8px; }
  .error-message { color: var(--color-error); }
  .workflows-list table { width: 100%; border-collapse: collapse; background-color: var(--color-surface); border-radius: 8px; overflow: hidden; }
  .workflows-list th, .workflows-list td { padding: 15px; text-align: left; border-bottom: 1px solid var(--color-border); }
  .workflows-list th { background-color: rgba(255, 255, 255, 0.05); font-size: 0.9em; text-transform: uppercase; color: var(--color-text-secondary); }
  .status-badge { padding: 5px 10px; border-radius: 12px; color: white; font-size: 0.85em; font-weight: bold; }
  .status-badge.enabled { background-color: var(--color-success); }
  .status-badge.disabled { background-color: #6c757d; }
  .actions-cell { display: flex; gap: 10px; }
  .action-btn { padding: 6px 12px; border: 1px solid var(--color-border); background-color: var(--color-background); color: var(--color-text-primary); cursor: pointer; border-radius: 4px; transition: background-color 0.2s, border-color 0.2s; }
  .action-btn.edit:hover { border-color: var(--color-accent); color: var(--color-accent); }
  .action-btn.delete:hover { border-color: var(--color-error); color: var(--color-error); }
</style>