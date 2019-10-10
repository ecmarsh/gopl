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
- Using `[#:#]`, unlike python, changes the reference, not create a copy. Use `copy` to copy.
- To initialize a stack, use `make([]T, initialLen, cap)` or define values with `[]T{vals...}`

### Maps

`m := make(map[kType]vType))`

- Go map is a reference to hash table.
- All values must be same type. But keys and values can be different types.
- Keys can be any comparable type.
- **Note:** Cannot take address of map element as tradeoff of dynamic map is new storage locations may be assigned to support growing or refreshing of elements.
- Zero value for map is nil.
- Map values are initialized to zero value, similar to python's `defaultdict`:
  ```go
  m := make(map[int]int)
  m[1] += 1
  m[1] += 1
  m[2] += 1
  // m: {1->2, 2->1}
  ```
- Not unusual to use map as a set. For example, a set of strings might be `map[string]bool`, but ensure its being used a set before assuming.
