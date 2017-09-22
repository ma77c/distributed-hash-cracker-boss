package main

import (
    "fmt"
    "./base62"
    "time"
)

func main() {

    start := []byte("0")
    end := []byte("4C92")

    current := start
    begin := time.Now()
    for string(current) != string(end) {
        // fmt.Printf("Current: %s\n", current)
        current = base62.Add(current, []byte("1"))
    }
    elapsed := time.Since(begin)
    fmt.Printf("Base62: %s\n", elapsed)
}
