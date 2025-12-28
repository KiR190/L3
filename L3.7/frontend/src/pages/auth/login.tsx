import { DemoLoginForm } from '@/features/auth/demo-login-form/demo-login-form';
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/shared/ui/card';

export function LoginPage() {
  return (
    <div className="max-w-md mx-auto py-12">
      <Card>
        <CardHeader>
          <CardTitle>Login to WarehouseControl</CardTitle>
          <CardDescription>
            Pick a demo role or edit credentials manually
          </CardDescription>
        </CardHeader>
        <CardContent>
          <DemoLoginForm />
        </CardContent>
        <CardFooter className="flex justify-center">
          <p className="text-sm text-muted-foreground">Demo accounts are pre-seeded in the DB.</p>
        </CardFooter>
      </Card>
    </div>
  );
}

