package cronutils

import "fmt"

// Splitter which will return incremented cron schedules.
type Splitter struct {
	diff    int
	current int
}

// NewSplitter for created incemented cron schedules.
func NewSplitter(diff int) *Splitter {
	return &Splitter{
		diff: diff,
	}
}

// Increment and return the latest cron schedule.
func (s *Splitter) Increment() string {
	if s.current > 60 {
		s.current = 0
	}

	cron := fmt.Sprintf("%d * * * *", s.current)

	s.current = s.current + s.diff

	return cron
}
