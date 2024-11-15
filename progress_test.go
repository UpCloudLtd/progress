package progress_test

import (
	"bytes"
	"regexp"
	"runtime"
	"testing"
	"time"

	"github.com/UpCloudLtd/progress"
	"github.com/UpCloudLtd/progress/messages"
	"github.com/stretchr/testify/assert"
)

func removeColorsOnWindows(expected string) string {
	if runtime.GOOS == "windows" {
		re := regexp.MustCompile("\x1b\\[[0-9]+m")
		return re.ReplaceAllString(expected, "")
	}
	return expected
}

func TestProgress_Push_ErrorChannel(t *testing.T) {
	t.Parallel()
	taskLog := progress.NewProgress(nil)
	taskLog.Start()
	defer taskLog.Stop()

	err := taskLog.Push(messages.Update{Message: "No status"})
	assert.EqualError(t, err, `can not push message with invalid status ""`)

	err = taskLog.Push(messages.Update{Message: "Valid update", Status: messages.MessageStatusStarted})
	assert.NoError(t, err)
}

func TestProgress_Start_PanicsIfCalledTwice(t *testing.T) {
	t.Parallel()
	taskLog := progress.NewProgress(nil)
	taskLog.Start()
	defer taskLog.Stop()
	assert.PanicsWithValue(t, "can not start progress log more than once", taskLog.Start)
}

func TestProgress_Push_ErrorsIfCalledBeforeStart(t *testing.T) {
	t.Parallel()
	taskLog := progress.NewProgress(nil)
	err := taskLog.Push(messages.Update{})
	assert.EqualError(t, err, "can not push updates into progress log that has not been started")
}

func TestProgress_Push_PanicsIfCalledAfterStop(t *testing.T) {
	t.Parallel()
	taskLog := progress.NewProgress(nil)
	taskLog.Start()
	taskLog.Stop()
	assert.Panics(t, func() { _ = taskLog.Push(messages.Update{}) })
}

func TestProgress_Stop_PanicsIfCalledBeforeStart(t *testing.T) {
	t.Parallel()
	taskLog := progress.NewProgress(nil)
	assert.PanicsWithValue(t, "can not stop progress log that has not been started", taskLog.Stop)
}

func TestProgress_Stop_PanicsIfCalledTwice(t *testing.T) {
	t.Parallel()
	taskLog := progress.NewProgress(nil)
	taskLog.Start()
	taskLog.Stop()
	assert.Panics(t, taskLog.Stop)
}

func TestProgress_Output(t *testing.T) {
	t.Parallel()
	cfg := progress.GetDefaultOutputConfig()
	buf := bytes.NewBuffer(nil)
	cfg.Target = buf

	taskLog := progress.NewProgress(cfg)
	taskLog.Start()

	err := taskLog.Push(messages.Update{Message: "Test update", Status: messages.MessageStatusStarted})
	assert.NoError(t, err)

	time.Sleep(time.Millisecond * 150) // Wait for the first render

	err = taskLog.Push(messages.Update{Message: "Test update", Status: messages.MessageStatusSuccess})
	assert.NoError(t, err)

	taskLog.Stop()

	output := buf.String()

	expected := removeColorsOnWindows("\x1b[34m> \x1b[0mTest update                                                                                       \n\x1b[32m✓ \x1b[0mTest update                                                                                       \n")
	assert.Equal(t, expected, output)
}

func TestProgress_NoProgressMessage(t *testing.T) {
	t.Parallel()
	cfg := progress.GetDefaultOutputConfig()
	buf := bytes.NewBuffer(nil)
	cfg.Target = buf

	taskLog := progress.NewProgress(cfg)
	taskLog.Start()
	defer taskLog.Stop()

	err := taskLog.Push(messages.Update{Message: "Test update", Status: messages.MessageStatusStarted})
	assert.NoError(t, err)
	err = taskLog.Push(messages.Update{Message: "Test update", ProgressMessage: "(50 %)"})
	assert.NoError(t, err)

	time.Sleep(time.Millisecond * 150) // Wait for the first render
	output := buf.String()

	expected := removeColorsOnWindows("\x1b[34m> \x1b[0mTest update                                                                                       \n")
	assert.Equal(t, expected, output)
}

func TestProgress_ClosesInProgressMessagesOnStop(t *testing.T) {
	t.Parallel()
	cfg := progress.GetDefaultOutputConfig()
	buf := bytes.NewBuffer(nil)
	cfg.Target = buf

	taskLog := progress.NewProgress(cfg)
	taskLog.Start()

	err := taskLog.Push(messages.Update{Message: "Test pending 1", Status: messages.MessageStatusPending})
	assert.NoError(t, err)
	time.Sleep(time.Microsecond * 25) // Ensure time difference on Windows
	err = taskLog.Push(messages.Update{Message: "Test started", Status: messages.MessageStatusStarted})
	assert.NoError(t, err)
	time.Sleep(time.Microsecond * 25) // Ensure time difference on Windows
	err = taskLog.Push(messages.Update{Message: "Test pending 2", Status: messages.MessageStatusPending})
	assert.NoError(t, err)

	time.Sleep(time.Millisecond * 150) // Wait for the first render

	taskLog.Stop()

	output := buf.String()

	expected := removeColorsOnWindows("\x1b[34m> \x1b[0mTest started                                                                                      \n\x1b[35m- \x1b[0mTest pending 1                                                                                    \n\x1b[35m- \x1b[0mTest pending 2                                                                                    \n\x1b[37m? \x1b[0mTest started                                                                                      \n")
	assert.Equal(t, expected, output)
}
