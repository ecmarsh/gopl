# 1. Introduction

## Setup Utils

Auto import: `$ go get golang.org/xtools/cmd/goimports`

## Printf common verbs

Verb | Spec
--- | ---
`%d` | decimal integer
`%x, %o, %b` | integer in hexadecimal, octal, binary
`%f, %g, %e` | floating-point number: 3.141593, 3.141592653589793 3.141593e+00
`%t` | boolean
`%c` | rune (Unicode code point)
`%s` | string
`%q` | quoted string "abc" or rune 'c'
`%v` | any value in a natural format
`%T` | type of any value
`%%` | literal percent sign (no operand)

- Note new lines and tabs not written by default: use escape sequences or use `%v.`

## Packages

- Packages in stdlib at [https://golang.org/pkg](https://golang.org/pkg)
- Packages from community at [https://godoc.org](https://godoc.org)

`go doc` is similar to `man`: `$ go doc http.ListenAndServe`
