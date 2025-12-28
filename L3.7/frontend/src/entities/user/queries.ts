import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { userApi } from './api';
import { authStorage } from '@/shared/lib/auth-storage';

export const useRegister = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: userApi.register,
    onSuccess: (data) => {
      authStorage.setToken(data.token);
      queryClient.setQueryData(['user'], data.user);
    },
  });
};

export const useLogin = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: userApi.login,
    onSuccess: (data) => {
      authStorage.setToken(data.token);
      queryClient.setQueryData(['user'], data.user);
    },
  });
};

export const useUser = () => {
  return useQuery({
    queryKey: ['user'],
    queryFn: userApi.getMe,
    enabled: authStorage.isAuthenticated(),
    retry: false,
  });
};

