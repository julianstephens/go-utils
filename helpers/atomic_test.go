package helpers

import (
	"os"
	"path/filepath"
	"testing"
)

// TestAtomicFileWrite tests basic atomic file write
func TestAtomicFileWrite(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")

	data := []byte("test data")

	err := AtomicFileWrite(filePath, data)
	if err != nil {
		t.Fatalf("AtomicFileWrite failed: %v", err)
	}

	// Verify file exists and contains correct data
	readData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if string(readData) != string(data) {
		t.Errorf("File content mismatch: got %q, want %q", readData, data)
	}
}

// TestAtomicFileWriteOverwrite tests overwriting existing file
func TestAtomicFileWriteOverwrite(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")

	// Write initial data
	initialData := []byte("initial")
	if err := AtomicFileWrite(filePath, initialData); err != nil {
		t.Fatalf("Initial write failed: %v", err)
	}

	// Overwrite with new data
	newData := []byte("new data that is longer")
	if err := AtomicFileWrite(filePath, newData); err != nil {
		t.Fatalf("Overwrite write failed: %v", err)
	}

	// Verify new data
	readData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if string(readData) != string(newData) {
		t.Errorf("File content mismatch: got %q, want %q", readData, newData)
	}
}

// TestAtomicFileWriteLargeData tests writing large data
func TestAtomicFileWriteLargeData(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "large.bin")

	// Create 1MB of data
	data := make([]byte, 1024*1024)
	for i := range data {
		data[i] = byte(i % 256)
	}

	err := AtomicFileWrite(filePath, data)
	if err != nil {
		t.Fatalf("AtomicFileWrite with large data failed: %v", err)
	}

	// Verify size
	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	if info.Size() != int64(len(data)) {
		t.Errorf("File size mismatch: got %d, want %d", info.Size(), len(data))
	}
}

// TestAtomicFileWriteWithPerm tests atomic write with custom permissions
func TestAtomicFileWriteWithPerm(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "perm_test.txt")

	data := []byte("test data")
	perm := os.FileMode(0600) // Read/write for owner only

	err := AtomicFileWriteWithPerm(filePath, data, perm)
	if err != nil {
		t.Fatalf("AtomicFileWriteWithPerm failed: %v", err)
	}

	// Verify permissions
	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	if info.Mode()&0777 != perm {
		t.Errorf("File permissions mismatch: got 0%o, want 0%o", info.Mode()&0777, perm)
	}

	// Verify content
	readData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if string(readData) != string(data) {
		t.Errorf("File content mismatch: got %q, want %q", readData, data)
	}
}

// TestSafeFileSync tests file syncing
func TestSafeFileSync(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "sync_test.txt")

	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer func() { _ = file.Close() }()

	_, err = file.WriteString("test data")
	if err != nil {
		t.Fatalf("Failed to write to file: %v", err)
	}

	// Should not error
	if err := SafeFileSync(file); err != nil {
		t.Errorf("SafeFileSync failed: %v", err)
	}
}

// TestSafeDirSync tests directory syncing
func TestSafeDirSync(t *testing.T) {
	tmpDir := t.TempDir()

	// Should not error
	if err := SafeDirSync(tmpDir); err != nil {
		t.Errorf("SafeDirSync failed: %v", err)
	}
}

// TestSafeDirSyncInvalidDir tests directory sync with invalid directory
func TestSafeDirSyncInvalidDir(t *testing.T) {
	// Non-existent directory should error
	err := SafeDirSync("/nonexistent/directory/path")
	if err == nil {
		t.Errorf("SafeDirSync should fail for non-existent directory")
	}
}

// TestAtomicFileWriteInvalidDir tests write to invalid directory
func TestAtomicFileWriteInvalidDir(t *testing.T) {
	// Directory doesn't exist
	err := AtomicFileWrite("/nonexistent/dir/file.txt", []byte("data"))
	if err == nil {
		t.Errorf("AtomicFileWrite should fail for non-existent directory")
	}
}

// TestAtomicFileWriteEmptyData tests writing empty data
func TestAtomicFileWriteEmptyData(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "empty.txt")

	data := []byte{}

	err := AtomicFileWrite(filePath, data)
	if err != nil {
		t.Fatalf("AtomicFileWrite with empty data failed: %v", err)
	}

	// Verify file exists and is empty
	readData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if len(readData) != 0 {
		t.Errorf("File should be empty, got %d bytes", len(readData))
	}
}

// TestAtomicFileWriteMultiple tests multiple consecutive writes
func TestAtomicFileWriteMultiple(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "multi.txt")

	testCases := [][]byte{
		[]byte("first"),
		[]byte("second longer data"),
		[]byte("short"),
		[]byte("another very long line of text for testing"),
	}

	for i, data := range testCases {
		err := AtomicFileWrite(filePath, data)
		if err != nil {
			t.Fatalf("Write %d failed: %v", i, err)
		}

		readData, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Read %d failed: %v", i, err)
		}

		if string(readData) != string(data) {
			t.Errorf("Write %d mismatch: got %q, want %q", i, readData, data)
		}
	}
}

