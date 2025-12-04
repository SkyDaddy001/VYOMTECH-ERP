export interface User {
  id: number
  email: string
  name: string
  role: string
  tenant_id: string
  created_at: string
  updated_at: string
}

export interface AuthResponse {
  token: string
  user: User
  message: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
  role: string
  tenant_id: string
  name?: string
}
