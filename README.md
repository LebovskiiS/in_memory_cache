# Go Cache Library

This library provides a simple cache implementation in Go with TTL (time-to-live) for each key.

✏️ **This project was created for learning purposes**, but you are free to use this library in your own projects without any restrictions. Feel free to modify and distribute it under fair use terms.

## Features

- 🧵 Thread-safe using `sync.RWMutex`.
- 🗑️ Automatic eviction of the least recently used items when the capacity is reached.
- ⏳ Key-based TTL (time-to-live) handling.

## Installation

To install the library, use the following command:

```sh
go get github.com/USERNAME/REPOSITORY_NAME
```  

Replace `USERNAME` and `REPOSITORY_NAME` with your GitHub username and repository name.

## Usage

Here’s how you can use the library in your Go projects:

```go
package main

import (
	"fmt"
	"time"
	"github.com/USERNAME/REPOSITORY_NAME/lesson" // Replace with your actual repository path
)

func main() {
	cache := lesson.NewCash()

	// Add items to the cache with a time-to-live (TTL)
	cache.Set("key1", "value1", 5*time.Second)

	// Retrieve items from the cache
	value, err := cache.Get("key1")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Key1:", value)
	}

	// Wait for the TTL to expire
	time.Sleep(6 * time.Second)

	// Try to retrieve the item again
	value, err = cache.Get("key1")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Key1:", value)
	}
}
```  

## License

This project was written for educational purposes, but you can freely use, modify, and distribute the library in your own projects under fair use conditions.

---  

⭐ Feel free to star the repository and contribute if you find this helpful!
