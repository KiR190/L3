export interface Item {
  id: string;
  sku: string;
  name: string;
  description?: string;
  quantity: number;
  location?: string;
  created_at: string;
  updated_at: string;
}

export interface ItemCreate {
  sku: string;
  name: string;
  description?: string;
  quantity: number;
  location?: string;
}

export interface ItemUpdate {
  sku?: string;
  name?: string;
  description?: string;
  quantity?: number;
  location?: string;
}

export interface ItemHistory {
  id: string;
  item_id: string;
  action: 'INSERT' | 'UPDATE' | 'DELETE' | string;
  old_data?: unknown;
  new_data?: unknown;
  user_id?: string;
  username?: string;
  role?: string;
  created_at: string;
}




