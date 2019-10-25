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

