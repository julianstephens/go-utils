# File Lock Package

The `filelock` package provides cross-platform file locking utilities for coordinating single-writer access between processes. It supports Linux, macOS, and Windows with a unified API.

## Features

- **Cross-platform**: Unix-like systems and Windows with unified API
- **Simple API**: `Lock()`, `TryLock()`, `Unlock()` methods
- **Non-blocking**: Optional non-blocking lock attempts
- **Error handling**: Clear error messages for common scenarios

## Installation

```bash
go get github.com/julianstephens/go-utils/filelock
```

## Usage

### Basic Locking

```go
package main

import (
	"fmt"
	"github.com/julianstephens/go-utils/filelock"
)

func main() {
	lock := filelock.New("/tmp/myapp.lock")

	// Acquire the lock (blocking)
	if err := lock.Lock(); err != nil {
		fmt.Printf("Failed to acquire lock: %v\n", err)
		return
	}
	defer lock.Unlock()

	// Perform protected operations
	fmt.Println("Lock acquired, performing operations...")
}
```

### Non-blocking Lock

```go
lock := filelock.New("/tmp/myapp.lock")

acquired, err := lock.TryLock()
if err != nil {
	fmt.Printf("Lock error: %v\n", err)
	return
}

if !acquired {
	fmt.Println("Lock is held by another process")
	return
}

defer lock.Unlock()

// Perform protected operations
fmt.Println("Lock acquired, performing operations...")
```

### Checking Lock Status

```go
lock := filelock.New("/tmp/myapp.lock")

if lock.IsLocked() {
	fmt.Println("Lock is currently held")
} else {
	fmt.Println("Lock is not held")
}

fmt.Println(lock.String()) // Print lock status
```

## API Reference

**`New(path string) *Locker`** - Creates a new Locker. Parent directory must exist.

**`Lock() error`** - Acquires the lock, blocking until available.

**`TryLock() (bool, error)`** - Non-blocking lock attempt. Returns true if acquired, false if held by another process.

**`Unlock() error`** - Releases the lock.

**`IsLocked() bool`** - Returns true if lock is held by this process.

**`Path() string`** - Returns the lock file path.

**`String() string`** - Returns string representation of lock state.

## Platform-Specific Behavior

**Unix-like Systems (Linux, macOS):** Uses `flock()` syscalls for advisory locking. Automatic release on process termination. Works across NFS (with caveats).

**Windows:** Uses `LockFile` API for mandatory locking. Automatic release on file close or process termination. Does not work across network shares.

## Use Cases

- **Single-Writer Databases**: Enforce one writer at a time (e.g., waldb)
- **PID Files**: Prevent multiple application instances
- **Resource Coordination**: Serialize cross-process access
- **Config Protection**: Prevent concurrent file updates

## Error Handling

Common errors:
- `lock is already held by this process` - Lock already acquired by this Locker
- `lock is not held by this process` - Unlock called without holding lock
- I/O errors - File system issues (permission denied, disk full, etc.)

## Example: Single-Writer Database

```go
type DB struct {
	path string
	lock *filelock.Locker
}

func Open(path string) (*DB, error) {
	db := &DB{
		path: path,
		lock: filelock.New(path + ".lock"),
	}
	if err := db.lock.Lock(); err != nil {
		return nil, fmt.Errorf("failed to acquire database lock: %w", err)
	}
	return db, nil
}

func (db *DB) Close() error {
	return db.lock.Unlock()
}
```

## Testing

Run tests with:

```bash
go test ./filelock
```

Run tests with race detection:

```bash
go test -race ./filelock
```

## Limitations

- **Advisory Locking**: All processes must cooperate to respect locks
- **Not Reentrant**: Cannot lock same Locker instance twice
- **File Persistence**: Lock files remain on disk after unlock
- **NFS Unreliable**: Unix file locking may not work over NFS

## Thread Safety

The `Locker` type is NOT thread-safe. Each goroutine should use its own `Locker` instance, or external synchronization must be used.

## Related Packages

- `github.com/julianstephens/go-utils/helpers` - Additional file system utilities
