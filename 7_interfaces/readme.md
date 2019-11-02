# Interfaces

- Interface types express generalizations or abstractions about the behaviors of other types.
- This allows us to write functions that are more flexible and adaptable because they are not tied to the details of one particular implementation.
- What makes Go's interfaces stand out is that they are _satisfied implicitly_, meaning there's no need to declare all the interfaces that a given concrete type satisfies; possessing the necessary methods is enough.
- This allows us to create new interfaces that are satisfied by existing concrete types without changing the existing types (useful for packages you don't control).

## Interfaces as Contracts

- Previous types explored are _concrete_, which specifies the exact representation of its values.
- Interface types are _abstract_; it doesn't expose representation or internal structure of its values or set of operations they support; it only reveals some of their methods.
- An interface types only allows you to know what it _can do_, or what behaviors are provided by its methods.
- The interface defines a contract between concrete type and callers.
  - Caller is required to provide a value of the concrete type with a method with appropriate signature.
  - Guarantees that method will do its job if satisfies interface.
- Interfaces offer _substitutability_, the freedom to substitute one type for another that satisfies the same interface.

## Interface Types

- An interface type specifies a set of methods that a concrete type must possess to be considered an instance of that interface.
- Example, the `io` package contains many interfaces:
```go
// Represents anything in which you can read bytes from
type Reader interface {
  Read(p []byte) (n int, err error)
}
// Represents anything you can close.
type Closer interface {
  Close() error
}
```
- Interfaces can also be declared as combinations of existing ones, called embedding an interface. Note order does _not_ matter for combined declarations.
```go
type ReadWriter interface {
  Reader
  Writer
}
type ReadWriteCloser interface {
  Reader
  Writer
  Closer
}
```

## Interface Satisfaction

- A type satisfies an interface if it possesses **all** the methods the interface requires.
- As short hand, terminology in Go is that a concrete type "is a" particular interface type if it satisfies the interface.
- An expression may be _assigned_ to an interface only if its type satisfies the interface:
```go
var w io.Writer
w = Os.Stdout         // OK: *os.File has a Write method
w = new(bytes.Buffer) // OK: *bytes.Buffer has Write method
w = time.Second       // compile error: time.Duration lacks Write method

var rwc io.ReadWriteCloser
rwc = os.Stdout         // OK: *os.File has Read, Write, Close methods
rwc = new(bytes.Buffer) // compile error: *bytes.Buffer lacks Close method
```
- The interface assignability rule applies even when right-hadn side is itself an interface:
```go
w = rwc  // OK: io.ReadWriteCloser has Write method
rwc = w  // compile error: io.Writer lacks Close method
```
- Although it is legal to call a *T method on an argument type T (as long as the argument is a variable) as Go implicitly takes the address, type T still does not have all the methods that a *T pointer has, so may satisfy fewer interfaces:
```go
type IntSet struct { /* ... */ }
func (*IntSet) String() string
var _ = IntSet{}.String() // compile error: String requires *IntSet receiver
// But if IntSet has an addressable value, we can call it
var s IntSet
var _ = s.String() // OK: s is a variable and &s has a String method
// Only &s will satisfy the Stringer interface though
var _ fmt.Stringer = &s // OK
var _ fmt.Stringer = s  // compile error: IntSet lacks String method (*IntSet) does
```
- An interface wraps and conceals the concrete type and value that it holds. Only methods revealed by the interface type may be called, even if the concrete type has others:
```go
os.Stdout.Write([]byte("hello")) // OK: *os.File has Write method
os.Stdout.Close()                // OK: *os.File has Close method

var w io.Writer
w = os.Stdout
w.Write([]byte("hello")) // OK: io.Writer has Write method
w.Close()                // compile error: io.Writer lacks Close method
```
- Interfaces with more methods place more demands on the types that implement it.
  - Thus, anything can satisfy the empty interface type: `type interface{}`
  - Examples of empty interfaces include `fmt.Println`, `errorf`, etc. which take any type.
  - Remember that to get the value back out though, we must use _type assertsions_.
- Since it is not necessary to declare the relationship between a concrete type and interface that satisfies it, it is occasionally useful to document and assert the relationship when it is intended but not otherwise enforced by the program:
```go
/* explicit declaration */
// *bytes.Buffer must satisfy io.Writer
var w io.Writer = new(bytes.Buffer)
/* same thing, but no need for w declaration */
// *bytes.Buffer must satisfy io.Writer
var _ io.Writer = (*bytes.Buffer)(nil)
```
- A concrete type may satisfy many unrelated interfaces, but it is sometimes useful to express abstractinos as an interface by grouping common methods to types.
- Unlike class-based languages where the set of interfaces satisfied by a class is explicit, in Go we can define new abstractions or groupings of interest when we need them, without modifying the declaration of the concrete type.
  - This becomes extremely handy when the concrete type comes from a package written by a different author.

## Parsing Flags with `flag.Value`

- `flag.Value` is a standard interface that helps define new notations for command-line flags:
```go
// sleep prints the time period and then sleeps for that time period
var period = flag.Duration("period", 1*time.Second, "sleep period")

func main() {
  flag.Parse()
  fmt.Printf("Sleeping for %v...", *period)
  time.Sleep(*period)
  fmt.Println()
}
// $ go build ./sleep
// $ ./sleep
// Sleeping for 1s...
// $ ./sleep -period 50ms
// Sleeping for 50ms...
```
- It's easy to define our own flag notations/data types by defining a type that satisfies the `flag.Value` interface:
```go
package flag

// Value is the interface to the value stored in a flag.
type Value interface {
  String() string
  Set(string) error
}
```
- See [tempflag](./tempflag) for example of satisfying `flag.Value` interface by defining necessary methods.

## Interface Values

- An interface type, or _interface value_, has two components:
  - concrete type, called the interface's _dynamic type_
  - a value of concrete type, called the interface's _dynamic value_
- In go (and all statically typed languages), types are a compile-time concept, so a type is not a value.
  - A set of values called _type descriptors_ provide more information about each type, such as names and methods.
  - In an interface value, the type component is represented by the appropriate type descriptor.
- Zero-value for interface is both type and value components `nil`.
- Upon assignment, interface value's dynamic type are set to the type descriptor.
- Calls through interfaces must use _dynamic dispatch_ since we cannot know at compile time what the dynamic type of an interface value will be.
  - Instead of a direct call, compiler must generate code to obtain address of the method from the type descriptor, and then make an indirect call to that address, where the receiver is a copy of the interface's dynamic value.
- Interface values can hold arbitrarily large dynamic values.
- Interface values can be compared, and will be equal if both are nil or if their dynamic types are identical and their dynamic values are equal according to the behavior for the dynamic type.
- Comparability creates possibility of using interfaces as keys in maps, switch statement operand, etc.
- NOTE: If an interface's dynamic type is not comparable (eg slice, maps, funcs), this causes a panic, so only use as keys,operands, etc. appropriately.
  - For debugging, handling errors, can be helpful to print the interface's dynamic type, using `%T` verb. (`fmt` uses reflection to obtain the name of the interface's dynamic type)

