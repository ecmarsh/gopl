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

- 


