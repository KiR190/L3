import { createRouter, createRoute, createRootRoute, redirect } from '@tanstack/react-router';
import { authStorage } from '@/shared/lib/auth-storage';
import { Header } from '@/widgets/header/header';
import { HomePage } from '@/pages/home';
import { LoginPage } from '@/pages/auth/login';
import { RegisterPage } from '@/pages/auth/register';
import { EventsListPage } from '@/pages/events/list';
import { EventDetailPage } from '@/pages/events/detail';
import { BookingsListPage } from '@/pages/bookings/list';
import { ProfilePage } from '@/pages/profile';
import { AdminEventsPage } from '@/pages/admin/events';

// Root route with layout
const rootRoute = createRootRoute({
  component: () => (
    <div className="min-h-screen bg-background">
      <Header />
      <main>
        <Outlet />
      </main>
    </div>
  ),
});

// Home route
const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  component: HomePage,
});

// Auth routes
const authLoginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/auth/login',
  component: LoginPage,
  beforeLoad: () => {
    if (authStorage.isAuthenticated()) {
      throw redirect({ to: '/events' });
    }
  },
});

const authRegisterRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/auth/register',
  component: RegisterPage,
  beforeLoad: () => {
    if (authStorage.isAuthenticated()) {
      throw redirect({ to: '/events' });
    }
  },
});

// Event routes (public)
const eventsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/events',
  component: EventsListPage,
});

const eventDetailRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/events/$eventId',
  component: EventDetailPage,
});

// Protected routes
const bookingsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/bookings',
  component: BookingsListPage,
  beforeLoad: () => {
    if (!authStorage.isAuthenticated()) {
      throw redirect({ to: '/auth/login' });
    }
  },
});

const profileRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/profile',
  component: ProfilePage,
  beforeLoad: () => {
    if (!authStorage.isAuthenticated()) {
      throw redirect({ to: '/auth/login' });
    }
  },
});

// Admin routes
const adminEventsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/admin/events',
  component: AdminEventsPage,
  beforeLoad: () => {
    if (!authStorage.isAuthenticated()) {
      throw redirect({ to: '/auth/login' });
    }
    // Note: Role check will happen on API calls, UI will show based on user.role
  },
});

// Create route tree
const routeTree = rootRoute.addChildren([
  indexRoute,
  authLoginRoute,
  authRegisterRoute,
  eventsRoute,
  eventDetailRoute,
  bookingsRoute,
  profileRoute,
  adminEventsRoute,
]);

// Create router
export const router = createRouter({ routeTree });

// Import Outlet separately to avoid circular dependency
import { Outlet } from '@tanstack/react-router';

// Type declaration for TypeScript
declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router;
  }
}

