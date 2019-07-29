package cronutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitter(t *testing.T) {
	wants := []string{
		"0 * * * *",
		"5 * * * *",
		"10 * * * *",
		"15 * * * *",
		"20 * * * *",
		"25 * * * *",
		"30 * * * *",
		"35 * * * *",
		"40 * * * *",
		"45 * * * *",
		"50 * * * *",
		"55 * * * *",
		"60 * * * *",
		"0 * * * *",
		"5 * * * *",
	}

	splitter := NewSplitter(5)

	for _, want := range wants {
		assert.Equal(t, want, splitter.Increment())
	}
}
