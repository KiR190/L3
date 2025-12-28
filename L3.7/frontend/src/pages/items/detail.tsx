import { useEffect, useMemo, useState } from 'react';
import { useNavigate, useParams } from '@tanstack/react-router';
import dayjs from 'dayjs';

import { useDeleteItem, useItem, useUpdateItem } from '@/entities/item/queries';
import { useUser } from '@/entities/user/queries';
import { Button } from '@/shared/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card';
import { Input } from '@/shared/ui/input';
import { Label } from '@/shared/ui/label';
import { ItemHistoryTable } from '@/widgets/item-history/item-history-table';

export function ItemDetailPage() {
  const navigate = useNavigate();
  const { itemId } = useParams({ strict: false }) as { itemId: string };

  const { data: user } = useUser();
  const role = user?.role;
  const canEdit = useMemo(() => role === 'admin' || role === 'manager', [role]);
  const canDelete = useMemo(() => role === 'admin', [role]);

  const { data: item, isLoading, isError, error } = useItem(itemId);
  const updateMutation = useUpdateItem();
  const deleteMutation = useDeleteItem();

  const [sku, setSku] = useState('');
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [quantity, setQuantity] = useState<number>(0);
  const [location, setLocation] = useState('');

  useEffect(() => {
    if (!item) return;
    setSku(item.sku);
    setName(item.name);
    setDescription(item.description || '');
    setQuantity(item.quantity);
    setLocation(item.location || '');
  }, [itemId, item]);

  const onSave = async () => {
    if (!item) return;
    await updateMutation.mutateAsync({
      id: item.id,
      data: {
        sku,
        name,
        description: description || undefined,
        quantity,
        location: location || undefined,
      },
    });
  };

  const onDelete = async () => {
    if (!item) return;
    await deleteMutation.mutateAsync(item.id);
    navigate({ to: '/items' });
  };

  return (
    <div className="container mx-auto py-8 space-y-6">
      <div className="flex items-start justify-between gap-4">
        <div>
          <div className="text-sm text-muted-foreground">Item</div>
          <h1 className="text-2xl font-bold font-mono">{itemId}</h1>
        </div>
        <Button variant="outline" onClick={() => navigate({ to: '/items' })}>
          Back to list
        </Button>
      </div>

      {isLoading ? (
        <div className="text-sm text-muted-foreground">Loading item…</div>
      ) : isError ? (
        <div className="text-sm text-destructive">Failed to load item: {(error as Error)?.message}</div>
      ) : item ? (
        <>
          <Card>
            <CardHeader>
              <CardTitle>Details</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-1">
                  <Label>SKU</Label>
                  <Input value={sku} onChange={(e) => setSku(e.target.value)} disabled={!canEdit} />
                </div>
                <div className="space-y-1">
                  <Label>Name</Label>
                  <Input value={name} onChange={(e) => setName(e.target.value)} disabled={!canEdit} />
                </div>
                <div className="space-y-1">
                  <Label>Description</Label>
                  <Input value={description} onChange={(e) => setDescription(e.target.value)} disabled={!canEdit} />
                </div>
                <div className="space-y-1">
                  <Label>Quantity</Label>
                  <Input
                    type="number"
                    value={quantity}
                    onChange={(e) => setQuantity(Number(e.target.value))}
                    disabled={!canEdit}
                  />
                </div>
                <div className="space-y-1">
                  <Label>Location</Label>
                  <Input value={location} onChange={(e) => setLocation(e.target.value)} disabled={!canEdit} />
                </div>
              </div>

              <div className="text-sm text-muted-foreground">
                Created: {dayjs(item.created_at).format('YYYY-MM-DD HH:mm')} · Updated:{' '}
                {dayjs(item.updated_at).format('YYYY-MM-DD HH:mm')}
              </div>

              <div className="flex gap-2">
                <Button onClick={onSave} disabled={!canEdit || updateMutation.isPending}>
                  {updateMutation.isPending ? 'Saving…' : 'Save'}
                </Button>
                <Button variant="destructive" onClick={onDelete} disabled={!canDelete || deleteMutation.isPending}>
                  {deleteMutation.isPending ? 'Deleting…' : 'Delete'}
                </Button>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>History</CardTitle>
            </CardHeader>
            <CardContent>
              <ItemHistoryTable itemId={item.id} />
            </CardContent>
          </Card>
        </>
      ) : null}
    </div>
  );
}


