import { useNavigate } from '@tanstack/react-router';
import { useQueryClient } from '@tanstack/react-query';
import { Button } from '@/shared/ui/button';
import { authStorage } from '@/shared/lib/auth-storage';

export const LogoutButton = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const handleLogout = () => {
    authStorage.removeToken();
    
    queryClient.resetQueries();
    queryClient.clear();
    
    navigate({ to: '/auth/login' });
  };

  return (
    <Button variant="ghost" onClick={handleLogout}>
      Logout
    </Button>
  );
};

