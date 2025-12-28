import type { Booking as BookingBase } from '../event/types';

export interface BookingRequest {
  seats_count: number;
}

export type Booking = BookingBase;

