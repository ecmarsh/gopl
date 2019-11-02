// Package palin provides functions for word games.
package palin

import "unicode"

// IsPalindrome reports whether s reads the same forward and backward.
// Letter case is ignored, as are non-letters.
func IsPalindrome(s string) bool {
	var letters []rune
	for _, r := range s {
		if unicode.IsLetter(r) {
			letters = append(letters, unicode.ToLower(r))
		}
	}
	for i, j := 0, len(letters)-1; i < j; i, j = i+1, j-1 {
		if letters[i] != letters[j] {
			return false
		}
	}
	return true
}
