package palin

import "testing"

func TestPalindrome(t *testing.T) {
	if !IsPalindrome("detartrated") {
		t.Error(`IsPalindrome("detartrated") = false`)
	}
	if !IsPalindrome("kayak") {
		t.Error(`IsPalindrome("kayak") = false`)
	}
}

func TestNonPalindrome(t *testing.T) {
	if IsPalindrome("palindrome") {
		t.Error(`IsPalindrome("palindrome") = true`)
	}
	if IsPalindrome("racecars") {
		t.Error(`IsPalindrome("racecars") = true`)
	}
}

func TestFrenchPalindrome(t *testing.T) {
	if !IsPalindrome("ètè") {
		t.Error(`IsPalindrome("ètè") = false`)
	}
}
func TestCanalPalindrome(t *testing.T) {
	input := "A man, a plan, a canal: Panama"
	if !IsPalindrome(input) {
		t.Errorf(`IsPalindrome(%q) = false`, input)
	}
}

// TestIsPalindrome is more comprehensive consisting of all cases
func TestIsPalindrome(t *testing.T) {
	var tests = []struct {
		input string
		want  bool
	}{
		{"", true},
		{"a", true},
		{"aa", true},
		{"ab", false},
		{"kayak", true},
		{"detartrated", true},
		{"A man, a plan, a canal: Panama", true},
		{"Evil I did dwell; lewd did I live.", true},
		{"Able was I ere I saw Elba", true},
		{"ètè", true},
		{"Et se resservir, ivresse reste.", true},
		{"palindrome", false}, // non-palindrome
		{"desserts", false},   // semi-palindrome
	}
	for _, test := range tests {
		if got := IsPalindrome(test.input); got != test.want {
			t.Errorf("IsPalindrome(%q) = %v", test.input, got)
		}
	}
}

/*
// USAGE
$ cd $GOPATH/11_testing/palin
$ go test
ok   go/path/11_testing/palin   0.0006s

*/
