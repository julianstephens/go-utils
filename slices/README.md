# Slices Package (Deprecated)

⚠️ **This package is deprecated.** Use the [generic](../generic) package instead, which provides all functionality from slices plus comprehensive functional programming utilities, advanced slice operations, and map utilities.

## Migration Guide

| slices Function | generic Equivalent |
|---|---|
| `If[T any]` | `generic.If[T any]` |
| `Difference` | `generic.Difference` |
| `DeleteElement` | `generic.DeleteElement` |
| `ContainsAll` | `generic.ContainsAll` |

See the [generic package](../generic) README for comprehensive documentation and advanced usage patterns.

---

The `slices` package provides generic slice utility functions for conditional selection, set operations, and element manipulation that are commonly needed across Go projects.

## Features

- **Conditional Logic**: Ternary operator implementation
- **Set Operations**: Slice subtraction (difference)
- **Element Manipulation**: Safe element deletion by index
- **Subset Validation**: Check if all elements in one slice exist in another

## Installation

```bash
go get github.com/julianstephens/go-utils/slices
```

## Usage

### Conditional Selection

```go
package main

import (
    "fmt"
    "github.com/julianstephens/go-utils/slices"
)

func main() {
    // If function mimics the ternary operator: cond ? vtrue : vfalse
    age := 25
    status := slices.If(age >= 18, "adult", "minor")
    fmt.Printf("Status: %s\n", status) // "adult"
    
    // Works with any type
    score := 85
    passed := slices.If(score >= 60, true, false)
    fmt.Printf("Passed: %t\n", passed) // true
    
    // Nested conditionals
    grade := slices.If(score >= 90, "A", 
                slices.If(score >= 80, "B", 
                    slices.If(score >= 70, "C", "F")))
    fmt.Printf("Grade: %s\n", grade) // "B"
}
```

### Slice Difference

```go
package main

import (
    "fmt"
    "github.com/julianstephens/go-utils/slices"
)

func main() {
    // Difference returns elements in slice a that are not in slice b
    allUsers := []string{"alice", "bob", "charlie", "diana", "eve"}
    activeUsers := []string{"alice", "charlie", "eve"}
    
    inactiveUsers := slices.Difference(allUsers, activeUsers)
    fmt.Printf("Inactive users: %v\n", inactiveUsers) // ["bob", "diana"]
    
    // Order is preserved from the first slice
    permissions := []string{"read", "write", "delete", "admin"}
    userPermissions := []string{"read", "admin"}
    
    missingPermissions := slices.Difference(permissions, userPermissions)
    fmt.Printf("Missing permissions: %v\n", missingPermissions) // ["write", "delete"]
    
    // Empty result when all elements are present
    subset := []string{"alice", "bob"}
    superset := []string{"alice", "bob", "charlie"}
    
    result := slices.Difference(subset, superset)
    fmt.Printf("Difference result: %v\n", result) // []
}
```

### Element Deletion

```go
package main

import (
    "fmt"
    "github.com/julianstephens/go-utils/slices"
)

func main() {
    // DeleteElement removes an element at the specified index
    numbers := []int{10, 20, 30, 40, 50}
    
    // Remove element at index 2 (value 30)
    result := slices.DeleteElement(numbers, 2)
    fmt.Printf("After deletion: %v\n", result) // [10, 20, 40, 50]
    
    // Remove first element
    fruits := []string{"apple", "banana", "cherry", "date"}
    result2 := slices.DeleteElement(fruits, 0)
    fmt.Printf("After removing first: %v\n", result2) // ["banana", "cherry", "date"]
    
    // Remove last element
    colors := []string{"red", "green", "blue"}
    result3 := slices.DeleteElement(colors, len(colors)-1)
    fmt.Printf("After removing last: %v\n", result3) // ["red", "green"]
}
```

### Subset Validation

```go
package main

import (
    "fmt"
    "github.com/julianstephens/go-utils/slices"
)

func main() {
    // ContainsAll checks if all elements in subset are present in mainSlice
    availableFeatures := []string{"auth", "logging", "metrics", "caching", "monitoring"}
    
    // Check required features
    requiredFeatures := []string{"auth", "logging"}
    hasRequired := slices.ContainsAll(availableFeatures, requiredFeatures)
    fmt.Printf("Has required features: %t\n", hasRequired) // true
    
    // Check optional features
    optionalFeatures := []string{"metrics", "analytics"}
    hasOptional := slices.ContainsAll(availableFeatures, optionalFeatures)
    fmt.Printf("Has optional features: %t\n", hasOptional) // false
    
    // Empty subset is always contained
    emptyFeatures := []string{}
    hasEmpty := slices.ContainsAll(availableFeatures, emptyFeatures)
    fmt.Printf("Has empty subset: %t\n", hasEmpty) // true
    
    // Works with numbers too
    availableNumbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    requiredNumbers := []int{2, 4, 6}
    hasNumbers := slices.ContainsAll(availableNumbers, requiredNumbers)
    fmt.Printf("Has required numbers: %t\n", hasNumbers) // true
}
```

### Practical Examples

