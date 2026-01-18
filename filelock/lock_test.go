package filelock_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/julianstephens/go-utils/filelock"
	tst "github.com/julianstephens/go-utils/tests"
)

func TestNew(t *testing.T) {
	lock := filelock.New("/tmp/test.lock")
	tst.AssertNotNil(t, lock, "New() should return non-nil locker")
	tst.AssertDeepEqual(t, lock.Path(), "/tmp/test.lock")
	tst.AssertFalse(t, lock.IsLocked(), "new lock should not be locked")
}

func TestLock(t *testing.T) {
	tmpDir := t.TempDir()
	lockPath := filepath.Join(tmpDir, "test.lock")

	lock := filelock.New(lockPath)

	// Acquire the lock
	err := lock.Lock()
	tst.AssertNoError(t, err, "Lock() should succeed")
	tst.AssertTrue(t, lock.IsLocked(), "lock should be locked after Lock()")

	// Clean up
	err = lock.Unlock()
	tst.AssertNoError(t, err, "Unlock() should succeed")
	tst.AssertFalse(t, lock.IsLocked(), "lock should not be locked after Unlock()")
}

func TestLockTwice(t *testing.T) {
	tmpDir := t.TempDir()
	lockPath := filepath.Join(tmpDir, "test.lock")

	lock := filelock.New(lockPath)

	// Acquire the lock first time
	err := lock.Lock()
	tst.AssertNoError(t, err, "First Lock() should succeed")

	// Try to lock again from same process
	err = lock.Lock()
	tst.AssertNotNil(t, err, "Second Lock() from same process should fail")

	// Clean up
	_ = lock.Unlock()
}

func TestUnlockWithoutLock(t *testing.T) {
	lock := filelock.New("/tmp/test.lock")

	// Try to unlock without locking
	err := lock.Unlock()
	tst.AssertNotNil(t, err, "Unlock() without Lock() should fail")
}

func TestTryLock(t *testing.T) {
	tmpDir := t.TempDir()
	lockPath := filepath.Join(tmpDir, "test.lock")

	lock := filelock.New(lockPath)

	// Try to acquire the lock (should succeed)
	acquired, err := lock.TryLock()
	tst.AssertNoError(t, err, "TryLock() should not error")
	tst.AssertTrue(t, acquired, "TryLock() should acquire lock on first attempt")
	tst.AssertTrue(t, lock.IsLocked(), "lock should be locked after TryLock()")

	// Clean up
	_ = lock.Unlock()
}

func TestTryLockTwice(t *testing.T) {
	tmpDir := t.TempDir()
	lockPath := filepath.Join(tmpDir, "test.lock")

	lock := filelock.New(lockPath)

	// Try to acquire the lock first time
	acquired, err := lock.TryLock()
	tst.AssertNoError(t, err, "First TryLock() should not error")
	tst.AssertTrue(t, acquired, "First TryLock() should succeed")

	// Try to acquire again from same process
	_, err = lock.TryLock()
	tst.AssertNotNil(t, err, "Second TryLock() from same process should error")

	// Clean up
	_ = lock.Unlock()
}

func TestMultipleInstances(t *testing.T) {
	tmpDir := t.TempDir()
	lockPath := filepath.Join(tmpDir, "test.lock")

	lock1 := filelock.New(lockPath)
	lock2 := filelock.New(lockPath)

	// First instance acquires the lock
	err := lock1.Lock()
	tst.AssertNoError(t, err, "First lock should succeed")

	// Second instance tries to acquire (should fail in non-blocking mode)
	acquired, err := lock2.TryLock()
	tst.AssertNoError(t, err, "TryLock should not error")
	tst.AssertFalse(t, acquired, "Second instance should not acquire locked file")

	// Release first lock
	err = lock1.Unlock()
	tst.AssertNoError(t, err, "Unlock should succeed")

	// Now second instance should be able to acquire
	acquired, err = lock2.TryLock()
	tst.AssertNoError(t, err, "TryLock should not error after release")
	tst.AssertTrue(t, acquired, "Second instance should acquire after release")

	// Clean up
	_ = lock2.Unlock()
}

func TestLockCreatesFile(t *testing.T) {
	tmpDir := t.TempDir()
	lockPath := filepath.Join(tmpDir, "test.lock")

	// Ensure the file doesn't exist yet
	_, err := os.Stat(lockPath)
	tst.AssertTrue(t, os.IsNotExist(err), "Lock file should not exist initially")

	lock := filelock.New(lockPath)
	err = lock.Lock()
	tst.AssertNoError(t, err, "Lock should succeed")

	// File should now exist
	_, err = os.Stat(lockPath)
	tst.AssertNoError(t, err, "Lock file should exist after Lock()")

	_ = lock.Unlock()
}

func TestLockWithExistingFile(t *testing.T) {
	tmpDir := t.TempDir()
	lockPath := filepath.Join(tmpDir, "test.lock")

	// Create the file beforehand
	f, err := os.Create(lockPath)
	tst.AssertNoError(t, err, "Create should succeed")
	_ = f.Close()

	lock := filelock.New(lockPath)
	err = lock.Lock()
	tst.AssertNoError(t, err, "Lock should succeed on existing file")

	_ = lock.Lock()
}

func TestString(t *testing.T) {
	lock := filelock.New("/tmp/test.lock")
	str := lock.String()
	tst.AssertTrue(t, len(str) > 0, "String() should return non-empty string")

	_ = lock.Lock()
	strLocked := lock.String()
	tst.AssertTrue(t, len(strLocked) > 0, "String() should return non-empty string when locked")
	_ = lock.Unlock()
}

func TestConcurrentLocking(t *testing.T) {
	tmpDir := t.TempDir()
	lockPath := filepath.Join(tmpDir, "test.lock")

	lock1 := filelock.New(lockPath)
	lock2 := filelock.New(lockPath)

	// Channel to signal when lock2 acquires the lock
	acquired := make(chan bool)

	// Lock with lock1
	err := lock1.Lock()
	tst.AssertNoError(t, err, "First lock should succeed")

	// Goroutine tries to lock with lock2 (should block)
	go func() {
		acquired <- true // Signal that we're about to try
		err := lock2.Lock()
		if err != nil {
			t.Logf("Lock2 failed: %v", err)
			acquired <- false
			return
		}
		acquired <- true
		_ = lock2.Unlock()
	}()

	// Wait for goroutine to start trying
	<-acquired

	// Give goroutine time to attempt lock (should be blocked)
	time.Sleep(100 * time.Millisecond)

	// Release lock1
	err = lock1.Unlock()
	tst.AssertNoError(t, err, "Unlock should succeed")

	// Now goroutine should acquire the lock
	result := <-acquired
	tst.AssertTrue(t, result, "Goroutine should acquire lock after release")
}
