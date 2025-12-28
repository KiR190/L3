package worker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"event-booker/internal/models"
	"event-booker/internal/queue"
	"event-booker/internal/repository"

	"github.com/segmentio/kafka-go"
	"github.com/wb-go/wbf/retry"
)

// NotificationService interface for sending notifications
type NotificationService interface {
	NotifyBookingCancelled(ctx context.Context, booking *models.BookingWithUser, event *models.Event) error
}

type Worker struct {
	repo                *repository.Repository
	queue               queue.QueueInterface
	notificationService NotificationService
}

func NewWorker(repo *repository.Repository, queue queue.QueueInterface, notificationService NotificationService) *Worker {
	return &Worker{
		repo:                repo,
		queue:               queue,
		notificationService: notificationService,
	}
}

// Start begins processing expiration tasks from the queue
func (w *Worker) Start(ctx context.Context) {
	log.Println("Worker started: monitoring booking expirations...")

	// Configure retry strategy for message processing
	strategy := retry.Strategy{
		Attempts: 3,
		Delay:    5 * time.Second,
		Backoff:  2,
	}

	// Create message channel
	msgChan := make(chan kafka.Message)

	// Subscribe to the queue
	go w.queue.Subscribe(ctx, msgChan, strategy)

	// Start two background goroutines:
	// 1. Process messages from queue
	go w.processQueueMessages(ctx, msgChan)

	// 2. Periodic check for expired bookings (backup mechanism)
	go w.periodicExpirationCheck(ctx)
}

func (w *Worker) processQueueMessages(ctx context.Context, msgChan <-chan kafka.Message) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Worker: stopping queue message processing")
			return
		case msg, ok := <-msgChan:
			if !ok {
				log.Println("Worker: message channel closed")
				return
			}

			// Decode the Kafka message
			var task models.ExpirationTask
			if err := json.Unmarshal(msg.Value, &task); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				// Commit the message anyway to avoid reprocessing bad messages
				if err := w.queue.Commit(ctx, msg); err != nil {
					log.Printf("Error committing message: %v", err)
				}
				continue
			}

			// Handle the expiration task
			w.handleExpirationTask(ctx, task)

			// Commit the message
			if err := w.queue.Commit(ctx, msg); err != nil {
				log.Printf("Error committing message: %v", err)
			}
		}
	}
}

func (w *Worker) handleExpirationTask(ctx context.Context, task models.ExpirationTask) {
	log.Printf("Processing expiration task for booking: %s (scheduled expiry: %s)", task.BookingID, task.ExpiresAt.Format(time.RFC3339))

	// Get the booking with user info
	bookingWithUser, err := w.repo.GetBookingWithUser(ctx, task.BookingID)
	if err != nil {
		if err == repository.ErrBookingNotFound {
			log.Printf("Booking %s not found (possibly already cancelled)", task.BookingID)
			return
		}
		log.Printf("Error fetching booking %s: %v", task.BookingID, err)
		return
	}

	// Check if booking is still unpaid
	if bookingWithUser.Status != models.BookingStatusUnpaid {
		log.Printf("Booking %s is %s, skipping cancellation", task.BookingID, bookingWithUser.Status)
		return
	}

	// Check if it's actually expired
	if time.Now().Before(bookingWithUser.ExpiresAt) {
		log.Printf("Booking %s not yet expired (expires at %s)", task.BookingID, bookingWithUser.ExpiresAt.Format(time.RFC3339))
		return
	}

	// Cancel the booking
	err = w.repo.CancelBooking(ctx, task.BookingID)
	if err != nil {
		log.Printf("Failed to cancel booking %s: %v", task.BookingID, err)
		return
	}

	log.Printf("Successfully cancelled expired booking: %s (event: %s, %d seats released)", task.BookingID, bookingWithUser.EventID, bookingWithUser.SeatsCount)

	// Get event info for notification
	event, err := w.repo.GetEventByID(ctx, bookingWithUser.EventID)
	if err != nil {
		log.Printf("Failed to get event info for notification: %v", err)
		// Don't return - cancellation was successful
	} else {
		// Send notification
		err = w.notificationService.NotifyBookingCancelled(ctx, bookingWithUser, event)
		if err != nil {
			log.Printf("Failed to send notification for booking %s: %v", task.BookingID, err)
		}
	}
}

// periodicExpirationCheck is a backup mechanism that runs every minute
// to check for any expired bookings that might have been missed
func (w *Worker) periodicExpirationCheck(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	log.Println("Worker: started periodic expiration check (every 1 minute)")

	for {
		select {
		case <-ctx.Done():
			log.Println("Worker: stopping periodic expiration check")
			return
		case <-ticker.C:
			w.checkAndCancelExpiredBookings(ctx)
		}
	}
}

func (w *Worker) checkAndCancelExpiredBookings(ctx context.Context) {
	// Get all expired unpaid bookings with user info
	expiredBookings, err := w.repo.GetExpiredUnpaidBookingsWithUser(ctx)
	if err != nil {
		log.Printf("Error fetching expired bookings: %v", err)
		return
	}

	if len(expiredBookings) == 0 {
		return
	}

	log.Printf("Found %d expired unpaid bookings to cancel", len(expiredBookings))

	for _, bookingWithUser := range expiredBookings {
		err := w.repo.CancelBooking(ctx, bookingWithUser.ID)
		if err != nil {
			log.Printf("Failed to cancel expired booking %s: %v", bookingWithUser.ID, err)
			continue
		}

		log.Printf("Cancelled expired booking: %s (event: %s, %d seats released)", bookingWithUser.ID, bookingWithUser.EventID, bookingWithUser.SeatsCount)

		// Get event info for notification
		event, err := w.repo.GetEventByID(ctx, bookingWithUser.EventID)
		if err != nil {
			log.Printf("Failed to get event info for notification: %v", err)
			continue
		}

		// Send notification
		err = w.notificationService.NotifyBookingCancelled(ctx, &bookingWithUser, event)
		if err != nil {
			log.Printf("Failed to send notification for booking %s: %v", bookingWithUser.ID, err)
		}
	}
}


