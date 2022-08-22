package messages

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageRenderer_moveToInProgressStartText(t *testing.T) {
	t.Parallel()
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
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
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
