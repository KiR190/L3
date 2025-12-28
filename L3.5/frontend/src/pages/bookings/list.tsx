import { useMyBookings } from '@/entities/booking/queries';
import { BookingCard } from '@/widgets/booking-card/booking-card';
import { Alert, AlertDescription } from '@/shared/ui/alert';

export function BookingsListPage() {
  const { data: bookings, isLoading, error } = useMyBookings();

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-8">
        <h1 className="text-4xl font-bold mb-2">My Bookings</h1>
        <p className="text-muted-foreground">
          View and manage your event bookings
        </p>
      </div>

      {isLoading && (
        <div className="space-y-4">
          {[1, 2, 3].map((i) => (
            <div key={i} className="h-48 bg-muted animate-pulse rounded-lg"></div>
          ))}
        </div>
      )}

      {error && (
        <Alert variant="destructive">
          <AlertDescription>Failed to load bookings</AlertDescription>
        </Alert>
      )}

      {bookings && bookings.length === 0 && (
        <Alert>
          <AlertDescription>
            You don't have any bookings yet. <a href="/events" className="underline">Browse events</a> to get started!
          </AlertDescription>
        </Alert>
      )}

      {bookings && bookings.length > 0 && (
        <div className="space-y-4">
          {bookings.map((booking) => (
            <BookingCard key={booking.id} booking={booking} />
          ))}
        </div>
      )}
    </div>
  );
}