```go
package main

import (
    "fmt"
    "github.com/julianstephens/go-utils/slices"
)

func filterUsersByRole(users []string, admins []string) (regularUsers []string) {
    // Find users who are not admins
    return slices.Difference(users, admins)
}

func validatePermissions(userPerms, requiredPerms []string) (bool, []string) {
    // Check if user has all required permissions
    hasAll := slices.ContainsAll(userPerms, requiredPerms)
    
    if hasAll {
        return true, nil
    }
    
    // Return missing permissions
    missing := slices.Difference(requiredPerms, userPerms)
    return false, missing
}

func removeUserFromList(users []string, userToRemove string) []string {
    // Find the user and remove them
    for i, user := range users {
        if user == userToRemove {
            return slices.DeleteElement(users, i)
        }
    }
    return users // User not found, return original slice
}

func getUserStatus(user string, activeUsers []string) string {
    // Use If for conditional status determination
    return slices.If(slices.ContainsAll([]string{user}, activeUsers), "active", "inactive")
}

func main() {
    allUsers := []string{"alice", "bob", "charlie", "diana", "eve"}
    adminUsers := []string{"alice", "eve"}
    activeUsers := []string{"alice", "bob", "charlie"}
    
    // Filter regular users
    regularUsers := filterUsersByRole(allUsers, adminUsers)
    fmt.Printf("Regular users: %v\n", regularUsers)
    
    // Validate permissions
    userPermissions := []string{"read", "write"}
    requiredPermissions := []string{"read", "write", "delete"}
    
    hasAll, missing := validatePermissions(userPermissions, requiredPermissions)
    fmt.Printf("Has all permissions: %t\n", hasAll)
    if !hasAll {
        fmt.Printf("Missing permissions: %v\n", missing)
    }
    
    // Remove a user
    updatedUsers := removeUserFromList(allUsers, "charlie")
    fmt.Printf("After removing charlie: %v\n", updatedUsers)
    
    // Check user status
    for _, user := range allUsers {
        status := getUserStatus(user, activeUsers)
        fmt.Printf("User %s is %s\n", user, status)
    }
}
```

### Advanced Usage with Custom Types

```go
package main

import (
    "fmt"
    "github.com/julianstephens/go-utils/slices"
)

type User struct {
    ID   int
    Name string
    Role string
}

func main() {
    users := []User{
        {ID: 1, Name: "Alice", Role: "admin"},
        {ID: 2, Name: "Bob", Role: "user"},
        {ID: 3, Name: "Charlie", Role: "user"},
        {ID: 4, Name: "Diana", Role: "admin"},
    }
    
    // Extract user names
    userNames := make([]string, len(users))
    for i, user := range users {
        userNames[i] = user.Name
    }
    
    // Use conditional logic with custom types
    alice := users[0]
    status := slices.If(alice.Role == "admin", "Administrator", "Regular User")
    fmt.Printf("%s is a %s\n", alice.Name, status)
    
    // Find difference in user names
    targetUsers := []string{"Alice", "Bob"}
    otherUsers := slices.Difference(userNames, targetUsers)
    fmt.Printf("Other users: %v\n", otherUsers)
    
    // Remove user by index
    usersAfterRemoval := slices.DeleteElement(users, 1) // Remove Bob
    fmt.Printf("Users after removal: %+v\n", usersAfterRemoval)
}
```

## API Reference

### Functions

- `If[T any](cond bool, vtrue T, vfalse T) T`
  - Ternary operator implementation
  - Returns `vtrue` if `cond` is true, otherwise returns `vfalse`
  - Works with any type

- `Difference(a []string, b []string) []string`
  - Returns elements in slice `a` that are not in slice `b`
  - Preserves order from slice `a`
  - Returns empty slice if all elements of `a` are in `b`

- `DeleteElement[T any](slice []T, index int) []T`
  - Removes element at specified index
  - Returns new slice without the element
  - **Note**: Does not check bounds - ensure index is valid

- `ContainsAll[T comparable](mainSlice, subset []T) bool`
  - Returns true if all elements in `subset` are present in `mainSlice`
  - Empty subset always returns true
  - Works with any comparable type

## Type Constraints

### Generic Types
- `If[T any]` - Works with any type
- `DeleteElement[T any]` - Works with any type
- `ContainsAll[T comparable]` - Works with comparable types only

### Comparable Types
Types that can be used with `ContainsAll`:
- Basic types: `string`, `int`, `float64`, `bool`, etc.
- Arrays with comparable elements
- Structs where all fields are comparable
- **Not**: slices, maps, functions, channels

## Performance Considerations

- `Difference`: O(n + m) time complexity, O(m) space complexity
- `ContainsAll`: O(n + m) time complexity, O(n) space complexity  
- `DeleteElement`: O(n) time complexity for copying elements
- `If`: O(1) constant time

## Safety Notes

1. **Index Bounds**: `DeleteElement` does not check array bounds. Ensure the index is valid:
   ```go
   if index >= 0 && index < len(slice) {
       result = slices.DeleteElement(slice, index)
   }
   ```

2. **Slice Modification**: `DeleteElement` returns a new slice; it does not modify the original:
   ```go
   original := []int{1, 2, 3}
   modified := slices.DeleteElement(original, 1)
   // original is still [1, 2, 3]
   // modified is [1, 3]
   ```

## Best Practices

1. **Use `If` for simple conditionals** instead of if-else blocks
2. **Check bounds before using `DeleteElement`**
3. **Use `Difference` for set-like operations** when you need elements unique to one slice
4. **Use `ContainsAll` for validation** to ensure required elements are present
5. **Consider performance implications** for large slices

## Integration

Works well with other go-utils packages:

```go
// Use with helpers package
result := slices.If(condition, helpers.Default(value, fallback), alternative)

// Use with generic package for more complex operations
filtered := generic.Filter(users, func(u User) bool {
    return slices.If(u.Active, u.Role == "admin", false)
})
```