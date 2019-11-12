# Channels and Goroutines

This chapter covers gorutines and channels, which support _communicating sequential processes_, CSPs, which is a model of concurrency in which values are passed between independent activities (goroutines), but variables are for the most part, confined to a single activity. For more on aspects of the traditional model of _shared memory multithreading_, see [Concurrency with Shared Variables](../9_concurrency/readme.md).

## Goroutines

- A goroutine is just a concurrently executing activity.
- Given there are two functions, a sequential program may call one function, then another, while in a concurrent program, calls to both functions may be active at the same time.
- Although they are different, for archetypal thinking, a gourtine is similar to a thread. Differences between the two are only quantitative and explored more in the concurrency section.
- When a program starts, its first goroutne is the main goroutine.
- Then new goroutines are created by the `go` statement`. A go statement causes the function to be called in a newly created goroutine and the go statement itself completes immediately.
- See [clock](./clock) for first example using one goroutine per connection, and introduces the [net package](https://golang.org/pkg/net) which provides components for building networked client and server programs that communicate over TCP, UDP, or Unix domain sockets.
- See [reverb](./reverb) for second example using multiple goroutines per connection.
- Be sure to concurrency chapter for consideration if it is safe to call methods of net.Conn currently, which is not true for most type.

## Channels

- If goroutines are the activities of a concurrent Go program, _channels_ are the connections between them.
- A channel is a communication mechanism that lets one goroutine send values to another goroutine.
- Each channel is a conduit for values of a particular type, which we call the channel's _element type_. So if a channel whose elements have type int, would be `chan int`.
- To create a channel, use the built-in `make` function:

```go
ch := make(chan int) // ch has type 'chan int'
```

- A channel is a reference to the data structure created by `make`, similar to maps, so when we copy a channel or pass one as an argument to a function, we are copying a reference to refer to the same structure. As a reminder, the zero-value is `nil`.
- Channels of the same type may be compared (`==`) and may also be compared to `nil`.
- The two principal operations of channels are `send` and `receive`, collectively called _communications_.
- Send transmits a value from one goroutine, through the channel, to another goroutine executing a corresponding receive expression. Both expressions are written using the `<-` operator.
- In the send statement, the `<-` separates the channel and value operands. In receive, `<-` precedes the channel operand. Note a receive expression whose result is not used is a valid statement.

```go
ch <- x  // a send statement
x = <-ch // a receive statement in an assignment statement
<-ch     // a receive statement; result is discarded
```

- The third operation of channel is _close_ (`close(ch)`), which sets a flag indicating that no more values will ever be sent on this channel and subsequent attempts to send will panic.
- Receive operations on a closed channel will yield the values already sent until no more values are left; any receive operations after that complete immediately and yield the zero value of the channel's element type.
- A channel created with a simple call to `make` is called an _unbuffered_ channel, but `make` accepts an optional second argument, int _capacity_. When capacity is non-zero, `make` creates a _buffered_ channel.

### Unbuffered Channels

- A send operation on an unbuffered channel blocks the sending goroutine until another goroutine executes a corresponding receive on the same channel, then the value is transmitted and both goroutines may continues.
- On the other hand, if the receive operation was attempted first, receiving goroutine is blocked until another goroutine performs a send on the same channel.
- Communication over unbuffered channel causes sending and receiving goroutines to _synchronize_, so we sometimes call unbuffered channels _synchronous_.
  - When a value is sent on an unbuffered channel, receipt of value happens *before* the reawakening of the sending goroutine.
- Side note on currency: "x happens before y" doesn't mean merely earlier in time, it means that it is guaranteed to do so and you may rely on the fact that all its prior effects, such as updates to variables, are complete.
  - When "x is concurrent with y", x doesn't happen either before or after; they aren't necesarily simultaneous, but we just can't assume anything about ordering.
- Each message sent over a channel has a value, but sometimes communication and moment it occurs is just as important; we call messages _events_ to stress the timing aspect.
- When an event has no additional information (its sole purpose is synchronization), its emphasized by using a channel whose element type is `struct{}` through its common use of a channel `bool` or `int` for same purpose since `<- 1` is shorter than `done <- struct{}{}`.

#### Example

Example of netcat except we make the program wait for the background goroutine to complete before exiting by using an unbuffered channel to sncyrhonize the two goroutines:

```go
// netcat3
func main() {
  conn, err := net.Dial("tcp", "localhost:8000")
  if err != nil {
    log.Fatal(err)
  }
  done := make(chan struct{})
  go func() {
    io.Copy(os.Stdout, conn) // NOTE: ignoring errors
    log.Println("done")
    done <- struct{}{} // signal the main goroutine
  }()
  mustCopy(conn, os.Stdin) // see clock or reverb for impl of mustCopy
  conn.Close()
  <-done // wait for background goroutine to finish
}
```

- In the above example, when user closes stdin stream, `mustCopy` returns and the main goroutine calls `conn.Close()`, which closes both halves of network connection.
  - Closing the write half causes server to see an EOF condition.
  - Closing the read half causes background goroutine's call to `io.Copy` to return a "read from a closed connection" error, which is why errors are ignored.
- Before returning, background goroutine logs a message and sends a value on the `done` channel, then the main goroutine waits until it has received the value before routining; as a result, program always logs the "done" message before exiting.

### Pipelines

- A pipeline is a channel used to connect goroutines together so that the output of one is the input to another.

**Example:** 3-state pipeline: Counter --naturalnums--> Squarer --squares-->Printer

Counter generates integers and sends them over channel to second goroutine `squarer` which receives each value, squares it, and sends the result over _another_ channel to a third goroutine, `printer`, which receives values and prints them.

```go
// Simple example prints infinite series of squares 0, 1, 4, 9,...
func main() {
  naturals := make(chan int)
  squares := make(chan int)

  // Counter
  go func() {
    for x := 0; ; x++ {
      naturals <- x
    }
  }()

  // Squarer
  go func() {
    for {
      x := <-naturals
      squares <- x * x
    }
  }()

  // Printer (in main goroutine)
  for {
    fmt.Println(<-squares)
  }
}
```

- A similar layout is used in long-running server programs wher channels are used for lifelong communication between goroutines containing infinite loops.
- If the sender knows that no further values will ever be sent on a channel, it is useful to communicate to receiver goroutines to stop waiting by _closing the channels with `close(naturals)`.
- After channel is closed, further send operations cause panic.
- After closed channel is _drained_ (last sent element received), all subscequent operations proceed without blocking but yield zero value.
- There is no way to test directly whether a channel has been closed, so we use a variant of receive operation that produces two results: received channel element, plus a boolean value, conventionally called `ok`, which is `true` for successful receive and `false` for a receive on a closed and drained channel:

```go
// Exemplifying variant of testing for closed/drained channel
go func Squarer() {
  for {
    x, ok := <-naturals
    if !ok {
      break // channel was closed and drained
    }
    squares <- x * x
  }
  close(squres)
}()
```

- We can also use a range loop, which is a more convenient syntax for receiving all the values sent on a channel and terminating the loop after the last one. Example after receiving 100 items:

```go
// Alternate version of Counter in pipeline example
go func Counter() {
  for x:= 0; x < 100; x++ {
    naturals <- x
  }
  close(naturals)
}()
go func Squarer() {
  for x := range naturals {
    squares <- x * x
  }
  close(squares)
}()
// Printer with range
for x := range squares {
  fmt.Println(x)
}
```

- Note it's only necessary to close a channel when it is important to tell the receiving goroutines that all data has been sent. A channel that the garbage collector determines to be unreachable will have its resources reclaimed whether or not its closed. Note: this is not the same as `Close()` for files -- files need to always be closed.
- Attempting to close an already-closed channel or a nil channel will cause a panic.
- See [cancellation](./readme.md#Cancellation) for another closing use as a broadcast mechanism.

### Unidirectional Channel Types

- The `squarer` function in the middle of the above examples takes input and output of the same type, but their intended use cases are opposite (in recieved from, out sent to). Note `in` and `out` are used by convention to convey that intention, but nothing actuallly prevents squarer from sending to in or receiving from out.
- When a channel is a suppplied as a funciton parameter, it is nearly always with the intent that it be used exclusively for sending or exclusively for receiving.
- To document primary in and out, Go provides _unidirectional_ channel types that expose only one of either the send or receive operations, where violations are detected at compile time:

Syntax | Type | Desc
--- | --- | ---
`chan<- T` | send-only | Allows sends, but not receives.
`<-chan T` | receive-only | Allows receives, but not sends.

```go
func counter(out chan<- int) {
  for x := 0; x < 100; x++ {
    out <- x
  }
  close(out)
}

