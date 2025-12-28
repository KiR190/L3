import apiClient from '@/shared/api/client';
import type { AuthResponse, LoginRequest, RegisterRequest, User } from './types';

export const userApi = {
  async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await apiClient.post('/auth/register', data);
    return response.data;
  },

  async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await apiClient.post('/auth/login', data);
    return response.data;
  },

  async getMe(): Promise<User> {
    const response = await apiClient.get('/auth/me');
    return response.data;
  },
};

