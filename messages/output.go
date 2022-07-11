package messages

import (
	"fmt"
	"io"
	"math"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
	"golang.org/x/term"
)

type RenderState int

const (
	RenderStateDone RenderState = -1
)

type OutputConfig struct {
	DefaultTextWidth    int
	ShowStatusIndicator bool
	StatusIndicatorMap  map[MessageStatus]string
	StatusColorMap      map[MessageStatus]text.Color
	InProgressAnimation []string
	UnknownColor        text.Color
	UnknownIndicator    string
	DetailsColor        text.Color
	ColorMessage        bool
	StopWatchcolor      text.Color
	ShowStopwatch       bool
	Target              io.Writer
}

var DefaultOutputConfig = OutputConfig{
	DefaultTextWidth:    100,
	ShowStatusIndicator: true,
	StatusIndicatorMap: map[MessageStatus]string{
		MessageStatusSuccess: "✓",
		MessageStatusWarning: "!",
		MessageStatusError:   "✗",
		MessageStatusStarted: ">",
		MessageStatusPending: "#",
		MessageStatusSkipped: "-",
	},
	StatusColorMap: map[MessageStatus]text.Color{
		MessageStatusSuccess: text.FgGreen,
		MessageStatusWarning: text.FgYellow,
		MessageStatusError:   text.FgRed,
		MessageStatusStarted: text.FgBlue,
		MessageStatusPending: text.FgCyan,
		MessageStatusSkipped: text.FgMagenta,
	},
	InProgressAnimation: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	UnknownColor:        text.FgWhite,
	UnknownIndicator:    "?",
	DetailsColor:        text.FgHiBlack,
	ColorMessage:        false,
	StopWatchcolor:      text.FgHiBlack,
	ShowStopwatch:       true,
	Target:              os.Stderr,
}

func (cfg OutputConfig) getStatusColor(status MessageStatus) text.Color {
	if color, ok := cfg.StatusColorMap[status]; ok {
		return color
	}
	return cfg.UnknownColor
}

func (cfg OutputConfig) getStatusIndicator(status MessageStatus) string {
	if indicator, ok := cfg.StatusIndicatorMap[status]; ok {
		return indicator
	}
	return cfg.UnknownIndicator
}

func (cfg OutputConfig) getInProgressAnimationFrame(renderState RenderState) string {
	animation := cfg.InProgressAnimation
	i := int(renderState) % len(animation)

	return animation[i]

}

func elapsedString(elapsedSeconds float64) string {
	if elapsedSeconds < 1 {
		return ""
	}

	if elapsedSeconds >= 999 {
		return "> 999 s"
	}

	return fmt.Sprintf("%3d s", int(elapsedSeconds))
}

func (cfg OutputConfig) getDimensions() (int, int) {
	file, ok := cfg.Target.(*os.File)
	if !ok {
		return cfg.DefaultTextWidth, 0
	}

	width, height, err := term.GetSize(int(file.Fd()))
	if err != nil {
		return cfg.DefaultTextWidth, 0
	}
	return width, height
}

// GetMaxWidth returns target terminals width or, if determining terminal dimensions failed, default value from OutputConfig
func (cfg OutputConfig) GetMaxWidth() int {
	width, _ := cfg.getDimensions()
	return width
}

// GetMaxHeight returns target terminals height or, if determining terminal dimensions failed, zero
func (cfg OutputConfig) GetMaxHeight() int {
	_, height := cfg.getDimensions()
	return height
}

func (cfg OutputConfig) formatDetails(msg *Message) string {
	wrapWidth := cfg.GetMaxWidth() - 2
	details := text.WrapSoft(cfg.DetailsColor.Sprint(msg.Details), wrapWidth)

	if cfg.ShowStatusIndicator {
		return strings.ReplaceAll("\n"+details, "\n", "\n  ")
	}
	return "\n" + details
}

