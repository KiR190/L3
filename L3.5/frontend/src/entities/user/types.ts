export interface User {
  id: string;
  email: string;
  role: 'user' | 'admin';
  telegram_username?: string;
  telegram_registered: boolean;
  preferred_notification: string;
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
  telegram_username?: string;
}

export interface LinkTelegramRequest {
  telegram_username: string;
}

export interface UpdateProfileRequest {
  email?: string;
  preferred_notification?: 'email' | 'telegram';
  telegram_username?: string;
}

