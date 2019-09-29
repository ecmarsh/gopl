# 2. Program Structure

## Predeclared Names

### Keywords

- `break`
- `case`
- `chan`
- `const`
- `continue`
- `default`
- `defer`
- `else`
- `fallthrough`
- `for`
- `func`
- `go`
- `goto`
- `if`
- `import`
- `interface`
- `map`
- `package`
- `range`
- `return`
- `select`
- `struct`
- `switch`
- `type`
- `var`

### Constants

- `true`
- `false`
- `iota`
- `nil`

### Types

- `int`
- `int8`
- `int16`
- `int32`
- `int64`
- `unit`
- `unit8`
- `unit16`
- `unit32`
- `uint64`
- `uintptr`
- `float32`
- `float64`
- `complex128`
- `complex64`
- `bool`
- `byte`
- `rune`
- `string`
- `error`

### Functions

- `make`
- `len`
- `cap`
- `new`
- `append`
- `copy`
- `close`
- `delete`
- `complex`
- `real`
- `imag`
- `panic`
- `recover`

## Variables

- There are no uninitialized variables (no assignment defaults to zero-value).
- Note do not confuse short, multi-variable declarations with tuple assignments (`:=` vs `=`)
- Short multi-variable assignments will assign to variables already declared. (scope is important)
  - but a short variable dec must have at least one declaration for this to work.
  - so just use ordinary assignment (`=`) if no new declarations
- Style: avoid tuple assignments if complicated expressions/hinders readability

### Pointers

- Pointers and addresses (* and &) same as in C lang. Except cannot do pointer arithmetic.

### Var Lifetimes

- Package-level variable lifetime is entire execution of program
- Local variables have dynamic lifetimes: new instance is created each time declaration statement is executed and lives until out of scope/unreachable

### Assertions

1. Map lookup: `v, ok = m[key]`
2. Type assertion: `v, ok = x.(T)`
3. Channel receive: `v, ok = <-ch`

### Assignability


## Types

- Assignability and comparisons vary by type.
- type declarations: `type name underlying-type`

## Packages & Files

- Each package serves as a separate name space for its declaration.
- Exported (public) identifiers start with uppercase letter.
- Doc comments are placed immediately before package declarations.
- By convention, import name should match the leaf of import path.
- Package-level variables are initialized in order they were declared in (note package names sorted before compiling)
- Packages resolve dependencies before importing. (if p imports q, q must be fully resolved first).
- `init` functions in packages are executed when the program starts. Thus, all packages init functions are run before `main`. 

## Scope

- Lexical blocks (`{}`)
- Style: handle error in if block, so successful execution path is not indented.










