# Packages and the _Go_ Tool

## Packages

### Introduction

- A package is just a distinct name space that encloses its identifiers.
- Go's compiler works quickly because:
  - imports are explicitly listed.
  - there are no cycles in dependencies, so can be compiled separately and even in parallel.
  - each package exports info about dependencies, so only needs to explore one depency level out.

### Import Paths

- Import paths should be globally unique.
- For packages not in standard library, should start with hoster or organization host/domain name and path to package, e.g., `github.com/go-sql-driver/mysql`.

### The Package Declaration

- Required at the start of every Go source file.
- Purpose is to determine the default identifier for that package (called the package name) when it is imported by another package.
- Conventionally, package name is the last segment of the import path. For example, every file in `math/rand` package starts with `package rand`. Note different packages may have the same name if their import paths differ.
- Import exceptions to the 'last segment' rule are:
  - A package defining a command (an executable Go program) always has the name `main`, regardless of the package's import in order to signal to `go build` that it must invoke the linker to make an executable file.
  - Some files in directory may have suffix `_test` if file name ends in `_test.go` where this package is the external test package.
  - Some tools for dependency management append the version number to package import paths, such as `gopkg.in/yaml.v2`, but the package name would just be `yaml`.

### Import Declarations

- Typically import declarations are grouped by domain, then alphabetically, groups separated with blank lines.
  - Note `gofmt` and `goimports` will group and sort for you.
- If two package names are the same, you can and must specify an alternative name to avoid a conflict using _renaming import_ syntax:
```go
import (
  "crypto/rand"
  mrand "math/rand" // alternative name mrand avoids conflict
)
```
- Renaming imports may be useful even when there is no conflict to either provide a better name or a shorter name, but the alternate name should be used consistently throughout the project.
- Renaming can also be used to avoid conflict with local variable names (e.g if source file has many local variables `path`, we can import the `path` package as `pathpkg`).

### Blank Imports

- To surpress the unused import error, we can rename the import to `_`, but the blank identifier can never be referenced.
```go
import _ "image/png" // register PNG decoder
```
- Blank imports are most often used to implement a compile-time mechanism where the main program can enable optional features by blank-importing additional packages.
- See [jpegconv](../jpegconv) for example of necessity for blank import, for a decoder to read the type of image input.
- Another example can be seen in the `database/sql` package to allow users to install just the database drivers they need:
```go
import (
  "database/sql"
  _ "github.com/lib/pq"                // Postgres support
  _ "github.com/go-sql-driver/mysql"   // MySQL support
)
db, err = sql.Open("postgres", dbname) // OK
db, err = sql.Open("mysql", dbname)    // OK
db, err = sql.Open("sqllite3", dbname) // unknown driver 'sqllite3'
```

### Packages and Naming

- Choose package names that are as concise as possible without being cryptic.
- Avoid choosing package names that are commonly used for related local variables.
- Package names usually take singular form, with exceptions if singular form conflicts with another variable keyword (eg bytes, errors).
- Avoid package names that have other connotations (eg "temp" for temperature coincides with "temporary").
- Use go packages as example. If similarities, create similar names to indicate parallelness.
- Consider method names with the identifier, e.g., `bytes.Equal`, `flag.Int`, `http.Get`, `json.Marshal`.
- If packages have the type def and the New method, keep the name short as can lead to lots of repitition. e.g., `rand.Rand`.

## The Go Tool

- Go tool is used for downloading, querying, formatting, building, testing, and installing packages of Go code.

### Most Commonly Used Commands

Command | Description
--- | ---
build | compile packages and dependencies
clean | remove object files
doc | show documentation for package or symbol
env | print Go environment information
fmt | run gofmt on package sources
get | download and install packages/deps
install | compile and install packages/deps
list | list packages
run | compile and run Go program
test | test packages
version | print Go version
vet | run go tool vet on packages

### Workspace Organization

- Only configuration that most users ever need to update is the GOPATH environment variable, which specifies the root of the workspace.
- When switching to a different workspace, update the go path. E.g for this project:
```sh
export GOPATH=$HOME/gopl
go get ...
``` 
- After downloading programs, will see 3 subdirectories in GOPATH:

**GOPATH Subdirectories**

dir | Purpose | Notes
--- | --- | ---
src | Holds source code. | Each package's import is relative to $GOPATH/src. Note will also see version-control repos beneath src.
pkg | Where build tools store compiled packages. | Will also contain multiple subdirs for each version.
bin | Holds binary executables |

- Another env variable is GOROOT, which specifies the root directory of the Go distribution, providing all the packages of the standard library. GOROOT's structure is similar to GOPATH.
  - Users don't need to set GOROOT because the go tool uses the location where it was installed. 

### Downloading Packages

- A package's import path does not indicate where to find it locally -- its where to find it on the Internet.
- `go get` can download a single pckage or an entire subtree or repository using the `...` notation: `go get domain/...`
- After downloading, it builds them and installs the libraries and commands.
- `go get` works best with popular code-hosting sites but for more-obscure domains, run `go help importpath` for more help.
- To specify a different domain name rather than the actual repository url, you can specify with metadata in the html from page name. Eg golang.org includes:
```html
<meta name="go-import"
  content="golang.org/x/net git https://go.googlesource.com/net">
```
- Use the `-u` option, `go get -u`, to retrieve the latest version of each package, which is convenient for getting started, but typically not appropriate for deployed projects.
- For deployed projects, we need to _vendor_ the code, or make it a persistent local copy of all the necessary dependencies and to update the copy carefully and deliberately. For more information, see the Vendor Directories in the output of the go help gopath command, which is supported directly as of Go 1.5.

