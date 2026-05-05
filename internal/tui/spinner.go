package tui

import (
	"fmt"
	"sync"
	"time"
)

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// graceDelay is how long to wait before showing the spinner.
// Fast commands finish within this window and never display it.
const graceDelay = 120 * time.Millisecond

// tickInterval controls spinner animation speed.
const tickInterval = 80 * time.Millisecond

// spinnerHandle lets the goroutine running inside RunWithSpinner pause and
// resume the spinner loop for interactive prompts.
type spinnerHandle struct {
	pauseReq  chan struct{}
	pauseAck  chan struct{} // spinner sends after it has cleared the line
	resumeReq chan struct{}
}

var (
	activeHandleMu sync.Mutex
	activeHandle   *spinnerHandle
)

// PauseSpinner pauses the running spinner (if any) and blocks until the
// spinner line has been cleared. Safe to call when no spinner is active.
func PauseSpinner() {
	activeHandleMu.Lock()
	h := activeHandle
	activeHandleMu.Unlock()
	if h == nil {
		return
	}
	h.pauseReq <- struct{}{}
	<-h.pauseAck
}

// ResumeSpinner resumes the spinner after PauseSpinner.
func ResumeSpinner() {
	activeHandleMu.Lock()
	h := activeHandle
	activeHandleMu.Unlock()
	if h == nil {
		return
	}
	h.resumeReq <- struct{}{}
}

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

	// f is still running — register the handle so PauseSpinner can find it.
	h := &spinnerHandle{
		pauseReq:  make(chan struct{}),
		pauseAck:  make(chan struct{}),
		resumeReq: make(chan struct{}),
	}
	activeHandleMu.Lock()
	activeHandle = h
	activeHandleMu.Unlock()
	defer func() {
		activeHandleMu.Lock()
		activeHandle = nil
		activeHandleMu.Unlock()
	}()

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
		case <-h.pauseReq:
			fmt.Print("\r\033[K") // clear before handing control back
			h.pauseAck <- struct{}{}
			select {
			case <-h.resumeReq:
			case <-done:
				fmt.Print("\r\033[K")
				return
			}
		case <-ticker.C:
		}
	}
}
