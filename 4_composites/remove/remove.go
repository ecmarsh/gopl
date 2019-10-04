// Package remove removes an element from the middle of a slice
package remove

// Slides higher elements into removed slot
func remove(slice []int, i int) []int {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}
