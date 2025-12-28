import { useState } from 'react';
import { useEvents } from '@/entities/event/queries';
import { CreateEventForm } from '@/features/event/create-event-form/create-event-form';
import { EditEventForm } from '@/features/event/edit-event-form/edit-event-form';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/shared/ui/card';
import { Button } from '@/shared/ui/button';
import { Badge } from '@/shared/ui/badge';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from '@/shared/ui/dialog';
import { Alert, AlertDescription } from '@/shared/ui/alert';
import type { Event } from '@/entities/event/types';
import { format } from 'date-fns';
import { PlusIcon } from 'lucide-react';

export function AdminEventsPage() {
  const { data: events, isLoading, error } = useEvents();
  const [showCreateDialog, setShowCreateDialog] = useState(false);
  const [editingEvent, setEditingEvent] = useState<Event | null>(null);

  if (isLoading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="h-8 bg-muted animate-pulse rounded mb-4"></div>
        <div className="h-64 bg-muted animate-pulse rounded"></div>
      </div>
    );
  }

  if (error) {
    return (
      <Alert variant="destructive">
        <AlertDescription>Failed to load events</AlertDescription>
      </Alert>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-4xl font-bold mb-2">Event Management</h1>
          <p className="text-muted-foreground">Create and manage events</p>
        </div>
        <Button onClick={() => setShowCreateDialog(true)}>
          <PlusIcon className="w-4 h-4" /> Create Event
        </Button>
      </div>

      {events && events.length === 0 && (
        <Alert>
          <AlertDescription>
            No events yet. Create your first event to get started!
          </AlertDescription>
        </Alert>
      )}

      {events && events.length > 0 && (
        <div className="space-y-4">
          {events.map((event) => (
            <Card key={event.id}>
              <CardHeader>
                <div className="flex items-start justify-between">
                  <div>
                    <CardTitle>{event.name}</CardTitle>
                    <CardDescription>
                      {format(new Date(event.event_date), 'PPP p')}
                    </CardDescription>
                  </div>
                  <Button 
                    variant="outline" 
                    size="sm"
                    onClick={() => setEditingEvent(event)}
                  >
                    Edit
                  </Button>
                </div>
              </CardHeader>
              <CardContent>
                <div className="flex gap-4 text-sm">
                  <div>
                    <span className="text-muted-foreground">Seats:</span>{' '}
                    <span className="font-semibold">{event.total_seats}</span>
                  </div>
                  <div>
                    <span className="text-muted-foreground">Timeout:</span>{' '}
                    <Badge variant="secondary">{event.payment_timeout_minutes}min</Badge>
                  </div>
                  <div>
                    <span className="text-muted-foreground">Created:</span>{' '}
                    {format(new Date(event.created_at), 'PP')}
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {/* Create Event Dialog */}
      <Dialog open={showCreateDialog} onOpenChange={setShowCreateDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Create New Event</DialogTitle>
            <DialogDescription>
              Add a new event for users to book
            </DialogDescription>
          </DialogHeader>
          <CreateEventForm onSuccess={() => setShowCreateDialog(false)} />
        </DialogContent>
      </Dialog>

      {/* Edit Event Dialog */}
      {editingEvent && (
        <Dialog open={!!editingEvent} onOpenChange={() => setEditingEvent(null)}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Edit Event</DialogTitle>
              <DialogDescription>
                Update event details
              </DialogDescription>
            </DialogHeader>
            <EditEventForm 
              event={editingEvent}
              onSuccess={() => setEditingEvent(null)}
              onCancel={() => setEditingEvent(null)}
            />
          </DialogContent>
        </Dialog>
      )}
    </div>
  );
}

