import { Link } from '@tanstack/react-router';
import { Button } from '@/shared/ui/button';
import { Card, CardHeader, CardTitle, CardDescription } from '@/shared/ui/card';
import { authStorage } from '@/shared/lib/auth-storage';

export function HomePage() {
  return (
    <div className="max-w-4xl mx-auto space-y-8 py-12">
      <div className="text-center space-y-4">
        <h1 className="text-5xl font-bold">WarehouseControl</h1>
        <p className="text-xl text-muted-foreground">
          Inventory CRUD with trigger-based audit log 
        </p>
        <div className="flex justify-center gap-4">
          <Link to="/items">
            <Button size="lg">Open Inventory</Button>
          </Link>
          {!authStorage.isAuthenticated() && (
            <Link to="/auth/login">
              <Button variant="outline" size="lg">Login</Button>
            </Link>
          )}
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-12">
        <Card>
          <CardHeader>
            <CardTitle>RBAC</CardTitle>
            <CardDescription>
              Viewer: read-only, Manager: edit, Admin: full access
            </CardDescription>
          </CardHeader>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Audit history</CardTitle>
            <CardDescription>
              Every change is tracked with before/after values
            </CardDescription>
          </CardHeader>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>DB triggers</CardTitle>
            <CardDescription>
              Trigger captures changes automatically (anti-pattern demo)
            </CardDescription>
          </CardHeader>
        </Card>
      </div>
    </div>
  );
}

