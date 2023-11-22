package messages_test

import (
	"bytes"
	"runtime"
	"testing"
	"time"

	"github.com/UpCloudLtd/progress/messages"
	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/stretchr/testify/assert"
)

const loremIpsum = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

func TestMessageRenderer_RenderMessageStore(t *testing.T) {
	t.Parallel()
	defaultConfig := messages.GetDefaultOutputConfig()

	disableColors := messages.GetDefaultOutputConfig()
	disableColors.DisableColors = true

	noIndicatorColorMessage := messages.GetDefaultOutputConfig()
	noIndicatorColorMessage.ColorMessage = true
	noIndicatorColorMessage.ShowStatusIndicator = false

	for _, test := range []struct {
		name          string
		config        messages.OutputConfig
		skipOnWindows bool
	}{
		{
			name:          "Default configuration",
			config:        defaultConfig,
			skipOnWindows: true,
		},
		{
			name:          "Disable colors",
			config:        disableColors,
			skipOnWindows: false,
		},
		{
			name:          "No indicator colored message",
			config:        noIndicatorColorMessage,
			skipOnWindows: true,
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if test.skipOnWindows && runtime.GOOS == "windows" {
				t.Skip("Skipping snapshot test on Windows, as output will not include ANSI codes included in the snapshot.")
			}

			cfg := test.config
			buf := bytes.NewBuffer(nil)
			cfg.Target = buf

			renderer := messages.NewMessageRenderer(cfg)
			store := messages.NewMessageStore()

			err := store.Push(messages.Update{
				Message: "Test pending (0s)",
				Status:  messages.MessageStatusPending,
			})
			assert.NoError(t, err)

			time.Sleep(time.Microsecond * 25) // Ensure time difference on Windows
			err = store.Push(messages.Update{
				Message: "Test started (0s)",
				Status:  messages.MessageStatusStarted,
			})
			assert.NoError(t, err)

			time.Sleep(time.Microsecond * 25) // Ensure time difference on Windows
			err = store.Push(messages.Update{
				Message: "Test skipped (0s, long message) - " + loremIpsum,
				Status:  messages.MessageStatusSkipped,
			})
			assert.NoError(t, err)

			time.Sleep(time.Microsecond * 25) // Ensure time difference on Windows
			err = store.Add(messages.Message{
				Message:  "Test success (100s)",
				Status:   messages.MessageStatusSuccess,
				Started:  time.Now().Add(time.Second * -100),
				Finished: time.Now(),
			})
			assert.NoError(t, err)

			err = store.Add(messages.Message{
				Message:  "Test error, 10 % (1000s, % char in message)",
				Status:   messages.MessageStatusError,
				Details:  "Error: Short dummy error message",
				Started:  time.Now().Add(time.Second * -1000),
				Finished: time.Now(),
			})
			assert.NoError(t, err)

			err = store.Add(messages.Message{
				Message:  "Test\tinvalid\nmessage\twith\ntabs\tand\nnewlines (5s, \\n and \\t chars in message)",
				Status:   messages.MessageStatusError,
				Started:  time.Now().Add(time.Second * -5),
				Finished: time.Now(),
			})
			assert.NoError(t, err)

			err = store.Add(messages.Message{
				Message:  "Test warning (10s, long message) - " + loremIpsum,
				Status:   messages.MessageStatusWarning,
				Details:  "Error: Long dummy error message - " + loremIpsum,
				Started:  time.Now().Add(time.Second * -1000),
				Finished: time.Now(),
			})
			assert.NoError(t, err)

			renderer.RenderMessageStore(store)

			err = store.Push(messages.Update{
				Key:     "progress-example",
				Message: "Test ProgressMessage is not visible in non-TTY output",
				Status:  messages.MessageStatusStarted,
			})
			assert.NoError(t, err)
			renderer.RenderMessageStore(store)

			err = store.Push(messages.Update{
				Key:             "progress-example",
				ProgressMessage: "(50%)",
			})
			assert.NoError(t, err)
			renderer.RenderMessageStore(store)

			err = store.Push(messages.Update{
				Key:    "progress-example",
				Status: messages.MessageStatusSuccess,
			})
			assert.NoError(t, err)
			store.Close()
			renderer.RenderMessageStore(store)

			output := buf.String()
			cupaloy.SnapshotT(t, output)
		})
	}
}
