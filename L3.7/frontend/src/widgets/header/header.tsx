import { Link } from '@tanstack/react-router';
import { useUser } from '@/entities/user/queries';
import { LogoutButton } from '@/features/auth/logout-button/logout-button';
import { authStorage } from '@/shared/lib/auth-storage';
import { Button } from '@/shared/ui/button';

export const Header = () => {
  const { data: user } = useUser();
  const hasToken = authStorage.isAuthenticated();

  return (
    <header className="border-b">
      <div className="container mx-auto px-4 py-4 flex items-center justify-between">
        <Link to="/" className="text-2xl font-bold">
          WarehouseControl
        </Link>

        <nav className="flex items-center gap-4">
          <Link to="/items" className="hover:underline">
            Items
          </Link>

          {hasToken ? (
            <>
              {user ? (
                <span className="text-sm text-muted-foreground">
                  {user.email} ({user.role})
                </span>
              ) : null}
              <LogoutButton />
            </>
          ) : (
            <>
              <Link to="/auth/login">
                <Button variant="ghost">Login</Button>
              </Link>
            </>
          )}
        </nav>
      </div>
    </header>
  );
};

