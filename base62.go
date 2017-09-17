package main

import (
    "fmt"
)

func reverse(array []byte) []byte {
    length := len(array)
    reversed := make([]byte, length)
    pos := 0
    for i := length - 1; i >= 0; i-- {
        reversed[pos] = array[i]
        pos++
    }
    return reversed
}
func add(augend []byte, addend []byte, base62dec map[byte]int, decbase62 map[int]byte) []byte {
    //
    var long []byte
    var short []byte

    if len(augend) > len(addend) {
        long = reverse(augend)
        short = reverse(addend)
    } else {
        long = reverse(addend)
        short = reverse(augend)
    }
    resultant := make([]byte, len(long))
    remainder := 0
    //
    // fmt.Printf("Long: %s\nShort: %s\n", long, short)
    for pos, char := range long {
        longDec := base62dec[char]
        shortDec := base62dec[short[pos]]
        singleSum := longDec + shortDec + remainder
        remainder = 0
        if singleSum >= 62 {
            remainder = singleSum - 62
            resultant[pos] = decbase62[remainder]
        } else {
            resultant[pos] = decbase62[singleSum]
        }

    }
    resultant = reverse(resultant)
    return resultant

}
func main() {

    base62dec := make(map[byte]int)
    decbase62 := make(map[int]byte)
    charset := []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

    for pos, char := range charset {
        base62dec[char] = pos
        decbase62[pos] = char
    }

    a := []byte("ABC")
    b := []byte("123")
    c := add(a, b, base62dec, decbase62)
    fmt.Printf("Resultant: %s", c)
    // b := []byte("abc")

    // c := add(a, b)

    // for pos, char := range a {
    //     append(char, a)
    // }
    // c := reverse(a)
    // d := reverse(c)
    // for pos, char := range c {
    //     fmt.Printf("%s = %s\n", pos, string(char))
    // }
    // fmt.Printf("%+s", d)



}

// var m map[string]int
//
//
//
// func add(augend string, addend string) string {
//
//
//     return sum
// }
