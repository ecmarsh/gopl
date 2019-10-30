// Package eval provides an expression evaluator
package eval

import (
	"fmt"
	"math"
)

// Env is the environment that maps variable names to values
type Env map[Var]float64

// --- Concrete Eval methods ---

// Eval on Var performs an environment lookup,
// returns a zero if variable is not defined.
func (v Var) Eval(env Env) float64 {
	return env[v]
}

// Method for literal simply returns the value
func (l literal) Eval(_ Env) float64 {
	return float64(l)
}

// Eval methods for unary and binary recursively evaluate their operands
// then apply the operation `op` to them.
func (u unary) Eval(env Env) float64 {
	switch u.op {
	case '+':
		return +u.x.Eval(env)
	case '-':
		return -u.x.Eval(env)
	}
	panic(fmt.Sprintf("unsupported unary operator: %q", u.op))
}
func (b binary) Eval(env Env) float64 {
	switch b.op {
	case '+':
		return b.x.Eval(env) + b.y.Eval(env)
	case '-':
		return b.x.Eval(env) - b.y.Eval(env)
	case '*':
		return b.x.Eval(env) * b.y.Eval(env)
	case '/':
		return b.x.Eval(env) / b.y.Eval(env)
	}
	panic(fmt.Sprintf("unsupported binary operator: %q", b.op))
}

// Method for call evaluates arguments to the `pow`, `sin`, or `sqrt` function.
func (c call) Eval(env Env) float64 {
	switch c.fn {
	case "pow":
		return math.Pow(c.args[0].Eval(env), c.args[1].Eval(env))
	case "sin":
		return math.Sin(c.args[0].Eval(env))
	case "sqrt":
		return math.Sqrt(c.args[0].Eval(env))
	}
	panic(fmt.Sprintf("unsupported function call: %s", c.fn))
}
