import { useState } from 'react';
import { useCreateEvent } from '@/entities/event/queries';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { Label } from '@/shared/ui/label';
import { Alert, AlertDescription } from '@/shared/ui/alert';
import type { ApiError } from '@/shared/types/common';

interface CreateEventFormProps {
  onSuccess?: () => void;
}

export const CreateEventForm = ({ onSuccess }: CreateEventFormProps) => {
  const createEventMutation = useCreateEvent();
  const [name, setName] = useState('');
  const [eventDate, setEventDate] = useState('');
  const [totalSeats, setTotalSeats] = useState(100);
  const [paymentTimeout, setPaymentTimeout] = useState(30);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    
    try {
      await createEventMutation.mutateAsync({
        name,
        event_date: new Date(eventDate).toISOString(),
        total_seats: totalSeats,
        payment_timeout_minutes: paymentTimeout,
      });
      setSuccess('Event created successfully!');
      setName('');
      setEventDate('');
      setTotalSeats(100);
      setPaymentTimeout(30);
      onSuccess?.();
    } catch (err: unknown) {
      setError((err as ApiError).error || 'Failed to create event');
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}
      
      {success && (
        <Alert>
          <AlertDescription>{success}</AlertDescription>
        </Alert>
      )}
      
      <div className="space-y-2">
        <Label htmlFor="name">Event Name</Label>
        <Input
          id="name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Concert Night"
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
        <p className="text-xs text-muted-foreground">
          How long users have to confirm payment (max 24 hours)
        </p>
      </div>

      <Button type="submit" className="w-full" disabled={createEventMutation.isPending}>
        {createEventMutation.isPending ? 'Creating...' : 'Create Event'}
      </Button>
    </form>
  );
};

