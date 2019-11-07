package eval

import (
	"fmt"
	"math"
	"testing"
)

// For each entry in the table, the test parses the expression,
// evaluates it in the environment, and prints the result.
func TestEval(t *testing.T) {
	tests := []struct {
		expr string
		env  Env
		want string
	}{
		{"sqrt(A / pi)", Env{"A": 87616, "pi": math.Pi}, "167"},
		{"pow(x, 3) + pow(y, 3)", Env{"x": 12, "y": 1}, "1729"},
		{"pow(x, 3) + pow(y, 3)", Env{"x": 9, "y": 10}, "1729"},
		{"5 / 9 * (F - 32)", Env{"F": -40}, "-40"},
		{"5 / 9 * (F - 32)", Env{"F": 32}, "0"},
		{"5 / 9 * (F - 32)", Env{"F": 212}, "100"},
		// additional tests not in the book
		{"-1 + -x", Env{"x": 1}, "-2"},
		{"-1 - x", Env{"x": 1}, "-2"},
	}
	var prevExpr string
	for _, test := range tests {
		// Print expr only when it changes.
		if test.expr != prevExpr {
			fmt.Printf("\n%s\n", test.expr)
			prevExpr = test.expr
		}
		expr, err := Parse(test.expr)
		if err != nil {
			t.Error(err) // parse error
			continue
		}
		got := fmt.Sprintf("%.6g", expr.Eval(test.env))
		fmt.Printf("\t%v => %s\n", test.env, got)
		if got != test.want {
			t.Errorf("%s.Eval() in %v = %q, want %q\n",
				test.expr, test.env, got, test.want)
		}
	}
}

/*
// output of testing
$ go test -v PATHFROMROOT/eval
	sqrt(A / pi)
		map[A:87616 pi:3.141592653589793] => 167
	pow(x, 3) + pow(y, 3)
		map[x:12 y:1] => 1729
		map[x:9 y:10] => 1729
	5 / 9 * (F - 32)
		map[F:-40] => -40
		map[F:32] => 0
		map[F:212] => 100
	-1 - x
		map[x:1] => -2
	-1 + -x
		map[x:1] => -2
*/

func TestCoverage(t *testing.T) {
	var tests = []struct {
		input string
		env   Env
		want  string // expected error from Parse/Check or result from Eval
	}{
		{"x % 2", nil, "unexpected '%'"},
		{"!true", nil, "unexpected '!'"},
		{"log(10)", nil, `unknown function "log"`},
		{"sqrt(1, 2)", nil, "call to sqrt has 2 args, want 1"},
		{"sqrt(A / pi)", Env{"A": 87616, "pi": math.Pi}, "167"},
		{"pow(x, 3) + pow(y, 3)", Env{"x": 9, "y": 10}, "1729"},
		{"5 / 9 * (F - 32)", Env{"F": -40}, "-40"},
	}

	for _, test := range tests {
		expr, err := Parse(test.input)
		if err == nil {
			err = expr.Check(map[Var]bool{})
		}
		if err != nil {
			if err.Error() != test.want {
				t.Errorf("%s: got %q, want %q", test.input, err, test.want)
			}
			continue
		}
		got := fmt.Sprintf("%.6g", expr.Eval(test.env))
		if got != test.want {
			t.Errorf("%s: %v => %s, want %s", test.input, test.env, got, test.want)
		}
	}
}

/*
// Usage
$ go test -v -run=Coverage $GOPATH/7_interfaces/eval
=== RUN TestCoverage
--- PASS: TestCoverage (0.00s)
PASS
ok $GOPATH/7_interfaces/eval  .0009s

*/
