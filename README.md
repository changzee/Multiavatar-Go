# multiavatar-go

[![Go Reference](https://pkg.go.dev/badge/github.com/changzee/multiavatar-go.svg)](https://pkg.go.dev/github.com/changzee/multiavatar-go)

A Go library for generating multicultural avatars. This is a port of the original JavaScript library [Multiavatar](https://github.com/multiavatar/Multiavatar).

This library generates unique, deterministic avatars based on a string input. The core algorithm is consistent with the original version, ensuring that the same input string will always produce the same avatar.

![Logo](https://raw.githubusercontent.com/multiavatar/Multiavatar/main/logo.png)

## Features

- **Deterministic Algorithm**: Always generates the same avatar for the same input.
- **No Dependencies**: The library is self-contained and requires no external dependencies.
- **Functional Options**: Easy-to-use API with support for functional options.
- **Customizable**: Supports generating avatars with or without a background.
- **Thread-Safe**: All public functions are safe for concurrent use.

## Installation

To install the library, use `go get`:

```bash
go get github.com/changzee/multiavatar-go
```

## Usage

Import the library into your project:

```go
import "github.com/changzee/multiavatar-go"
```

### Basic Example

To generate a default avatar, simply call the `Generate` function with a string.

```go
package main

import (
	"log"
	"os"

	"github.com/changzee/multiavatar-go"
)

func main() {
	// Generate an avatar for the string "Binx Bond"
	svgCode := multiavatar.Generate("Binx Bond")

	// Save the SVG to a file
	err := os.WriteFile("avatar.svg", []byte(svgCode), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
```

### Example with Options

The library uses the functional options pattern to configure the avatar generation. For example, you can generate an avatar with a transparent background using the `WithoutBackground()` option.

```go
package main

import (
	"log"
	"os"

	"github.com/changzee/multiavatar-go"
)

func main() {
	// Generate an avatar with a transparent background
	svgCode := multiavatar.Generate("John Doe", multiavatar.WithoutBackground())

	// Save the SVG to a file
	err := os.WriteFile("avatar_transparent.svg", []byte(svgCode), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
```

## API Reference

### `Generate(input string, options ...Option) string`

This is the main function that generates the SVG avatar.

- `input string`: A UTF-8 string that serves as the seed for the avatar. The same input will always produce the same avatar.
- `options ...Option`: A variadic set of functional options to customize the generation.

Returns a string containing the complete, well-formed SVG code for the avatar.

### Options

#### `WithoutBackground() Option`

This option removes the colored background from the avatar, making it transparent.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. The original Multiavatar project has its own license that should be respected.
