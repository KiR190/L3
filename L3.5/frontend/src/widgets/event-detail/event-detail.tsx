import { useEvent } from '@/entities/event/queries';
import { CreateBookingForm } from '@/features/booking/create-booking-form/create-booking-form';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/shared/ui/card';
import { Badge } from '@/shared/ui/badge';
import { Alert, AlertDescription } from '@/shared/ui/alert';
import { authStorage } from '@/shared/lib/auth-storage';
import { format } from 'date-fns';

interface EventDetailProps {
  eventId: string;
}

export const EventDetail = ({ eventId }: EventDetailProps) => {
  const { data: eventDetail, isLoading, error } = useEvent(eventId);
  const isAuthenticated = authStorage.isAuthenticated();

  if (isLoading) {
    return (
      <div className="max-w-3xl mx-auto space-y-6">
        <Card className="animate-pulse">
          <CardHeader>
            <div className="h-8 bg-muted rounded w-1/2"></div>
            <div className="h-4 bg-muted rounded w-1/3"></div>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              <div className="h-4 bg-muted rounded w-full"></div>
              <div className="h-4 bg-muted rounded w-3/4"></div>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  if (error) {
    return (
      <Alert variant="destructive">
        <AlertDescription>Failed to load event details</AlertDescription>
      </Alert>
    );
  }

  if (!eventDetail) {
    return null;
  }

  const { event, available_seats, active_bookings } = eventDetail;
  const seatPercentage = (available_seats / event?.total_seats) * 100;

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      <Card>
        <CardHeader>
          <div className="flex items-start justify-between">
            <div>
              <CardTitle className="text-3xl">{event.name}</CardTitle>
              <CardDescription className="text-lg mt-2">
                {format(new Date(event.event_date), 'PPP p')}
              </CardDescription>
            </div>
            <Badge 
              variant={seatPercentage > 50 ? 'default' : seatPercentage > 20 ? 'secondary' : 'destructive'}
              className="text-lg px-4 py-2"
            >
              {available_seats} / {event.total_seats} seats
            </Badge>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center gap-4 text-sm text-muted-foreground">
            <div>⏱️ Payment timeout: {event.payment_timeout_minutes} minutes</div>
            <div>📅 Created: {format(new Date(event.created_at), 'PP')}</div>
          </div>

          <div className="h-2 bg-muted rounded-full overflow-hidden">
            <div 
              className={`h-full transition-all ${
                seatPercentage > 50 ? 'bg-green-500' : 
                seatPercentage > 20 ? 'bg-yellow-500' : 
                'bg-red-500'
              }`}
              style={{ width: `${seatPercentage}%` }}
            />
          </div>
        </CardContent>
      </Card>

      {isAuthenticated && available_seats > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Book Your Seats</CardTitle>
            <CardDescription>
              Reserve your seats now. You'll have {event.payment_timeout_minutes} minutes to confirm payment.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <CreateBookingForm 
              eventId={event.id} 
              availableSeats={available_seats}
            />
          </CardContent>
        </Card>
      )}

      {!isAuthenticated && (
        <Alert>
          <AlertDescription>
            Please <a href="/auth/login" className="underline">login</a> or <a href="/auth/register" className="underline">register</a> to book seats.
          </AlertDescription>
        </Alert>
      )}

      {available_seats === 0 && (
        <Alert variant="destructive">
          <AlertDescription>
            This event is sold out. All seats have been booked.
          </AlertDescription>
        </Alert>
      )}

      {active_bookings?.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Active Bookings ({active_bookings?.length})</CardTitle>
            <CardDescription>Current reservations for this event</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              {active_bookings?.map((booking) => (
                <div key={booking.id} className="flex items-center justify-between py-2 border-b last:border-0">
                  <div className="text-sm">
                    {booking.seats_count} seats
                  </div>
                  <Badge variant={booking.status === 'paid' ? 'default' : booking.status === 'unpaid' ? 'secondary' : 'destructive'}>
                    {booking.status}
                  </Badge>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
};

