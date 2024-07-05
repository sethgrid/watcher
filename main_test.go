package main

import (
	"sync"
	"testing"
	"time"
)

// could set up table driven tests to compare different newMeta values and assert given events are sent
func TestContentUpdatedTracked(t *testing.T) {

	meta := Meta{
		FilePath:     "/path/to/file",
		LatestUpdate: time.Now(),
		Permissions:  "rw-rw-rw-",
		Content:      []byte("Hello, world"),
	}
	dontPoll := 5 * time.Minute

	w, err := NewWatcher(meta, dontPoll)
	if err != nil {
		t.Fatalf("Error creating watcher: %s", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		select {
		case <-w.done:
			t.Error("done channel closed, expected event instead")
			return
		case event := <-w.Event:
			if event != "Content changed" {
				t.Errorf("Expected 'Content updated', got %s", event)
			}
		}
	}()

	go func() {
		meta.Content = []byte("Goodbye, world")
		meta.LatestUpdate = time.Now().Add(1 * time.Minute)
		w.metaCh <- meta
	}()

	wg.Wait()
	t.Logf("Test complete")
}