### Gotchas: An Interface Containing a Nil Pointer is Non-Nil

- A nil interface value, which contains no value at all, is not the same as an interface value containing a pointer that happens to be nil.
- For example:
```go
const debug = true
// main collects the output of the function f in a bytes.Buffer
func main() {
  var buf *bytes.Buffer
  if debug {
    buf = new(bytes.Buffer) // enable collection of output
  }
  f(buf) // NOTE: subtly incorrect!
  if (debug) {
    // ...use buf...
  }
}
// If out is non-nil, output will be written to it.
func f(out io.Writer) {
  // ...do something...
  if out != nil {
    out.Write([]byte("done!\n"))
  }
  // If debug is changed to false, panics here because
  // buf's value is nil, but its dynamic type is still *bytes.Buffer
  // so out is not nil and tries to write to a nil value.
}
```
- Note solution to above example is to change the type of `buf` in `main` to `io.Writer` to avoid assignment of dysfunctional value to the interface; `(*bytes.Buffer).Write` has implicit precondition that receiver is not nil so shouldn't have assigned nil pointer.

## Important Interfaces in Go's STL

### Sorting with `sort.Interface`

[https://golang.org/pkg/sort](https://golang.org/pkg/sort)

- `sort` package provides in-place sorting of any sequence according to any ordering function.
- An in-place sort algorithm needs the length, a comparison function, and a swap function:
```go
// from package sort
type Interface interface {
  Len() int
  Less(i, j int) bool // i, j are indices of sequence elements
  Swap(i, j int)
}
// Example fulfillment - note structs are a common type for custom sorts
type StringSlice []string
func (p StringSlice) Len() int { return len(p) }
func (p StringSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p StringSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
// Usage
sort.Sort(StringSlice(names))
// included go package for this equivalent is sort.Strings(s []string)
```
- Keep in mind `sort.Interface` can be adopted to other uses, such as `IsPalindrome(s sort.Interface) bool` with equality comparison if `!s.Less(i, j) && !s.Less(j, i)` 

### The `http.Handler` Interface

[https://golang.org/pkg/net/http/#Handler](https://golang.org/pkg/net/http/#Handler)

```go
// net/http
package http

type Handler interface {
  ServeHTTP(w ResponseWriter, r *Request)
}
// ListenAndServe requires a server address, such as "localhost:8000",
// and an instance of the Handler interface to dispatch all requests to
// It runs forever, or until the server fails with an error.
// Upon failure, the error is always non-nil and returned.
func ListenAndServe(address string, h Handler) error
```

- See [`http`](./http) for examples of use case with `http.Handler` interface.
- To simply the association between URLs and handlers, `net/http` also provides `ServeMux`, a request multiplexer which aggregates a collection of `http.Handlers` into a single `http.Handler`, using the fact that different types satisfying the same interface are _substitutable_; the web server can dispatch requests to any `http.Handler`, regardless of which concrete type is behind it.
  - See [example with ServeMux](./http/http3).
- Go doesn't provide frameworks analogous to Rails/Django, etc. but STL is flexible enough that frameworks are often unnecessary.
- Remember the web server invokes each handler in a new goroutine, so handlers must take precautions such as _locking_ when accessing variables that other goroutines, including other requests to the same handler, may be accessing. See [concurrency section](../8_concurrency/) for more.

### The `error` Interface

[https://golang.org/pkg/errors/](https://golang.org/pkg/errors/)

- The `error` type is just an interface type with a single method on it that returns an error message:
```go
type error interface {
  Error() string
}
``` 
- The simplest way to create an error is by calling `errors.New`, which returns a new `error` for a given error message. Here's the rest of the `errors` package:
```go
package errors

func New(text string) error { return &errorString{text} }

type errorString struct { text string }

func (e *errorString) Error() string { return e.text }
```
- It's underlying type is a struct as opposed to string to protect its representation from inadvertent (or premeditated) updates.
- The reason that pointer type `*errorString`, not `errorString` alone, satisfies the error interface is so that every call to `New` allocates a distinct error instance that is equal to no other.
- Note that calls to `errors.New` are relatively infrequent because we use `Errorf` as a wrapper function:
```go
package fmt

import "errors"

func errorF(format string, args ...interface{}) error {
  return errors.New(Sprintf(format, args...))
}
```
- The [syscall package](https://golang.org/pkg/syscall/) provides many different types of errors, including `Errno` for defining errors numerically as lookup keys; these map to POSIX errors.
- We can discriminate errors using type assertions.

## Type Assertions

Type assertions check that the dynamic type of its operand matches the asserted type.

```go
x.(T)
// x is expression of interface type
// t is a type, called the asserted type
```

- If T is concrete, then checks whether x's dynamic type is identical to T.
  - If the check fails, type assertion causes a panic.
- If T is an interface type, type assertion checks whether x's dyanmic type satisfies T.
  - A type assertion to an interface type changes the type of the expression, making a different (and usually larger) set of methods accessible, while preserving the dynamic type and value components inside the interface value.
- Example:
```go
var w io.Writer
w = os.Stdout
rw := w.(io.ReadWriter) // success: *os.File has both Read and Write

w = new(ByteCounter)
rw = w.(io.ReadWriter) // panic: *ByteCounter has no Read method
```
- Type assertion always fail if operand is a nil interface value.
- In order to not cause a panic, use in an assignment to handle:
```go
var w io.Writer = os.Stdout
f, ok := w.(*os.File)      // success:  ok, b == os.Stdout
b, ok := w.(*bytes.Buffer) // failure: !ok, b == nil
// Often it is used to decide next steps...
// Note it is common to use the original variable name rather than a new name
if w, ok := w.(*os.File); ok {
  // ...use w as *os.File...
}
```

### Discriminating Errors with Type Assertions

- Discriminating error types allows us to provide robust error messages across different systems.
- The [`os` package](https://golang.org/pkg/os/) provides helper functions to classify the failure indicated by a given error value.
- There are 3 common types of errors where helper functions are provided accordingly:
```go
package os

func IsExist(err error) bool      // file already exists (create operations)
func IsNotExist(err error) bool   // file not found (read operations)
func IsPermission(err error) bool // permission denied
```
- To handle errors most reliably, represent structured error values using a dedicated type.
  - For example, `os.PathError` can be used to describe failures in an operation on a file path (like Open or Delete).
  - Another variant is `os.LinkError` which can be used to describe failures involving two file paths, like `Symlink` and `Rename`.
- Example usage to distinguish error type and display appropriate error message:
```go
import (
  "errors"
  "syscall"
)

var ErrNotExist = errors.New("file does not exist")
// IsNotExist returns a boolean indicating whether the error is known to
// report that a file or directory does not exist. It is satisfied by
// ErrNotExist as well as some syscall errors.
func IsNotExist(err error) bool {
  if pe, ok := err.(*PathError); ok {
    err = pe.Err
  }
  return err == syscall.ENOENT || err == ErrNotExist
}
```
- Remember to do error discrimination immediately after the failig operation, before an error is propogated to the caller so the built error is not lost.

### Querying Behaviors with Interface Type Assertions

- We can use interface type assertions to only add additional methods (which often require a copy) when necessary. Instead of allocating a copy of every time, we can use the already satisfied interface:

```go
// WriteString writes s to w.
// If w has a WriteString method, it is invoked instead of w.Write.
func writeString(w io.Writer, s string) (n int, err error) {
  type stringWriter interface {
    WriteString(string) (n int, err error)
  }
  if sw, ok := w.(stringWriter); ok {
    return sw.WriteString(s) // avoid allocating memory for a copy
  }
  return w.Write([]byte(s)) // allocate the temporary copy
}

func writeHeader(w io.Writer, contentType string) error {
  if _, err := writeString(w, "Content-Type: "); err != nil {
    return err
  }
  if _, err := writeString(w, contentType); err != nil {
    return err
  }
  // ...
}
```

## Type Switches

- Interfaces are used in two distinct styles:
  - 1. An interface's methods express the similarities of concrete types that satisfy the interface but hide the representation details and intrinsic operations of those concrete types; emphasis on the methods, not on the concrete types. This style is exmplified in `io.Reader/Writer`, `fmt.Stringer`, `sort.Interface`, `http.Handler`, and `error`.
  - 2. Discriminated unions: explotation of ability of interface value to hold values of a variety of concrete types by considering interface to be the _union_ of those types; emphasis on concrete types that satisfy the interface, not on the interface's methods (if it has any). ie. there is no hiding of information.
- Type switches focused mainly on the second style.

```go
// Type switch
// Case order is significant as possibility of two cases matching.
// NOTE: no fallthrough is allowed
switch x.(type) {
case nil:       //...
case int, uint: //...
case bool:      //...
case string:    //...
default:        //... (Typically a panic is used here)
}
// Extended type switch form that binds extracted variable to new variable
switch x := x.(type) { /* ... */ }
```

- As with ordinary switches, typeswitches create a new lexical scope. Each case also creates its own lexical scope.
- Typical panic format for default case looks like `panic(fmt.Sprintf("unexpected type %T: %v", x, x))`
- Best uses for cases where a function might except an "any" empty interface type (`interface{}`), but it actually must be one of many types, called a _discriminated union_; in the example above, it's a discriminated union of int, uint, bool, string, and nil.

See [token-based XML decoding example](./xmlselect) for using typeswitches in practice.

## Interface Design Tips

- Remember interfaces are only needed when there are two or more concrete types that must be dealt with in a uniform way.
  - Take advantage of controlling what is/isn't exported.
  - Avoid creating a set of interfaces at the beginning of creating a new package and then only later defining the concrete types that satisfy them as this results in many unnessary abstractions with only a single implementation.
  - The exception to this rule is when an interface is satisfied by a single concrete type but that type cannot live in the same package as the interface because of its dependencies; in that case it's a good way to decouple two packages.
- Good rule of thumb for interface design is only ask for what you need; smaller interfaces are easier to satisfy.
