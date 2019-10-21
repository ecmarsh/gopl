// Package intset implements a bit vector type
package intset

import (
	"bytes"
	"fmt"
	"math/bits"
)

// w, word size, is the effective size of uint in bits,
// either 32 or 64, depending on platform.
const w = 32 << (^uint(0) >> 63)

// IntSet is a set of small non-negative intgers.
// Its zero value represents the empty set.
type IntSet struct {
	words []uint
}

// Len returns the number of elements in the set
func (s *IntSet) Len() int {
	var l int
	for _, word := range s.words {
		l += bits.OnesCount(word)
	}
	return l
}

// Remove removes x from the set.
func (s *IntSet) Remove(x int) {
	word, bit := x/w, uint(x%w)
	s.words[word] &^= 1 << bit
}

// Clear removes all elements from the set.
func (s *IntSet) Clear() {
	s.words = s.words[:0]
}

// Copy returns a copy of the set.
func (s *IntSet) Copy() *IntSet {
	t := &IntSet{}
	copy(t.words, s.words)
	return t
}

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
	word, bit := x/w, uint(x%w)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add adds the non-negative value x to the set.
// Each word has 64 bits, so we use the quotient x/64 as the word
// index and the remainder x%64 as the bit index within that word.
func (s *IntSet) Add(x int) {
	word, bit := x/w, uint(x%w)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// AddAll adds multiple non-negative values to set.
func (s *IntSet) AddAll(nums ...int) {
	for _, x := range nums {
		s.Add(x)
	}
}

// UnionWith sets s to the union of s and t.
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// IntersectWith sets s to the intersection of s and t.
func (s *IntSet) IntersectWith(t *IntSet) {
	for i, tword := range t.words {
		if i > len(s.words) {
			break
		}
		s.words[i] &= tword
	}
}

// DifferenceWith sets s to the set-theoretic difference of s and t.
// If an element in s is also in t, the element is removed from s.
func (s *IntSet) DifferenceWith(t *IntSet) {
	for i, tword := range t.words {
		if i > len(s.words) {
			break
		}
		s.words[i] &^= tword
	}
}

// SymmetricDifference sets s to elements in s or t, but not both.
func (s *IntSet) SymmetricDifference(t *IntSet) {
	for i, tword := range t.words {
		if i > len(s.words) {
			break
		}
		s.words[i] ^= tword
	}
}

// Elems returns a slice containing elements of the set.
// Useful when needing to iterate over with a range loop.
func (s *IntSet) Elems() []int {
	a := []int{}
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < w; j++ {
			if word&(1<<uint(j)) != 0 {
				a = append(a, w*i+j)
			}
		}
	}
	return a
}

// String returns the set as a string of the form "{1 2 3}"
// Implements a similar method as intsToString
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < w; j++ {
			if word&(1<<uint(j)) != 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprintf(&buf, "%d", w*i+j)
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

/*
// Usage
var x.Add(1)
x.Add(144)
x.Add(9)
fmt.Println(x.String()) // "{1 9 144}"

y.Add(9)
y.Add(42)
fmt.Println(y.String()) // "{9 42}"

x.UnionWith(&y)
fmt.Println(x.String()) // "{1 9 42 144}"

fmt.Println(x.Has(9), x.Has(123)) // "true false"

// NOTE: Since we declard String and Has as methods of pointer
// type *Intset for consistenstency with other tow methods,
// an IntSet value does not have a String method,
// which may lead to surprising cases when printing:
fmt.Println(x) // "{[1439804651168 0 65536]}"
*/
