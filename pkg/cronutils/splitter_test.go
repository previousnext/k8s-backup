package cronutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitter(t *testing.T) {
	wants := []string{
		"0 0 * * *",
		"5 0 * * *",
		"10 0 * * *",
		"15 0 * * *",
		"20 0 * * *",
		"25 0 * * *",
		"30 0 * * *",
		"35 0 * * *",
		"40 0 * * *",
		"45 0 * * *",
		"50 0 * * *",
		"55 0 * * *",
		"0 0 * * *",
		"5 0 * * *",
	}

	splitter := NewSplitter(5)

	for _, want := range wants {
		assert.Equal(t, want, splitter.Increment())
	}
}
