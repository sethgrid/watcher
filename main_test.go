package main

import (
	"fmt"
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
			return
		case event := <-w.Event:
			fmt.Printf("fmt Event: %s\n", event)
			t.Errorf("t Event: %s", event)
		}

	}()

	go func() {
		newMeta := meta
		meta.Content = []byte("Goodbye, world")
		newMeta.LatestUpdate = time.Now().Add(1 * time.Minute)
		w.metaCh <- newMeta
	}()

	wg.Wait()
	t.Logf("Test complete")
}