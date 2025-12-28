import apiClient from '@/shared/api/client';
import type { Item, ItemCreate, ItemHistory, ItemUpdate } from './types';

function expectArray<T>(data: unknown, context: string): T[] {
  if (Array.isArray(data)) return data as T[];
  const preview =
    typeof data === 'string'
      ? data.slice(0, 200)
      : (() => {
          try {
            return JSON.stringify(data).slice(0, 200);
          } catch {
            return String(data);
          }
        })();
  throw new Error(`Invalid response for ${context}: expected array, got ${typeof data}. ${preview}`);
}

export const itemApi = {
  async getItems(limit = 50, offset = 0): Promise<Item[]> {
    const response = await apiClient.get('/items/', { params: { limit, offset } });
    return expectArray<Item>(response.data, 'GET /items');
  },

  async getItem(id: string): Promise<Item> {
    const response = await apiClient.get(`/items/${id}`);
    return response.data;
  },

  async createItem(data: ItemCreate): Promise<Item> {
    const response = await apiClient.post('/items/', data);
    return response.data;
  },

  async updateItem(id: string, data: ItemUpdate): Promise<Item> {
    const response = await apiClient.put(`/items/${id}`, data);
    return response.data;
  },

  async deleteItem(id: string): Promise<void> {
    await apiClient.delete(`/items/${id}`);
  },

  async getHistory(itemId: string, limit = 50, offset = 0): Promise<ItemHistory[]> {
    const response = await apiClient.get(`/items/${itemId}/history`, { params: { limit, offset } });
    return expectArray<ItemHistory>(response.data, 'GET /items/:id/history');
  },

  async downloadHistoryCSV(itemId: string): Promise<{ blob: Blob; filename: string }> {
    const response = await apiClient.get(`/items/${itemId}/history/export.csv`, {
      responseType: 'blob',
    });

    const disposition = response.headers?.['content-disposition'] as string | undefined;
    const match = disposition?.match(/filename=([^;]+)/i);
    const filename = match ? match[1].replace(/"/g, '').trim() : `item-history-${itemId}.csv`;

    return { blob: response.data as Blob, filename };
  },
};




