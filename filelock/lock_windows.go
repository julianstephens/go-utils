//go:build windows
// +build windows

package filelock

import (
	"fmt"
	"os"
	"syscall"
)

// lockFile acquires an exclusive lock on the file (blocking).
// Uses Windows LockFile API.
func lockFile(f *os.File) error {
	handle := syscall.Handle(f.Fd())

	// Lock the entire file (offset 0, length 0xffffffff)
	// Note: On Windows, locking is mandatory, not advisory
	err := lockFileRange(handle, 0, 0xffffffff)
	if err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}
	return nil
}

// tryLockFile attempts to acquire an exclusive lock without blocking.
// Returns true if the lock was acquired, false if it's held by another process.
func tryLockFile(f *os.File) (bool, error) {
	handle := syscall.Handle(f.Fd())

	// Try to lock without blocking
	err := lockFileRangeNonBlocking(handle, 0, 0xffffffff)
	if err == syscall.ERROR_LOCK_VIOLATION {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to try lock: %w", err)
	}
	return true, nil
}

// unlockFile releases the lock on the file.
func unlockFile(f *os.File) error {
	handle := syscall.Handle(f.Fd())
	return unlockFileRange(handle, 0, 0xffffffff)
}

// lockFileRange locks a range of bytes in a file
func lockFileRange(handle syscall.Handle, offset, length uint32) error {
	ret, _, err := syscall.SyscallN(
		syscall.NewLazyDLL("kernel32.dll").NewProc("LockFile").Addr(),
		uintptr(handle),
		uintptr(offset),
		uintptr(0), // offset high
		uintptr(length),
		uintptr(0), // length high
	)
	if ret == 0 {
		return err
	}
	return nil
}

// lockFileRangeNonBlocking locks a range of bytes without blocking
func lockFileRangeNonBlocking(handle syscall.Handle, offset, length uint32) error {
	ret, _, err := syscall.SyscallN(
		syscall.NewLazyDLL("kernel32.dll").NewProc("LockFileEx").Addr(),
		uintptr(handle),
		uintptr(2), // LOCKFILE_EXCLUSIVE_LOCK | LOCKFILE_FAIL_IMMEDIATELY
		uintptr(0), // reserved
		uintptr(length),
		uintptr(0), // length high
	)
	if ret == 0 {
		return err
	}
	return nil
}

// unlockFileRange unlocks a range of bytes in a file
func unlockFileRange(handle syscall.Handle, offset, length uint32) error {
	ret, _, err := syscall.SyscallN(
		syscall.NewLazyDLL("kernel32.dll").NewProc("UnlockFile").Addr(),
		uintptr(handle),
		uintptr(offset),
		uintptr(0), // offset high
		uintptr(length),
		uintptr(0), // length high
	)
	if ret == 0 {
		return err
	}
	return nil
}
