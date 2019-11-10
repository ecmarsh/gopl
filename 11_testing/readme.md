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

- To test a command, it is helpful to break out the essential part of the function, and use main as a driver.
  - During testing the main function is ignored.
  - A good strategy is organizing test cases a table to test different types of input.
  - Then we "fake" implementation by replacing other parts of the production and reading output; faking implementations makes configuration simpler, more reliable and easier to observe as well as avoid undesirable side effects. 
- Note that the `*_test.go` package for an executable, can also be named package `main`.
- If panics occur during tests, the test driver recovers, but the test is considered a failure.
- Expected errors occuring from bad user input, missing files, or imporper configuration should be reported by returning a non-nil error value.
- See [echoargs](../1_intro/echoargs/) for example of testing a command with table test-cases.

## White Box Testing

- White vs black box testing is categorized by level of knowledge they require of the internal workings of the tested package.
- Black box tests assumes nothing about the package other than what is exposed by its API and the documentation.
- A _white-box_ test has privelaged access to internal functions and data structures of a package so it can make observes and changes that an ordinary client cannot.
- An example of white-box is checking that data types are maintained after every operation.
- Black box updates are typically more robust and require fewer updates.
- White box helps to provide detailed coverage of the trickier parts of the implementation.
- In previous examples, the test for `IsPalindrome` is a black box test, simply calling the exported function, while the `EchoArgs` test uses a global variable of the package, making it a white-box test.
- Typically with white-box testing we fake an implementation for simpler configuration and better reliability. This is why its important to move the alogrithmic part of the function and the driver to separate functions for testing.
- Typically we add a private package-level variable to use for output depending on testing or production environment.
  - Remember to restore the original global variable if overriding by keeping a reference, then deferring a rest to the original function/package global.
  - Using global variables in this way is safe because `go test` does not normally run multiple tests concurrently.

## External Test Packages

- Since cyclic dependencies are forbidden in Go, resolve by declaring the package in test suffixed with `_test`.
  - Note another package is created, but cannot be imported or used/imported by that name.
- The external testing package is logically higher than the other packages.
- External test packages are especially useful for integration tests of several components since we can import packages freely exactly as an application would.
- To see which packages are included in production code (ie packages that go build will use):
```sh
# using fmt package as an example
$ go list -f={{.GoFiles}} fmt
[doc.go format.go print.go scan.go]
# to see testing packages
$ go list -f={{.TestGoFiles}} fmt
[export_test.go] # note usually fmt does not have any
# to see external testing packages included only for testing
$ go list -f={{.XTestGoFiles}} fmt
[fmt_test.go scan_test.go stringer_test.go]
```
- If external test packages need access to unexported items, create an in-package `_test.go` file and export the variables you need as a back door. Conventionally, this file is called `export_test.go`
  - Example with fmt (note no tests, just redeclaration of needed variables):
  ```go
  package fmt

  var IsSpace = isSpace
  ```
  - Now external tests can use `IsSpace` with techniques of white-box testing.

## Writing Effective Tests

- Go requires the user to implement functions for most of the testing features, by design.
- A good dtest doesn't explode on failure, but prints a clear and succint description of the symption of the problem and any other relevant facts regarding context.
- Ideally, maintainers shouldn't need to read source code to decipher a test failure.
- A good test shouldn't give up after one failure but try to report several errors in a single run since pattern of failures may be self revealing.
- Example of _BAD_ test which provides almost useless information:
```go
import (
  "fmt"
  "strings"
  "testing"
)

// A poor assertion function
func assertEqual(x, y int) {
  if x != y {
    panic(fmt.Sprintf("%d is %d", x, y))
  }
}
func TestSplit(t *testing.T) {
  words := strings.Split("a:b:c", ":")
  assertEqual(len(words), 3)
  // ...
}
```
- Assertion functions suffer from premature abstraction; by treating the failure of a particular tst as a mere difference, it forfeits the opportunity to provide meaningful contest.
- Example of improved test report that shows function that was called, its input, and the significance of the result; it explicitly identifies the actual value and the expectation, then continues to execute even if the assertion failures:
```go
func TestSplit(t *testing.T) {
  s, sep := "a:b:c", ":"
  word := strings.Split(s, sep)
  if got, want := len(words), 3; got != want {
    t.Errorf("Split(%q, %q) returned %d words, want %d",
      s, sep, got, want)
  }
}
```
- When needed, it is appropriate to use utility functions to make the testing simpler. (one example of a good utility function for this is `reflect.deepEqual`)
- Key to good test is to start by implementing the concrete behavior you want and only then use functions to simplify the code and eliminate repitition.
- Best results are rarely obtained by starting with a library of abstract, generic testing functions.

## Avoiding Brittle Tests

- An application that fails when it encounters new but valid inputs is called _buggy_.
- A test that spuriously fails when a sound chnge was made to the program is called _britle_. Brittle tests can exasperate its maintainer.
- Brittle tests fail for almost any change to the production code, good or bad and are sometimes called _change detector_ or _status quo_ tests. Time spent resolving these tests often depletes any benefit they may have once provided.
- Following are some best practices for avoiding brittle tests:
  - Test program's simpler and more stable interfaces in preference to its internal functions.
  Don't check for exact string matches, for example, but look for relevant substrings that will remain unchanged as the program evolves.
  - Note it is often worth writing a function that will distill complex output down to its essense so that assertions are reliable.