### Building Packages

-`go build packages...` compiles each argument package.
- If the package is a library, result is discarded and just checks that the package is free of compile errors.
- If package is named `main`, the linker is invoked to create an executable in the current directory; where the name of the executable is taken from the last segment of the package's import path.
- Each directory should contain one package since each command requires its own directory.
- Package arguments can be specified by import path or by relative directory name (`..` or `.`).
- No argument provided assumes the current directory.
- Examples:
```sh
$ cd $GOPATH/src/domain/subdir/package
$ go build # OK
# from anywhere
$ go build gopl.io/ch1/hello world # OK
# from gopath
$ go build ./src/gopl.io/ch1/helloworld # OK
# following is NOT ok
$ cd $GOPATH
$ go build src/gopl.io/ch1/helloworld
Error: cannot find package "src/gopl.io/ch1/helloworld".
```

- `go build` also accepts a list of files (usually for temporary purposes). If package name of file is `main`, then executable name comes from the basename of the first `.go` file.

#### `go run`

- For quick throwaway programs use `go run files...` with the specific `.go` file as an argument.
  - If a list of files is given, the first argument that isn't suffixed with `.go` are assumed to be the beginnning of the list of arguments to the go executable(s).
- When using go build, the program is built with all of its dependencies and throws away all of the compiled code except for the final executable, if any.

#### `go install`

- When projects grow large, time to recompile can become noticiable, which is why the `go install` command exists.
- `install` is similar to `build`, except it saves the compiled code for each package and command instead of throwing it away.
- Compiled packages are saved beneath `$GOPATH/pkg` corresponding to the `src` directory where `src` code resides, and then the executables are saved in the `bin` directory.
  - Note many users add `$GOPATH/bin` to path to avoid typing out the path each time.
- We can also use `go build -i [pkg]` to install packages that are dependencies.
- When `go build` or `go install` is run, packages that haven't changed aren't recompiled to save time.

#### Cross Compiling

- In order to _cross-compile_ a go program (compile an executable intended for a different OS or CPU), set the `GOOS` AND `GOARCH` variables during the build:
```sh
# GOARCH=386 go build ...
```
- Example to print the operating system and architecture:
```go
func main() {
  fmt.Println(runtime.GOOS, runtime.GOARCH)
}
```
- Special comments called _build tags_ provide fine-grained control for build targets.
- For example, if file contains the comment `// +build linux darwin`, before package declaration (and its doc comment), `go build` will compile it only when building for Linux or MacOS X.
- To never compile a file, use `// +build ignore`.
- For more details on build constraints, run `go doc go/build`.

### Documenting Packages

- Every exported package member and declaration should be immediatedly preceded by a comment explaining its purpose and usage.
- Use complete sentences for doc comments.
- If doc comments are long (eg hundreds of lines). These cases may warrant a file of their own, which is usually called `doc.go`.
- Go's convention favors brevity and simplicity, which includes and is emphasized for documentation; many declarations can be explained in one well-worded sentence, and self-explanatory behavior doesn't necessitate a comment.

#### `go doc`

- `go doc [pkg]` prints the declaration and doc comment of the specified entitity, which may be a package, package member, or a method. 
  - Note this tool does not need complete import paths or correct identifier case.
- A related, but different tool is `godoc` which serves cross-linked HTML pages that provide the same information as `go doc` plus more:
```sh
$ godoc -http :8000 # runs an instance of godoc for workspace to browse your own packages
```
- Also see flags -analysis=type and -analysis=pointer flags to agument the documentation and source code with the results of advanced static analysis.

### Internal Packages

- There is a middle ground between completely hidden identifiers (unexported and only visible within same package), and exported, visible identifiers.
- These become useful for cases such as breaking up a large package into more manageable parts and may not want to reveal interfaces between those parts to other packages, or sharing utility functions across several packages of a project without exposing them widely, or just experimenting with a new package without prematurely committing to its API by putting it "on probation" with a limited set of clients.
- If a go import path contains `internal`, then it may be imported only by another package that is inside the tree rooted at the parent of the `internal` directory.
- For example, given:
```txt
net/http
net/http/internal/chunked
net/http/httputil
net/url
```
  - `net/http/internal/chunked` can be imported from `net/http/httputil` or `net/http`, but not from `net/url.
  - But `net/url` may import `net/http/httputil`.

### Querying Packages with `go list`

- `go list` reports information about available packages.
- Simplest usage is to test whether a package is present in the workspace and print its import path if so:
```sh
$ go list github.com/go-sql-driver/mysql
github.com/go-sql-driver/mysql
```
- Note argument to `go list` may contain the ellipses wildcard which matches any substring of the package's import package; also useful to enumerate all packages within a Go workspace or specific subtree, or even just related to a particular topic using `go list ...topic...`.
- To get print the entire record of a package, do `go list -json [entity]`
- The `-f` flag lets users customize output format using the template language of package `text/template`.
  - For example, to print the transitive dependencies of the `strconv` package, separated by spaces:
  ```sh
  $ go list -f '{{join.Deps " "}}' strconv
  errors math runtime unicode/utf8 unsafe
  ```
  - Or to print direct imports of each package:
  ```sh
  go list -f '{{.ImportPath}} -> {{join .Imports " "}}' [pkg]/...
  importpath/pkg -> dependency dependency ...
  ...
  ```
