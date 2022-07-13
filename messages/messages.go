package messages

import (
	"fmt"
	"sort"
	"time"
)

type Update struct {
	Key     string
	Message string
	Status  MessageStatus
	Details string
}

type Message struct {
	Key      string
	Message  string
	Status   MessageStatus
	Details  string
	Created  time.Time
	Started  time.Time
	Finished time.Time
}

func getMessageKey(key, message string) string {
	if key != "" {
		return key
	}
	return message
}

func validateMessage(message string) error {
	if message == "" {
		return fmt.Errorf("can not push message with empty message")
	}
	return nil
}

func validateStatus(status MessageStatus) error {
	if !status.IsValid() {
		return fmt.Errorf(`can not push message with invalid status "%s"`, status)
	}
	return nil
}

func (msg *Message) update(update Update) {
	if msg.Created.IsZero() {
		msg.Created = time.Now()
	}
	if update.Status != MessageStatusPending && msg.Started.IsZero() {
		msg.Started = time.Now()
	}

	if update.Status.IsFinished() {
		msg.Finished = time.Now()

		if msg.Started.IsZero() {
			msg.Started = msg.Finished
		}
	}

	if msg.Key == "" {
		msg.Key = getMessageKey(update.Key, update.Message)
	}

	if update.Message != "" {
		msg.Message = update.Message
	}
	if update.Status != "" {
		msg.Status = update.Status
	}
	if update.Details != "" {
		msg.Details = update.Details
	}
}

func (msg Message) ElapsedSeconds() float64 {
	if msg.Started.IsZero() {
		return 0
	}

	end := time.Now()
	if !msg.Finished.IsZero() {
		end = msg.Finished
	}

	return end.Sub(msg.Started).Seconds()
}

type MessageStore struct {
	inProgress map[string]*Message
	finished   []*Message
}

func NewMessageStore() *MessageStore {
	return &MessageStore{
		inProgress: make(map[string]*Message),
	}
}

func (ms *MessageStore) storeMessage(msg *Message) {
	if msg.Status.IsFinished() {
		delete(ms.inProgress, msg.Key)
		ms.finished = append(ms.finished, msg)
	} else {
		ms.inProgress[msg.Key] = msg
	}
}

// Add existing Message to Message store. Useful for adding, for example, historical data to MessageStore. For live data, prefer Push.
func (ms *MessageStore) Add(msg Message) error {
	if err := validateMessage(msg.Message); err != nil {
		return err
	}
	if err := validateStatus(msg.Status); err != nil {
		return err
	}

	msg.Key = getMessageKey(msg.Key, msg.Message)
	ms.storeMessage(&msg)

	return nil
}

// Push update to MessageStore. Automatically creates/updates Message and sets it timestamps based on update content.
func (ms *MessageStore) Push(update Update) error {
	key := getMessageKey(update.Key, update.Message)
	if key == "" {
		return fmt.Errorf("can not push message without key or message")
	}

	var msg *Message
	if prev, ok := ms.inProgress[key]; !ok {
		if err := validateMessage(update.Message); err != nil {
			return err
		}
		if err := validateStatus(update.Status); err != nil {
			return err
		}

		msg = &Message{}
	} else {
		msg = prev
	}

	msg.update(update)
	ms.storeMessage(msg)
	return nil
}

// ListInprogress lists in-progress messages in MessageStore sorted by started time.
func (ms *MessageStore) ListInProgress() []*Message {
	messages := []*Message{}
	for _, msg := range ms.inProgress {
		messages = append(messages, msg)
	}

	sort.Slice(messages, func(i, j int) bool {
		// Sort zero before any value
		if messages[i].Started.IsZero() && !messages[j].Started.IsZero() {
			return true
		}
		if !messages[i].Started.IsZero() && messages[j].Started.IsZero() {
			return false
		}
		// For not started, sort by created time
		if messages[i].Started.IsZero() && messages[j].Started.IsZero() {
			return messages[i].Created.Before(messages[j].Created)
		}
		// Sort by started time
		return messages[i].Started.Before(messages[j].Started)
	})

	return messages
}

// ListFinished lists finished messages in MessageStore in order they were marked finished
func (ms *MessageStore) ListFinished() []*Message {
	return ms.finished
}

// Close sets status of pending messages to skipped and started message to unknown.
func (ms *MessageStore) Close() {
	for _, msg := range ms.ListInProgress() {
		if msg.Status == MessageStatusPending {
			ms.Push(Update{ //nolint:errcheck
				Key:    msg.Key,
				Status: MessageStatusSkipped,
			})
		}
		if msg.Status == MessageStatusStarted {
			ms.Push(Update{ //nolint:errcheck
				Key:    msg.Key,
				Status: MessageStatusUnknown,
			})
		}
	}
}
