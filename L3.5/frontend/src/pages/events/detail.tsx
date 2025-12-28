import { useParams, Link } from '@tanstack/react-router';
import { EventDetail } from '@/widgets/event-detail/event-detail';
import { Button } from '@/shared/ui/button';

export function EventDetailPage() {
  const { eventId } = useParams({ strict: false }) as { eventId: string };

  return (
    <div className="container mx-auto px-4 py-8">
      <Link to="/events">
        <Button variant="ghost" className="mb-4">
          ← Back to Events
        </Button>
      </Link>
      
      <EventDetail eventId={eventId} />
    </div>
  );
}

