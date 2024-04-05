# SID

The main purpose of SID is to seamlessly convert between a JSON-formatted string representation of an ID and its uint64 counterpart when interacting with a database.

## Installation

You can install SID using `go get`:

```bash
go get github.com/deloz/sid/v1
```

## Usage

Here's a basic example of how to use SID:

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/deloz/sig/v1"
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

## Features

- Create IDs from uint64, int64, and string representations.
- Marshal and unmarshal IDs to/from JSON and text.
- Compare IDs.
- Sort slices of IDs.

## Contributing

Contributions are welcome! If you find a bug or want to propose a new feature, feel free to open an issue or submit a pull request.

## License

SID is licensed under the MIT License. See [LICENSE](LICENSE) for details.
