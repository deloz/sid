# SID

The main purpose of SID is to seamlessly convert between a JSON-formatted string representation of an ID and its int64 counterpart when interacting with a database.

## Installation

You can install SID using `go get`:

```bash
go get -u -v github.com/deloz/sid
```

## Usage

Here's a basic example of how to use SID:

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/deloz/sid"
)

func main() {
	// Create a new ID
	id := sid.New(18283195028008204)

	// Print the ID as a string
	fmt.Println("ID:", id.String())

	// Check if the ID is zero
	fmt.Println("IsZero:", id.IsZero())

	// Marshal the ID to JSON
	jsonData, _ := json.Marshal(id)
	fmt.Println("JSON:", string(jsonData))

	// Unmarshal the ID from JSON
	var newID sid.ID
	_ = json.Unmarshal(jsonData, &newID)
	fmt.Println("Unmarshaled ID:", newID)
}
```

[[play](https://go.dev/play/p/Fie93MkAzmq)]

## Features

- Support creation of IDs from `uint64` (up to `math.MaxInt64`), `int64`, and numeric string representations.
- Implement JSON (de)serialization, comparison, and sorting for custom ID types in Go.

## Contributing

Contributions are welcome! If you find a bug or want to propose a new feature, feel free to open an issue or submit a pull request.

## License

SID is licensed under the MIT License. See [LICENSE](LICENSE) for details.
