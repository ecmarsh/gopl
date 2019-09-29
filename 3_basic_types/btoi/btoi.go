// Package btoi converts a boolean to its integer value
// returns 1 if b is true and 0 if false
package btoi

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}
