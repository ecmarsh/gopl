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

- TODO
