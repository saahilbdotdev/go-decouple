# go-decouple (v0.1.0-alpha1)

Go decouple your settings from your application!

This project is a Go fork of the Python project [python-decouple](https://github.com/HBNetwork/python-decouple).

## Installation

```bash
go get github.com/saahilbdotdev/go-decouple
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/saahilbdotdev/go-decouple"
)

func main() {
	// Loads in order of priority: settings.ini file, .env file; both files are assumed to be in the current directory
	// If the settings.ini file is not found, it will be ignored
	// If the .env file is not found, it will be ignored
	// Highest priority will be given to the environment variables

	config := config.Default()

	// Get the value of the key "HOME" from the environment variables
	home := config.Get("HOME", nil, nil)

	// Get the value of the key "DEBUG" as a bool from the environment variables; here default value is false
	debug := config.Get("DEBUG", false, "bool")
}
```

## License

MIT License

Copyright (c) 2024 Saahil Bhavsar

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
