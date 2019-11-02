# Testing

## `go test` tool

- Using `go test [pkg/file]` builds files ending in `_test.go`.
- Within test files, three kinds of functions are treated specially:
  - tests (name begins with `Test`)
  - benchmarks (name begins with `Benchmark`)
  - examples (name starts with `Example`)
- After invoking, the tool scans the `*_test.go` files for special functions, generates a temporary `main` package that calls them all in the proper way, builds and runs it, reports the results, and then cleans up.

## Test Functions

```go
// Optional "Name" portion must begin with a capitl letter
func TestName(t *testing.T) {
  // ...
}
```
- The `t` parameter provides methods for reporting test failures and logging.
- See [./palin](./palin) for a simple example of testing if a word is a palindrome.
- If tests are performing slow, use the -v flag to see the times of individual tests and then the `-run` flag to run only tests that match a pattern:
```sh
go test -v -run="regex"
```
- Example also includes a table-driven style. We could improve on table in [the example](./palin_test.go) by grouping the things we are testing into different tables and providing a more helpful error message for each case.
- Note that one failure does not cause an entire stack print and following cases are still run (tests are independent of each other).
- To stop a test case after a failure or if there's a cascade, use `t.Fatal` or `t.Fatalf`.
- Test failure methods are usually of the form: `F(x) = y, want z`, where `F(x)` explains attempted operation and its input, `y` is the actual result, and `z` is the expected result.
- As best practices:
  - Avoid boilerplate and redundant information
  - When testing a boolean, omit the want z since we know what z should be.
  - If `x`, `y`, or `z` is lengthy, print a concise summary of the relevant parts instead.

## Randomized Testing

- Randomized testing consists of exploring a broader range by constructing inputs at random.
- Strategies include:
  - Writing an alternative implementation of the function that uses a less efficient but simpler and clearer algorithm, then checking that both implementations give the same result.
  - Create input variables according to a pattern so we know what output to expect.
- Since randomized tests are nondeterministic, it's important to log the failing test record with sufficient information to reproduce the failure.
- Using the current time as a source of randomness is a good way to explore novel inputs each time a test is run over its lifetime; especially valuable with automated system to run all tests periodically.

```go
import "math/rand"

// randomPalindrome returns a palindrome whose length and contents
// are derived from the psuedo-random number generator ring.
func randomPalindrome(rng *rand.Rand) string {
  n := rng.Intn(25) // random length up to 24
  runes := make([]rune, n)
  for i := 0; i < (n+1)/2; i++ {
    r := rune(rng.Intn(0x1000)) // random rune up to '\u0999'
    runes[i] = r
    runes[n-1-i] = r
  }
  return string(runes)
}
func TestRandomPalindromes(t *testing.T) {
  // Initialize a psuedo-random number generator.
  seed := Time.Now().UTC().UnixNano()
  t.Logf("Random seed: %d", seed)
  rng := rand.New(rand.NewSource(seed))
  for i := 0; i < 1000; i++ {
    p := randomPalindrome(rng)
    if !IsPalindrome(p) {
      t.Errorf("IsPalindrome(%q) = false", p)
    }
  }
}
```

## Testing a Command

