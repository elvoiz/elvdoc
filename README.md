# elvdoc

elvdoc is a Go package for reading and writing `.elv` files. It is designed to help developers parse, generate, and manipulate files with the `.elv` extension, making it easy to integrate `.elv` file support into your Go projects.

> **Note:** elvdoc is mainly used in the Elvoiz app, an application specialized in invoice management.

## Features
- Read `.elv` files
- Write `.elv` files
- Parse and generate structured data
- Easy-to-use API

## Installation

```
go get github.com/elvoiz/elvdoc
```

## Usage

```go
package main

import (
    "github.com/elvoiz/elvdoc"
)

func main() {
    // Example: Reading an .elv file
    doc, err := elvdoc.ReadFile("example.elv")
    if err != nil {
        panic(err)
    }
    // Work with doc...

    // Example: Writing an .elv file
    err = elvdoc.WriteFile("output.elv", doc)
    if err != nil {
        panic(err)
    }
}
```

## License

MIT License
