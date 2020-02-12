# Concurrency with Shared Variables

- Recap: When we cannot confidently say that one event _happens before_ the other, then the events _x_ and _y_ are _concurrent_.

## Race Conditions

- A function is _concurrency-safe_ if it is correct sequentially and it continues to work correctly even when called concurrently (ie called from two or more goroutines with no additional synchronization).
  - A type is concurrency-safe if all its accessible methods and operations are concurrency-safe.
- We avoid concurrent access to most variables either by _confining_ them to a single goroutine or by maintaining a higher-level invariant of _mutual exclusion_.
- Reasons a function might not work when called concurrently include deadlock, livelock, and resource starvation. The most important one for this chapter is _race conditions_.
- A **race condition** is a situation in which the program does not give the correct result for some interleavings of the operations of multiple gourtines. They are danagerous because they may remain latent in a program and appear infrequently (eg under heavy load or when using certain compilers, platforms, or archs) making them hard to reproduce/diagnose.
- The classic example of bank account handling transactions is a type of condition called a **data race**, which occurs when _two gourtines access the same variable concurrently and at least one of the accesses is a write_.
- Following this definition, there are three ways to avoid a data race:

1. Don't write the variable. In the example map below, which is lazily populated as each key is requested for the first time, if `Icon` is called sequentially, the program works fine, but if `Icon` is called concurrently, there is a data race accessing the map:

```go
var icons = make(map[string]image.Image)

func loadIcon(name string) image.Image

// NOTE: not concurrency-safe!
func Icon(name string) image.IMage {
  icon, ok := icons[name]
  if !ok {
    icon = loadIcon(name)
    icons[name] = icon
  }
  return icon
}
```

Instead, we can initialize the map with all necessary entries before creating additional goroutines and never modify it again, so any number of goroutines may safely call `Icon` concurrently since each one only reads the map:

```go
var icons = map[string]image.Image {
  "spades.png": loadIcon("spades.png"),
  "hearts.png": loadIcon("hearts.png"),
  "diamonds.png": loadIcon("diamonds.png"),
  "clubs.png": loadIcon("clubs.png"),
}

// Concurrency-safe.
func Icon(name string) image.Image { return icons[name] }
```

2. Avoid accessing the variable from multiple goroutines (the approach taken in may of the examples in the [channels section](../8_channels/readme.md)). In otherwords, variables are _confined_ to a single goroutine.

The Go mantra "Do not communicate by sharing memory; instead, share memory by communicating" applies to this method where goroutines must use a channel to send the confining gourtine a request to query or update the variable. A gourtine that brokers access to a confined variable using channel requests is called a _monitor goroutine_ for that variable.

Bank example rewritten with the `balance` variable confined to a monitor goroutine, `teller`:

```go
// Package bank provides a concurrency-safe bank with one account.
package bank

var deposits = make(chan int) // send amount to deposit
var balances = make(chan int) // receive balance

func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }

func teller() {
  var balance int // balance is confined to teller goroutine
  for {
    select {
      case amount := <-deposits:
        balance += amount
      case balances <- balance:
    }
  }
}

func init() {
  go teller() // start the monitor goroutine
}
```

We can even apply this through multiple stages; if a variable is confined to one stage of the pipeline, then confined to the next, and so on, then essentially all accesses to the variable are sequential. This discipline pipeline is sometimes called _serial confinement_. Example below where `Cakes` are serially confined, first to the `baker` goroutine, then to the `icer` goroutine.

```go
type Cake struct{ state string }

func baker(cooked chan<- *Cake) {
  for {
    cake := new(Cake)
    cake.state = "cooked"
    cooked <- cake // baker never touches this cake again
  }
}

func icer (iced chan<- *Cake, cooked <-chan *Cake) {
  for cake := range cooked {
    cake.state = "iced"
    iced <- cake // icer never touches this cake again
  }
}
```

3. Allow many goutines to access the variable, but only one at a time, also called _mutual exclusion_, explored in the next section.

### Mutual Exclusion: `sync.Mutex`

- A counting semaphore that counts only to one is called a _binary semaphore_. This shares the same idea except we limit the channel to capacity 1 to ensure that at most one goroutine accesses a shared variable at a time.

