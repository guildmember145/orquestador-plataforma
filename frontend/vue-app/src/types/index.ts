// Tipos de autenticación
export interface User {
  id: number
  username: string
  email: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
}

// Tipos de workflow
export interface Workflow {
  id: string; // El UUID siempre será un string en el frontend
  user_id: string;
  name: string;
  description: string;
  trigger: any;
  actions: any[];
  is_enabled: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateWorkflowRequest {
  name: string
  description: string
  is_enabled: boolean
}


export interface ExecutionLog {
  id: string;
  workflow_id: string;
  user_id: string;
  status: 'running' | 'completed' | 'failed';
  triggered_at: string;
  completed_at?: string | null;
  logs: any;
}