func squarer(out chan<-int, in <-chan int) {
  for v := range in {
    out <- v * v
  }
  close(out)
}

func printer(in <- chan int) {
  for v := range in {
    fmt.Println(v)
  }
}

func main() {
  naturals := make(chan int)
  squares := make(chan int)

  go counter(naturals)
  go squarer(squres, naturals)
  printer(squres)
}
// counter(naturals) implicitly converts naturals (type chan int) to the
// type of the paramer, chan<- int and the printer(squares) does a similar conversion.
```

- Note conversions from bidirectional to unidirectional channel types are permitted in any assignment, but once you have a value of a unidirectional type such as chan<- int, there is no way to obtain from it a value of type `chan int` that refers to the same channel data structure (no going back).

### Buffered Channels

- A buffered channel has a queue of elements, where the queue's maximum size is determined upon creation, by providing it as the second argument to `make`.

```go
// creates a buffered channel capable of holding 3 string values
ch = make(chan string, 3)
// ch -> |s, s, s|
```

- Send operations on buffered channel inserts element at the _back_ of queue, and receive pops from the front.
- If channel is full, send operation blocks goroutine until space is made available by another goroutine's receive.
- Filling the channel: 3x `ch <- "text"`. Receive one value: `fmt.Println(<-ch)`. Now channel is neither full nor empty ("partially full buffered channel") so send operation or receive proceeds without blocking. In this way, the channel's buffer decouples the sending and receiving goroutines.
- Note `len(ch)` is number of items in channel, and `cap(ch)` is capacity of channel. Note `len` is likely to be stale by time received in a concurrent program though.
- Normally send and receive operations are performed by different goroutines. Even for simple ones, this should be done. If all you need is a simple queue, then use a slice instead.

Example application of a buffered channel that makes parallel requests to three _mirrors_, or equivalent but geographically distributed servers. It sends responses over a buffered channel, then receives and returns only the first response, which is the quickest one to arrive, so `mirroredQuery` returns a result even before the two slower servers have responded.

```go
func mirroredQuery() string {
  responses := make(chan string, 3)
  go func() { responses <- request("asia.gopl.io") }()
  go func() { responses <- request("americas.gopl.io") }()
  go func() { responses <- request("europe.gopl.io") }()
  return <-responses // return the quickest response
}
func request(hostname string) (response string) { /* ... */ }
```

- In above example, if used an unbuffered channel, two slower gourtines would have been stuck trying to send their responses on a channel where no goroutine will ever receive, called a _goroutine leak_.
  - Note leaked goroutines, unlike garbage variables, are not automatically collected so remember to make sure they terminate themselves when no longer needed.
- Unbuffered channels give stronger synchronization guarantees because every send operation is synchronized with its corresponding receive.
- Buffered channels are decoupled. When an upper bound on number of values that will be sent on a channel is known, it's not unusual to create a buffered channel and perform all the secds before the first value is receive.
- Failure to allocate sufficient buffer capacity causes a program to deadlock.
- Keep in mind performance for channel buffering. If one item ahead of another is slower, it will slow down the ones behind it, so may be useful to introduce a second to get second up to speed of the first (assembly-line metaphor).

## Looping in Parallel

- Problems where order does not matter, consisting entirely of subproblems that are completely independent of each other are called _embarassingly parallel_.
- Embarassingly parallel problems are easiest kind to implement concurrently and enjoy performance that scales linearly with amount of parallelism.
- At its simplest, remember if using variables, to give the goroutine its own block scope:

```go
// example that looks to create thumbnails from filenames in parallel
for _, f := range filenames {
    go func() {
        thumbnail.ImageFile(f) // Note: We need to handle errors
    }()
}
```

- In order to know the error, we need to return values of each goroutine to the main one.
- To prevent blocking, we can use a buffered channel to return names of generated image files
 along with any errors. However, this makes it difficult to predict number of loop iterations.
- The solution below uses `sync.WaitGroup` which allows us to know when the last goroutine has
 finished (which may not be the last one to start) acting as a special counter that increments
  before each goroutine starts and decrements it as it finishes. This structure is common
   idiomatic pattern for looping in parallel when we don't know the number of iterations:

```go
// makeThumbnails makes thumbnails for each file received from the channel.
// It returns the number of bytes occupied by the files it creates.
func makeThumbnails(filenames <-chan string) int64 {
    sizes := make(chan int64)
    var wg sync.WaitGroup // number of working goroutines
    for f := range filenames {
        wg.Add(1) // increment count of active goroutines
        // worker
        go func(f string) {
            defer wg.Done() // decrement counter when goroutine finished
            thumb, err := thumbnail.ImageFile(f)
            if err != nil {
                log.Println(err)
                return
            }
            info _ := os.Stat(thumb) // OK to ignore error
            sizes <- info.Size() 
        }(f)
    }
    // closer
    go func() {
        wg.Wait()
        close(sizes) // close channel after all finished
    }()
    var total int64
    for size := range sizes {
        total += size
    }
    return total
}
```

- Above, `Add` (wait group incrementer) must be called before goroutine starts.
- We defer `Done` to ensure counter is decremented even in the error case.
- The sizes channel carries each file size back to main goutine to compute the sum of bytes.
- Notice that the closer gourtine waits for the workers to finish _before_ closing the `sizes
` channel; wait and close must be concurrent with the loop over `sizes`.
    - If wait operation were placed in main goroutine before the loop, it would never end; if
     placed after, it would be unreachable since loop would never terminate since nothing closing
      the channel.
- See [concurrent web crawler](./crawl/), a common concurrency pattern (often asked in
 interviews) for more on looping in parallel.
 
 ## Multiplexing with `select`
 
 - 
