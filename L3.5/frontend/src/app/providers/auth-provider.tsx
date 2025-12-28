import { createContext, type ReactNode } from 'react';
import { useUser } from '@/entities/user/queries';
import { authStorage } from '@/shared/lib/auth-storage';
import { useQueryClient } from '@tanstack/react-query';
import type { User } from '@/entities/user/types';

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | null>(null);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const queryClient = useQueryClient();
  const { data: user, isLoading } = useUser();

  const logout = () => {
    authStorage.removeToken();
    queryClient.clear();
    window.location.href = '/auth/login';
  };

  return (
    <AuthContext.Provider
      value={{
        user: user ?? null,
        isLoading,
        isAuthenticated: authStorage.isAuthenticated(),
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

// export const useAuth = () => {
//   const context = useContext(AuthContext);
//   if (!context) throw new Error('useAuth must be used within AuthProvider');
//   return context;
// };
