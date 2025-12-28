import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { bookingApi } from './api';
import type { BookingRequest } from './types';

export const useCreateBooking = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ eventId, data }: { eventId: string; data: BookingRequest }) =>
      bookingApi.createBooking(eventId, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['events', variables.eventId] });
      queryClient.invalidateQueries({ queryKey: ['bookings'] });
    },
  });
};

export const useConfirmBooking = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: bookingApi.confirmBooking,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['bookings'] });
      queryClient.invalidateQueries({ queryKey: ['events'] });
    },
  });
};

export const useMyBookings = () => {
  return useQuery({
    queryKey: ['bookings'],
    queryFn: bookingApi.getMyBookings,
  });
};

export const useBooking = (bookingId: string) => {
  return useQuery({
    queryKey: ['bookings', bookingId],
    queryFn: () => bookingApi.getBooking(bookingId),
    enabled: !!bookingId,
  });
};

