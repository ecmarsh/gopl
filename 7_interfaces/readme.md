# Interfaces

- Interface types express generalizations or abstractions about the behaviors of other types.
- This allows us to write functions that are more flexible and adaptable because they are not tied to the details of one particular implementation.
- What makes Go's interfaces stand out is that they are _satisfied implicitly_, meaning there's no need to declare all the interfaces that a given concrete type satisfies; possessing the necessary methods is enough.
- This allows us to create new interfaces that are satisfied by existing concrete types without changing the existing types (useful for packages you don't control).

## Interfaces as Contracts

