import { useState, useEffect } from 'react';
import { useUpdateEvent } from '@/entities/event/queries';
import type { Event } from '@/entities/event/types';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { Label } from '@/shared/ui/label';
import { Alert, AlertDescription } from '@/shared/ui/alert';
import type { ApiError } from '@/shared/types/common';
import { format } from 'date-fns';

interface EditEventFormProps {
  event: Event;
  onSuccess?: () => void;
  onCancel?: () => void;
}

export const EditEventForm = ({ event, onSuccess, onCancel }: EditEventFormProps) => {
  const updateEventMutation = useUpdateEvent();
  const [name, setName] = useState(event.name);
  const [eventDate, setEventDate] = useState('');
  const [totalSeats, setTotalSeats] = useState(event.total_seats);
  const [paymentTimeout, setPaymentTimeout] = useState(event.payment_timeout_minutes);
  const [error, setError] = useState('');

  useEffect(() => {
    // Format date for datetime-local input
    const date = new Date(event.event_date);
    const formatted = format(date, "yyyy-MM-dd'T'HH:mm");
    setEventDate(formatted);
  }, [event.event_date]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    
    try {
      await updateEventMutation.mutateAsync({
        id: event.id,
        data: {
          name,
          event_date: new Date(eventDate).toISOString(),
          total_seats: totalSeats,
          payment_timeout_minutes: paymentTimeout,
        },
      });
      onSuccess?.();
    } catch (err: unknown) {
      setError((err as ApiError).error || 'Failed to update event');
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}
      
      <div className="space-y-2">
        <Label htmlFor="name">Event Name</Label>
        <Input
          id="name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="date">Event Date & Time</Label>
        <Input
          id="date"
          type="datetime-local"
          value={eventDate}
          onChange={(e) => setEventDate(e.target.value)}
          required
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="seats">Total Seats</Label>
        <Input
          id="seats"
          type="number"
          min="1"
          max="100000"
          value={totalSeats}
          onChange={(e) => setTotalSeats(parseInt(e.target.value))}
          required
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="timeout">Payment Timeout (minutes)</Label>
        <Input
          id="timeout"
          type="number"
          min="1"
          max="1440"
          value={paymentTimeout}
          onChange={(e) => setPaymentTimeout(parseInt(e.target.value))}
          required
        />
      </div>

      <div className="flex gap-2">
        <Button type="submit" className="flex-1" disabled={updateEventMutation.isPending}>
          {updateEventMutation.isPending ? 'Updating...' : 'Update Event'}
        </Button>
        {onCancel && (
          <Button type="button" variant="outline" onClick={onCancel}>
            Cancel
          </Button>
        )}
      </div>
    </form>
  );
};

