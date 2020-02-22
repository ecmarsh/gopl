# Reflection

Go provides a mechanism called _reflection_ to update variables and inspect their values at run time, to call their methods, and to apply the operations intrinsic to their representation, which also lets us treat types themselves as first-class values.

Reflection increases the expressiveness of the language and are crucial for implementation of important APIs, specifically string formatting provided by fmt, protocol encoding provided by packages like `encoding/json` and `encoding/xml`, and the mechanism provided by `text/template` and `html/template` packages.

## Why Reflection

- To write a function capable of dealing uniformly with values of types that don't satisfy a common interface, don't have a known representation, or don't exist at the same time we design a function (or all of the above).
- A good example is the logic within `fmt.Fprintf`, which can usefully print values of any type, even user-defined ones.
- To implement this ourselves, we'll create a function called `Sprint`, as it will return the result as a string like `fmt.Sprintf` does. It starts with a type switch that tests whether argument defines a `String` method, and call it if so. Then add switch cases to test the value's dynamic type against each of the basic types (string, int bool, etc), then perform the appropriate formatting operation.

```go
func Spring(x interface{}) string {
  type stringer interface {
    String() string
  }

  switch x := x.(type) {
    case stringer:
      return x.String()
    case string:
      return x
    case int:
      return strconv.Itoa(x)
    // ...similar cases for int16, uint32, so on...
    case bool:
      if x {
        return "true"
      }
      return "false"
    default:
      // array, chan, func, map, pointer, slice, struct
      return "???"
  }
}
```

- So how do we deal with the composite types. We could add more cases but number of types is infinite. And what about named types like `url.Values`. Even for underlying values, it wouldn't match because two types are not identical and the type switch can't include a case for each type like `url.Values` because that requires the library to depend on its client.
- To handle this scenario, we need to inspect the representation of values of unknown types, which is why we need reflection.

## `reflect.Type` and `reflect.Value`

