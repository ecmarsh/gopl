# 3. Basic Data Types

1. **Basic types**
2. Aggregate types (arrays and structs)
3. Reference types (incl pointers, slices, maps, fns, channels all refer to program variables indirectly)
4. Interface types

## Integers

- 4 Sizes: 8, 16, 32, 64 bit integers
- Signed and unsigned (unit)
- Use normal int or unit to let compiler make choice based on 32-bit/64-bit hardware.
- `rune` is synonym for `int32`
- `byte` is synonym for `uint8`, but used to emphasize value is piece of raw data rather than numeric value.

### Binary Operators

In order of precedence:

1. `* / % << >> & &^`
2. `+ - | ^`
3. `== != < <= > >=`
4. `&&`
5. `||`

- Operations with same level of precedence associate to left.
- Note the mod operator can only be applied to int operations.
- Mod operator does not use true mod like python. (just flips sign).
- `/` operator similar to java (depends on operand type).
- Note overflown bits are **silently** discarded.
- Rune literals can be written as a character within **single** quotes. (can be printed with %c or %q)

## Floating-Point Numbers

- Prefer float64 vs float32 since arithmetic can accumulate quickly. 32 is ~6 decimal points and 64 is ~15.  
- Scientific notation is allowed (e.g 6.26e-34)
- Common artihmetic error results (z is float64(0):
  - z: 0
  - -z: -0
  - 1/z: +Inf
  - -1/z: -Inf
  - z/z: NaN

Note: Complex numbers are available.

### NaN

- Check with math.IsNaN()
- NaN comparisons always yield false except with `NaN != NaN`

## Booleans

- Short circuiting is valid in Go.
- Remember && has higher precedence than ||.
- No implicit conversion from booleans to numeric variables, or vice versa.
  - See `../btoi` for impelementation if needed.

## Strings

- Immutable
- Use slices to access characters/substrings (a copy is created).

### String Literals

- Double quotes `"`
- Go source files are encoded in UTF-8, so can use escape characters within string literals.

### Raw String Literals

- Backticks <code>`<code>
- Everything is taken literally (new lines, backslashes, etc.). No escape sequences.
- Useful for regex, documentation, HTML templates, JSON literals, etc.

### Unicode (Review)

- US-ASCII uses 7 bits to represent 128 "characters" incl upper/lower/digits/punc/control chars.
- Unicode v 8 defines code points for over 120,000 characters. Each one has a standard number called a unicode code point. In Go, Unicode Code Point = rune.

### UTF-8

- UTF-8 is a variable length encoding of Unicode points.
- Uses 1-4 bytes for each rune (most 2-3), but only 1 byte for ASCII characters. 
- No embedded NUL bytes.
- `unicode/utf8` package provides functions for encoding and decoding runes as bytes.
- In go, can use `\uhhhh` (16-bit) or `\uhhhhhhhh` (32-bit) where `h` is a hexidecimal digit.

### Strings and Byte Slices

- Relevant packages include `bytes`, `strings`, `strconv` and `unicode`
- `path/filepath` provides package for manipulating hierarchal names.
- While strings immutable, elements of a byte slice can be freely modified: `b := []byte(str)`
- For string builder, use `bytes.buffer`. Methods parallel strings lib.
  - Use `bb.WriteByte(char)` for ASCII and `bb.WriteRune(rune)` for Unicode. `bb.String()` to complete.

### Conversion with Numbers

- `strconv` package.
- We can use `fmt.Sprintf` to format, or the library function like `strconv.Itoa(n)`
- Another useful is `strconv.ParseInt(str, base, maxbase)`

## Constants

- Every constant's underlying type must be basic type (boolean, string, or number)
- const declarations prevent values from being changed.
- const values are known to compiler at run time so can use in types.
- Can emit an assignment and previous assignment will fall through:
```go
const (
  a = 1
  b
  c = 2
  d
)
fmt.Println(a, b, c, d) // "1 1 2 2"
```
### `iota` Constant Generator

- Used to create sequence of values w/o spelling out directly.
- Initialize first one to `iota` which will be 0 and rest get automatically incremented by 1.
- Note: this is essentially what `enums` (ts/java) is.
- Useful for creating flags / bit masks.

### Untyped Constants

- Untyped constants can have up to 256 bits of precision.
- Six types: untyped boolean, untyped integer, untyped rune, untyped floating-point, untyped complex, and untyped string.
- Only constants can be untyped. When used in program, they are implicitly converted to the other type, as long as target type can represent the original value (eg rounding for real and complex floats)






  


