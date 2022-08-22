//nolint:funlen // Does not make sense to split example usage into multiple functions
package main

import (
	"fmt"
	"time"

	"github.com/UpCloudLtd/progress"
	"github.com/UpCloudLtd/progress/messages"
)

const loremIpsum = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

func main() {
	cfg := progress.GetDefaultOutputConfig()

	taskLog := progress.NewProgress(cfg)
	taskLog.Start()
	defer taskLog.Stop()

	_ = taskLog.Push(messages.Update{
		Key:     "first-example",
		Message: "Progress is a library for communicating CLI app progress to the user",
		Status:  messages.MessageStatusStarted,
	})

	time.Sleep(time.Millisecond * 1500)

	_ = taskLog.Push(messages.Update{
		Key:     "parallel-example",
		Message: "There can be multiple active progress messages at once",
		Status:  messages.MessageStatusStarted,
	})

	_ = taskLog.Push(messages.Update{
		Key:     "progress-example",
		Message: "Progress message can include part that is only outputted to TTY terminals",
		Status:  messages.MessageStatusStarted,
	})

	time.Sleep(time.Millisecond * 300)
	for i := 1; i < 10; i++ {
		_ = taskLog.Push(messages.Update{
			Key:             "progress-example",
			ProgressMessage: fmt.Sprintf("(%d%%)", i*10),
		})
		time.Sleep(time.Millisecond * 300)
	}

	_ = taskLog.Push(messages.Update{
		Key:    "progress-example",
		Status: messages.MessageStatusSuccess,
	})

	_ = taskLog.Push(messages.Update{
		Key:     "first-example",
		Message: "Progress messages can be updated while they are in pending or started state",
		Status:  messages.MessageStatusStarted,
	})

	time.Sleep(time.Millisecond * 1500)

	_ = taskLog.Push(messages.Update{
		Key:     "parallel-example",
		Status:  messages.MessageStatusError,
		Details: "Error: Message details can be used, for example, to communicate error messages to the user.",
	})

	time.Sleep(time.Millisecond * 1500)

	_ = taskLog.Push(messages.Update{
		Key:     "unknown-example",
		Status:  messages.MessageStatusStarted,
		Message: "If message has started status when log is closed, its status is set to unknown",
	})

	_ = taskLog.Push(messages.Update{
		Key:     "pending-example",
		Status:  messages.MessageStatusPending,
		Message: "Pending tasks are not written to output. If message has pending status when log is closed, its status is set to skipped",
	})

	time.Sleep(time.Millisecond * 1500)

	_ = taskLog.Push(messages.Update{
		Key:     "first-example",
		Message: "Progress messages are, by default, written to stderr",
		Status:  messages.MessageStatusSuccess,
	})

	time.Sleep(time.Millisecond * 1500)

	_ = taskLog.Push(messages.Update{
		Key:     "long",
		Message: "Long messages are truncated - " + loremIpsum,
		Details: loremIpsum,
		Status:  messages.MessageStatusStarted,
	})

	time.Sleep(time.Second * 3)

	_ = taskLog.Push(messages.Update{
		Key:    "long",
		Status: messages.MessageStatusWarning,
	})

	time.Sleep(time.Second * 3)
}
