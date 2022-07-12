package progress

import (
	"bytes"
	"regexp"
	"runtime"
	"testing"
	"time"

	"github.com/kangasta/progress/messages"
	"github.com/stretchr/testify/assert"
)

func TestProgress_Push_ErrorChannel(t *testing.T) {
	taskLog := NewProgress(DefaultOutputConfig)
	taskLog.Start()
	defer taskLog.Stop()

	err := taskLog.Push(messages.Update{Message: "No status"})
	assert.EqualError(t, err, `can not push message with invalid status ""`)

	err = taskLog.Push(messages.Update{Message: "Valid update", Status: messages.MessageStatusStarted})
	assert.NoError(t, err)
}

func TestProgress_Start_PanicsIfCalledTwice(t *testing.T) {
	taskLog := NewProgress(DefaultOutputConfig)
	taskLog.Start()
	defer taskLog.Stop()
	assert.PanicsWithValue(t, "can not start progress log more than once", taskLog.Start)
}

func TestProgress_Push_PanicsIfCalledAfterStop(t *testing.T) {
	taskLog := NewProgress(DefaultOutputConfig)
	taskLog.Start()
	taskLog.Stop()
	assert.Panics(t, func() { taskLog.Push(messages.Update{}) }) //nolint:errcheck
}

func TestProgress_Stop_PanicsIfCalledBeforeStart(t *testing.T) {
	taskLog := NewProgress(DefaultOutputConfig)
	assert.PanicsWithValue(t, "can not stop progress log that has not been started", taskLog.Stop)
}

func TestProgress_Stop_PanicsIfCalledTwice(t *testing.T) {
	taskLog := NewProgress(DefaultOutputConfig)
	taskLog.Start()
	taskLog.Stop()
	assert.Panics(t, taskLog.Stop)
}

func TestProgress_Output(t *testing.T) {
	cfg := DefaultOutputConfig
	buf := bytes.NewBuffer(nil)
	cfg.Target = buf

	taskLog := NewProgress(cfg)
	taskLog.Start()

	err := taskLog.Push(messages.Update{Message: "Test update", Status: messages.MessageStatusStarted})
	assert.NoError(t, err)

	time.Sleep(time.Millisecond * 100) // Wait for the first render

	err = taskLog.Push(messages.Update{Message: "Test update", Status: messages.MessageStatusSuccess})
	assert.NoError(t, err)

	taskLog.Stop()

	output := buf.String()

	expected := "\x1b[34m> \x1b[0mTest update                                                                                       \n\x1b[32m✓ \x1b[0mTest update                                                                                       \n"
	if runtime.GOOS == "windows" {
		re := regexp.MustCompile("\x1b\\[[0-9]+m")
		expected = re.ReplaceAllString(expected, "")
	}
	assert.Equal(t, expected, output)
}