## Coverage

- "Testing shows the presence, not the absence of bugs."
- The degree to which a suite exercises the package under test is called a test's _coverage_.
- _Statement coverage_ is the simplest and most commonly used which is the fraction of source statements that are executed at least once during the test.
- To see coverage, go has a tool `go coverage` which is integrated into the `go test` tool.
- See [eval_test.go(TestCoverage)](../7_interfaces/eval/eval_test.go) for example of coverage test.
- To run the coverage, use the `-coverProfile` flag with go test:
```sh
# prints the summary of function statements covered to c.out
$ go test -run=Coverage -coverprofile=c.out $GOPATH/path/pkg
ok      $GOPATH/path/pkg     .0032s    coverage: 68.5% of statements
# to see a count of times ran use -cover-mode=count
# if you just need summary use
$ go test -cover
# in order to view output of c.out as html file:
$ go tool cover-html=c.out
```
- Note that some statements should always be red (e.g panics as default switch statements in testing); achieving 100% coverage usually isn't feasable.
- Other cases that make it unfeasable is handling esoteric errors, but always need to decide on tradeoff of cost of failures and cost of writing those tests.
- Coverage tools help identify the weakest spots, but need to use same good programming sense when writing programming tests.

## Benchmark Functions

- Uses `*testing.B` which is similar to `*testing.T` but includes extra features related to performance measurement and exposes an integer field N which specifies number of times to perform the operation being measurement.
- Simple benchmark test:

```go
import "testing"

func BenchmarkIsPalindrome(b *testing.B) {
  for i := 0; i < b.N; i++ {
    IsPalindrome("A man, a plan, a canal: Panama")
  }
}
```

- To run above:

```sh
$ cd $GOPATH/path/to/paldindrome
$ go test -bench=.
# tells us IsPalindrome took about 1.035 microseconds averaged over 1m runs
PASS BenchmarkIsPalindrome-8 1000000  1035 ns/op
ok    $GOPATH/path/to/paldindrome   2.179s
# note the suffix ...-8 is the value of GOMAXPROCS,
# important for concurrent benchmarks
```

- To see memory allocation statistics in report, use the `-benchmem` command-line flag.
```sh
$ go test -bench=. -benchmem
PASS
BenchmarkIsPalindrome 1000000 1026 ns/op 304 B/op  4 allocs/op
```

- As a side note, the quickest program for IsPalindrome will be the one that makes the fewest number of memory allocations.
- For benchmarking, the important measurements are the __relative__ timings of two different operations.
- For comparative benchmarks, typically just different functions with different parameters:

```go
func benchmark(b *testing.B, size int) { /* ... */ }
func benchmark10(b *testing.B, size int) { benchmark(b, 10) }
func benchmark100(b *testing.B, size int) { benchmark(b, 100) }
func benchmark1000(b *testing.B, size int) { benchmark(b, 1000) }
```

- Parameter size specifies size of input varies across benchmarks but is constant within each benhmark.
- Don't use the parameter `b.N` as the input size since unless you interpret it as an interation count for a fixed input, results will be meaninless.
- Benchmarks are important to keep around as program evolves or input grows, new operarating systems, etc.

## Profiling

- Profiling is an automated approach to performance measurement based on sampling a number of profile events during execution, then extrapolating from them during a post-processing step.
- The go tool supports different types of profiling:

Type | Desc | Command
--- | --- | ---
CPU Profile | identifies functions whose execution requires most CPU time. | `go test -cpuprofile=cpu.out` 
Heap Profile | identifies statements responsible for allocating the most memory. | `go test -memprofile=mem.out` 
Block Profiling | identifies operations responsible for blocking goroutines the longest such as sys calls, channels sends and receives, and acquisitions of locks. | `go test -blockprofile=block.out`

- Profiling becomes especially useful for long-running applications, so Go runtime's profiling features can be enabled under programmer control using the [runtime API](https://golang.org/pkg/runtime/).
- After gathering a profile, analyze it using the _pprof_ tool. Not an every day tool and basic use requires only two arguments: the executable that produced the profile and the profile log.

Example shell session to gather and display a simple CPU profile:

```bash
$ go test -run=NONE -bench=ClientServerParallelTLS64\
[PS2] -cpuprofile=cpu.log net/http
PASS
BenchmarkClientServerParallelTLS64-8 1000
  3141325 ns/op 143010 B/op 1747 allocs/op
ok   net/http   3.395s

# -text specifies output format (in this case a table with 10 hottest functions)
$ go tool pprof -text -nodecount=10 .http.test cpu.log
2570ms of 3590ms total (71.59%)
Dropped 129 nodes (cum <= 17.95ms)
Showing top 10 nodes out of 166 (cum >= 60ms)
  flat    flat%  sum %     cum    cum%
  1730ms 48.19% 48.19%  1750ms  48.75%  crypto/elliptic.p256ReduceDegree
   230ms  6.41% 54.60%   250ms   6.96%  crypto/elliptic.p256Diff
   ......
   ......

```

- Text may be enough to find cause of some issues, but for subtler issues, it's easier to use one of `pprof's` graphical displays which require GraphViz, download available [here](www.graphviz.org).
- For more on Go's profiling tools, read the Go blog's ["Profiling Go Programs"](https://blog.golang.org/profiling-go-programs).
