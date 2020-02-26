package equal

import (
	"reflect"
	"unsafe"
)

func equal(x, y reflect.Value, seen map[comparison]bool) bool {
	if !x.IsValid() || !y.IsValid() {
		return x.IsValid() == y.IsValid()
	}
	if x.Type() != y.Type() {
		return false
	}

	// Cycle check: to ensure termination even for
	// cyclic structures, completed comparisons must be recorded.
	if x.CanAddr() && y.CanAddr() {
		xptr := unsafe.Pointer(x.UnsafeAddr())
		yptr := unsafe.Pointer(y.UnsafeAddr())
		if xptr == yptr {
			return true // identical references
		}
		c := comparison{xptr, yptr, x.Type()}
		if seen[c] {
			return true // already seen
		}
		seen[c] = true
	}

	switch x.Kind() {
	case reflect.Bool:
		return x.Bool() == y.Bool()
	case reflect.String:
		return x.String() == y.String()

		// ... numeric cases omitted for brevity...

	case reflect.Chan, reflect.UnsafePointer, reflect.Func:
		return x.Pointer() == y.Pointer()

	case reflect.Ptr, reflect.Interface:
		return equal(x.Elem(), y.Elem(), seen)

	case reflect.Array, reflect.Slice:
		if x.Len() != y.Len() {
			return false
		}
		for i := 0; i < x.Len(); i++ {
			if !equal(x.Index(i), y.Index(i), seen) {
				return false
			}
		}
		return true

		// .. struct and map cases omitted for brevity...

	}
	panic("unreachable")
}

// Equal reports whether x and y are deeply equal.
// Impl: it must call `reflect.ValueOf` on its arguments.
func Equal(x, y interface{}) bool {
	seen := make(map[comparison]bool)
	return equal(reflect.ValueOf(x), reflect.ValueOf(y), seen)
}

type comparison struct {
	x, y unsafe.Pointer
	t    reflect.Type
}

/*
Example
fmt.println(Equal([]int{1,2,3}, []int{1, 2, 3})        // "true"
fmt.Println(Equal([]string{"foo"}, []string{"bar"}))   // "false"
fmt.Println(Equal([]string(nil), []string{}))          // "true"
fmt.Println(Equal(map[string]int(nil), map[string]{})) // "true"

// Even works on cyclic inputs without infinite recursion:

// circular linked lists a -> b -> a and c -> c
type link struct {
	value string
	tail *link
}
a, b, c := &link{"a"}, &link{"b"}, &link{value: "c"}
a.tail, b.tail, c.tail = b, a, c
fmt.Println(Equal(a, a)) // "true"
fmt.Println(Equal(b, b)) // "true"
fmt.Println(Equal(c, c)) // "true"
fmt.Println(Equal(a, b)) // "false"
fmt.Println(Equal(a, c)) // "false"
*/
