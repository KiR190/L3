export interface Event {
  id: string;
  name: string;
  event_date: string;
  total_seats: number;
  payment_timeout_minutes: number;
  created_at: string;
  updated_at: string;
}

export interface Booking {
  id: string;
  event_id: string;
  user_id: string;
  seats_count: number;
  status: 'unpaid' | 'paid' | 'cancelled';
  created_at: string;
  paid_at?: string;
  expires_at: string;
}

export interface EventDetail {
  event: Event;
  available_seats: number;
  active_bookings: Booking[];
}

export interface CreateEventRequest {
  name: string;
  event_date: string;
  total_seats: number;
  payment_timeout_minutes?: number;
}

export interface UpdateEventRequest {
  name?: string;
  event_date?: string;
  total_seats?: number;
  payment_timeout_minutes?: number;
}

