package awg

import (
	"errors"
	"fmt"
	"runtime"
	"time"
)

const (
	// StatusIdle means that WG did not run yet
	StatusIdle int = iota
	// StatusSuccess means successful execution of all tasks
	StatusSuccess
	// StatusTimeout means that job was broken by timeout
	StatusTimeout
	// StatusError means that job was broken by error in one task (if stopOnError is true)
	StatusError
)

// WaitgroupFunc func
type WaitgroupFunc func() error

// AdvancedWaitGroup enhanced wait group struct
type AdvancedWaitGroup struct {
	stack       []WaitgroupFunc
	timeout     time.Duration
	stopOnError bool
	status      int
	errors      []error
}

// SetTimeout defines timeout for all tasks
func (wg *AdvancedWaitGroup) SetTimeout(t time.Duration) *AdvancedWaitGroup {
	wg.timeout = t
	return wg
}

// SetStopOnError make wiatgroup stops if any task returns error
func (wg *AdvancedWaitGroup) SetStopOnError(b bool) *AdvancedWaitGroup {
	wg.stopOnError = b
	return wg
}

// Add adds new task in waitgroup
func (wg *AdvancedWaitGroup) Add(f WaitgroupFunc) *AdvancedWaitGroup {
	wg.stack = append(wg.stack, f)
	return wg
}

// AddSlice adds new tasks in waitgroup
func (wg *AdvancedWaitGroup) AddSlice(s []WaitgroupFunc) *AdvancedWaitGroup {
	for _, f := range s {
		wg.stack = append(wg.stack, f)
	}
	return wg
}

// Start runs tasks in separate goroutines
func (wg *AdvancedWaitGroup) Start() *AdvancedWaitGroup {
	wg.status = StatusSuccess

	if taskCount := len(wg.stack); taskCount > 0 {
		failed := make(chan error, taskCount)
		done := make(chan bool, taskCount)
		timer := time.NewTimer(wg.timeout)

		for _, f := range wg.stack {
			go func(f WaitgroupFunc, failed chan<- error, done chan<- bool) {
				// Handle panic and pack it into stdlib error
				defer func() {
					if r := recover(); r != nil {
						buf := make([]byte, 1000)
						runtime.Stack(buf, false)
						failed <- errors.New(fmt.Sprintf("Panic handeled\n%v\n%s", r, string(buf)))
					}
				}()

				if err := f(); err != nil {
					failed <- err
				} else {
					done <- true
				}
			}(f, failed, done)
		}

	ForLoop:
		for taskCount > 0 {
			select {
			case err := <-failed:
				wg.errors = append(wg.errors, err)
				taskCount--
				if wg.stopOnError {
					wg.status = StatusError
					break ForLoop
				}
			case <-done:
				taskCount--
			case <-timer.C:
				if wg.timeout > 0 {
					wg.status = StatusTimeout
					break ForLoop
				}
			}
		}
	}

	return wg
}

// Reset performs cleanup task queue and reset state
func (wg *AdvancedWaitGroup) Reset() {
	wg.stack = []WaitgroupFunc{}
	wg.timeout = 0
	wg.stopOnError = false
	wg.status = StatusIdle
	wg.errors = []error{}
}

// GetLastError returns last error that caught by execution process
func (wg *AdvancedWaitGroup) GetLastError() error {
	if l := len(wg.errors); l > 0 {
		return wg.errors[l-1]
	}
	return nil
}

// GetAllErrors returns all errors that caught by execution process
func (wg *AdvancedWaitGroup) GetAllErrors() []error {
	return wg.errors
}

// Status return result state string
func (wg *AdvancedWaitGroup) Status() int {
	return wg.status
}
