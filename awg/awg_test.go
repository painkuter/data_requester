package awg

import (
	"errors"
	"runtime"
	"testing"
	"time"
)

func slowFunc() error {
	time.Sleep(time.Second)
	return nil
}

func fastFunc() error {
	// do nothing
	return nil
}

func errorFunc() error {
	return errors.New("Test error")
}

func panicFunc() error {
	panic("Test panic")
	return nil
}

// TestAdvancedWorkGroupTimeout test for timeout
func Test_AdvancedWorkGroupTimeout(t *testing.T) {
	var wg AdvancedWaitGroup

	wg.Add(fastFunc)
	wg.Add(fastFunc)
	wg.Add(fastFunc)
	wg.Add(slowFunc)
	wg.Add(slowFunc)
	wg.Add(slowFunc)

	wg.SetTimeout(time.Nanosecond * 10).
		Start()

	if wg.Status() != StatusTimeout {
		t.Error("AWG should stops by timeout!")
	}
}

// TestAdvancedWorkGroupError test for error
func Test_AdvancedWorkGroupError(t *testing.T) {
	var wg AdvancedWaitGroup

	wg.Add(errorFunc)
	wg.Add(fastFunc)
	wg.Add(fastFunc)
	wg.Add(slowFunc)
	wg.Add(slowFunc)
	wg.Add(slowFunc)

	wg.SetStopOnError(true).
		Start()

	if wg.Status() != StatusError {
		t.Error("AWG should stops by error!")
	}

}

// TestAdvancedWorkGroupSuccess test for success case
func Test_AdvancedWorkGroupSuccess(t *testing.T) {
	var wg AdvancedWaitGroup

	wg.Add(fastFunc)
	wg.Add(fastFunc)
	wg.Add(fastFunc)
	wg.Add(fastFunc)
	wg.Add(slowFunc)
	wg.Add(slowFunc)

	wg.SetStopOnError(true).
		Start()

	if wg.Status() != StatusSuccess {
		t.Error("AWG result should be 'success'!")
	}
}

// TestAdvancedWorkGroupPanic test for success case
func Test_AdvancedWorkGroupPanic(t *testing.T) {
	var wg AdvancedWaitGroup

	wg.Add(slowFunc)
	wg.Add(panicFunc)

	wg.SetStopOnError(true).
		Start()

	if wg.Status() != StatusError {
		t.Error("AWG result should be 'error'!")
	}
}

// TestAWGStopOnError tests AdvancedWaitGroup with StopOnError set to false
// and with failing task.
func TestAWGStopOnError(t *testing.T) {
	var wg AdvancedWaitGroup
	wg.Add(fastFunc)
	wg.Add(errorFunc)
	wg.Add(fastFunc)
	wg.SetStopOnError(false).
		Start()
}

// Test_AdvancedWorkGroup_NoLeak tests for goroutines leaks
func Test_AdvancedWorkGroup_NoLeak(t *testing.T) {
	var wg AdvancedWaitGroup

	wg.Add(errorFunc)

	wg.SetStopOnError(true).
		Start()

	time.Sleep(2 * time.Second)

	numGoroutines := runtime.NumGoroutine()

	var wg2 AdvancedWaitGroup

	wg2.Add(errorFunc)
	wg2.Add(slowFunc)

	wg2.SetStopOnError(true).
		Start()

	time.Sleep(2 * time.Second)

	numGoroutines2 := runtime.NumGoroutine()

	if numGoroutines != numGoroutines2 {
		t.Fatalf("We leaked %d goroutine(s)", numGoroutines2-numGoroutines)
	}
}
