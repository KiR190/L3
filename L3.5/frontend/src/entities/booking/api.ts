import apiClient from '@/shared/api/client';
import type { Booking, BookingRequest } from './types';

export const bookingApi = {
  async createBooking(eventId: string, data: BookingRequest): Promise<Booking> {
    const response = await apiClient.post(`/events/${eventId}/book`, data);
    return response.data;
  },

  async confirmBooking(bookingId: string): Promise<Booking> {
    const response = await apiClient.post(`/bookings/${bookingId}/confirm`);
    return response.data;
  },

  async getMyBookings(): Promise<Booking[]> {
    const response = await apiClient.get('/bookings');
    return response.data;
  },

  async getBooking(bookingId: string): Promise<Booking> {
    const response = await apiClient.get(`/bookings/${bookingId}`);
    return response.data;
  },
};

