# Functions

## Function Declarations

```txt
func name(parameter-list) (result-list) {
  body
}
```

Function declaration variations:
```go
func add(x int, y int) int   { return x + y }
func sub(x, y int) (z int)   { z = x - y; return }
func first(x int, _ int) int { return x }
func zero(int, int) int      { return 0 }
```

- Parameter list specifies names and types of function's parameters, which are the local variables whose values, or arguments, are supplied by the caller. The result list specifies the types of the values that the function returns.
- Two functions have the same type or signature if they have the same sequence of parameter types and the same sequence of result types.
- All arguments are passed by value, so function receives a copy of each argument; modifications to copy do not affect the caller.
  - If argument contains a dereference (like pointer, slice, map, function or channel) then caller may be affected by modifications the function makes to variables.
- If a function does not have a body, this means it is implemented in a language other than Go:
  - Example: `func Sin(x float64) float64`

## Recursion

- Functions may be recursive: they may call themselves either directly or indirectly.
- Recursive call must be within a function that is defined (not func expression).
- When adding values to slice, etc, remember that everything is passed by values, so each calee will receive a copy of the slice. The function does not modify the version that the caller has.
  - For example, in backtracking algorithms, normally, you'd `push/append()` recurse, then `pop()` but in Go, there's no need to `pop()` after.

## Multiple Return Values

- A function can return more than one result. 
- The most common situation is when an operation may fail and a function returns (successful T, error).
- Callers of functions such as (success, err) sneed to explicitly assign the _all return values_ or use the blank identifier:
  - `res, err := fn(args)`
  - `res, _ := fn(ars)`
- When returning multiple success values (especially of same type), it is best practice to use a named return to help 'self-document' what is expected:
  - `func HourMinSec(t time.Time) (hour, minute, second int) { ... }`
- For a function with _named results_, a naked/bare return (return operands are omitted) is ok
  - Assign the variables a value (or leave them as zero-value) and with naked return, they will be returned in order listed in signature.

## Errors

- Some functions may always succeed (such as bool return like `Contains()`), others succeed as long as precondition is met (all arguments included, etc.), but many functions cannot always assure success due to factors outside of programmer's control.
- For example, I/O has many error possibilities, even in simple read/writes.
- Conventially, errors where failure is an expected behavior return additional parameter (usually) as the last result:
  - such as a lookup in maps return a bool if key does not exist.
  - or an `error` result that reveals success by having a value of `nil` or not.
- Go errors differ by most commonly returning errors as ordinary values, not _exceptions_.
  - Note there are _exceptions_ for non-routine errors that occur.
- As a tradeoff, more control-flow is needed to handle these returns appropriately, which is the point of using ordinary values as the mechanism.

### Error-Handling Strategies


5 common strategies for error-handling:

1. Propogate the error: make the failure in a subroutine the failure of the calling routine:
```go
resp, err := html.Get(url)
if err != nil {
  return nil, err
}
// Alternatively, propogate a more specific error.
doc, err := html.Parse(resp.Body)
resp.Body.Close()
if err != nil {
  return nil, fmt.Errorf("parsing %s as Html: %v", url, err)
}
```
  - Note `fmt.Errorf` formats an error message using `fmt.Sprintf` to return a new error value - very useful for providing a clear causal chain from root problem to overall failure.
  - As a best practice, message strings should _not_ be capital and _not_ include newlines as error logs often contain many errors. This makes easier to self contain messages when searching using `grep`, etc.
  - As an example of specificity, if reading many links, files, include the name of link/file etc where the failure occured in the error message.

2. Retry the failed operation, mainly for errors that represent trasient or unpredicatable problems. 
  - Of course, set a time or attempt limit using this strategy.
  - See [`waitforserver.go`](./waitforserver/waitforserver.go) for example.

3. Stop the program gracefully, if function is essential to program progress.
  - Generally, should be reserved for the `main` package of a program.
  - Go packages usually propogate errors, unless error is sign of internal inconsistency (aka a bug).
  - Example:
  ```go
  if err := WaitForServer(url); err != nil {
    fmt.FPrintf(os.Stderr, "Site is down: %v\n", err)
    os.Exit(1)
  }
  // More conveniently, use log.Fatalf for same results
  if err := WaitForServer(url); err != nil {
    log.Fatalf("Site is down: %v\n", err)
  }
  // By default, Fatalf prints date and time (useful for long running)
  // To surpress dait and time, we can use
  log.Setprefix("wait: ") // Add package as prefix
  log.SetFlags(0)         // Surpress date/time
  ```

