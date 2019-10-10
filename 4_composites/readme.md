# 4. Composite types

## Arrays

- Arrays are fixed length. Initial length must be a constant.

```go
var a [3] int
n := [3]int{1, 2, 3} // cannot assign different length arrays now
r := [...]int{99: -1} // assigns 100 element array with 99 0's and the last -1
```

- If array's element types are comparable, arrays can be compared with `==` or `!=`.
- This can be useful such as comparing arrays of bytes.
- Except for special cases such as SHA256's fixed-size hash, arrays are less preferrable as function parameters or results than slices.

## Slices

- Slice is like a dynamic array, which gives access to its underlying array.
- If slice were a strict, would resemble:

```go
type IntSlice struct {
  ptr      *int
  len, cap int
}
```

- Slices capacity must be able to handle amount of new elements before appending.
- Using `[#:#]`, unlike python, changes the reference, not create a copy. Use `copy` to copy.
- To initialize a stack, use `make([]T, initialLen, cap)` or define values with `[]T{vals...}`

## Maps

`m := make(map[kType]vType))`

- Go map is a reference to hash table.
- All values must be same type. But keys and values can be different types.
- Keys can be any comparable type.
- **Note:** Cannot take address of map element as tradeoff of dynamic map is new storage locations may be assigned to support growing or refreshing of elements.
- Zero value for map is nil.
- Map values are initialized to zero value, similar to python's `defaultdict`:
  ```go
  m := make(map[int]int)
  m[1] += 1
  m[1] += 1
  m[2] += 1
  // m: {1->2, 2->1}
  ```
- Not unusual to use map as a set. For example, a set of strings might be `map[string]bool`, but ensure its being used a set before assuming.

## Structs

- Groups together zero or more named values (called _fields_) of aritrary types a single entity.
- To initialize:
```go
type StructName struct {
  Field1   type
  Field2   type
  F3, F4   type // Fields may be combined for same type, but typically combined only if related.
}
var objName StructName
```
- Note field order is significant to identify. That is, changing order of fields defines different struct type.
- A struct *field* is exported if it begins with capital letter (following Go's control mechanism), and may contain a mixture of private and public fields. 
- Structs may not declare a field of same type, (similar to arrays) but may declare a field of pointer type *Struct:
```go
// Recursive structure may be defined with pointer fields
type S struct {
  ID      int 
  sibling *S
}
```
- Zero value for struct is composition of zero values for its fields.
- Note: some go programmers use empty structs as set similar to map to indicate only keys are significant, but tradeoff between savings and syntax is arguable so typically not best practice.
- Structs are comparable with `==` or `!=` if all fields of a struct are comparable.
- Comparable struct types may also be used as the key of a map.

### Struct Assignment/Access

- Fields can be accessed/assigned using dot notation `objName.Field1 ...` or taking address and using pointer:
```go
field2 := &objName.Field2
*field2 = "Prepended " + *field2 
```
- Dot notation may also be used with a pointer to struct:
```go
var objAlias *StructName = &objName
// Following statements are equivalent
objAlias.Field1 = ...
(*objAlias).Field1 = ...
```

- Note to assign a field, lefthand side must either be a pointer or a variable. For example:
```go
// Finds struct by id. Returns a pointer to struct.
func StructByID(id int) *StructName { ... }
id := objName.ID // assume ID is a field in StructName and objName has been initialized
// This works because fn returns a pointer. If returned just StructName, compile error since StructName is not a variable.
StructByID(id).Field1 = ...
```

### Struct Literals

Two forms of struct literals, which _cannot_ be mixed (incl imports).

1) Order-based:  
  ```go
  type Point struct{ X, Y int }
  p := Point{1, 2}
  ```
  - requires a value be specified for every field, in the right order.
  - can become cumbersome of number of fields grow or order changes.
  - only use within the package that defines the struct type or in smaller types where their is an ordering convention such as `color.RGBA{red, green, blue, alpha}`

2) Name-based: `anim := gif.Gif{LoopCount: nframes}`
  - if field is omitted, set to zero value
  - with names provided, order does not matter (similar to `**kwargs in python`)

- Struct values can be passed as arguments to functions and returned from them:
```go
func Scale(p Point, factor int) Point {
  return Point{p.x * factor, p.Y * factor}
}
fmt.Println(Scale(Point{1, 2}, 5))
```
- Use pointer to struct for args/return if it is a larger struct for efficiency:
```go
func CalcBonus(e *Employee, percent int) {
  return e.Salary * percent / 100
}
// Pointer required to modify its argument so fn recives ref, not copy:
func AwardRaise(e *Employee) {
  e.Salary = e.Salary * 105 / 100
}
```

- Because structs are often handled with pointers, we can use shorthand to define and get address:
```go
pp := &Point{1, 2}
```

### Struct Embedding and Anonymous Fields

- Struct embedding lets us use one named struct type as an anonymouse field of another struct type.
- Can be useful as a shortcut like `x.f` to represent a chain of fields like `x.d.e.f`.

#### Example

Consider 2-D drawing programs that has library of shapes, including a Circle and Wheel:

```go
type Circle struct {
  X, Y, Radius int
}
type Wheel struct {
  X, Y, Radius, Spokes int
}

// Create a wheel
var w Wheel
w.X = 8
w.Y = 8
w.Radius = 5
w.Spokes = 20
```

Notice Circle and Wheel share common fields, and other shapes may as well, so it can be useful to factor our commonalities:

```go
// Refactor point set
type Point struct {
  X, Y int
}

type Circle struct {
  Center Point
  Radius int
}

type Wheel struct {
  Circle Circle // Note circle is embedded
  Spokes int
}

// Makes app clearer, but makes accessing embedded fields more verbose:
var w Wheel
w.Circle.Center.X = 8 
w.Circle.Center.Y = 8 
w.Circle.Radius = 5
w.Spokes = 20
```

