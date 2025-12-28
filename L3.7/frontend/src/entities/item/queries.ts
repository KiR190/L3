import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { itemApi } from './api';

export const useItems = (limit = 50, offset = 0) => {
  return useQuery({
    queryKey: ['items', limit, offset],
    queryFn: () => itemApi.getItems(limit, offset),
  });
};

export const useItem = (id: string) => {
  return useQuery({
    queryKey: ['item', id],
    queryFn: () => itemApi.getItem(id),
    enabled: !!id,
  });
};

export const useCreateItem = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: itemApi.createItem,
    onSuccess: (created) => {
      qc.invalidateQueries({ queryKey: ['items'] });
      qc.invalidateQueries({ queryKey: ['itemHistory', created.id] });
    },
  });
};

export const useUpdateItem = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Parameters<typeof itemApi.updateItem>[1] }) =>
      itemApi.updateItem(id, data),
    onSuccess: (_data, variables) => {
      qc.invalidateQueries({ queryKey: ['items'] });
      qc.invalidateQueries({ queryKey: ['item'] });
      qc.invalidateQueries({ queryKey: ['itemHistory', variables.id] });
    },
  });
};

export const useDeleteItem = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: itemApi.deleteItem,
    onSuccess: (_data, id) => {
      qc.invalidateQueries({ queryKey: ['items'] });
      qc.invalidateQueries({ queryKey: ['item'] });
      qc.invalidateQueries({ queryKey: ['itemHistory', id] });
    },
  });
};

export const useItemHistory = (itemId: string, limit = 50, offset = 0) => {
  return useQuery({
    queryKey: ['itemHistory', itemId, limit, offset],
    queryFn: () => itemApi.getHistory(itemId, limit, offset),
    enabled: !!itemId,
  });
};




