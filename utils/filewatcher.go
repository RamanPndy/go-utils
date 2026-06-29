package goutils

import "fmt"

type Watcher struct {
	Events chan Event
	Errors chan error
}

type Event struct {
	Name string
	Op   Op
}

type Op uint32

const (
	Create Op = 1 << iota
	Write
	Remove
	Rename
	Chmod
)

func NewWatcher() (*Watcher, error) {
	return &Watcher{
		Events: make(chan Event),
		Errors: make(chan error),
	}, nil
}

type FileWatcher struct {
	*Watcher
}

func NewFileWatcher() (*FileWatcher, error) {
	watcher, err := NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create file watcher: %w", err)
	}
	return &FileWatcher{watcher}, nil
}

func (fw *FileWatcher) Close() error {
	return fw.Close()
}

func (fw *FileWatcher) Add(name string) error {
	// In a real implementation, you would add the file to the watch list here.
	// For this example, we'll just simulate that the file is being watched.
	return nil
}

func NotifyOnFileChange(filePath string, callback func()) error {
	watcher, err := NewFileWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = watcher.Add(filePath)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&Write == Write || event.Op&Create == Create || event.Op&Remove == Remove {
					callback()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("error:", err)
			}
		}
	}()

	return nil
}
