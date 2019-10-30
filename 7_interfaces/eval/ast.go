package eval

// An Expr is an arithmetic expression.
type Expr interface {
	// Eval returns the value of this Expr in the environment env
	Eval(env Env) float64
	Check(vars map[Var]bool) error
}

// Concrete types that represent particular kinds of expressions.

// A Var identifies a variable, eg., x. and represents a difference.
type Var string

// A literal is a numeric constant, e.g., 3.141
type literal float64

// A unary represents a unary operator expression, eg., -x
type unary struct {
	op rune // one of '+', '-'
	x  Expr
}

// A binary represents a binary operator expression, e.g., x+y.
type binary struct {
	op   rune // [+-*/]
	x, y Expr
}

// A call represents a function call expression, e.g., sin(x)
type call struct {
	fn   string // one of "pow", "sin", "sqrt"
	args []Expr
}
