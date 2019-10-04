# 4. Composite types

## Arrays

- Arrays are fixed length. Initial length must be a constant.

```go
var a [3] int
n := [3]int{1, 2, 3} // cannot assign different length arrays now
r := [...]int{99: -1} // assigns 100 element array with 99 0's and the last -1
```

- If array's element types are comparable, arrays can be compared with `==` or `!=`.
- This can be useful such as comparing arrays of bytes.
- Except for special cases such as SHA256's fixed-size hash, arrays are less preferrable as function parameters or results than slices.

## Slices

- Slice is like a dynamic array, which gives access to its underlying array.
- If slice were a strict, would resemble:

```go
type IntSlice struct {
  ptr      *int
  len, cap int
}
```

- Slices capacity must be able to handle amount of new elements before appending.

### Stack Implementation

```go
stack := make([]int, initialLen, capacity)
push := func(stack []int, v int) {
  stack = append(stack, v)
}
pop := func(stack []int) int {
  top := stack[len(stack)-1]
  stack = stack[:len(stack)-1]
  return top
}
```
