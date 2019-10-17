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

### 5 Common Error-Handling Strategies

1. Propogate the error: make the failure in a subroutine the failure of the calling routine:
    - Note `fmt.Errorf` formats an error message using `fmt.Sprintf` to return a new error value - very useful for providing a clear causal chain from root problem to overall failure.
    - As a best practice, message strings should _not_ be capital and _not_ include newlines as error logs often contain many errors. This makes easier to self contain messages when searching using `grep`, etc.
    - As an example of specificity, if reading many links, files, include the name of link/file etc where the failure occured in the error message.
    - Example
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

2. Retry the failed operation, mainly for errors that represent trasient or unpredicatable problems. 
    - Of course, set a time or attempt limit using this strategy.
    - See [`waitforserver.go`](./waitforserver/main.go) for example.

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

5. Least commonly, safely an ignore the entire error
    - Best practice to always consider errors after function call.
    - If using this method, best to document intention of ignoring clearly.
    - Example:
      ```go
      dir, err := ioutil.TempDir("", "scratch")
      if err != nil {
        return fmt.Errof("failed to create temp dir: %v", err)
      }
      // ...use temp dir..
      os.RemoveAll(dir) // ignore errors; $TMPDIR is cleaned periodally
      ```

**Note**: Similar to _nodejs_ style, handle failures before success.

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
  // ...use r...
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

### Gotchas: Capturing Iteration Variables

- Sometimes Go's lexical scope rules can cause surprising results.
- To capture an iteration variable, assign it to a different variable in the enclosing scope first:
```go
var rmdirs []func()
for _, d := range tempDirs() {
  dir := d // The necessary assignment 
  os.MkdirAll(dir, 0755) // umask 022
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

- Create a deferred function by prefixing function call with `defer`:
  - `defer resp.Body.Close()`
- Function and argument expressions evaluated when statement is executed, but actual call is _deferred_ until the function that contains the `defer` statement has finished, whatever branch it took.
- Any number of calls may be deferred
  - Note they are executed in the _reverse_ order in which they were deferred.
- Often used with paired operations like open/close, disconnect/connect, lock/unlock, etc. to ensure resources are released in all cases, no matter how complex the control flow.
- See example in [`defer/title`](./defer/title.go)
  - First version requires a duplicate call to close connection in all paths.
  - Second calls deferred close to execute it after everything else has finished.
- Deferred functions have access to outer scope still.
  - For functions with many return statements, a handy trick is to defer an anonymous defer statement within a function with a named result:
  ```go
  func double(x int) (result int) {
    defer func() { fmt.Printf("double(%d) = %d\n", x, result) }()
    return x + x
  }
  _ = double(4) // Stdout: "double(4) = 8"
  // Note the result value is updated to x+x in return,
  // then it is called with new result.
  ```
- Be careful when using `defer` in a loop, since it won't be executed until all values in the loop have completed.
  - For example, if looping through files and deferring a close, could run out of file descriptors before others have closed.
  - One possible solution to this is to move the inner loop logic into another function, and call that function within the loop. Then defer will be called after each iteration.

## Panic

- When mistakes are detected at runtime (such as out-of-bounds array access or nil pointer dereference), it _panics_.
- Typical panic stops normal execution, _executes the deferred calls_, then crashes with log message which includes the _panic value_ and stack trace.
- There is also a built-in `panic` function that can be called directly.
  - This is most useful and best practice to use when an "impossible" situation happens and execution reaches a case that cannot logically happen:
  ```go
  switch s := suit(drawCard()); s {
  case "Spades":   //...
  case "Hearts":   //...
  case "Diamonds": //...
  case "Clubs":    //...
  default:
    panic(fmt.Sprintf("invalid suit %q", s)) // Joker?
  }
  ```
- Unless providing a meaningful message for panic that runtime panic will catch anyway, best to just let runtime panic.
- Panic is similar to an "exception", but since always causes a program crash, should only be used for grave errors. Otherwise use `error` values.
  - There are cases when may need to separate panic errors from errors. For example, if regular expression receives input that can't be compiled, it panics, but for an incorrect pattern, it gives error value.
- As a diagnostic tool, can defer a call to dump the stack ([`runtime pkg`](https://golang.org/pkg/runtime)) in main in order to dump the stack in case of a panic:
```go
func main() {
  defer printStack() {
      defer printStack()
      f(3)
  }
  func printStack() {
    var buf [4096]byte
    n := runtime.Stack(buf[:], false) // Returns bytes written
    os.Stdout.Write(buf[:n])
  }
```
- Giving up in a panic is usually correct, but there are some cases when we can recover, or at least clean up.

## Recover

- Recovering from panic is most useful for cleaning up or getting some useful error messages rather than an immediate crash. Cleanup examples include closing a connection that wasn't deferred, etc.
- Another example is if panic occurs in some odd corner case, we can convert panic into ordinary error value, print an extra message to alert, and continue:
```go
func Parse(input string) (s *Syntax, err error) {
  defer func() {
    if p := recover(); p != nil {
      err = fmt.Errof("internal error: %v", p)
    }
  }()
  // ...parser...
}
```
- Be careful when converting panics (be sure to cause some kind of alert), as it might cause lurking bugs to go unnoticed.
- Some best practices:
  - Don't try to recover from another package's panics (public APIs should report failures as `errors`)
  - Don't recover from a apanic that may pass through a function you don't maintain since you cannot reason about its safety.
  - Overall, `recover` _very_ selectively, if at all.
- See [recover example](./recover/soletitle.go) for how you might check a panic value and use recover, although this example isn't a best practice.
- There are some panics that cannot be recovered from, fatal errors, such as running out of memory.
