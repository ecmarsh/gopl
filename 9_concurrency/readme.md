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

- TODO
