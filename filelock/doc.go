package filelock

// Package filelock provides cross-platform file locking utilities for
// coordinating access between processes. It supports both Unix-like systems
// (Linux, macOS) and Windows with a unified API.
//
// The package implements advisory file locking, meaning processes must
// cooperate to respect locks. It is suitable for ensuring single-writer
// access to resources like databases.
//
// Basic usage:
//
//	lock := filelock.New("/path/to/lockfile")
//	if err := lock.Lock(); err != nil {
//		// Handle lock acquisition failure
//	}
//	defer lock.Unlock()
//
//	// Perform protected operations
