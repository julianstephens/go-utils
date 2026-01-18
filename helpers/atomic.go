package helpers

import (
	"fmt"
	"os"
	"path/filepath"
)

// AtomicFileWrite performs an atomic file write by:
// 1. Writing to a temporary file in the same directory
// 2. Syncing the file to disk
// 3. Atomically renaming it to the target path
// This prevents partial/corrupted writes if the process crashes.
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

	// Clean up temp file if something goes wrong
	defer func() {
		if err != nil {
			_ = os.Remove(tmpPath)
		}
	}()

	// Write data to temp file
	if _, err := tmpFile.Write(data); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("failed to write data: %w", err)
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
// Similar to AtomicFileWrite but allows setting file mode.
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

	// Clean up temp file if something goes wrong
	defer func() {
		if err != nil {
			_ = os.Remove(tmpPath)
		}
	}()

	// Write data to temp file
	if _, err := tmpFile.Write(data); err != nil {
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

	return nil
}
