# Saika Language Specification

Version 1.0.0

## Introduction

Saika is a programming language designed to provide a modern syntax while being fully compatible with the Go ecosystem. Saika code transpiles directly to Go, enabling seamless integration with existing Go libraries and tools. The language maintains most of Go's features and semantics while introducing syntactic enhancements and quality-of-life improvements.

## Key Features

- **Go Compatibility**: Saika transpiles to Go, ensuring compatibility with Go's ecosystem.
- **Chinese Character Function Syntax**: Functions can be declared using the Chinese character '數' (shù) instead of "func".
- **Enhanced Syntactic Sugar**: Simplified variable declarations with 'let' and 'var'.
- **Type System**: Strong, static typing similar to Go.
- **Familiar Control Flow**: Standard control flow structures enhanced with syntactic improvements.
- **Structs and Interfaces**: Full support for Go's type system.

## Lexical Elements

### Comments

Saika supports the same comment styles as Go:

```
// Line comment

/* Block comment
   spanning multiple lines */
```

### Identifiers

Identifiers follow the same rules as Go:

- Must begin with a letter or underscore
- Subsequent characters can be letters, digits, or underscores
- Case-sensitive

```
validIdentifier
_validIdentifier
ValidIdentifier123
```

### Keywords

Saika includes the following keywords:

```
break       case        chan        const       continue
default     defer       else        fallthrough for
go          goto        if          import      interface
let         map         package     range       return
select      struct      switch      type        var
數           // Chinese character for function declarations
```

### Operators and Punctuation

Saika supports all Go operators:

```
+    &     +=    &=     &&    ==    !=    (    )
-    |     -=    |=     ||    <     <=    [    ]
*    ^     *=    ^=     <-    >     >=    {    }
/    <<    /=    <<=    ++    =     :=    ,    ;
%    >>    %=    >>=    --    !     ...   .    :
```

## Types

Saika supports all Go types:

### Basic Types

```
bool
string
int  int8  int16  int32  int64
uint uint8 uint16 uint32 uint64 uintptr
byte // alias for uint8
rune // alias for int32, represents a Unicode code point
float32 float64
complex64 complex128
```

### Composite Types

```
array     // Fixed-size collection of elements
slice     // Dynamic-size collection of elements
struct    // User-defined type with named fields
map       // Key-value collection
channel   // Communication mechanism between goroutines
interface // Method set specification
function  // First-class function type
```

## Declarations and Scope

### Package Declaration

Every Saika file begins with a package declaration:

```
package packagename
```

### Import Declaration

Imports follow Go's style:

```
import "packagename"
import "path/to/package"

// Multiple imports
import (
    "fmt"
    "strings"
    "math"
)

// Aliased import
import fmt "fmt"
```

### Variable Declarations

Variables can be declared using `var`, `let`, or `:=`:

```
// Using var (like Go)
var x int
var y int = 10
var z = 20

// Using let (Saika extension)
let a int
let b int = 30
let c = 40

// Short variable declaration
i := 50
```

### Constant Declarations

Constants are declared using the `const` keyword:

```
const PI = 3.14159
const Version = "1.0.0"

// Multiple constants
const (
    StatusOK = 200
    StatusNotFound = 404
    StatusError = 500
)
```

### Function Declarations

Functions can be declared using the `數` character or `func` keyword:

```
// Using the Saika 數 character
數 add(a int, b int) int {
    return a + b
}

// Using traditional Go func keyword (also supported)
func multiply(a int, b int) int {
    return a * b
}

// With receiver (method)
數 (p Person) name() string {
    return p.firstName + " " + p.lastName
}
```

### Type Declarations

Types are declared using the `type` keyword:

```
type MyInt int

type Person struct {
    firstName string
    lastName string
    age int
}

type Stringer interface {
    String() string
}
```

## Expressions

Saika supports all Go expressions:

### Operators

The precedence of operators is the same as in Go:

1. `*`, `/`, `%`, `<<`, `>>`, `&`, `&^`
2. `+`, `-`, `|`, `^`
3. `==`, `!=`, `<`, `<=`, `>`, `>=`
4. `&&`
5. `||`

### Function Calls

```
result := add(5, 10)
fmt.Println("Hello, Saika!")
```

### Method Calls

```
length := name.length()
strings.ToUpper("hello")
```

## Statements

### Assignment

```
x = 10
y, z = z, y  // Swap values
```

### Conditional Statements

```
if x > 10 {
    fmt.Println("x is greater than 10")
} else if x < 0 {
    fmt.Println("x is negative")
} else {
    fmt.Println("x is between 0 and 10")
}
```

### Switch Statements

```
switch day {
case "Monday":
    fmt.Println("Start of the week")
case "Friday":
    fmt.Println("End of the workweek")
default:
    fmt.Println("Another day")
}

// Switch with no expression
switch {
case score >= 90:
    fmt.Println("A grade")
case score >= 80:
    fmt.Println("B grade")
default:
    fmt.Println("Lower grade")
}
```

### Loops

```
// Traditional for loop
for i := 0; i < 10; i++ {
    fmt.Println(i)
}

// While-style loop
for count > 0 {
    count--
}

// Infinite loop
for {
    if condition {
        break
    }
}

// Range-based loop
for index, value := range items {
    fmt.Printf("%d: %v\n", index, value)
}

// Just the value
for _, value := range items {
    fmt.Println(value)
}
```

### Control Flow

```
// Break statement
for i := 0; i < 10; i++ {
    if i == 5 {
        break
    }
}

// Continue statement
for i := 0; i < 10; i++ {
    if i % 2 == 0 {
        continue
    }
    fmt.Println(i)  // Print odd numbers
}

// Return statement
數 getValue() int {
    return 42
}
```

## Packages and Imports

Saika follows Go's package system. Each Saika file begins with a package declaration and can import other packages.

```
package main

import "fmt"
import "strings"
```

The main package is the entry point for a Saika executable program. Other packages provide libraries that can be imported.

## Examples

### Hello World

```
package main

import "fmt"

數 main() {
    fmt.Println("Hello, Saika!")
}
```

### Functions and Parameters

```
package main

import "fmt"

數 add(a int, b int) int {
    return a + b
}

數 greet(name string) {
    fmt.Printf("Hello, %s!\n", name)
}

數 main() {
    result := add(5, 7)
    fmt.Println("5 + 7 =", result)
    
    greet("Saika")
}
```

### Structures and Methods

```
package main

import "fmt"

struct Person {
    name string
    age int
}

數 (p Person) describe() string {
    return fmt.Sprintf("%s is %d years old", p.name, p.age)
}

數 main() {
    p := Person{name: "John", age: 30}
    fmt.Println(p.describe())
}
```

## Differences from Go

While Saika is largely compatible with Go, there are some syntactic differences:

1. Functions can be declared with `數` instead of `func`
2. Variables can be declared with `let` as an alternative to `var`
3. Structs can be defined using the `struct` keyword directly (without requiring a `type` declaration first)

## Implementation Notes

The Saika transpiler converts Saika code to Go code, which can then be compiled with the standard Go compiler. This approach allows Saika to leverage the entire Go ecosystem, including libraries, tools, and the runtime.

## Future Directions

The Saika language is under active development. Future versions may include:

- Additional syntactic sugar for common operations
- Enhanced error messages and diagnostics
- Integration with popular IDEs and text editors
- Language server protocol support
- More extensive standard library wrappers

## Conclusion

Saika offers a familiar yet refreshed experience for Go developers, maintaining Go's philosophy of simplicity and efficiency while providing syntactic enhancements. The language's full compatibility with Go ensures that developers can leverage their existing knowledge and the rich Go ecosystem.