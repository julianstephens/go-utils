## v0.4.2

- logger: improves rotating file output robustness and adds context integration
    - makes `MaxBackups` and `MaxAge` optional constraints - clients can use one, both, or neither
    - adds `Close()` method for proper resource cleanup on both Logger instances and global logger
    - closes previous file output when `SetFileOutput()` or `SetFileOutputWithConfig()` is called again (prevents file handle leaks)
    - moves default rotation values to package-level constants for clarity and reusability
    - comprehensive test coverage for file rotation with various constraint combinations
    - adds `SetOutput(io.Writer)` method that automatically closes previous file output to prevent leaks
    - adds `Sync()` method to flush pending logs to disk during graceful shutdown
    - adds `WithContext(ctx context.Context)` for structured logging with request/trace IDs
        - automatically extracts trace ID, request ID, user ID from context with multiple key format support
        - compatible with common context patterns (trace-id, traceId, traceID, etc.)
    - improves `Close()` to call `Sync()` before closing to ensure all logs are written
    - adds comprehensive initialization patterns and context integration documentation
    - wraps errors with context in `SetLogLevel()`, `SetFileOutput()`, and `SetFileOutputWithConfig()` for better debugging
    - validates configuration parameters (non-empty filepath, positive MaxSize and MaxAge values, non-negative MaxBackups)
    - adds `SafeLog()` method to recover from panics in logging operations for critical code paths
    - adds comprehensive graceful shutdown patterns in documentation with HTTP server examples
    - improved error handling throughout configuration methods
- validator: adds production-ready improvements for reliability and performance
    - adds `SafeValidate()` method to `CustomValidator` to recover from panics in validation chains
    - prevents validation errors from crashing the application in critical code paths
    - adds comprehensive performance documentation with benchmarks and optimization tips
    - includes performance tips for validator ordering and reuse patterns
    - documents efficiency characteristics (zero allocations for most validators, lazy validation)
    - adds safe validation example in documentation showing panic recovery usage

## v0.4.1

