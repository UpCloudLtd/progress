package messages

import (
	"bytes"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/stretchr/testify/assert"
)

const loremIpsum = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

func TestMessageRenderer_RenderMessageStore(t *testing.T) {
	defaultConfig := GetDefaultOutputConfig()

	disableColors := GetDefaultOutputConfig()
	disableColors.DisableColors = true

	noIndicatorColorMessage := GetDefaultOutputConfig()
	noIndicatorColorMessage.ColorMessage = true
	noIndicatorColorMessage.ShowStatusIndicator = false

	for _, test := range []struct {
		name          string
		config        OutputConfig
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
		t.Run(test.name, func(t *testing.T) {
			if test.skipOnWindows && runtime.GOOS == "windows" {
				t.Skip("Skipping snapshot test on Windows, as output will not include ANSI codes included in the snapshot.")
			}

			cfg := test.config
			buf := bytes.NewBuffer(nil)
			cfg.Target = buf

			renderer := NewMessageRenderer(cfg)
			store := NewMessageStore()

			err := store.Push(Update{
				Message: "Test pending (0s)",
				Status:  MessageStatusPending,
			})
			assert.NoError(t, err)

			time.Sleep(time.Microsecond * 25) // Ensure time difference on Windows
			err = store.Push(Update{
				Message: "Test started (0s)",
				Status:  MessageStatusStarted,
			})
			assert.NoError(t, err)

			time.Sleep(time.Microsecond * 25) // Ensure time difference on Windows
			err = store.Push(Update{
				Message: "Test skipped (0s, long message) - " + loremIpsum,
				Status:  MessageStatusSkipped,
			})
			assert.NoError(t, err)

			time.Sleep(time.Microsecond * 25) // Ensure time difference on Windows
			err = store.Add(Message{
				Message:  "Test success (100s)",
				Status:   MessageStatusSuccess,
				Started:  time.Now().Add(time.Second * -100),
				Finished: time.Now(),
			})
			assert.NoError(t, err)

			err = store.Add(Message{
				Message:  "Test error (1000s)",
				Status:   MessageStatusError,
				Details:  "Error: Short dummy error message",
				Started:  time.Now().Add(time.Second * -1000),
				Finished: time.Now(),
			})
			assert.NoError(t, err)

			err = store.Add(Message{
				Message:  "Test warning (10s, long message) - " + loremIpsum,
				Status:   MessageStatusWarning,
				Details:  "Error: Long dummy error message - " + loremIpsum,
				Started:  time.Now().Add(time.Second * -1000),
				Finished: time.Now(),
			})
			assert.NoError(t, err)

			renderer.RenderMessageStore(store)

			err = store.Push(Update{
				Key:     "progress-example",
				Message: "Test ProgressMessage is not visible in non-TTY output",
				Status:  MessageStatusStarted,
			})
			assert.NoError(t, err)
			renderer.RenderMessageStore(store)

			err = store.Push(Update{
				Key:             "progress-example",
				ProgressMessage: "(50%)",
			})
			assert.NoError(t, err)
			renderer.RenderMessageStore(store)

			err = store.Push(Update{
				Key:    "progress-example",
				Status: MessageStatusSuccess,
			})
			assert.NoError(t, err)
			store.Close()
			renderer.RenderMessageStore(store)

			output := buf.String()
			cupaloy.SnapshotT(t, output)
		})
	}
}

func TestMessageRenderer_moveToInProgressStartText(t *testing.T) {
	for _, test := range []struct {
		name             string
		inProgressWidth  int
		inProgressHeight int
		terminalWidth    int
		moveUp           int
	}{
		{
			name:             "Terminal width increases",
			inProgressWidth:  30,
			inProgressHeight: 3,
			terminalWidth:    31,
			moveUp:           3,
		},
		{
			name:             "Terminal width stays the same",
			inProgressWidth:  30,
			inProgressHeight: 3,
			terminalWidth:    30,
			moveUp:           3,
		},
		{
			name:             "Terminal width decreases from 30 to 29",
			inProgressWidth:  30,
			inProgressHeight: 3,
			terminalWidth:    29,
			moveUp:           6,
		},
		{
			name:             "Terminal width decreases from 30 to 10",
			inProgressWidth:  30,
			inProgressHeight: 3,
			terminalWidth:    10,
			moveUp:           9,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			cfg := GetDefaultOutputConfig()
			cfg.DefaultTextWidth = test.terminalWidth

			r := NewMessageRenderer(cfg)
			r.inProgressHeight = test.inProgressHeight
			r.inProgressWidth = test.inProgressWidth

			// Cursor up = \x1b[A
			// Erase from cursor to EOL = \x1b[K
			moveUpCount := strings.Count(r.moveToInProgressStartText(), "\x1b[A\x1b[K")
			assert.Equal(t, test.moveUp, moveUpCount)
		})
	}
}