See [`reflect` package](https://golang.org/pkg/reflect) for more documentation.

- Perhaps the two most important types are `Type` and `Value`.

### Type and Typeof

- `Type` represents a Go type, an interface with many methods for discriminating among types and inspecting their components, like the fields of a struct or the parameters of a function.
- The sole implementation of `reflect.Type` is the type descriptor which is the same entity that identifies the dynamic type of an interface value.
- `reflect.TypeOf` is a function that accepts any type (`interface{}`)and returns its dynamic type as a `reflect.Type`:

```go
t := reflect.Typeof(3) // a reflect.Type
fmt.Println(t.String()) // "int"
fmt.Println(t)          // "int
```

- Recall that an assignment from a concrete value to an interface type performs an implicit interface conversion, which creates an interface value consisting of two components: its _dynamic type_ its the operand's type (`int`) and the _dynamic value_ is the operand's value (3).
- It is also capable of representing interface types:

```go
var w io.Writer = os.Stdout
fmt.Println(reflect.TypeOf(w)) // "*os.File"
```

- Note that `reflect.Type` satisifes `fmt.Stringer`. The dynamic type of an interface value is useful for debugging and logging, and `fmt.Printf` provides the short hand `%T` that uses `reflect.TypeOf` internally: `fmt.Printf("%T\n", 3) // int`

### Value and ValueOf

- A `reflect.Value` can hold a value of any type.
- `reflect.ValueOf` accepts any `interface{}` and returns a `reflect.Value` containing the interface's dynamic value. The results are always concrete, but `reflect.Value` can hold interfaces values too:

```go
v := reflect.ValueOf(3) // a reflect.Value
fmt.Println(v) // "3"
fmt.Printf("%v\n", v) // "3"
fmt.Println(v.String()) // NOTE: "<int Value>"
```

- Calling the type method on a Value returns its type as a `reflect.Type`:

```go
t := v.Type() // a reflect.Type
fmt.Println(t.String()) // "int"
```

- The inverse operation to `reflect.ValueOf` is the `reflect.Value.Interface` method which returns an `interface{}` holding the same concrete value as the `reflect.Value`:

```go
v := reflect.ValueOf(3) // a reflect.Value
x := v.Interface() // an interface{}
i := x.(int) // an int
fmt.Printf("%d\n", i) // "3"
```

- The difference between `reflect.Value` and an `interface{}` is that an empty interface hides the representation and intrinsic operations of the value it holds and exposes none of its methods, so unless the dynamic type is known and we use a type assertion to peer inside (as above), there is little we can do to the value within. But with a `Value`, it has many methods for inspecting its contents, regardless of its type.
- See [`format` example](./format) for a second attempt at a formatting function. Instead of a type switch, it uses `reflect.Value`'s `Kind` method to discriminate cases as their are only a finite number of _kinds_: the basic types `Bool`, `String`, and all the numbers; the aggretate types `Array` and `Struct`; the reference types `Chan`, `Func`, `Ptr`, `Slice`, and `Map`; `Interface` types; and finally `Invalid` meaning no value at all (the zero value kind of `reflect.Value` is reflect `Invalid`).

## `Display` Example: a Recursive Value Printer

- See [`display` example](./display) for example of improving the display of composite types.
- Best practice is to avoid exposing reflection in the API of a package by wrapping it:

```go
func Display(name string, x interface {}) {
  fmt.Printf("Display %s (%T):\n", name, x)
  display(name, reflect.ValueOf(x))
}
```

- Explanation of each case in examples:
  - **Slices and arrays**: `display` recursively invokes itself on each element of the sequence, appending the subscript notation "[i]" to the path. Only a few of `reflect.Value`'s methods are safe to call on any given value, e.g., `Index()` is only safe to call on `Slice`, `Array`, or `String`.
  - **Structs**: `Field(i)` returns the `i`-th field as a `reflect.Value` including fields promoted from anonymous fields. We use `reflect.Type` of the struct to append the field selector notation ".f" to the path and access the name of its `i`-th field.
  - **Maps**: The subscript notation "[key]" is appended to path as a shortcut because the type of map key isn't restricted to the types of the `format` example. Extending cases to other composite types requires more effort since these can also be valid map keys.
  - **Pointers**: The `reflect.Value` operation is safe even if the pointer value is nil (where more appropriate case is `Invalid`, but `isNil` is used to detect nil pointers explicitly. We prefix the path with an asterisk and parenthesize it to avoid ambiguity.
  - **Interfaces**: Use `isNil` to determine if interface is nil and if not, retrieve the value using the `Elem()` method to simply print its type and value.

- The example represents some cycles, and many Go programs contain at least some cyclic data. Additional bookkeeping is required to record the set of references that have been followed so far (this is also costly). The general solution requires using `unsafe` language features, exemplified in [low level programming notes](../13_lowlevel/readme.md).
  - When `fmt.Sprint` encounters a pointer, it breaks the recursion by printing the pointer's numeric value. Occasionally it gets stuck trying to print a slice or map that contains itself as an element, but rare cases don't warrant the considerable extra trouble of handling cycles.

- See [Encoding S-Expressions](./sexpr) example for another way of handling additional constructs such as `integer`, string with Go style quotations, symbols with unquoted names, and a list (zero or more items enclosed in parentheses).

## Setting Variables with `reflect.Value`

- So far, we've seen how to interpret values, but the point of knowing how to do this is to _change_ them.
- A variable is an _addressable_ storage location that contains a value, and its value may be updated through that address.
- With `reflect.Values`, some are addressable while others are not. For example:

```go
x := 2                   // value type  variable?
a := reflect.ValueOf(2)  // 2     int   no
b := reflect.ValueOf(x)  // 2     int   no
c := reflect.ValueOf(&x) // &x    *int  no
d := c.Elem()            // 2     int   yes (x)
```

- No `reflect.Value` returned by `reflect.ValueOf(x)` is addressable. But `d`, derived from `c` by dereferencing the pointer within it, refers to a variable and and is thus an addressable `Value` for any variable x.
- To determine if a `reflect.Value` is addressable, we can use the method `CanAddr`:

```go
fmt.Println(a.CanAddr()) // "false"
fmt.Println(b.CanAddr()) // "false"
fmt.Println(c.CanAddr()) // "false"
fmt.Println(d.CanAddr()) // "true"
```

- We obtain an addressable `reflect.Value` whenever we indirect through a pointer, even if we started from a non-addressable `Value`.
  - `reflect.ValueOf(e).Index(i)` refers to a variable, and is thus addressableeven if `reflect.ValueOf(e)` is not.
- To recover the variable from an addressable `reflect.Value` requires three steps.
  - First, call `Addr()`, which returns a `Value` holding a pointer to the variable.
  - Second, call `Interface()` on the `Value`, which returns an `interface{}` value containing the pointer.
  - Finally, if we know type of variable, we can use type assertion to retrieve the contents of the interface as an ordinary pointer and update the variable through the pointner:
  
  ```go
  x := 2 
  d := reflect.ValueOf(&x).Elem()   // d refers to the variable x
  px := d.Addr().Interface().(*int) // px := &x
  *px = 3                           // x = 3
  fmt.Println(x)                    // "3

  // or update the variable referred to by an
  // addressable `reflect.Value` directly,
  // without using a pointer by calling `reflect.Value.Set`
  d.Set(reflect.ValueOf(4))
  fmt.Println(x) // "4"

  // It is crucial to make sure the value is assignable
  // to the type of the variable or causes a panic:
  d.Set(reflect.ValueOf(int64(5))) // panic: int64 is not assignable to int

  // calling `Set` on a non-addressable reflect.Value panics too
  x := 2
  b := reflect.ValueOf(x)
  b.Set(reflect.ValueOf(x))
  b.Set(reflect.ValueOf(3)) // panic: Set using unaddr val

  // Variants of Set exist for certain groups of basic types: 
  d := reflect.ValueOf(&x).Elem()
  d.SetInt(3)
  fmt.Println(x) // "3"
  ```

- The specialized methods are somehwat more forgiving. For example, `SetInt` will succeed so long as the variable's type is some kind of signed integer, or even a named type whose underlying type is a signed integer. Even if the value is too large, ti will be truncated to fit.
- An addressable `reflect.Value` records whether it was obtained by traversing an unexported struct field, and if so, disallows modification. `CanAddr` is not usually the right check to use before setting a variable. `CanSet` is better since it can report whether a variable is addressable _and_ settable:

```go
fmt.Println(fd.CanAddr(), fd.CanSet()) // "true false"
```
