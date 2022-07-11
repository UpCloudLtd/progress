package progress

import (
	"time"

	"github.com/kangasta/progress/messages"
)

type OutputConfig messages.OutputConfig

var DefaultOutputConfig OutputConfig = OutputConfig(messages.DefaultOutputConfig)

type Progress struct {
	store      *messages.MessageStore
	renderer   *messages.MessageRenderer
	updateChan chan messages.Update
	stopChan   chan bool
	doneChan   chan bool
}

// NewProgress creates new Progress instance.
func NewProgress(config OutputConfig) *Progress {
	return &Progress{
		store:      messages.NewMessageStore(),
		renderer:   messages.NewMessageRenderer(messages.OutputConfig(config)),
		updateChan: make(chan messages.Update),
		stopChan:   make(chan bool),
		doneChan:   make(chan bool),
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
			break
		case update := <-p.updateChan:
			p.store.Push(update)
		case <-ticker.C:
			p.renderer.RenderMessageStore(p.store)
		}
	}
}

// Start the progress logging in a new goroutine.
func (p Progress) Start() {
	go p.run()
}

// Push updates to the progress log.
func (p Progress) Push(update messages.Update) {
	p.updateChan <- update
}

// Stop the goroutine handling progress logging and render the final progress state.
func (p Progress) Stop() {
	p.stopChan <- true
	// Block until stop is handled
	<-p.doneChan
}
