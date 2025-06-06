import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

interface Workflow {
  id: number
  name: string
  description: string
  is_enabled: boolean
  created_at?: string
  updated_at?: string
}

export const useWorkflowStore = defineStore('workflow', () => {
  // State
  const workflows = ref<Workflow[]>([])
  const isWorkflowsLoading = ref(false)
  const getWorkflowError = ref<string | null>(null)

  // Getters
  const allWorkflows = computed(() => workflows.value)

  // Actions
  const fetchWorkflows = async () => {
    isWorkflowsLoading.value = true
    getWorkflowError.value = null

    try {
      // Aquí harías la llamada a tu API
      // const response = await api.getWorkflows()
      
      // Datos de ejemplo
      const mockWorkflows: Workflow[] = [
        {
          id: 1,
          name: 'Workflow de Ejemplo 1',
          description: 'Este es un workflow de prueba',
          is_enabled: true
        },
        {
          id: 2,
          name: 'Workflow de Ejemplo 2',
          description: 'Otro workflow de ejemplo',
          is_enabled: false
        }
      ]
      
      workflows.value = mockWorkflows
      
    } catch (err) {
      getWorkflowError.value = 'Error al cargar workflows'
      console.error('Error fetching workflows:', err)
    } finally {
      isWorkflowsLoading.value = false
    }
  }

  const createWorkflow = async (workflowData: Omit<Workflow, 'id'>) => {
    try {
      // Llamada a la API para crear workflow
      // const response = await api.createWorkflow(workflowData)
      
      // Simulación
      const newWorkflow: Workflow = {
        ...workflowData,
        id: Date.now() // ID temporal
      }
      
      workflows.value.push(newWorkflow)
      return newWorkflow
      
    } catch (err) {
      getWorkflowError.value = 'Error al crear workflow'
      throw err
    }
  }

  const deleteWorkflow = async (id: number) => {
    try {
      // await api.deleteWorkflow(id)
      workflows.value = workflows.value.filter(w => w.id !== id)
    } catch (err) {
      getWorkflowError.value = 'Error al eliminar workflow'
      throw err
    }
  }

  return {
    workflows,
    isWorkflowsLoading,
    getWorkflowError,
    allWorkflows,
    fetchWorkflows,
    createWorkflow,
    deleteWorkflow
  }
})