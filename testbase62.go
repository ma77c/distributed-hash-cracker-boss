package main

import (
    "fmt"
    "./base62"
)

func main() {

    a := []byte("xzz")
    b := []byte("11")
    c := base62.Add(a, b)

    fmt.Printf("Resultant: %s", c)

}
