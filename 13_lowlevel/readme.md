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

- TODO