func (cfg OutputConfig) GetMessageText(msg *Message, renderState RenderState) string {
	termWidth := cfg.GetMaxWidth()

	status := ""
	color := cfg.getStatusColor(msg.Status)
	if cfg.ShowStatusIndicator {
		indicator := cfg.getStatusIndicator(msg.Status)
		if msg.Status.IsInProgress() && cfg.GetMaxHeight() > 0 {
			indicator = cfg.getInProgressAnimationFrame(renderState)
		}

		status = color.Sprintf("%s ", indicator)
	}

	elapsed := elapsedString(msg.ElapsedSeconds())
	if elapsed != "" {
		elapsed = cfg.StopWatchcolor.Sprintf(" %s", elapsed)
	}

	lenFn := text.RuneWidthWithoutEscSequences
	message := msg.Message
	maxMessageWidth := termWidth - lenFn(status) - lenFn(elapsed)
	if len(message) > maxMessageWidth {
		message = fmt.Sprintf("%s…", message[:maxMessageWidth-1])
	} else {
		message = text.Pad(message, maxMessageWidth, ' ')
	}
	if cfg.ColorMessage {
		message = color.Sprintf(message)
	}

	details := ""
	if msg.Details != "" && msg.Status.IsFinished() {
		details = cfg.formatDetails(msg)
	}

	return fmt.Sprintf("%s%s%s%s\n", status, message, elapsed, details)
}

type MessageRenderer struct {
	finishedMap      map[string]bool
	config           OutputConfig
	renderState      RenderState
	finishedIndex    int
	inProgressWidth  int
	inProgressHeight int
}

func NewMessageRenderer(config OutputConfig) *MessageRenderer {
	return &MessageRenderer{
		finishedMap: make(map[string]bool),
		config:      config,
	}
}

func (mr MessageRenderer) write(args ...any) {
	fmt.Fprint(mr.config.Target, args...)
}

func (mr MessageRenderer) prepareMessage(msg *Message, keyPostfix ...string) string {
	key := fmt.Sprint(msg.Key, keyPostfix)

	if mr.finishedMap[key] {
		return ""
	}

	mr.finishedMap[key] = true
	return mr.config.GetMessageText(msg, mr.renderState)
}

func (mr *MessageRenderer) RenderMessageStore(ms *MessageStore) {
	text := mr.moveToInProgressStartText()

	// Render finished messages
	finished := ms.ListFinished()[mr.finishedIndex:]
	for _, msg := range finished {
		if msg.Status.IsFinished() {
			text += mr.config.GetMessageText(msg, mr.renderState)
		}
	}
	mr.finishedIndex += len(finished)

	// Render in-progress messages
	inProgress := ms.ListInProgress()
	count := 0
	for _, msg := range inProgress {
		if !msg.Status.IsInProgress() {
			continue
		}
		maxHeight := mr.config.GetMaxHeight()
		if maxHeight == 0 {
			// Print message when it is started and when its message changes to new value
			text += mr.prepareMessage(msg, msg.Message, "started")
		} else {
			if count >= maxHeight {
				break
			}
			text += mr.config.GetMessageText(msg, mr.renderState)
			count++
		}
	}
	if text != "" {
		mr.write(text)
	}

	mr.inProgressHeight = count
	mr.inProgressWidth = mr.config.GetMaxWidth()
	mr.renderState++
}

func (mr *MessageRenderer) moveToInProgressStartText() string {
	if mr.inProgressHeight == 0 {
		return ""
	}

	// Check if terminal width has decreased and get height to clear
	currentWidth := mr.config.GetMaxWidth()
	currentHeight := mr.inProgressHeight
	if currentWidth < mr.inProgressWidth {
		contentLength := mr.inProgressWidth * mr.inProgressWidth
		currentHeight = int(math.Ceil(float64(contentLength) / float64(currentWidth)))
	}

	// Move to first column and move cursor to beginning of in-progress area
	return "\r" + strings.Repeat(text.CursorUp.Sprint(), currentHeight)
}
