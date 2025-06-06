<template>
  <div class="workflow-create-view">
    <header>
      <h1>Crear Nuevo Workflow</h1>
      <router-link to="/dashboard/workflows" class="back-link">Volver al Dashboard</router-link>
    </header>
    <WorkflowForm />
  </div>
</template>

<script setup lang="ts">
import WorkflowForm from '@/components/Workflows/WorkflowForm.vue';
</script>

<style scoped>
.workflow-create-view header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  border-bottom: 1px solid #d48e0b;
  padding-bottom: 10px;
}
.back-link {
  text-decoration: none;
  color: #007bff;
}
.back-link:hover {
  text-decoration: underline;
}
</style>