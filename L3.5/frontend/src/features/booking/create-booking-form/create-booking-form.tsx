import { useState } from 'react';
import { useCreateBooking } from '@/entities/booking/queries';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { Label } from '@/shared/ui/label';
import { Alert, AlertDescription } from '@/shared/ui/alert';
import type { ApiError } from '@/shared/types/common';

interface CreateBookingFormProps {
  eventId: string;
  availableSeats: number;
  onSuccess?: () => void;
}

export const CreateBookingForm = ({ eventId, availableSeats, onSuccess }: CreateBookingFormProps) => {
  const createBookingMutation = useCreateBooking();
  const [seatsCount, setSeatsCount] = useState(1);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    
    if (seatsCount > availableSeats) {
      setError(`Only ${availableSeats} seats available`);
      return;
    }
    
    try {
      await createBookingMutation.mutateAsync({
        eventId,
        data: { seats_count: seatsCount },
      });
      setSuccess('Booking created! Please confirm payment to secure your seats.');
      setSeatsCount(1);
      onSuccess?.();
    } catch (err: unknown) {
      setError((err as ApiError).error || 'Failed to create booking');
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
        <Label htmlFor="seats">Number of Seats</Label>
        <Input
          id="seats"
          type="number"
          min="1"
          max={availableSeats}
          value={seatsCount}
          onChange={(e) => setSeatsCount(parseInt(e.target.value))}
          required
        />
        <p className="text-xs text-muted-foreground">
          {availableSeats} seats available
        </p>
      </div>

      <Button type="submit" className="w-full" disabled={createBookingMutation.isPending || availableSeats === 0}>
        {createBookingMutation.isPending ? 'Booking...' : availableSeats === 0 ? 'Sold Out' : 'Book Seats'}
      </Button>
    </form>
  );
};

