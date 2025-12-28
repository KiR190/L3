package sender

import (
	"fmt"
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TelegramSender handles Telegram notifications
type TelegramSender struct {
	bot   *tgbotapi.BotAPI
	users map[string]int64 // username -> chatID mapping
	lock  sync.RWMutex
	
	// Callback for when a user registers
	onUserRegister func(username string, chatID int64)
}

// NewTelegramSender creates a new Telegram sender
func NewTelegramSender(token string) (*TelegramSender, error) {
	if token == "" {
		return nil, fmt.Errorf("telegram bot token is empty")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w", err)
	}

	log.Printf("Telegram bot authorized as @%s", bot.Self.UserName)

	return &TelegramSender{
		bot:   bot,
		users: make(map[string]int64),
	}, nil
}

// SetUserRegisterCallback sets the callback for user registration
func (t *TelegramSender) SetUserRegisterCallback(callback func(username string, chatID int64)) {
	t.onUserRegister = callback
}

// ListenAndServe starts listening for Telegram updates
func (t *TelegramSender) ListenAndServe() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.bot.GetUpdatesChan(u)

	log.Println("Telegram bot started listening for commands...")

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			t.handleCommand(update.Message)
		}
	}
}

// handleCommand processes bot commands
func (t *TelegramSender) handleCommand(msg *tgbotapi.Message) {
	switch msg.Command() {
	case "register", "start":
		username := msg.From.UserName
		chatID := msg.Chat.ID

		if username == "" {
			reply := tgbotapi.NewMessage(chatID, "❌ You need to set a Telegram username first to use this bot.")
			_, _ = t.bot.Send(reply)
			return
		}

		t.lock.Lock()
		t.users[username] = chatID
		t.lock.Unlock()

		// Call the registration callback if set
		if t.onUserRegister != nil {
			t.onUserRegister(username, chatID)
		}

		reply := tgbotapi.NewMessage(chatID, fmt.Sprintf("✅ Successfully registered!\n\nYour Telegram username: @%s\nYou will now receive booking notifications here.", username))
		_, err := t.bot.Send(reply)
		if err != nil {
			log.Printf("Failed to send registration confirmation: %v", err)
		}

		log.Printf("Telegram user registered: @%s (chatID: %d)", username, chatID)

	default:
		reply := tgbotapi.NewMessage(msg.Chat.ID, "Unknown command. Use /register to register for notifications.")
		_, _ = t.bot.Send(reply)
	}
}

// Send sends a message to a user by username
func (t *TelegramSender) Send(username, subject, body string) error {
	t.lock.RLock()
	chatID, ok := t.users[username]
	t.lock.RUnlock()

	if !ok {
		return fmt.Errorf("user @%s not registered in Telegram bot", username)
	}

	// Combine subject and body
	message := fmt.Sprintf("*%s*\n\n%s", subject, body)

	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "Markdown"

	_, err := t.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send Telegram message to @%s: %w", username, err)
	}

	log.Printf("📱 Telegram message sent to @%s: %s", username, subject)
	return nil
}

// SendToChat sends a message directly to a chat ID
func (t *TelegramSender) SendToChat(chatID int64, subject, body string) error {
	message := fmt.Sprintf("*%s*\n\n%s", subject, body)

	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "Markdown"

	_, err := t.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send Telegram message to chat %d: %w", chatID, err)
	}

	log.Printf("📱 Telegram message sent to chat %d: %s", chatID, subject)
	return nil
}

// RegisterUser manually registers a user mapping
func (t *TelegramSender) RegisterUser(username string, chatID int64) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.users[username] = chatID
}

// Stop stops the Telegram bot
func (t *TelegramSender) Stop() {
	log.Println("Stopping Telegram bot...")
	t.bot.StopReceivingUpdates()
	log.Println("Telegram bot stopped")
}

