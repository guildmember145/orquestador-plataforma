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
  id: number
  name: string
  description: string
  is_enabled: boolean
  created_at?: string
  updated_at?: string
}

export interface CreateWorkflowRequest {
  name: string
  description: string
  is_enabled: boolean
}