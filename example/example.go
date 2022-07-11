package main

import (
	"time"

	"github.com/kangasta/progress"
	"github.com/kangasta/progress/messages"
)

const loremIpsum = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

func main() {
	cfg := progress.DefaultOutputConfig

	taskLog := progress.NewProgress(cfg)
	taskLog.Start()
	defer taskLog.Stop()

	taskLog.Push(messages.Update{ //nolint:errcheck
		Key:     "first-example",
		Message: "Example message",
		Status:  messages.MessageStatusStarted,
	})

	time.Sleep(time.Millisecond * 1500)

	taskLog.Push(messages.Update{ //nolint:errcheck
		Key:     "parallel-example",
		Message: "Example message running in parallel",
		Status:  messages.MessageStatusStarted,
	})

	time.Sleep(time.Second * 3)

	taskLog.Push(messages.Update{ //nolint:errcheck
		Key:     "first-example",
		Message: "Updated example message",
		Status:  messages.MessageStatusStarted,
	})

	time.Sleep(time.Millisecond * 1500)

	taskLog.Push(messages.Update{ //nolint:errcheck
		Key:     "parallel-example",
		Status:  messages.MessageStatusError,
		Details: "Error: Message details can be used, for example, to communicate error messages to the user.",
	})

	time.Sleep(time.Millisecond * 1500)

	taskLog.Push(messages.Update{ //nolint:errcheck
		Key:     "unknown-example",
		Status:  messages.MessageStatusStarted,
		Message: "If message has started status when log is closed, its status is set to unknown",
	})

	time.Sleep(time.Millisecond * 1500)

	taskLog.Push(messages.Update{ //nolint:errcheck
		Key:     "first-example",
		Message: "Finished example message",
		Status:  messages.MessageStatusSuccess,
	})

	time.Sleep(time.Millisecond * 1500)

	taskLog.Push(messages.Update{ //nolint:errcheck
		Key:     "long",
		Message: "Long messages are truncated - " + loremIpsum,
		Details: loremIpsum,
		Status:  messages.MessageStatusStarted,
	})

	time.Sleep(time.Second * 3)

	taskLog.Push(messages.Update{ //nolint:errcheck
		Key:    "long",
		Status: messages.MessageStatusWarning,
	})

	time.Sleep(time.Second * 3)
}
