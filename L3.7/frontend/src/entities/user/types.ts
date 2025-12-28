export interface User {
  id: string;
  email: string;
  role: 'viewer' | 'manager' | 'admin';
}

export interface AuthResponse {
  user: User;
  token: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
}


