// Package stack shows how to use a a slice to implement a stack.
package stack

var stack = make([]int, 0)

func push(val int) {
	stack = append(stack, val)
}

func pop() int {
	top := stack[len(stack)-1]
	stack = stack[:len(stack)-1]
	return top
}
