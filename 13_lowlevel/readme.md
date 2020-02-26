# Low-Level Programming

- Go programs guarantee a number of safety properties against failing programs, but ocassionally we may choose to forfeit some helpful guarantees to achieve the highest performance, interpoerate with libraries written in other languages, or to implement a function that cannot be expressed in pure Go.

## `unsafe.Sizeof`, `Alignof`, and `Offsetof`

- The [`unsafe` package](http://golang.org/pkg/unsafe) provides access to a number of built-in language features that are not ordinary available because they expose details of Go's memory layout. It is used extensively within low-level packages like `runtime`, `os`, `syscall`, and `net` that interact with the OS, but is almost never needed by ordinary packages and should be used frivously.
- Often these three functions are not very "unsafe", and may be helpful for understanding the layout of raw memory in a program when optimizing for space.

### `unsafe.Sizeof`

- `unsafe.Sizeof` reports the size in bytes of the representation of its operand, which may be an expression of any type; the expression is not evaluated:

```go
import "unsafe"
// Call to sizeof is a constant expr of type uintptr
// result may be used as the dimension of an
// array type to compute other constants
fmt.Println(unsafe.Sizeof(float64(0))) // "8"
```

- It reports only the size of the fixed part of each data structure (like the ptr and length of a string), but not the indirect parts (like the contents of a string).
- Computers load and store values from memory most efficiently when the values are properly _aligned_. Alignment requirements of higher multiples are unusual, even for larger data types such as complex 8.
- The size of a value of an aggretate type (struct or array) is at least the sum of the sizes of its fields or elements but may be greater due to presence of "holes", which are unused spaces added by the compiler to ensure the following field/element is properly aligned relative to the start of the struct or array.
- Some examples of typical sizes for non-aggregate Go-types:

Type | Size
---- | ----
`bool` | 1 byte
`intN, uintN, floatN, complexN` | N/8 bytes (e.g `float64`=8 bytes)
`int,uint,uintptr` | 1 word
`*T` | 1 word
`string` | 2 words (data, len)
`[]T` | 3 words (data,len,cap)
`map` | 1 word
`func` | 1 word
`chan` | 1 word
`interface` | 2 words (type, value)

- Language specs don't guarantee order in which fields are declared matches order in memory, so in theory, a compiler is free to rearrange them to pack them more efficiently (although none do as of now).

### `unsafe.Alignof`

- `unsafe.Alignof` reports the required alignment of its argument's type. It may be applied to an expression fo any type and yields a constant.
- Typically, boolean and nmeric types are aligned to their size (up to max of 8 bytes) and all other types are word-aligned.

### `unsafe.Offsetof`

- `unsafe.Offsetof`, whose operand must be a field selector `x.f`, computes the offset of field `f` relative to the start of its enclosing struct `x`, accounting for holes, if any.

## `unsafe.Pointer`

- Recall that a pointer type are written `*T`, meaning "a pointer to a variable of type T".
- `unsafe.Pointer` is a special pointer type that can hold the address of any variable. They are also comparable and may be compared with nil, its zero type.
- `unsafe.Pointer` may be converted from and to ordinary pointers (not necessarily of the same type `*T`). For example, we can inspect the bit pattern of a floating-point variable by converting a `*float64` ptr to a `*uint64`:

```go
package math

func Float64bits(f float64) uint64 {
  return *(*uint64)(unsafe.Pointer(&f))
}

fmt.Printf("%#016x\n", Float64bits(1.0)) // "0x3ff0000000000000"
```

- Many `unsafe.Pointer` values are intermediaries for converting ordinary pointers to raw numeric addresses and back again. This example takes the address of a variable `x`, adds the offset of its `b` field, converts the resulting address to `*int16`, and through that pointer updates x.b:


```go
var x struct {
  a bool
  b int16
  c []int
}

// equivalent to pb := &x.b
pb := (*int16)(unsafe.Pointer(
        uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)))
*pb := 42

fmt.Println(x.b) // "42"
```

- Many things can go wrong with conversions, especially after an `unsafe.Pointer` to `uintptr` conversion so best practice is to assume the bare minimum and treat all `uintptr` values as if they contain the _former_ address of a variable, and minimize the number of operations between converting an `unsafe.Pointer` to a `uintptr` and using that `uintptr`.
- When calling a library function that returns a `uintptr`, the result should be immediately converted to an `unsafe.Pointer` to ensure that it continues to point to the same variable:

```go
package reflect

func (Value) Pointer() uintptr
func (Value) UnsafeAddr() uintptr
func (Value) InterfaceData() [2]uintptr // (index 1)
```

### Example: Deep Equivalence

- `reflect.DeepEqual` reports whether two values are deeply equal. Basic values are compared via `==` and composite values are traversed recursively.
- For example: to compare two `[]string` values:

```go
func TestSplit(t *testing.T) {
  got := strings.Split("a:b:c", ":")
  want := []string{"a", "b", "c"}
  if !reflect.DeepEqual(got, want) { /* ... */ }
}
```

- `DeepEqual` does not consider a nil map equal to a non-nil empty map, nor a nil slice equal to a non-nil empty one:

```go
var a, b []string = nil, []string{}
fmt.Println(reflect.DeepEqual(a, b)) // "false"

var c, d map[string]int = nil, make(map[string]int)
fmt.Println(reflect.DeepEqual(c, d)) // "false"
```

- See [equal example](./equal) for a modified equivalence function that compares arbitrary values.
