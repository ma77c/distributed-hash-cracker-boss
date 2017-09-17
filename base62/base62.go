package base62

var (
    base62dec = make(map[byte]int)
    decbase62 = make(map[int]byte)
)

func init() {
    charset := []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
    for pos, char := range charset {
        base62dec[char] = pos
        decbase62[pos] = char
    }
}
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
func Add(augend []byte, addend []byte) []byte {
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
    carry := 0
    var shortDec int
    //
    // fmt.Printf("Long: %s\nShort: %s\n", long, short)
    for pos, char := range long {
        longDec := base62dec[char]
        if pos >= len(short) {
            shortDec = 0
        } else {
            shortDec = base62dec[short[pos]]
        }

        singleSum := longDec + shortDec + carry
        remainder = 0
        carry = 0
        if singleSum >= 62 {
            carry = 1
            remainder = singleSum - 62
            resultant[pos] = decbase62[remainder]
        } else {
            resultant[pos] = decbase62[singleSum]
        }
    }
    if carry == 1 {
        r := make([]byte, len(resultant) + 1)
        for pos, char := range resultant {
            r[pos] = char
        }
        r[len(resultant)] = decbase62[1]
        resultant = r
    }
    resultant = reverse(resultant)
    return resultant

}
// func main() {
//
//     base62dec := make(map[byte]int)
//     decbase62 := make(map[int]byte)
//     charset := []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
//
//     for pos, char := range charset {
//         base62dec[char] = pos
//         decbase62[pos] = char
//     }
//
//     a := []byte("xzz")
//     b := []byte("11")
//     c := add(a, b, base62dec, decbase62)
//     fmt.Printf("Resultant: %s", c)
//
// }
