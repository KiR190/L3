import { useMemo, useState } from 'react';
import { useNavigate } from '@tanstack/react-router';
import dayjs from 'dayjs';
import { useUser } from '@/entities/user/queries';
import { useCreateItem, useDeleteItem, useItems, useUpdateItem } from '@/entities/item/queries';
import type { Item } from '@/entities/item/types';
import { authStorage } from '@/shared/lib/auth-storage';
import { Button } from '@/shared/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/shared/ui/dialog';
import { Input } from '@/shared/ui/input';
import { Label } from '@/shared/ui/label';
import { DataTable } from '@/shared/ui/data-table';
import type { ColumnDef } from '@tanstack/react-table';
import { MoreHorizontal } from 'lucide-react';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/shared/ui/dropdown-menu';

export function ItemsListPage() {
  const navigate = useNavigate();
  const { data: user, isLoading } = useUser();
  const { data: items, isLoading: itemsLoading, isError: itemsIsError, error: itemsError } = useItems();
  const createMutation = useCreateItem();
  const updateMutation = useUpdateItem();
  const deleteMutation = useDeleteItem();
  const role = user?.role;

  const canEdit = useMemo(() => role === 'admin' || role === 'manager', [role]);
  const canDelete = useMemo(() => role === 'admin', [role]);
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [isEditOpen, setIsEditOpen] = useState(false);
  const [editing, setEditing] = useState<Item | null>(null);

  const [sku, setSku] = useState('');
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [quantity, setQuantity] = useState<number>(0);
  const [location, setLocation] = useState('');

  if (!authStorage.isAuthenticated()) {
    return null;
  }

  const onCreate = async () => {
    await createMutation.mutateAsync({
      sku,
      name,
      description: description || undefined,
      quantity,
      location: location || undefined,
    });
    setIsCreateOpen(false);
  };

  const openEdit = (it: Item) => {
    setEditing(it);
    setSku(it.sku);
    setName(it.name);
    setDescription(it.description || '');
    setQuantity(it.quantity);
    setLocation(it.location || '');
    setIsEditOpen(true);
  };

  const onUpdate = async () => {
    if (!editing) return;
    await updateMutation.mutateAsync({
      id: editing.id,
      data: {
        sku,
        name,
        description: description || undefined,
        quantity,
        location: location || undefined,
      },
    });
    setIsEditOpen(false);
  };

  const onDelete = async (it: Item) => {
    await deleteMutation.mutateAsync(it.id);
  };

  const columns: ColumnDef<Item>[] = [
    {
      accessorKey: 'sku',
      header: 'SKU',
      cell: ({ row }) => <span className="font-mono">{row.original.sku}</span>,
    },
    {
      accessorKey: 'name',
      header: 'Name',
      cell: ({ row }) => row.original.name,
    },
    {
      accessorKey: 'description',
      header: 'Description',
      cell: ({ row }) => row.original.description || '-',
    },
    {
      accessorKey: 'quantity',
      header: 'Qty',
      cell: ({ row }) => row.original.quantity,
    },
    {
      accessorKey: 'location',
      header: 'Location',
      cell: ({ row }) => row.original.location || '-',
    },
    {
      accessorKey: 'updated_at',
      header: 'Updated',
      cell: ({ row }) => dayjs(row.original.updated_at).format('YYYY-MM-DD HH:mm'),
    },
    {
      id: 'actions',
      header: '',
      cell: ({ row }) => {
        const it = row.original;
        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="h-8 w-8 p-0">
                <span className="sr-only">Open menu</span>
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem onClick={() => navigate({ to: `/items/${it.id}` })}>Open</DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem disabled={!canEdit} onClick={() => openEdit(it)}>
                Edit
              </DropdownMenuItem>
              <DropdownMenuItem
                disabled={!canDelete || deleteMutation.isPending}
                className="text-destructive focus:text-destructive"
                onClick={() => onDelete(it)}
              >
                Delete
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        );
      },
    },
  ];

  return (
    <div className="container mx-auto py-8 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">Inventory</h1>
          <div className="text-sm text-muted-foreground">
            {isLoading ? 'Loading user…' : `Logged in as ${user?.email} (${user?.role})`}
          </div>
        </div>
        <div className="flex gap-2">
          <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
            <DialogTrigger asChild>
              <Button disabled={!canEdit}>Add item</Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Create item</DialogTitle>
              </DialogHeader>
              <div className="space-y-3">
                <div className="space-y-1">
                  <Label>SKU</Label>
                  <Input value={sku} onChange={(e) => setSku(e.target.value)} />
                </div>
                <div className="space-y-1">
                  <Label>Name</Label>
                  <Input value={name} onChange={(e) => setName(e.target.value)} />
                </div>
                <div className="space-y-1">
                  <Label>Description</Label>
                  <Input value={description} onChange={(e) => setDescription(e.target.value)} />
                </div>
                <div className="space-y-1">
                  <Label>Quantity</Label>
                  <Input
                    type="number"
                    value={quantity}
                    onChange={(e) => setQuantity(Number(e.target.value))}
                  />
                </div>
                <div className="space-y-1">
                  <Label>Location</Label>
                  <Input value={location} onChange={(e) => setLocation(e.target.value)} />
                </div>
                <Button onClick={onCreate} disabled={createMutation.isPending}>
                  {createMutation.isPending ? 'Creating…' : 'Create'}
                </Button>
              </div>
            </DialogContent>
          </Dialog>
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Items</CardTitle>
        </CardHeader>
        <CardContent>
          {itemsLoading ? (
            <div className="text-sm text-muted-foreground">Loading items…</div>
          ) : itemsIsError ? (
            <div className="text-sm text-destructive">
              Failed to load items: {(itemsError as Error)?.message || 'unknown error'}
            </div>
          ) : (
            <DataTable columns={columns} data={Array.isArray(items) ? items : []} filterColumnId="sku" filterPlaceholder="Filter SKU..." />
          )}
        </CardContent>
      </Card>

      <Dialog open={isEditOpen} onOpenChange={setIsEditOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit item</DialogTitle>
          </DialogHeader>
          <div className="space-y-3">
            <div className="space-y-1">
              <Label>SKU</Label>
              <Input value={sku} onChange={(e) => setSku(e.target.value)} />
            </div>
            <div className="space-y-1">
              <Label>Name</Label>
              <Input value={name} onChange={(e) => setName(e.target.value)} />
            </div>
            <div className="space-y-1">
              <Label>Description</Label>
              <Input value={description} onChange={(e) => setDescription(e.target.value)} />
            </div>
            <div className="space-y-1">
              <Label>Quantity</Label>
              <Input type="number" value={quantity} onChange={(e) => setQuantity(Number(e.target.value))} />
            </div>
            <div className="space-y-1">
              <Label>Location</Label>
              <Input value={location} onChange={(e) => setLocation(e.target.value)} />
            </div>
            <Button onClick={onUpdate} disabled={!canEdit || updateMutation.isPending}>
              {updateMutation.isPending ? 'Saving…' : 'Save'}
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}


