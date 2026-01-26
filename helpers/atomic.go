package helpers

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// AtomicFileWrite performs an atomic file write by:
// 1. Writing to a temporary file in the same directory
// 2. Syncing the file to disk
// 3. Atomically renaming it to the target path
// This prevents partial/corrupted writes if the process crashes.
//
// Note: os.Rename on Unix is atomic but on Windows it fails if the destination exists.
// This implementation is suitable for Unix/Linux systems. For cross-platform support
// with Windows, consider wrapping this or providing a platform-specific alternative.
// If the target file exists, its permissions are preserved in the new file.
func AtomicFileWrite(path string, data []byte) error {
	// Get the directory of the target file
	dir := filepath.Dir(path)
	if dir == "" {
		dir = "."
	}

	// Create temp file in the same directory (ensures same filesystem)
	tmpFile, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	// Track cleanup requirement explicitly to avoid subtle refactoring bugs
	shouldCleanup := true
	defer func() {
		if shouldCleanup {
			_ = os.Remove(tmpPath)
		}
	}()

	// Write data to temp file using io.Copy to handle short writes correctly
	// (os.File.Write can legally return n < len(data) with err == nil)
	if _, err := io.Copy(tmpFile, bytes.NewReader(data)); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("failed to write data: %w", err)
	}

	// If target exists, copy its permissions to the temp file
	if stat, err := os.Stat(path); err == nil {
		if err := tmpFile.Chmod(stat.Mode().Perm()); err != nil {
			_ = tmpFile.Close()
			return fmt.Errorf("failed to set file permissions: %w", err)
		}
	}

	// Sync to disk
	if err := SafeFileSync(tmpFile); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("failed to sync temp file: %w", err)
	}

	// Close the file before renaming
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// Atomically replace the target file
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("failed to rename temp file to target: %w", err)
	}

	// Sync the directory to ensure the rename is durable
	if err := SafeDirSync(dir); err != nil {
		return fmt.Errorf("failed to sync directory: %w", err)
	}

	// Success: prevent cleanup
	shouldCleanup = false
	return nil
}

// SafeFileSync syncs file data and metadata to disk.
// Returns nil on success, or error if sync fails.
func SafeFileSync(f *os.File) error {
	if err := f.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}
	return nil
}

// SafeDirSync syncs a directory to ensure directory operations (renames, deletes) are durable.
// This is important for atomic file operations where the rename must be persisted.
func SafeDirSync(dirPath string) error {
	dir, err := os.Open(dirPath)
	if err != nil {
		return fmt.Errorf("failed to open directory for sync: %w", err)
	}
	defer func() { _ = dir.Close() }()

	if err := dir.Sync(); err != nil {
		return fmt.Errorf("failed to sync directory: %w", err)
	}
	return nil
}

// AtomicFileWriteWithPerm performs atomic file write with specific permissions.
// Similar to AtomicFileWrite but allows explicitly setting file mode.
// The provided perm overrides any existing target file permissions.
//
// Note: os.Rename on Unix is atomic but on Windows it fails if the destination exists.
// This implementation is suitable for Unix/Linux systems. For cross-platform support
// with Windows, consider wrapping this or providing a platform-specific alternative.
func AtomicFileWriteWithPerm(path string, data []byte, perm os.FileMode) error {
	// Get the directory of the target file
	dir := filepath.Dir(path)
	if dir == "" {
		dir = "."
	}

	// Create temp file in the same directory
	tmpFile, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	// Track cleanup requirement explicitly to avoid subtle refactoring bugs
	shouldCleanup := true
	defer func() {
		if shouldCleanup {
			_ = os.Remove(tmpPath)
		}
	}()

	// Write data to temp file using io.Copy to handle short writes correctly
	// (os.File.Write can legally return n < len(data) with err == nil)
	if _, err := io.Copy(tmpFile, bytes.NewReader(data)); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("failed to write data: %w", err)
	}

	// Set permissions before closing
	if err := tmpFile.Chmod(perm); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("failed to set file permissions: %w", err)
	}

	// Sync to disk
	if err := SafeFileSync(tmpFile); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("failed to sync temp file: %w", err)
	}

	// Close the file before renaming
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// Atomically replace the target file
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("failed to rename temp file to target: %w", err)
	}

	// Sync the directory
	if err := SafeDirSync(dir); err != nil {
		return fmt.Errorf("failed to sync directory: %w", err)
	}

	// Success: prevent cleanup
	shouldCleanup = false
	return nil
}