```go
var (
  sema = make(chan struct{}, 1) // a binary sempahore guarding balance
  balance int
)

func Deposit(amount int) {
  sema <- struct{}{} // acquire token
  balance = balance + amount
  <- sema // release token
}

func Balance() int {
  sema <- struct{}{} // acquire token
  b := balance
  <-sema // release token
  return b
}
```

- Pattern of _mutual exclusion_ is useful that it is supported directly by the `Mutex` type from the `sync` package. Its `Lock` method acquires the token (called a _lock_) and its `Unlock` method releases it:

```go
import "sync"

var (
  mu      sync.Mutex // guards balance
  balance int
)

func Deposit(amount int) {
  mu.Lock()
  balance = balance + amount
  mu.Unlock()
}

func Balance() int {
  mu.Lock()
  b := balance
  mu.Unlock()
  return b
}
```

- In the above, the balance variable must call mutex's `Lock` in order to acquire an exclusive lock. If some other go routine acquired the lock, the operation will block until the other goroutine unlocks the variable again.
  - By convention, variables guarded by a mutex are declared immediatedly after the declaration of the mutex itself. (as in the var statement above).
- The region between lock and unlock (where the goroutine can freely read and modify the shared variables) is called a _critical section_.
- This arrangement of functions, mutex lock, and variables ic alled a _monitor_, meaning a broker that ensures variables are accessed sequentially.
- In more complex critical sections (especially ones in which errors must be dealt with by returning early), it can be hard to tell that calls to Lock and Unlock are strictly paired on all paths. In these scenarios, a deferred call to `Unlock` implicitly extends to the end of the current function so we don't remember to have to free the lock:
  - Another benefit of using the deferred `Unlock`, is that it will run even if the critical section panics, which may be important in programs that make use of `recover`.
  - Note that a defer call is marginally more expensive than just a call to `Unlock`, but arguably not enough to justify less clear code.

```go
func Balance() int {
  mu.Lock()
  defer mu.Unlock() // called after the return statement has read the value of balance
  return balance // note we no longer need the variable b anymore
}
```

- Note that mutex locks are not _re-entrant_: it's not possible to lock a mutex that's already locked as it leads to deadlock where nothing can proceed.
- For instance, with a non-atomic function below, it is tempting to try and lock an entire sequence (because if an excessive withdrawl is attempted and the balance dips below zero, may cause a concurrent withdrawal for modest sum to be be rejected). The problem with the function is that there are three separate operations, each of which acquires and then releases the mutex lock.

```go
// NOTE: incorrect! causes deadlock
func Withdraw(amount int) bool {
  mu.Lock()
  defer mu.Unlock()
  // NOT atomic!
  Deposit(-amount)
  if Balance() < 0 {
    Deposit(amount)
    return false // insufficient funds
  }
  return true
}
```

- A common solution is to divide a function such as `Deposit` into two; an unexported function `deposit`, that assumes the lock is already held and does the real work, and an exported function `Deposit` that acquires the lock before calling `deposit`:

```go
func Withdraw(amount int) bool {
  mu.Lock()
  defer mu.Unlock()
  deposit(-amount)
  if balance < 0 {
    deposit(amount)
    return false // insufficient funds
  }
  return true
}

func Deposit(amount int) {
  mu.Lock()
  defer mu.Unlock()
  deposit(amount)
}

func Balance() int {
  mu.Lock()
  defer mu.Unlock()
  return balance
}

// This function requires that the lock be held.
func deposit(amount int) { balance += amount }
```

- When you use a mutex, make sure that both it and the variables it guards are not exported, whether they are package-level variables or the fields of a struct. Encapsulation helps us maintain concurrency invariants.

### Read/Write Mutexes: `sync.RWMutex`

- If a function only needs to read (e.g if multiple requests to read a balance are coming in), we can use a special kind of lock that allows read-only operations to proceed in parallel with each other, but write operations to have fully exclusive access. This lock is called _multiple readers, single writer_ lock and is provided by `sync.RWMutex`:

```go
var mu sync.RWMutex
var balance int

func Balance() int {
  mu.RLock() // readers lock
  defer mu.RUnlock()
  return balance
}
```

