# Methods

For this section (and generally in Go terminology):

- An _object_ is a value or variable that has methods
- A _method_ is a function associated with a particular type.

## Method Declarations

```go
func (r Receiver) MethodName(args T) resT { }
```
- Declare methods by adding an extra parameter before the function name.
  - The parameter attaches the function to the type of that parameter
```go
type Point struct{ X, Y float64 }
// method of the Point type
func (p Point) Distance(q Point) float64 {
  return math.Hypot(q.X-p.X, q.Y-p.Y)
}
// p Point is the receiver
// p.Distance is called a selector
```
- Don't use `this` or `self` to name receivers; use the receiver name just like you would for any other parameter.
  - Common choice for names is the first letter of the reciever type.
- Note there is no conflict between method declarations and function declarations of the same name. For example, if `func Distance()` was also defined in above, no error. But a method and field in same struct cannot have the same name (compile error).
- Methods may be declared for any type, using a named type (as long as its underlying point isn't a pointer or interface):
```go
// Give slice path a method called Distance
type Path []Point
func (path Path) Distance() float64 {
  sum := 0.0
  for i := range path {
    if i > 0 {
      sum += path[i-1].Distance(path[i])
    }
  }
  return sum
}
```
- Methods with same signature/type need different names, but different signatures can use the same method name.

## Methods with a Pointer Receiver

```go
func (pr *PtrReceiver) Method(arg T) { ... } 
//   ^ Parentheses necessary
// To invoke a method on a pointer receiver:
(*pr).Method(arg)
// or shorthand, and compiler will perform implicit (&pr)
pr.Method(arg)
// Note variable must be defined first since must be addressable
PtrReceiver{field:x}.Method(arg) // NOT allowed since no address yet
```
- Named types and pointers to named types are the only types that may appear in a receiver declaration.
  - The named type must not have an underlying type of a pointer:
  ```go
  type P *int
  func (P) f() { ... } // compile error
  ```
- In every valid method call, exactly one of these staements is true:
  - Either receiver argument has same type as receiver parameter (both type T or both type *T):
  ```go
  NamedType{fields...}.Method(arg) // T
  nt.Method(arg) // *T
  ```
  - Or receiver argument is a _variable_ of type T and receiver parameter has type *T, where compiler implicitly takes the address of the variable:
  ```go
  p.Method(arg) // implicit (&p)
  ```
  - Or receiver argument type *T and receiver parameter has type T, where compiler implicitly dereferences receiver (loads the value):
  ```go
  ptr.Method(arg) // implicit (*ptr)
  ```
- If all methods of named type T have receiver of type T (not *T), it is safe to copy all instances, but calling any of its methods makes a copy.
  - Avoid copying instances of T if the method has a pointer receiver since may involate internal variants.

## Nil Is a Valid Receiver Value

- As functions allow nil pointers as arguments, so do some methods for their receiver as sometimes `nil` is a meaninful zero value of the type (maps, slices, etc):
```go
// An IntList is a linked list of integers.
// A nil *IntList represents the empty list.
type IntList struct {
  Value int
  Tail *IntList
}
// Sum returns the sum of the list elements
func (list *IntList) Sum() int {
  if list == nil {
    return 0
  }
  return list.Value + list.Tail.Sum()
}
```
- When you define type whose methods allow nil as receiver value, best practice is to point this out explictly in the documentation as in example above.

## Composing Types by Struct Embedding

- As we can refer indirectly to [embedded structs](../4_composites/readme.md#Struct_Embedding_and_Anonymous_Fields), we can call methods in the same way.
- When a method is called indirectly, the method has been _promoted_ to type we're calling it on.
- This is the mechanism that allows many methods to be built up by composition of several fields.
- An important note is that the containing struct with embedded struct is _not_ similar to inheritance (where embedded would be base class and containing would be derived). It is a closer relationship to "has-a", so would be an "implements" relationship.
  - As a takeaway, we cannot use containing structs in place of say a function that takes a type of the embedded struct. The embedded struct must be explicitly selected through the container.
  - Example:
  ```go
  import "image/color"

  type Point struct { X, Y float64 }
  type ColoredPoint struct {
    Point
    Color color.RGBA
  }
  // When we promote a method, compiler implictly generates
  // wrappers that would function similarly to this:
  func (p ColoredPoint) Distance(q Point) float64 {
    return p.Point.Distance(q) // Method is called explicitly on p.Point
  }
  func (p *ColoredPoint) ScaleBy(factor float64) {
    p.Point.ScaleBy(factor) // Same, receiver value is p.Point, not *p
  }
  ```
- We can reduce the explicit call by using a pointer as the anonymous field, so fields and methods are promoted indirectly from pointed-to object:
  ```go
  type ColoredPoint struct {
    *Point
    Color color.RGBA
  }
  p := ColoredPoint{&Point{1, 1}, red}
  q := ColoredPoint&Point{{5, 4}, blue}
  fmt.Println(p.Distance((*q.Point)) // "5"
  q.Point = p.Point                  // p and q now share same point
  p.ScaleBy(2)                       // ScaleBy is promoted indirectly
  fmt.Println(*p.Point, *q,Point)    // "{2 2} {2 2}"
  ```
- A struct type may have more than one anonymous field (we could have defined `Color` as just `color.RGBA` above).
  - Then `ColoredPoint` would _have_ (not inherit, but be able to use) any of the additional methods of `Point` and `RGBA`, plus any other methods declared by `ColoredPoint`.
  - When a method is called in containing struct, compiler resolves by looking first for directly declared method, then for methods promoted once from embedded fields, then methods promoted twice, etc. If call is ambiguous (eg two methods with same name promoted from same rank), the compiler reports an error.
- It can sometimes be useful for _unnamed struct types_ to have methods too by allowing for more expressive names and self-explanatory syntax:
  ```go
  // Shows part of a simple cache implemented with two pkg-level vars
  var (
    mu sync.Mutex // guards mapping
    mapping = make(map[string]string)
  )
  func Lookup(key string) string {
    mu.Lock()
    v := mapping[key]
    mu.Unlock()
    return v
  }
  // Below is equivalent to above but groups together related
  // variables by defining methods on unnamed struct types
  var cache = struct {
    sync.Mutex
    mapping map[string]string {
      mapping: make(map[string]string)
    }
  }
  // Lookup becomes self-explanatory now
  func Lookup(key string) string {
    cache.Lock()
    v := cache.mapping[key]
    cache.Unlock()
    return v
  }
  ```

## Method Values and Expressions