- logger: adds rotating file output support
    - adds `FileRotationConfig` type for configurable rotation parameters
    - adds `SetFileOutput(filepath)` for rotating file logs with sensible defaults (100MB, 3 backups, 28 days, compression enabled)
    - adds `SetFileOutputWithConfig(config)` for custom rotation configuration
    - adds support for global and instance-level file output configuration
    - integrates [lumberjack](https://github.com/natefinch/lumberjack) for robust log rotation
- improves code formatting consistency across validator and test packages

## v0.4.0

- checksum: adds fast cryptographic checksum utilities
    - adds `CRC32C()` for Castagnoli checksums optimized for storage
    - adds `CRC32()` for IEEE and Koopman variants
    - adds `Verify()` and `VerifyWithAlgorithm()` for data integrity checks
    - adds `AppendCRC32C()` and `StripCRC32C()` for self-checksumming
    - adds streaming support via `hash.Hash32` interface
- filelock: adds cross-platform file locking utilities
    - adds `Lock()` for exclusive file locks across processes
    - adds `TryLock()` for non-blocking lock attempts
    - adds `Unlock()` for releasing locks
    - supports Linux, macOS, and Windows with unified API
- health: adds health check and diagnostic utilities
    - adds `Checker` interface for custom health checks
    - adds `Repairer` interface for automated repair operations
    - adds standardized exit codes (0=OK, 1=Warning, 2=Error)
    - adds `Report` type with aggregated check results and timestamps
    - adds `RunChecks()` to execute multiple checkers
    - adds `RepairAll()` to attempt repairs on failed components
- helpers: removes deprecated slice functions and consolidates with generic package
    - **BREAKING**: removes `ContainsAll[T comparable]` - use `generic.ContainsAll` instead
    - **BREAKING**: removes `Difference` - use `generic.Difference` instead  
    - **BREAKING**: removes `DeleteElement` - use `generic.DeleteElement` instead
    - adds atomic file operations for crash-safe writes
        - adds `AtomicFileWrite()` for atomic file writes with fsync and rename
        - adds `AtomicFileWriteWithPerm()` for atomic writes with custom permissions
        - adds `SafeFileSync()` for syncing file data to disk
        - adds `SafeDirSync()` for syncing directory to ensure durability
    - removes unused `encoding/json` import
- jsonutil: adds file I/O operations for JSON marshaling/unmarshaling
    - adds `ReadFile()` for reading and unmarshaling JSON files
    - adds `ReadFileStrict()` for strict field matching when reading files
    - adds `ReadFileWithOptions()` for custom unmarshal options
    - adds `WriteFile()` for marshaling and writing JSON files
    - adds `WriteFileIndent()` for indented JSON file writes
    - adds `WriteFileWithOptions()` for custom marshal options
- slices: marks package as deprecated in favor of generic package
    - **DEPRECATED**: Package will be removed in v0.6.0, migrate to `generic` package
    - updates package documentation with deprecation notice and migration guide

## v0.3.1

- validator:
    - **BREAKING**: removes legacy `Validator` type and `New()` function (superceded by factory functions)
    - adds `ValidateMatchesField[T comparable]` for comparing two values (password confirmation)
    - adds `CustomValidator` with fluent builder pattern for composing validators
        - adds `NewCustomValidator()` factory function
        - adds `Add()` method for appending validators with chaining support
        - adds `Validate()` method for executing all validators with AND logic
    - adds `ErrFieldMismatch` error constant
- logger: improves thread-safety for concurrent applications
    - adds `sync.RWMutex` protection for global logger configuration changes
    - `SetLogLevel()`, `SetOutput()`, and `SetFormatter()` now thread-safe during concurrent logging
    - adds comprehensive concurrency tests for multi-threaded scenarios
    - verified with Go race detector for data race detection
- cliutil: fixes race condition in Spinner
    - adds `sync.RWMutex` protection for concurrent access to spinner state
    - `Start()`, `Stop()`, and `UpdateMessage()` now thread-safe
    - verified with Go race detector

## v0.3.0

- **BREAKING**: redesigns `validator` package with generic validators
    - adds `NumberValidator[T]` with 15 comprehensive numeric validation methods
    - adds `StringValidator[T]` with length validation methods
    - adds `ParseValidator` with 12 parsing validation methods
    - adds factory functions `Numbers[T]()`, `Strings[T]()`, `Parse()` for type-safe validator creation
    - adds type constraints `Number`, `StringLike`, `Emptyable` for generic validation
    - adds `ValidationError` type with detailed context and error chaining
    - adds generic `ValidateNonEmpty[T]` supporting strings, bytes, runes, maps, and slices
    - improves test coverage to 92.1% with comprehensive test suite
    - reorganizes tests into dedicated files (`number_test.go`, `string_test.go`, `parse_test.go`)
    - updates documentation with advanced usage examples and API reference
    - maintains backward compatibility for original validation functions

## v0.2.8

- adds MIT license

## v0.2.7

- adds optional custom error messages to all assertion functions in `tests`
- adds generic comparison functions to `tests`: `AssertGreaterThan`, `AssertLessThan`, `AssertGreaterThanOrEqual`, `AssertLessThanOrEqual`, `AssertEqual`

## v0.2.6

- removes `Error.Code` from `httputil.response`
- adds support for additional details in `Error.Details` in `httputil.response`
- adds `ValidateUUID` to `validator`

## v0.2.5

- adds `PromptPassword` and `PromptPasswordWithValidation` to `cliutil` 
- adds `validator`

## v0.2.4

- adds `Exists` to `helpers`
- improves package documentation

## v0.2.3

- adds `security` package

## v0.2.2

- adds support for long-lived refresh tokens to `httputil.auth.JWTManager`
- adds cookie security helpers to `httputil.auth.JWTManager`

## v0.2.1

- adds `TokenExpiration` helpers to `httputil.auth.JWTManager`
- adds `NewEmpty` constructor for responder without hooks
- adds JWT middleware to `httputil.middleware`

## v0.2.0

- adds `tests`
- improves documentation

## v0.1.6

- fixes `httputil.response.ErrorWithStatus`

## v0.1.5

- refactors documentation
- fixes `httputil.response` error statuses
- fixes `dbutil` field mapper names

## v0.1.4

- adds `slices`
- adds `logger`
- adds `config`
- adds `generic`
- adds `jsonutil`
- adds `cliutil`
- adds `dbutil`

## v0.1.3

- adds `httputil.auth.utils`
- adds `httputil.middleware`
- removes unnecessary user identifier claims
- refactors `httputil.auth.GenerateToken` for simpler user identification

## v0.1.2

- fixes superfluous header write in `response.OK`
- adds `helpers.Default`
- adds `httputil.auth`

## v0.1.1

- moves `helpers` to root
- adds `httputil.request`

## v0.1.0

- adds mux response package
- adds common helpers