To solve this, we can use _anonymous fields_, which must be a named type or a pointer to a named type:

```go
// Embedding with anonymous fields
type Circle struct {
  Point      // Point is type, NOT name
  Radius int
}

type Wheel struct {
  Circle    // Circle is type, NOT name
  Spokes int
}

// Now we can omit subfield names
// Note fields still have names, so we could use explicit forms (commented)
var w Wheel
w.X = 8       // w.Circle.Point.X = 8
w.Y = 8       // w.Circle.Point.Y = 8
w.Radius = 5  // w.Circle.Radius = 5
w.Spokes = 20
```

Note there is no corresponding syntax for struct literal syntax. We cannot use ordered or field based form.

- Compile errors for unknown fields:
  - `w = Wheel{8, 8, 5, 20}`
  - `w = Wheel{X: 8, Y: 8, Radius: 5, Spokes: 20}`

See [package embed](./structs/embed/) for correct use examples.

Lastly, anonymous fields do not _need_ to be struct types; any type or pointer to a named type will do. This is how we can compose simpler methods into complex object behavior. See [method section](../6_methods/) for more.

## JSON (JavaScript Object Notation)

- JSON is most widely used standard notation for sending and receiving structured information.
  - XML and ASN.1 (Google's Protocol Buffers) are others that serve similar purpose.
- Encodes JavaScript values as Unicode text. Can also represent Go's basic data types and composite types (above).
- Basic JSON types are:
  - numbers (decimal or scientific notation)
  - booleans 
  - strings (sequence of Unicode code point in _double_ quotes)
    - strings have similar backslash escapes (`\uhhhh`), which denote UTF-16 codes, not runes.
- Basic types can be combined recursively using JSON arrays and objects.
- Examples:

Type | Syntax
---- | ------
boolean | `true`
number | `-3.145` 
string | `"She said \"Hello, 你好\""`
array | `["gold", "silver", "bronze"]`
object | `{"year": 2020,`
 ..    |   `"event": "shotput",`
 ..    |   `"medals": ["gold", "silver", "bronze"]}`

### Marshaling

- Go data structures may converted to JSON through _marshaling_ via:
  - `json.Marshal(obj)` (struct instances are ideal for `obj`)
  - `Marshal` eliminates whitespace, producing compact, but unreadable output.
- To produce formatted output, use:
  - `json.MarshalIndent(obj, linePrefix(eg ""), indentString(eg "    "))`
- Marshal obtains JSON key names through struct field names via _reflection_. **Note only exported fields are marshaled**.
- JSON may be modified through _field tags_, a string of metada associated at compile time with the field of a struct, for example:
```go
type Movie struct {
  Title  string
  Year   int `json:"released"` // Changes JSON key to "released"
  Color  bool `json:"color,omitempty"` // Does not output color field if has types zero value
}
```
- The JSON key controls behavior of the [`encoding/json` package](https://golang.org/pkg/encoding/json/). Other `encoding/..` packages follow this convention.
- Inverse to marshaling is _unmarshaling_, via `json.Unmarshal(jsonData, &structfieldToUnmarshal)`
- Example:
```go
var titles []struct { Title string }
if err := json.Unmarshal(data, &titles); err != nil {
  log.Fatalf("JSON unmarshaling failed: %s", err)
}
fmt.Println(titles) // "[{Movie Title 1}, {Movie Title 2}, {Movie Title 3}]"
```
- This is useful for using data obtained via web APIs, which often return JSON responses.

### Text and HTML Templates

- [`text/template` package](https://golang.org/pkg/text/template/) and [`html/template` package](https://golang.org/pkg/html/template/) provide mechanisms for substituting values of variables into a text or HTML template (such as more advanced formatting of JSON).
- Can use template strings (similar to handlebars, jinja, etc) to format output, for example:
```go
const templ = `{{.TotalCount}} issues:
{{range .Items}}--------------------------------
Number: {{.Number}}
User:   {{.User.Login}}
Title:  {{.Title | printf "%.64s"}}
Age:    {{.CreatedAt | daysAgo}} days
{{end}}`
// dot initially refers to templates parameter (expands)
// range creates a loop (with internal dot bound to range input),
// and loop end denoted with last line
// note piping also works. In templates, printf is analagous to fmt.Sprintf
```
- To produce the output of a template string:
  1) Parse template into suitable internal representation. (only needs to be done once)
  2) Execute it on specific inputs.
  - Example:
  ```go
  report, err := template.New("report"). // Creates and returns template
    Funcs(template.FuncMap{"daysAgo": daysAgo}). // adds daysAgo fn to available fns
    Parse(templ)
  if err != nil {
    log.Fatal(err) // Template parse failure is fatal bug. Also see template.Must
  }
  ```

- The same logic can be used for HTML templates, which additionally has features for automatic and context-approriate escaping of strings within HTML, JS, CSS, or URLs to prevent against inj attack.
```go
import "html/template"

var issueList = template.Must(template.New("issuelist").Parse(`
  <h1>{{.TotalCount}} issues</h1>
  <table>
  <tr style='text-align: left'>
    <th>#</th>
    ...
  </tr>
  {{range .Items}}
  <tr>
    <td><a href='{{.HTMLURL}}'>{{.Title}}</a><td>
    ...
  </tr>
  {{end}}
  </table>
`))
```
- If variables contain (wanted) HTML syntax, we can surpress escaping by using `template.HTML` instead of `string`. See [`autoescape`](./autoescape/) for difference demonstration.
