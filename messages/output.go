package messages

import (
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"strings"

	"github.com/UpCloudLtd/progress/terminal"
	"github.com/jedib0t/go-pretty/v6/text"
	"golang.org/x/term"
)

type RenderState int

const (
	RenderStateDone RenderState = -1
)

var whitespace = regexp.MustCompile(`\s`)

type OutputConfig struct {
	DefaultTextWidth            int
	DisableColors               bool
	ForceColors                 bool
	ShowStatusIndicator         bool
	StatusIndicatorMap          map[MessageStatus]string
	FallbackStatusIndicatorMap  map[MessageStatus]string
	StatusColorMap              map[MessageStatus]Color
	InProgressAnimation         []string
	FallbackInProgressAnimation []string
	UnknownColor                Color
	UnknownIndicator            string
	DetailsColor                Color
	ColorMessage                bool
	StopWatchcolor              Color
	ShowStopwatch               bool
	Target                      io.Writer
}

func GetDefaultOutputConfig() OutputConfig {
	return OutputConfig{
		DefaultTextWidth:    100,
		DisableColors:       false,
		ShowStatusIndicator: true,
		StatusIndicatorMap: map[MessageStatus]string{
			MessageStatusSuccess: "✓", // Check mark: U+2713
			MessageStatusWarning: "!",
			MessageStatusError:   "✗", // Ballot X: U+2717
			MessageStatusStarted: ">",
			MessageStatusPending: "#",
			MessageStatusSkipped: "-",
		},
		FallbackStatusIndicatorMap: map[MessageStatus]string{
			MessageStatusSuccess: "√", // Square root: U+221A
			MessageStatusError:   "X",
		},
		StatusColorMap: map[MessageStatus]Color{
			MessageStatusSuccess: text.FgGreen,
			MessageStatusWarning: text.FgYellow,
			MessageStatusError:   text.FgRed,
			MessageStatusStarted: text.FgBlue,
			MessageStatusPending: text.FgCyan,
			MessageStatusSkipped: text.FgMagenta,
		},
		InProgressAnimation:         []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		FallbackInProgressAnimation: []string{"/", "-", "\\", "|"},
		UnknownColor:                text.FgWhite,
		UnknownIndicator:            "?",
		DetailsColor:                text.FgHiBlack,
		ColorMessage:                false,
		StopWatchcolor:              text.FgHiBlack,
		ShowStopwatch:               true,
		Target:                      os.Stderr,
	}
}

func (cfg OutputConfig) shouldUseFallback() bool {
	if terminal.IsWindowsTerminal(cfg.Target) && !terminal.IsUnicodeSafeWindowsTermProgram() {
		return true
	}

	return false
}

func (cfg OutputConfig) getColor(c Color) Color {
	if cfg.ForceColors {
		return c
	}
	if cfg.DisableColors || os.Getenv("NO_COLOR") != "" {
		return noColor{}
	}
	return c
}

func (cfg OutputConfig) getStatusColor(status MessageStatus) Color {
	if color, ok := cfg.StatusColorMap[status]; ok {
		return cfg.getColor(color)
	}
	return cfg.getColor(cfg.UnknownColor)
}

func (cfg OutputConfig) getDetailsColor() Color {
	return cfg.getColor(cfg.DetailsColor)
}

func (cfg OutputConfig) getStopWatchcolor() Color {
	return cfg.getColor(cfg.StopWatchcolor)
}

func (cfg OutputConfig) getStatusIndicator(status MessageStatus) string {
	indicator := cfg.UnknownIndicator
	if preferred, ok := cfg.StatusIndicatorMap[status]; ok {
		indicator = preferred
	}

	if cfg.shouldUseFallback() {
		if fallback, ok := cfg.FallbackStatusIndicatorMap[status]; ok {
			indicator = fallback
		}
	}

	return indicator
}

func (cfg OutputConfig) getInProgressAnimationFrame(renderState RenderState) string {
	animation := cfg.InProgressAnimation
	if cfg.shouldUseFallback() {
		animation = cfg.FallbackInProgressAnimation
	}

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

// GetMaxWidth returns target terminals width or, if determining terminal dimensions failed, default value from OutputConfig.
func (cfg OutputConfig) GetMaxWidth() int {
	width, _ := cfg.getDimensions()
	return width
}

// GetMaxHeight returns target terminals height or, if determining terminal dimensions failed, zero.
func (cfg OutputConfig) GetMaxHeight() int {
	_, height := cfg.getDimensions()
	return height
}

func (cfg OutputConfig) formatDetails(msg *Message) string {
	wrapWidth := cfg.GetMaxWidth() - 2

	var details string
	// If details contains newline characters, assume that details are preformatted (e.g., stack trace, console output, ...)
	if strings.Contains(msg.Details, "\n") {
		details = text.WrapText(cfg.getDetailsColor().Sprint(msg.Details), wrapWidth)
	} else {
		details = text.WrapSoft(cfg.getDetailsColor().Sprint(msg.Details), wrapWidth)
	}

	if cfg.ShowStatusIndicator {
		return strings.ReplaceAll("\n"+details, "\n", "\n  ")
	}
	return "\n" + details
}

func (cfg OutputConfig) GetMessageText(msg *Message, renderState RenderState) string {
	isInteractive := cfg.GetMaxHeight() > 0

	status := ""
	color := cfg.getStatusColor(msg.Status)
	if cfg.ShowStatusIndicator {
		indicator := cfg.getStatusIndicator(msg.Status)
		if msg.Status.IsInProgress() && isInteractive {
			indicator = cfg.getInProgressAnimationFrame(renderState)
		}

		status = color.Sprintf("%s ", indicator)
	}

	elapsed := elapsedString(msg.ElapsedSeconds())
	if elapsed != "" {
		elapsed = cfg.getStopWatchcolor().Sprintf(" %s", elapsed)
	}

	lenFn := text.RuneWidthWithoutEscSequences
	message := msg.Message
	if isInteractive && msg.ProgressMessage != "" {
		message += " " + msg.ProgressMessage
	}
	maxMessageWidth := cfg.GetMaxWidth() - lenFn(status) - lenFn(elapsed)
	// Some terminals initially return 0 width, skip rendering message in that case.
	if maxMessageWidth < 0 {
		return ""
	}
	message = whitespace.ReplaceAllString(message, " ")
	if len(message) > maxMessageWidth {
		message = fmt.Sprintf("%s…", message[:maxMessageWidth-1])
	} else {
		message = text.Pad(message, maxMessageWidth, ' ')
	}
	if cfg.ColorMessage {
		message = color.Sprint(message)
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
		// In-progress messages always span whole terminal width. Thus, when terminal width is decreased, every in-progress message is wrapped to the next line(s).
		currentHeight *= int(math.Ceil(float64(mr.inProgressWidth) / float64(currentWidth)))
	}

	// Move to first column, move cursor to beginning of in-progress area and erase all lines on the way.
	return "\r" + strings.Repeat(text.CursorUp.Sprint()+text.EraseLine.Sprint(), currentHeight)
}
