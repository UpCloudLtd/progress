//nolint:funlen // Does not make sense to split example usage into multiple functions
package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/UpCloudLtd/progress"
	"github.com/UpCloudLtd/progress/messages"
)

const loremIpsum = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

const (
	firstKey    = "first-example"
	parallelKey = "parallel-example"
	progressKey = "progress-example"
	unknownKey  = "unknown-example"
	pendingKey  = "pending-example"
	longKey     = "long-example"
)

func main() {
	taskLog := progress.NewProgress(nil)
	taskLog.Start()
	defer taskLog.Stop()

	parallelCount := 1
	if len(os.Args) > 1 {
		n, err := strconv.Atoi(os.Args[1])
		if err == nil && n > 0 {
			parallelCount = n
		}
	}

	_ = taskLog.Push(messages.Update{
		Key:     firstKey,
		Message: "Progress is a library for communicating CLI app progress to the user",
		Status:  messages.MessageStatusStarted,
	})

	time.Sleep(time.Millisecond * 1500)

	for i := 0; i < parallelCount; i++ {
		_ = taskLog.Push(messages.Update{
			Key:     fmt.Sprintf("%s-%02d", parallelKey, i),
			Message: "There can be multiple active progress messages at once",
			Status:  messages.MessageStatusStarted,
		})
	}

	_ = taskLog.Push(messages.Update{
		Key:     progressKey,
		Message: "Progress message can include part that is only outputted to TTY terminals",
		Status:  messages.MessageStatusStarted,
	})

	time.Sleep(time.Millisecond * 300)
	for i := 1; i < 10; i++ {
		_ = taskLog.Push(messages.Update{
			Key:             progressKey,
			ProgressMessage: fmt.Sprintf("(%d%%)", i*10),
		})
		time.Sleep(time.Millisecond * 300)
	}

	_ = taskLog.Push(messages.Update{
		Key:    progressKey,
		Status: messages.MessageStatusSuccess,
	})

	_ = taskLog.Push(messages.Update{
		Key:     firstKey,
		Message: "Progress messages can be updated while they are in pending or started state",
		Status:  messages.MessageStatusStarted,
	})

	time.Sleep(time.Millisecond * 1500)

	for i := 0; i < parallelCount; i++ {
		status := messages.MessageStatusSuccess
		details := ""
		if i == 0 {
			status = messages.MessageStatusError
			details = "Message details can be used, for example, to communicate error messages to the user."
		}

		_ = taskLog.Push(messages.Update{
			Key:     fmt.Sprintf("%s-%02d", parallelKey, i),
			Status:  status,
			Details: details,
		})
	}

	time.Sleep(time.Millisecond * 1500)

	_ = taskLog.Push(messages.Update{
		Key:     unknownKey,
		Status:  messages.MessageStatusStarted,
		Message: "If message has started status when log is closed, its status is set to unknown",
	})

	_ = taskLog.Push(messages.Update{
		Key:     pendingKey,
		Status:  messages.MessageStatusPending,
		Message: "Pending tasks are not written to output. If message has pending status when log is closed, its status is set to skipped",
	})

	time.Sleep(time.Millisecond * 1500)

	_ = taskLog.Push(messages.Update{
		Key:     firstKey,
		Message: "Progress messages are, by default, written to stderr",
		Status:  messages.MessageStatusSuccess,
	})

	time.Sleep(time.Millisecond * 1500)

	_ = taskLog.Push(messages.Update{
		Key:     longKey,
		Message: "Long messages are truncated - " + loremIpsum,
		Details: loremIpsum,
		Status:  messages.MessageStatusStarted,
	})

	time.Sleep(time.Second * 3)

	_ = taskLog.Push(messages.Update{
		Key:    longKey,
		Status: messages.MessageStatusWarning,
	})

	time.Sleep(time.Second * 3)
}
