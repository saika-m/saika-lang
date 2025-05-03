# Saika Language

<p align="center">
  <img src="https://via.placeholder.com/200x200?text=Saika" alt="Saika Logo" width="200" height="200">
</p>

<p align="center">
  A modern programming language that transpiles to Go
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#installation">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#examples">Examples</a> •
  <a href="#documentation">Documentation</a> •
  <a href="#contributing">Contributing</a> •
  <a href="#license">License</a>
</p>

## Features

- **Go Compatibility**: Seamlessly works with Go's ecosystem
- **Modern Syntax**: Streamlined function declarations with '數' (shù) character
- **Strong Typing**: Maintains Go's type system with enhanced readability
- **Powerful Toolchain**: Build, run, and manage Saika projects efficiently
- **Production-Ready**: Robust error handling and comprehensive testing
- **Cross-Platform**: Works on Windows, macOS, and Linux

## Installation

### Prerequisites

- Go 1.22 or higher

### From Source

```bash
# Clone the repository
git clone https://github.com/saika-m/saika-lang.git
cd saika-lang

# Build the Saika compiler
go build -o saika ./cmd/saika

# Add to your PATH (optional)
mv saika /usr/local/bin/  # Linux/macOS
# or
copy saika.exe C:\Windows\System32\  # Windows
```

### Using Go Install

```bash
go install github.com/saika-m/saika-lang/cmd/saika@latest
```

## Usage

### Basic Commands

```bash
# Compile a Saika file to an executable
saika build file.saika

# Run a Saika file
saika run file.saika

# View help information
saika help

# Show version information
saika version
```

### Advanced Options

```bash
# Specify output directory
saika build -o build/ file.saika

# Enable verbose output
saika build -v file.saika

# Add include paths for imports
saika build -I ./lib -I ./vendor file.saika

# Compile multiple files
saika build file1.saika file2.saika

# Use wildcards
saika build examples/*.saika
```

## Examples

### Hello World

```go
// hello.saika
package main

import "fmt"

數 main() {
    fmt.Println("Hello, Saika!")
}
```

### Functions and Types

```go
// functions.saika
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

### Structs and Methods

```go
// person.saika
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

## Project Structure

```
saika-lang/
├── cmd/
│   └── saika/           # Command-line interface
│       └── main.go
├── examples/            # Example Saika programs
│   ├── hello.saika
│   └── advanced.saika
├── internal/
│   ├── ast/             # Abstract Syntax Tree
│   ├── codegen/         # Code generation
│   ├── lexer/           # Lexical analysis
│   ├── parser/          # Parsing
│   └── transpiler/      # Transpilation engine
├── tests/               # Test suite
├── go.mod               # Go module file
├── LICENSE              # Project license
└── README.md            # Project documentation
```

## Documentation

- [Language Specification](docs/LANGUAGE_SPEC.md)
- [User Guide](docs/USER_GUIDE.md)
- [API Reference](docs/API_REFERENCE.md)
- [Development Guide](docs/DEVELOPMENT.md)

## Roadmap

- [ ] Language server protocol implementation
- [ ] IDE extensions for VS Code, JetBrains IDEs
- [ ] Package management system
- [ ] Interactive REPL
- [ ] Enhanced standard library
- [ ] Web framework support
- [ ] Documentation generator
- [ ] Performance optimizations

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the project
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please read [CONTRIBUTING.md](docs/CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- The Go team for creating an excellent language
- All contributors who have helped shape Saika
- The open-source community for their invaluable tools and libraries

---

<p align="center">
  Made with ❤️ by the Saika team
</p>