// TestAtomicFileWritePreservesContent tests that concurrent writes don't corrupt data
func TestAtomicFileWriteConsistency(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "consistency.txt")

	// Write specific test data
	testData := []byte("consistent test data")
	if err := AtomicFileWrite(filePath, testData); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Read multiple times to ensure consistency
	for i := 0; i < 5; i++ {
		readData, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Read %d failed: %v", i, err)
		}

		if string(readData) != string(testData) {
			t.Errorf("Read %d inconsistent: got %q, want %q", i, readData, testData)
		}
	}
}

// TestAtomicFileWriteNestedDir tests writing to nested directories
func TestAtomicFileWriteNestedDir(t *testing.T) {
	tmpDir := t.TempDir()
	// Note: nested path must already exist
	nestedDir := filepath.Join(tmpDir, "a", "b", "c")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatalf("Failed to create nested directory: %v", err)
	}

	filePath := filepath.Join(nestedDir, "nested.txt")
	data := []byte("nested file data")

	err := AtomicFileWrite(filePath, data)
	if err != nil {
		t.Fatalf("AtomicFileWrite to nested dir failed: %v", err)
	}

	readData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read nested file: %v", err)
	}

	if string(readData) != string(data) {
		t.Errorf("Nested file content mismatch: got %q, want %q", readData, data)
	}
}

// TestAtomicFileWritePreservesPermissions tests that existing file permissions are preserved
func TestAtomicFileWritePreservesPermissions(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "perm_preserve.txt")

	// Create initial file with specific permissions
	initialPerm := os.FileMode(0600)
	initialData := []byte("initial data")
	if err := AtomicFileWriteWithPerm(filePath, initialData, initialPerm); err != nil {
		t.Fatalf("Initial write failed: %v", err)
	}

	// Verify initial permissions
	stat, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("Failed to stat initial file: %v", err)
	}
	if stat.Mode().Perm() != initialPerm {
		t.Errorf("Initial permissions mismatch: got 0%o, want 0%o", stat.Mode().Perm(), initialPerm)
	}

	// Overwrite without specifying permissions - should preserve them
	newData := []byte("new data")
	if err := AtomicFileWrite(filePath, newData); err != nil {
		t.Fatalf("Overwrite failed: %v", err)
	}

	// Verify permissions were preserved
	stat, err = os.Stat(filePath)
	if err != nil {
		t.Fatalf("Failed to stat overwritten file: %v", err)
	}
	if stat.Mode().Perm() != initialPerm {
		t.Errorf("Permissions not preserved: got 0%o, want 0%o", stat.Mode().Perm(), initialPerm)
	}

	// Verify new data was written
	readData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if string(readData) != string(newData) {
		t.Errorf("File content mismatch: got %q, want %q", readData, newData)
	}
}

// TestAtomicFileWriteHandlesShortWrites tests that short writes are handled correctly
// by using io.Copy internally instead of relying on os.File.Write
func TestAtomicFileWriteHandlesShortWrites(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "short_write.txt")

	// Create a large data buffer to increase likelihood of short writes
	// (though on most systems this will still write all in one go)
	data := make([]byte, 10*1024*1024) // 10MB
	for i := range data {
		data[i] = byte(i % 256)
	}

	err := AtomicFileWrite(filePath, data)
	if err != nil {
		t.Fatalf("AtomicFileWrite with large data failed: %v", err)
	}

	// Verify all data was written
	readData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if len(readData) != len(data) {
		t.Errorf("File size mismatch: got %d bytes, want %d bytes", len(readData), len(data))
	}

	// Spot check some data
	for i, expected := range []int{0, 1000, 100000, 1000000, len(data) - 1} {
		if i >= len(data) {
			break
		}
		if readData[expected] != data[expected] {
			t.Errorf("Data mismatch at byte %d: got %d, want %d", expected, readData[expected], data[expected])
		}
	}
}

// TestAtomicFileWritePermOverride tests that AtomicFileWriteWithPerm overrides permissions
func TestAtomicFileWritePermOverride(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "perm_override.txt")

	// Create initial file with permissive permissions
	initialPerm := os.FileMode(0644)
	initialData := []byte("initial")
	if err := AtomicFileWriteWithPerm(filePath, initialData, initialPerm); err != nil {
		t.Fatalf("Initial write failed: %v", err)
	}

	// Overwrite with more restrictive permissions
	newPerm := os.FileMode(0600)
	newData := []byte("new restricted")
	if err := AtomicFileWriteWithPerm(filePath, newData, newPerm); err != nil {
		t.Fatalf("Overwrite with new perms failed: %v", err)
	}

	// Verify new permissions were set
	stat, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}
	if stat.Mode().Perm() != newPerm {
		t.Errorf("Permissions not overridden: got 0%o, want 0%o", stat.Mode().Perm(), newPerm)
	}

	// Verify new data
	readData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if string(readData) != string(newData) {
		t.Errorf("File content mismatch: got %q, want %q", readData, newData)
	}
}
