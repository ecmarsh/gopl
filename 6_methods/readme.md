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






