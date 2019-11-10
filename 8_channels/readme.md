# Channels and Goroutines

This chapter covers gorutines and channels, which support _communicating sequential processes_, CSPs, which is a model of concurrency in which values are passed between independent activities (goroutines), but variables are for the most part, confined to a single activity. For more on aspects of the traditional model of _shared memory multithreading_, see [Concurrency with Shared Variables](../9_concurrency/readme.md).

## Goroutines

- A goroutine is just a concurrently executing activity.
- Given there are two functions, a sequential program may call one function, then another, while in a concurrent program, calls to both functions may be active at the same time.
- Although they are different, for archetypal thinking, a gourtine is similar to a thread. Differences between the two are only quantitative and explored more in the concurrency section.
- When a program starts, its first goroutne is the main goroutine.
- Then new goroutines are created by the `go` statement`. A go statement causes the function to be called in a newly created goroutine and the go statement itself completes immediately.
- See [clock](./clock) for first example using one goroutine per connection.
- See [reverb](./reverb) for second example using multiple goroutines per connection.

## Channels

- TODO
