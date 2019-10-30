package eval

// Check methods checks for static errors in expression syntax tree.

import (
	"fmt"
	"strings"
)

// Check returns nil for evaluation of literal and Var since cannot fail
func (v Var) Check(vars map[Var]bool) error {
	vars[v] = true
	return nil
}
func (literal) Check(vars map[Var]bool) error {
	return nil
}

// Methods for unary and binary first check that operator is valid
func (u unary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-", u.op) {
		return fmt.Errorf("unexpected unary op %q", u.op)
	}
	return u.x.Check(vars)
}
func (b binary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-", b.op) {
		return fmt.Errorf("unexpected binary op %q", b.op)
	}
	if err := b.x.Check(vars); err != nil {
		return err
	}
	return b.x.Check(vars)
}
func (c call) Check(vars map[Var]bool) error {
	arity, ok := numParams[c.fn]
	if !ok {
		return fmt.Errorf("unknown function %q", c.fn)
	}
	if len(c.args) != arity {
		return fmt.Errorf("call to %s has %d args, want %d",
			c.fn, len(c.args), arity)
	}
	for _, arg := range c.args {
		if err := arg.Check(vars); err != nil {
			return err
		}
	}
	return nil
}

var numParams = map[string]int{"pow": 2, "sin": 1, "sqrt": 1}

/*
// selection of flawed inputs with errors
// Parse reports syntax errors
// Check reports semantic errors
x % 2 				unexpected '%'
math.Pi 			unexpected '.'
!true 				unexpected '!'
"hello" 			unexpected ':'
log(10) 			unknown function "log"
sqrt(1, 2)		call to sqrt has 2 args, want 1
*/
