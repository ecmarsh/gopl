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
