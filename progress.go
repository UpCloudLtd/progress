package progress

import (
	"fmt"
	"time"

	"github.com/UpCloudLtd/progress/messages"
)

type OutputConfig messages.OutputConfig

// GetDefaultOutputConfig returns pointer to a new instance of default output configuration.
func GetDefaultOutputConfig() *OutputConfig {
	config := OutputConfig(messages.GetDefaultOutputConfig())
	return &config
}

type Progress struct {
	store      *messages.MessageStore
	renderer   *messages.MessageRenderer
	updateChan chan messages.Update
	errorChan  chan error
	stopChan   chan bool
	doneChan   chan bool
}

// NewProgress creates new Progress instance. Use nil config for default output configuration.
func NewProgress(config *OutputConfig) *Progress {
	if config == nil {
		config = GetDefaultOutputConfig()
	}

	return &Progress{
		store:     messages.NewMessageStore(),
		renderer:  messages.NewMessageRenderer(messages.OutputConfig(*config)),
		errorChan: make(chan error),
		doneChan:  make(chan bool),
	}
}

func (p Progress) run() {
	ticker := time.NewTicker(time.Millisecond * 95)
	defer ticker.Stop()

	for {
		select {
		case <-p.stopChan:
			p.store.Close()
			p.renderer.RenderMessageStore(p.store)
			p.doneChan <- true
			return
		case update := <-p.updateChan:
			p.errorChan <- p.store.Push(update)
		case <-ticker.C:
			p.renderer.RenderMessageStore(p.store)
		}
	}
}

// Start the progress logging in a new goroutine. Panics if called more than once.
func (p *Progress) Start() {
	if p.stopChan != nil {
		panic("can not start progress log more than once")
	}

	p.stopChan = make(chan bool)
	p.updateChan = make(chan messages.Update)
	go p.run()
}

// Push updates to the progress log. Errors if called before start or if called with an invalid update. Panics if called after Close.
func (p Progress) Push(update messages.Update) error {
	if p.updateChan == nil {
		return fmt.Errorf("can not push updates into progress log that has not been started")
	}

	p.updateChan <- update
	return <-p.errorChan
}

// Stop the goroutine handling progress logging and render the final progress state. Panics if called before start or more than once.
func (p Progress) Stop() {
	if p.stopChan == nil {
		panic("can not stop progress log that has not been started")
	}

	p.stopChan <- true
	// Block until stop is handled
	<-p.doneChan

	close(p.stopChan)
	close(p.updateChan)
}