- The above differs from the normal lock/unlock methods to acquire and release a writer or _exclusive_ lock.
- `RLock` can only be used if there are no writes to shared variables in the critical section.
- Note that you must be positive that there is absolutely no way that you will cause some sort of update (e.g a method that appears to be a simple accessor might also increment an internal usage counter or update a cache to repeat so that calls are faster). When in doubt, use an exclusive lock.
- It is only profitable to use an RWMutex when most of goroutines that acquire the lock are readers, and the lock is under _contention_, meaning goroutines have to wait to acquire it. The internal bookkeeping of an `RWMutex` makes it slower than a regular mutex for uncontended locks.

## Memory Synchronization

- Synchronization is not merely just the order of execution of multiple goroutines - it also affects memory.
- With multi-processor computers that each have their own local cache of main memory, writes to memory may be buffered within each processor and flushed out to main memory only whe necessary - and this may be in a different order than originally written by the writing gourtine. Syncrhonization primitives like channels and mutex cause the processor to flush out and commit all of its accumulated writes to that the effects of gourtine execution up to that point are guaranteed to be visible to goroutines running on other processes.
- Example showing the issue:

```go
var x, y int
go func() {
  x = 1                     // A1
  fmt.Print("x:", x, " ") // A2
}
go func() {
  y = 1                   // B1
  fmt.Print("y:", y, " ") // B2
}

// Possible Output:
y:0 x:1
x:0 y:1
x:1 y:1
y:1 x:1
```

- The zeros may come as a surprise; although one goroutine must observe the write to its associated variable, it does not necessarily observe the write to the other gourtine, so it may print a stale value for the other variable (0).
- All of these concurrency problems can be avoided by consistent use of simple, established patterns.
  - Where possible, confine variables to a single gourtine.
  - For all other variables, use mutual exclusion.

## Lazy Initialization: `sync.Once`

- It is good practice to defer expensive initialization steps until it is needed, especially if there is a possibility that the initialized will not be needed.

Example of lazy initialization from Icon example:

```go
var icons map[string]image.Image

func loadIcons() {
  icons = map[string]image.Image{
    "spades.png": loadIcon("spades.png"),
    "hearts.png": loadIcon("hearts.png"),
    "diamonds.png": loadIcon("diamonds.png"),
    "clubs.png": loadIcon("clubs.png"),
  }
}

// NOTE: not concurrency safe!
func Icon(name string) image.Image {
  if icons == nil {
    loadIcons() // one time initialization
  }
  return icons[name]
}
```

- The above is not concurrently safe because if two routines try to get an icon, one routine may think the icons is not nill and try to access it before it had actually been assigned. The simplest way to ensure this does not happen is to synchronize them using a mutex:

```go
var mu sync.Mutex // guards icons
var icons map[string]image.Image

// Concurrency-safe
func Icon(name string) image.Image {
  mu.Lock()
  defer mu.Unlock()
  if icons == nil {
    loadIcons()
  }
  return icons[name]
}
```

- However, in the above, two goroutines cannot access the variable concurrently, even once the variable has been initialized and never modified again. We could use an `RWMutex`, but this adds complexity and is error-prone.
- As an alternative, `sync` provides a specialized solution specifically for this problem called `sync.Once` which consists of a mutex and a boolean variable that records whether initialization has already taken place. It's sole method, `Do`, accepts the initialization function as its argument:

```go
var loadIconsOnce sync.Once
var icons map[string]image.Image

// Concurrency-safe
func Icon(name string) image.Image {
  loadIconsOnce.Do(loadIcons)
  return icons[name]
}
```

- In the above, each call to `Do(loadIcons)` locks the mutex and checks the boolean variable. The first call sets it to true so subsequent calls do nothing and it becomes visible to all goroutines.

## The Race Detector

- Adding the `-race` flag to `go build, run, test` causes the compiler to build a modified verison of app or test that records all accesses to shared variables that occured during execution, along with the identify of the gourtine that read/wrote the variable. Syncrhonization events (go statements, channel ops, mutex locks/unlocks, waitgroups, etc) are also recorded.
- The race detector only reports all data races that were executed, but it cannot prove that no races will ever occur. As a best practice, tests should exercise packages using concurrency. Concurrency programs require more overhead and time/memory to run, but overhead is tolerable and race detector can save hours or even days of debugging.
- See [concurrent non-blocking cache example](./memo) that addresses the problem of _memoizing_ a function, or cahcing the result of a function so it only needs to be computed once.