4. Log the error and continue, if can continue with possible reduced functionality:
  - Example:
  ```go
  if err := Ping(); err != nil {
    log.Printf("ping failed: %v; networking disabled", err)
    // Note log functions automatically append '\n' if not present. 
  }
  // Not using log and printing directly to error stream 
  if err := Ping(); err != nil {
    fmt.Fprintf(os.Stderr, "ping failed: %v; networking disabled\n", err)
  }
  ```

5. As a last option and least commonly, safely an ignore the entire error
  - Example:
  ```go
  dir, err := ioutil.TempDir("", "scratch")
  if err != nil {
    return fmt.Errof("failed to create temp dir: %v", err)
  }
  // ...use temp dir..
  os.RemoveAll(dir) // ignore errors; $TMPDIR is cleaned periodally
  ```
  - Best practice to always consider errors after function call.
  - If using this method, best to document intention of ignoring clearly.

- Similar to node style, handle failures before success.

### End of File (EOF)

- EOF has its own error class to distinguish this error from others.
  - For example if `n` is length of file, and we don't read `n` bytes, then any error represents a failure.
  - But typically we need to respond to EOF differently than any others.
- The [`io` package](https://golang.org/pkg/io/) guarantees that any read failure caused by EOF is always reported by distinguished error, `io.EOF`, defined as `var EOF = errors.new("EOF")`
- To detect EOF condition:
```go
in := bufio.NewReader(os.Stdin)
for {
  r, _, err := in.RadRune()
  if err == io.EOF {
    break // finished reading
  }
  if err != nil {
    return fmt.Errof("read failed: %v", err)
  }
  // ...use r..
}
```

## Function Values

- Functions are first-class variables in Go.
- Function values have types and may be assigned to variables or passed to or returned from functions:
```go
func square (n int) { return n * n }
func negative(n int) int { return -n }
func product(m, n int) int { return m * n }

f := square
fmt.Println(f(3)) // "9"

f = negative
fmt.Println(f3)) // "-3"
fmt.Printf("%%T\n", f) // "func(int) int"
```
- Zero-value of a function type is `nil` and calling a `nil` function type causes a panic.
  - Accordingly, function values may be compared with nil: `if fn != nil { ... }`
  - But other than nil, they are not comparable, which also means they may not be used as keys in maps.
- Function values let us parameterize functions over data, and also behavior (passing in functions).
  - See [`funcvals.forEachNode`](./funcvals/funcvals.go) for an example of passing functions as parameter.

## Anonymous Functions

- While named functions can be declared only at the package level, _function literals_ can denote a function value without any expression.
- A function literal's value is called an _anonymous function_.
- Indicated by not using name following `func` keyword:
  - `strings.Map(func(r, rune) rune { return r + 1 }, "HAL-9000")`
  - `add := func(x, y int) int { return x + y }`
  - as a return: `return func() int { ... }`
    - Note the parent funcs return would also be int, type of return in return function.
- Function values also have state: anonymous inner functions can access and update the local variables of enclosing functions. (closures)
  - This is why functions are not comparable.
  - Accordingly, variable lifetimes are not determined by its scope.
- See [topological sort](./anon/toposort.go) for example.

### Caveats: Capturing Iteration Variables

- Sometimes Go's lexical scope rules can cause surprising results.
- To capture an iteration variable, assign it to a different variable in the enclosing scope first:
```go
var rmdirs []func()
for _, d := range tempDirs() {
  dir := d // The necessary assignment 
  os.MkdirAll(dir, 0755)
  rmdirs = append(rmdirs, func() {
    os.Removeall(dir)
  })
}
// Note sometimes you'll see programmers just do d := d
// to declare an inner variable of same name.
```
- We see this problem not just with loops, but also with `go` statement and with deferred functions.

## Variadic Functions

- One that can be called with varying numbers of arguments. 
- Declare a variadic function by preceding the type of the final parameter with an ellipses (`...`):
  - `func sum(vals ...int) {}`
  - `vals` behaves like slice within function body.
- Inversely, we can spread a slices values such as `sum(values...)`.
- Most common use case is for string formatting (`Printf` and variants).
- To accept any values at all for its final arguments, can use `(finalargs ...interface{})`

## Deferred Function Calls
