// src/types/auth.ts
export interface LoginRequest { email: string; password: string; }
export interface RegisterRequest { username: string; email: string; password: string; }
export interface UserProfile { id: string; username: string; email: string; }