package runner

import (
	"errors"
	"sync"
	"testing"
)

// TestPipeline_ConcurrentErrors mirrors the stdout and stderr copy goroutines both failing
// when a process exits: recording the error must be race-free and keep the first one.
func TestPipeline_ConcurrentErrors(t *testing.T) {
	p := &Pipeline{}
	first := errors.New("first")

	p.setErr(first)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() { defer wg.Done(); p.setErr(errors.New("later")) }()
		go func() { defer wg.Done(); _ = p.Err() }()
	}
	wg.Wait()

	if !errors.Is(p.Err(), first) {
		t.Fatalf("Err() = %v, want the first error", p.Err())
	}
}
