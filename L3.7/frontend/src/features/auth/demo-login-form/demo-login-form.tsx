import { useEffect, useMemo, useState } from 'react';
import { useNavigate } from '@tanstack/react-router';
import { useLogin, useUser } from '@/entities/user/queries';
import type { LoginRequest } from '@/entities/user/types';
import type { ApiError } from '@/shared/types/common';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { Label } from '@/shared/ui/label';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card';

type DemoRole = 'admin' | 'manager' | 'viewer';

const DEMO_USERS: Record<DemoRole, LoginRequest> = {
  admin: { email: 'admin@local', password: 'Admin12345' },
  manager: { email: 'bob@local', password: 'Manager12345' },
  viewer: { email: 'alice@local', password: 'Viewer12345' },
};

export const DemoLoginForm = () => {
  const navigate = useNavigate();
  const loginMutation = useLogin();
  const { data: user } = useUser();

  const [role, setRole] = useState<DemoRole>('viewer');
  const defaults = useMemo(() => DEMO_USERS[role], [role]);

  const [email, setEmail] = useState(defaults.email);
  const [password, setPassword] = useState(defaults.password);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setEmail(defaults.email);
    setPassword(defaults.password);
  }, [defaults.email, defaults.password]);

  useEffect(() => {
    if (user) {
      navigate({ to: '/items' });
    }
  }, [navigate, user]);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    try {
      await loginMutation.mutateAsync({ email, password });
      navigate({ to: '/items' });
    } catch (err) {
      setError((err as ApiError).error || 'Login failed');
    }
  };

  return (
    <Card className="w-full max-w-md">
      <CardHeader>
        <CardTitle>Login</CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={onSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="role">Demo role</Label>
            <select
              id="role"
              className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
              value={role}
              onChange={(e) => setRole(e.target.value as DemoRole)}
            >
              <option value="viewer">Viewer (read-only)</option>
              <option value="manager">Manager (edit)</option>
              <option value="admin">Admin (full)</option>
            </select>
          </div>

          <div className="space-y-2">
            <Label htmlFor="email">Email</Label>
            <Input id="email" value={email} onChange={(e) => setEmail(e.target.value)} />
          </div>

          <div className="space-y-2">
            <Label htmlFor="password">Password</Label>
            <Input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </div>

          {error && <div className="text-sm text-destructive">{error}</div>}

          <Button type="submit" className="w-full" disabled={loginMutation.isPending}>
            {loginMutation.isPending ? 'Logging in...' : 'Login'}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
};




