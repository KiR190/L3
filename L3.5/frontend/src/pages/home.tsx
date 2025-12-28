import { Link } from '@tanstack/react-router';
import { Button } from '@/shared/ui/button';
import { Card, CardHeader, CardTitle, CardDescription } from '@/shared/ui/card';

export function HomePage() {
  return (
    <div className="max-w-4xl mx-auto space-y-8 py-12">
      <div className="text-center space-y-4">
        <h1 className="text-5xl font-bold">Welcome to EventBooker</h1>
        <p className="text-xl text-muted-foreground">
          Book your seats for amazing events with our secure reservation system
        </p>
        <div className="flex justify-center gap-4">
          <Link to="/events">
            <Button size="lg">Browse Events</Button>
          </Link>
          <Link to="/auth/register">
            <Button variant="outline" size="lg">Get Started</Button>
          </Link>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-12">
        <Card>
          <CardHeader>
            <CardTitle>🎫 Easy Booking</CardTitle>
            <CardDescription>
              Reserve your seats in seconds with our simple booking flow
            </CardDescription>
          </CardHeader>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>⏱️ Secure Hold</CardTitle>
            <CardDescription>
              Your seats are held while you complete payment - no rush!
            </CardDescription>
          </CardHeader>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>🔔 Smart Notifications</CardTitle>
            <CardDescription>
              Get updates via email or Telegram - never miss an update
            </CardDescription>
          </CardHeader>
        </Card>
      </div>
    </div>
  );
}

