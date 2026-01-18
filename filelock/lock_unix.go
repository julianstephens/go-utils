//go:build unix
// +build unix

package filelock

import (
	"fmt"
	"os"
	"syscall"
)

// lockFile acquires an exclusive lock on the file (blocking).
// Uses flock-based locking on Unix systems.
func lockFile(f *os.File) error {
	return syscall.Flock(int(f.Fd()), syscall.LOCK_EX)
}

// tryLockFile attempts to acquire an exclusive lock without blocking.
// Returns true if the lock was acquired, false if it's held by another process.
func tryLockFile(f *os.File) (bool, error) {
	err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err == syscall.EWOULDBLOCK || err == syscall.EAGAIN {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to acquire lock: %w", err)
	}
	return true, nil
}

// unlockFile releases the lock on the file.
func unlockFile(f *os.File) error {
	return syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
}
