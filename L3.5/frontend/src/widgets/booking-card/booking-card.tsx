import { useEvent } from '@/entities/event/queries';
import type { Booking } from '@/entities/booking/types';
import { ConfirmPaymentButton } from '@/features/booking/confirm-payment-button/confirm-payment-button';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/shared/ui/card';
import { Badge } from '@/shared/ui/badge';
import { format } from 'date-fns';

interface BookingCardProps {
  booking: Booking;
}

export const BookingCard = ({ booking }: BookingCardProps) => {
  const { data: eventDetail } = useEvent(booking.event_id);

  return (
    <Card>
      <CardHeader>
        <div className="flex items-start justify-between">
          <div>
            <CardTitle>{eventDetail?.event.name || 'Loading...'}</CardTitle>
            <CardDescription>
              {eventDetail && format(new Date(eventDetail.event.event_date), 'PPP p')}
            </CardDescription>
          </div>
          <Badge 
            variant={
              booking.status === 'paid' ? 'default' : 
              booking.status === 'unpaid' ? 'secondary' : 
              'destructive'
            }
          >
            {booking.status}
          </Badge>
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid grid-cols-2 gap-4 text-sm">
          <div>
            <div className="text-muted-foreground">Seats</div>
            <div className="font-semibold">{booking.seats_count}</div>
          </div>
          <div>
            <div className="text-muted-foreground">Booking ID</div>
            <div className="font-mono text-xs">{booking.id.slice(0, 8)}...</div>
          </div>
          <div>
            <div className="text-muted-foreground">Created</div>
            <div>{format(new Date(booking.created_at), 'PP p')}</div>
          </div>
          {booking.paid_at && (
            <div>
              <div className="text-muted-foreground">Paid</div>
              <div>{format(new Date(booking.paid_at), 'PP p')}</div>
            </div>
          )}
          {booking.status === 'unpaid' && (
            <div>
              <div className="text-muted-foreground">Expires</div>
              <div>{format(new Date(booking.expires_at), 'PP p')}</div>
            </div>
          )}
        </div>

        {booking.status === 'unpaid' && (
          <ConfirmPaymentButton 
            bookingId={booking.id} 
            expiresAt={booking.expires_at}
          />
        )}
      </CardContent>
    </Card>
  );
};

