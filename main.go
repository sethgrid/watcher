package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	meta, err := GrabMeta("/path/set/by/env/variable/or/flag")
	// meta, err := GrabMeta("README.md")
	if err != nil {
		fmt.Printf("Error getting file info: %s", err)
		os.Exit(1)
	}

	w, err := NewWatcher(meta, 5*time.Second)
	if err != nil {
		fmt.Printf("Error creating watcher: %s", err)
		os.Exit(1)
	}

	go func() {
		for event := range w.Event {
			fmt.Printf("Event: %s\n", event)
		}
	}()

	// Block until killed
	<-w.Closed
}

func NewWatcher(meta Meta, pollTime time.Duration) (*Watcher, error) {
	w := &Watcher{
		Event:  make(chan string, 1),
		done:   make(chan struct{}),
		metaCh: make(chan Meta),
		Closed: make(chan struct{}),
	}

	go func() {
		ticker := time.NewTicker(pollTime)
		defer ticker.Stop()
		for {
			select {
			case <-w.done:
				return
			case <-ticker.C:
				newMeta, err := GrabMeta(meta.FilePath)
				if err != nil {
					w.Event <- fmt.Sprintf("Error getting file info: %s", err)
					continue
				}
				w.metaCh <- newMeta
				// fmt.Println("Polled")
			}
		}
	}()

	go func() {
		for {
			select {
			case <-w.done:
				close(w.Closed)
				return
			case newMeta := <-w.metaCh:
				event, err := GenEvent(meta, newMeta)
				if err != nil {
					w.Event <- fmt.Sprintf("Error generating event: %s", err)
				}
				w.Event <- event
				meta = newMeta
			}
		}
	}()

	return w, nil
}

func GenEvent(oldMeta, newMeta Meta) (string, error) {
	if string(oldMeta.Content) != string(newMeta.Content) {
		return "Content changed", nil
	}

	if oldMeta.Permissions != newMeta.Permissions {
		return "Permissions changed", nil
	}

	return "Unknown change", nil
}

type Watcher struct {
	Event  chan string
	done   chan struct{}
	metaCh chan Meta
	Closed chan struct{}
}

type Meta struct {
	FilePath     string
	LatestUpdate time.Time
	Permissions  string
	Content      []byte
}

func GrabMeta(filepath string) (Meta, error) {
	file, err := os.Stat(filepath)
	if err != nil {
		return Meta{}, fmt.Errorf("Error getting file info: %w", err)
	}

	content, err := os.ReadFile(filepath)
	if err != nil {
		return Meta{}, fmt.Errorf("Error reading file: %w", err)
	}

	return Meta{
		FilePath:     filepath,
		LatestUpdate: file.ModTime(),
		Permissions:  file.Mode().String(),
		Content:      content,
	}, nil
}
