import apiClient from '@/shared/api/client';
import type { Event, EventDetail, CreateEventRequest, UpdateEventRequest } from './types';

export const eventApi = {
  async getEvents(): Promise<Event[]> {
    const response = await apiClient.get('/events');
    return response.data;
  },

  async getEvent(id: string): Promise<EventDetail> {
    const response = await apiClient.get(`/events/${id}`);
    return response.data;
  },

  async createEvent(data: CreateEventRequest): Promise<Event> {
    const response = await apiClient.post('/events', data);
    return response.data;
  },

  async updateEvent(id: string, data: UpdateEventRequest): Promise<Event> {
    const response = await apiClient.put(`/events/${id}`, data);
    return response.data;
  },
};

