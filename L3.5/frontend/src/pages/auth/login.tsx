import { Link } from '@tanstack/react-router';
import { LoginForm } from '@/features/auth/login-form/login-form';
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/shared/ui/card';

export function LoginPage() {
  return (
    <div className="max-w-md mx-auto py-12">
      <Card>
        <CardHeader>
          <CardTitle>Login to EventBooker</CardTitle>
          <CardDescription>
            Enter your credentials to access your account
          </CardDescription>
        </CardHeader>
        <CardContent>
          <LoginForm />
        </CardContent>
        <CardFooter className="flex justify-center">
          <p className="text-sm text-muted-foreground">
            Don't have an account?{' '}
            <Link to="/auth/register" className="underline">
              Register here
            </Link>
          </p>
        </CardFooter>
      </Card>
    </div>
  );
}

