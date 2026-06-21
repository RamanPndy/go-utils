// Package dsnloader reads a database DSN from a file and watches it for
// changes using fsnotify, enabling hot-reload of the database connection
// without restarting the process.
package goutils

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
)

// OnChange is called when the DSN file contents change. The new DSN string
// is passed as the argument.
type OnChange func(newDSN string) error

// Loader watches a DSN file and invokes a callback when the DSN changes.
type Loader struct {
	path     string
	onChange OnChange

	mu      sync.RWMutex
	current string

	watcher *fsnotify.Watcher
	done    chan struct{}
	logger  *logrus.Logger
}

// NewDsnLoader creates a Loader that reads the DSN from path and watches it for changes.
// The initial DSN value is read synchronously; call Watch to begin file monitoring.
func NewDsnLoader(path string, onChange OnChange, logger *logrus.Logger) (*Loader, error) {
	dsn, err := readDSN(path)
	if err != nil {
		return nil, fmt.Errorf("dsnloader: initial read: %w", err)
	}

	return &Loader{
		path:     path,
		onChange: onChange,
		current:  dsn,
		done:     make(chan struct{}),
		logger:   logger,
	}, nil
}

// DSN returns the most recently loaded DSN string.
func (l *Loader) DSN() string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.current
}

// Watch starts watching the DSN file for changes in a background goroutine.
// It returns immediately. Call Close to stop watching.
func (l *Loader) Watch() error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("dsnloader: creating watcher: %w", err)
	}
	if err := w.Add(l.path); err != nil {
		_ = w.Close()
		return fmt.Errorf("dsnloader: watching %s: %w", l.path, err)
	}
	l.watcher = w

	go l.loop()
	return nil
}

// Close stops the file watcher and releases resources.
func (l *Loader) Close() error {
	if l.watcher != nil {
		err := l.watcher.Close()
		<-l.done
		return err
	}
	return nil
}

func (l *Loader) loop() {
	defer close(l.done)
	for {
		select {
		case event, ok := <-l.watcher.Events:
			if !ok {
				return
			}
			if event.Op&(fsnotify.Write|fsnotify.Create) == 0 {
				continue
			}
			l.reload()
		case err, ok := <-l.watcher.Errors:
			if !ok {
				return
			}
			l.logger.Warn("Watch error", "error", err)
		}
	}
}

func (l *Loader) reload() {
	newDSN, err := readDSN(l.path)
	if err != nil {
		l.logger.Warn("Failed to read DSN file", "path", l.path, "error", err)
		return
	}

	l.mu.RLock()
	same := newDSN == l.current
	l.mu.RUnlock()
	if same {
		return
	}

	l.logger.Info("DSN changed; reconnecting")
	if err := l.onChange(newDSN); err != nil {
		l.logger.Warn("Reconnect failed", "error", err)
		return
	}

	l.mu.Lock()
	l.current = newDSN
	l.mu.Unlock()
	l.logger.Info("Reconnect successful")
}

func readDSN(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	dsn := strings.TrimSpace(string(data))
	if dsn == "" {
		return "", fmt.Errorf("dsnloader: file %s is empty", path)
	}
	return dsn, nil
}

// ReadFile reads and returns the trimmed DSN string from a file.
func ReadFile(path string) (string, error) {
	return readDSN(path)
}
