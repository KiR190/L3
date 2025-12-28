import { createRouter, createRoute, createRootRoute, redirect } from '@tanstack/react-router';
import { authStorage } from '@/shared/lib/auth-storage';
import { Header } from '@/widgets/header/header';
import { HomePage } from '@/pages/home';
import { LoginPage } from '@/pages/auth/login';
import { ItemsListPage } from '@/pages/items/list';
import { ItemDetailPage } from '@/pages/items/detail';

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
      throw redirect({ to: '/items' });
    }
  },
});

// Protected routes
const itemsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/items',
  component: ItemsListPage,
  beforeLoad: () => {
    if (!authStorage.isAuthenticated()) {
      throw redirect({ to: '/auth/login' });
    }
  },
});

const itemDetailRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/items/$itemId',
  component: ItemDetailPage,
  beforeLoad: () => {
    if (!authStorage.isAuthenticated()) {
      throw redirect({ to: '/auth/login' });
    }
  },
});

// Create route tree
const routeTree = rootRoute.addChildren([
  indexRoute,
  authLoginRoute,
  itemsRoute,
  itemDetailRoute,
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

