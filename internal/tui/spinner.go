package tui

import (
	"fmt"
	"time"
)

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// graceDelay is how long to wait before showing the spinner.
// Fast commands finish within this window and never display it.
const graceDelay = 120 * time.Millisecond

// tickInterval controls spinner animation speed.
const tickInterval = 80 * time.Millisecond

// RunWithSpinner runs f in a goroutine and shows a spinner on stdout while it
// runs. The spinner only appears if f takes longer than graceDelay, so fast
// commands produce no visual noise.
func RunWithSpinner(f func()) {
	done := make(chan struct{})
	go func() {
		f()
		close(done)
	}()

	// Grace period — if f finishes quickly, skip the spinner entirely.
	select {
	case <-done:
		return
	case <-time.After(graceDelay):
	}

	// f is still running — start animating.
	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()

	i := 0
	for {
		fmt.Printf("\r%s", prefixStyle.Render(spinnerFrames[i%len(spinnerFrames)]))
		i++
		select {
		case <-done:
			fmt.Print("\r\033[K") // clear the spinner line
			return
		case <-ticker.C:
		}
	}
}
