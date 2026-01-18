# Helpers Package

The `helpers` package provides general utility functions including conditional helpers, file system utilities, and struct manipulation. For slice-specific operations, see the [slices](../slices) package.

## Features

- **Conditional Helpers**: Ternary operator and default value handling
- **File Operations**: Filesystem checks and atomic writes
- **Atomic Writes**: Crash-safe file operations with sync and rename
- **Type Utilities**: Pointer and struct conversion helpers

## Installation

```bash
go get github.com/julianstephens/go-utils/helpers
```

## Usage

### Conditional Helpers

```go
// Ternary operator
age := 25
status := helpers.If(age >= 18, "adult", "minor")  // "adult"

// Default value handling - returns default if val is zero value
emptyStr := ""
result := helpers.Default(emptyStr, "default")  // "default"
zeroInt := 0
intResult := helpers.Default(zeroInt, 42)  // 42
```

### File Operations

```go
// Check existence
exists := helpers.Exists("config.json")  // true/false

// Create file or directory if it doesn't exist
if err := helpers.Ensure("data.json", false); err != nil {  // file
    log.Fatal(err)
}
if err := helpers.Ensure("./logs", true); err != nil {  // directory
    log.Fatal(err)
}
```

### Generic Types

All generic functions work with custom types:

```go
type Person struct {
    Name string
    Age  int
}

// If works with any type
people := []Person{{Name: "Alice", Age: 30}, {Name: "Bob", Age: 25}}
older := helpers.If(people[0].Age > people[1].Age, people[0], people[1])

// Default recognizes zero values for all types
var p Person  // zero value
defaultPerson := helpers.Default(p, Person{Name: "Unknown", Age: 0})
```

## API Reference

### Conditional Helpers
- `If[T any](cond bool, vtrue T, vfalse T) T` - Ternary operator: returns vtrue if cond is true, otherwise vfalse
- `Default[T any](val T, defaultVal T) T` - Returns defaultVal if val is the zero value for its type, otherwise returns val

### File Operations
- `Exists(path string) bool` - Check if a file or directory exists at the given path
- `Ensure(path string, isDir bool) error` - Create a file or directory if it doesn't exist
- `ReadJSONFile(filename string, v interface{}) error` - Read JSON file and unmarshal into provided struct
- `WriteJSONFile(filename string, v interface{}) error` - Marshal struct to JSON and write to file

### Atomic File Operations
- `AtomicFileWrite(path string, data []byte) error` - Write data to file atomically with default permissions (0666)
- `AtomicFileWriteWithPerm(path string, data []byte, perm os.FileMode) error` - Atomic write with custom permissions
- `SafeFileSync(f *os.File) error` - Sync file data to disk
- `SafeDirSync(dir string) error` - Sync directory to ensure durability

### Type Utilities
- `StringPtr(s string) *string` - Return a pointer to the given string
- `StructToMap(obj any) map[string]any` - Convert struct to map using reflection

## Zero Values

The `Default` function recognizes zero values for all Go types: numbers (`0`, `0.0`), strings (`""`), booleans (`false`), pointers (`nil`), slices/maps/channels (`nil`), and structs with zero-valued fields.

### Atomic File Operations

Crash-safe file writes using temp file + rename pattern:

```go
// Atomic write (crash-safe)
data := []byte("critical data")
if err := helpers.AtomicFileWrite("record.json", data); err != nil {
    log.Fatal(err)
}

// With custom permissions
if err := helpers.AtomicFileWriteWithPerm("config.json", data, 0600); err != nil {
    log.Fatal(err)
}

// Manual sync for existing files
helpers.SafeFileSync(f)     // Sync file data
helpers.SafeDirSync(".")    // Sync directory (ensures durability)
```

Why atomic writes? All-or-nothing guarantee: process crash yields either old or new file, never partial writes.

## Thread Safety

- `If`, `Default` are thread-safe (pure functions)
- Atomic file operations are safe to call concurrently (each call is independent)
- `Exists`/`Ensure` are thread-safe for individual calls

## Best Practices

- **Use atomic writes** for critical files (configs, WAL records, manifests)
- **Temp + rename > truncate + write** for durability guarantees
- **Check sync errors** (disk full, permissions)
- **Use custom permissions** for sensitive files (0600 for secrets)
- **Verify zero values** with `Default` - sometimes zero is valid

## Integration

Works seamlessly with other go-utils packages:

```go
// With config package
port := helpers.Default(cfg.Server.Port, 8080)

// With filelock for exclusive writes
lock := filelock.New("mydata.lock")
lock.Lock()
helpers.AtomicFileWrite("data.json", data)
lock.Unlock()

// With checksum and jsonutil for integrity
record := checksum.AppendCRC32C(data)
if err := jsonutil.WriteFile("record.bin", record); err != nil {
    log.Fatal(err)
}

// With generic for collections
if generic.ContainsAll(features, required) {
    helpers.AtomicFileWrite("manifest.json", manifestData)
}
```

## Related Packages

- **[jsonutil](../jsonutil)** - JSON marshaling, unmarshaling, and file I/O
- **[generic](../generic)** - Comprehensive generic utilities for functional programming, slice/map operations
- **[filelock](../filelock)** - Cross-platform file locking
- **[checksum](../checksum)** - Fast CRC32 checksums for data integrity