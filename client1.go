package main
import (
    "fmt"
    "time"
    "net"
    "sync"
    "strconv"
)

func unHash(hash string) {
  hash_int, err := strconv.Atoi(hash)
  if err != nil {
    fmt.Printf("Some error %v\n", err)
    return
  }
  for i := 0; i < 100; i++ {
    if (i == hash_int) {
      fmt.Printf("Hash Found : %v\n", i)
      return
    }
  }
}
func xyz(conn net.Conn, p []byte) {
  for {
    n, err := conn.Read(p)
    if err == nil {
      fmt.Printf("SERVER : %s\n", p[:n])
      go unHash("36")
    } else {
        fmt.Printf("Some error %v\n", err)
    }
  }
}

func main() {
  var wg = &sync.WaitGroup{}
    p :=  make([]byte, 2048)
    conn, err := net.Dial("udp", "127.0.0.1:1234")
    if err != nil {
        fmt.Printf("Some error %v", err)
        return
    }
    wg.Add(1)
    go xyz(conn, p)
    time.Sleep(1 * time.Second);
    fmt.Fprintf(conn, "Give me a hash to work on ...")
    wg.Wait()
}
