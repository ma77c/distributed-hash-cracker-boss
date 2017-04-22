package main
import (
    "fmt"
    "time"
    "net"
    "sync"
    "strings"
    "strconv"
)

func keepReading(conn net.Conn, ch chan int) {
  var byteArray = make([]byte, 2048)
  for {
    n, err := conn.Read(byteArray)
    if err == nil {
      fmt.Printf("SERVER : %s\n", byteArray[:n])
      if (strings.Contains(string(byteArray[:n]), "hash")) {
        go unHash("36", ch)
      }
    } else {
        fmt.Printf("Some error %v\n", err)
    }
  }
}
func unHash(hash string, ch chan int) {
  hash_int, err := strconv.Atoi(hash)
  if err != nil {
    fmt.Printf("Some error %v\n", err)
    return
  }
  for i := 0; i < 100; i++ {
    time.Sleep(1 * time.Second)
    if (i == hash_int) {
      fmt.Printf("Hash Found : %v\n", i)
      ch <- 200;
    }
  }
}

func xyz(conn net.Conn, wg *sync.WaitGroup) {
  ch := make(chan int)
  var ch_clean int
  go keepReading(conn, ch)
  for {
    ch_clean = <- ch
    if (ch_clean == 200) {
      fmt.Printf("hellohello\n")
      wg.Done()
    }
  }
}

func main() {
  var wg = &sync.WaitGroup{}
  conn, err := net.Dial("udp", "127.0.0.1:1234")
  if err != nil {
      fmt.Printf("Some error %v", err)
      return
  }
  wg.Add(1)
  go xyz(conn, wg)
  time.Sleep(1 * time.Second);
  fmt.Fprintf(conn, "Give me a hash to work on ...")
  wg.Wait()
}
