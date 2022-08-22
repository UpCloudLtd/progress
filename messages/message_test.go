package messages_test

import (
	"testing"
	"time"

	"github.com/UpCloudLtd/progress/messages"
	"github.com/stretchr/testify/assert"
)

func TestMessageStore_Push_Errors(t *testing.T) {
	for _, test := range []struct {
		name          string
		update        messages.Update
		expectedError string
	}{
		{
			name:          "No key/message",
			update:        messages.Update{Status: messages.MessageStatusSuccess},
			expectedError: "can not push message without key or message",
		},
		{
			name:          "No message",
			update:        messages.Update{Key: "test", Status: messages.MessageStatusSuccess},
			expectedError: "can not push message with empty message",
		},
		{
			name:          "No status",
			update:        messages.Update{Message: "Testing"},
			expectedError: `can not push message with invalid status ""`,
		},
		{
			name:          "Invalid status",
			update:        messages.Update{Message: "Testing", Status: "invalid"},
			expectedError: `can not push message with invalid status "invalid"`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			ms := messages.NewMessageStore()
			err := ms.Push(test.update)
			assert.EqualError(t, err, test.expectedError)
		})
	}
}

func TestMessageGroup_ListInProgress_ListFinished(t *testing.T) {
	ms := messages.NewMessageStore()

	assert.NoError(t, ms.Push(messages.Update{Message: "2nd", Status: messages.MessageStatusPending}))
	assert.NoError(t, ms.Push(messages.Update{Message: "1st", Status: messages.MessageStatusStarted}))
	time.Sleep(time.Microsecond * 25) // Ensure time difference on Windows
	assert.NoError(t, ms.Push(messages.Update{Message: "2nd", Status: messages.MessageStatusStarted}))

	inProgress := ms.ListInProgress()
	assert.Len(t, ms.ListFinished(), 0)
	assert.Equal(t, "1st", inProgress[0].Message)
	assert.Equal(t, "2nd", inProgress[1].Message)

	assert.NoError(t, ms.Push(messages.Update{Message: "2nd", Status: messages.MessageStatusSuccess}))

	assert.Len(t, ms.ListInProgress(), 1)
	assert.Len(t, ms.ListFinished(), 1)
}

func TestMessageStore_Push_UpdatesMessage(t *testing.T) {
	ms := messages.NewMessageStore()

	assert.NoError(t, ms.Push(messages.Update{Key: "test", Message: "Testing", Status: messages.MessageStatusPending}))

	msg := ms.ListInProgress()[0]
	assert.Equal(t, "Testing", msg.Message)
	assert.True(t, msg.Started.IsZero())
	assert.True(t, msg.Finished.IsZero())

	tic := time.Now()
	time.Sleep(time.Microsecond * 25) // Ensure time difference on Windows
	assert.NoError(t, ms.Push(messages.Update{Key: "test", Message: "Still testing", Status: messages.MessageStatusStarted}))

	assert.Equal(t, "Still testing", msg.Message)
	assert.True(t, tic.Before(msg.Started))
	assert.True(t, msg.Finished.IsZero())

	time.Sleep(time.Microsecond * 25) // Ensure time difference on Windows
	toc := time.Now()
	time.Sleep(time.Microsecond * 25) // Ensure time difference on Windows
	assert.NoError(t, ms.Push(messages.Update{Key: "test", Status: messages.MessageStatusError, Details: "Test details"}))

	assert.Equal(t, "Still testing", msg.Message)
	assert.True(t, toc.After(msg.Started))
	assert.True(t, toc.Before(msg.Finished))
	assert.Equal(t, msg.Details, "Test details")
}
