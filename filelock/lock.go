package filelock

import (
	"fmt"
	"os"
)

// Locker provides file-based locking for cross-process synchronization.
type Locker struct {
	path   string
	file   *os.File
	locked bool
}

// New creates a new Locker for the specified file path.
// The lock file will be created if it doesn't exist.
// Note: The parent directory must exist.
func New(path string) *Locker {
	return &Locker{
		path:   path,
		locked: false,
	}
}

// Path returns the file path of the lock.
func (l *Locker) Path() string {
	return l.path
}

// IsLocked returns true if the lock is currently held.
func (l *Locker) IsLocked() bool {
	return l.locked
}

// TryLock attempts to acquire the lock without blocking.
// Returns true if the lock was acquired, false if it's already held by another process.
// Returns an error if an unexpected I/O error occurs.
func (l *Locker) TryLock() (bool, error) {
	if l.locked {
		return false, fmt.Errorf("lock is already held by this process")
	}

	// Open or create the lock file
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return false, fmt.Errorf("failed to open lock file: %w", err)
	}

	// Try to acquire the lock
	acquired, err := tryLockFile(f)
	if err != nil {
		_ = f.Close()
		return false, err
	}

	if !acquired {
		_ = f.Close()
		return false, nil
	}

	l.file = f
	l.locked = true
	return true, nil
}

// Lock acquires the lock, blocking until it's available.
// Returns an error if an unexpected I/O error occurs.
func (l *Locker) Lock() error {
	if l.locked {
		return fmt.Errorf("lock is already held by this process")
	}

	// Open or create the lock file
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to open lock file: %w", err)
	}

	// Acquire the lock (blocks until available)
	if err := lockFile(f); err != nil {
		_ = f.Close()
		return err
	}

	l.file = f
	l.locked = true
	return nil
}

// Unlock releases the lock.
// Returns an error if the lock is not held or if unlock fails.
func (l *Locker) Unlock() error {
	if !l.locked {
		return fmt.Errorf("lock is not held by this process")
	}

	if l.file == nil {
		l.locked = false
		return fmt.Errorf("lock file handle is nil")
	}

	err := unlockFile(l.file)
	if closeErr := l.file.Close(); closeErr != nil && err == nil {
		err = closeErr
	}

	l.file = nil
	l.locked = false
	return err
}

// String returns a string representation of the lock state.
func (l *Locker) String() string {
	status := "unlocked"
	if l.locked {
		status = "locked"
	}
	return fmt.Sprintf("Locker{path: %q, status: %s}", l.path, status)
}
