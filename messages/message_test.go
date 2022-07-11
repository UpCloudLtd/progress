package messages

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMessageStore_Push_Errors(t *testing.T) {
	for _, test := range []struct {
		name          string
		update        Update
		expectedError string
	}{
		{
			name:          "No key/message",
			update:        Update{Status: MessageStatusSuccess},
			expectedError: "can not push message without key or message",
		},
		{
			name:          "No message",
			update:        Update{Key: "test", Status: MessageStatusSuccess},
			expectedError: "can not push message with empty message",
		},
		{
			name:          "No status",
			update:        Update{Message: "Testing"},
			expectedError: `can not push message with invalid status ""`,
		},
		{
			name:          "Invalid status",
			update:        Update{Message: "Testing", Status: "invalid"},
			expectedError: `can not push message with invalid status "invalid"`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			ms := NewMessageStore()
			err := ms.Push(test.update)
			assert.EqualError(t, err, test.expectedError)
		})
	}
}

func TestMessageGroup_ListInProgress_ListFinished(t *testing.T) {
	ms := NewMessageStore()

	assert.NoError(t, ms.Push(Update{Message: "2nd", Status: MessageStatusPending}))
	assert.NoError(t, ms.Push(Update{Message: "1st", Status: MessageStatusStarted}))
	assert.NoError(t, ms.Push(Update{Message: "2nd", Status: MessageStatusStarted}))

	messages := ms.ListInProgress()
	assert.Len(t, ms.ListFinished(), 0)
	assert.Equal(t, "1st", messages[0].Message)
	assert.Equal(t, "2nd", messages[1].Message)

	assert.NoError(t, ms.Push(Update{Message: "2nd", Status: MessageStatusSuccess}))

	assert.Len(t, ms.ListInProgress(), 1)
	assert.Len(t, ms.ListFinished(), 1)
}

func TestMessageStore_Push_UpdatesMessage(t *testing.T) {
	ms := NewMessageStore()

	assert.NoError(t, ms.Push(Update{Key: "test", Message: "Testing", Status: MessageStatusPending}))

	msg := ms.ListInProgress()[0]
	assert.Equal(t, "Testing", msg.Message)
	assert.True(t, msg.Started.IsZero())
	assert.True(t, msg.Finished.IsZero())

	tic := time.Now()
	assert.NoError(t, ms.Push(Update{Key: "test", Message: "Still testing", Status: MessageStatusStarted}))

	assert.Equal(t, "Still testing", msg.Message)
	assert.True(t, tic.Before(msg.Started))
	assert.True(t, msg.Finished.IsZero())

	toc := time.Now()
	assert.NoError(t, ms.Push(Update{Key: "test", Status: MessageStatusError, Details: "Test details"}))

	assert.Equal(t, "Still testing", msg.Message)
	assert.True(t, toc.After(msg.Started))
	assert.True(t, toc.Before(msg.Finished))
	assert.Equal(t, msg.Details, "Test details")
}
