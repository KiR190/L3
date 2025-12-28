import { useState, useEffect } from 'react';
import { useConfirmBooking } from '@/entities/booking/queries';
import { Button } from '@/shared/ui/button';
import { Alert, AlertDescription } from '@/shared/ui/alert';
import type { ApiError } from '@/shared/types/common';

interface ConfirmPaymentButtonProps {
  bookingId: string;
  expiresAt: string;
  onSuccess?: () => void;
}

export const ConfirmPaymentButton = ({ bookingId, expiresAt, onSuccess }: ConfirmPaymentButtonProps) => {
  const confirmMutation = useConfirmBooking();
  const [error, setError] = useState('');
  const [timeLeft, setTimeLeft] = useState('');

  useEffect(() => {
    const calculateTimeLeft = () => {
      const now = new Date().getTime();
      const expiry = new Date(expiresAt).getTime();
      const diff = expiry - now;

      if (diff <= 0) {
        setTimeLeft('Expired');
        return;
      }

      const minutes = Math.floor(diff / 60000);
      const seconds = Math.floor((diff % 60000) / 1000);
      setTimeLeft(`${minutes}:${seconds.toString().padStart(2, '0')}`);
    };

    calculateTimeLeft();
    const timer = setInterval(calculateTimeLeft, 1000);

    return () => clearInterval(timer);
  }, [expiresAt]);

  const handleConfirm = async () => {
    setError('');
    
    try {
      await confirmMutation.mutateAsync(bookingId);
      onSuccess?.();
    } catch (err: unknown) {
      setError((err as ApiError).error || 'Failed to confirm payment');
    }
  };

  const isExpired = timeLeft === 'Expired';
  const isUrgent = !isExpired && timeLeft.split(':')[0] !== '' && parseInt(timeLeft.split(':')[0]) < 5;

  return (
    <div className="space-y-2">
      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}
      
      <div className="flex items-center gap-2">
        <Button 
          onClick={handleConfirm} 
          disabled={confirmMutation.isPending || isExpired}
          className="flex-1"
        >
          {confirmMutation.isPending ? 'Confirming...' : 'Confirm Payment'}
        </Button>
        <div className={`text-sm font-mono px-3 py-2 rounded-md ${isUrgent ? 'bg-red-100 text-red-900' : 'bg-muted'}`}>
          {timeLeft}
        </div>
      </div>
      
      {isExpired && (
        <p className="text-xs text-destructive">This booking has expired</p>
      )}
    </div>
  );
};

