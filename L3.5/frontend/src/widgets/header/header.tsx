import { Link } from '@tanstack/react-router';
import { useUser } from '@/entities/user/queries';
import { LogoutButton } from '@/features/auth/logout-button/logout-button';
import { Button } from '@/shared/ui/button';
import { Badge } from '@/shared/ui/badge';

export const Header = () => {
  const { data: user, isLoading } = useUser();
  const isAuthenticated = !!user && !isLoading;
  const isAdmin = user?.role === 'admin';

  return (
    <header className="border-b">
      <div className="container mx-auto px-4 py-4 flex items-center justify-between">
        <Link to="/" className="text-2xl font-bold">
          🎫 EventBooker
        </Link>

        <nav className="flex items-center gap-4">
          <Link to="/events" className="hover:underline">
            Events
          </Link>

          {isAuthenticated ? (
            <>
              <Link to="/bookings" className="hover:underline">
                My Bookings
              </Link>
              {isAdmin && (
                <Link to="/admin/events" className="hover:underline flex items-center gap-1">
                  Admin
                  <Badge variant="destructive" className="text-xs">Admin</Badge>
                </Link>
              )}
              <Link to="/profile" className="hover:underline">
                Profile
              </Link>
              <LogoutButton />
            </>
          ) : (
            <>
              <Link to="/auth/login">
                <Button variant="ghost">Login</Button>
              </Link>
              <Link to="/auth/register">
                <Button>Register</Button>
              </Link>
            </>
          )}
        </nav>
      </div>
    </header>
  );
};

