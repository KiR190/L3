import { EventList } from '@/widgets/event-list/event-list';

export function EventsListPage() {
  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-8">
        <h1 className="text-4xl font-bold mb-2">Upcoming Events</h1>
        <p className="text-muted-foreground">
          Browse and book tickets for amazing events
        </p>
      </div>
      
      <EventList />
    </div>
  );
}

