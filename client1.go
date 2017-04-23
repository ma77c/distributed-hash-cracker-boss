package main
import (
    "fmt"
    "time"
    "net"
    "sync"
    "strings"
    "strconv"
)

func keepReading(conn net.Conn, work_ch chan int) {
  dh := make(chan int)
  var byteArray = make([]byte, 2048)
  for {
    n, err := conn.Read(byteArray)
    if err == nil {
      fmt.Printf("SERVER : %s\n", byteArray[:n])
      if strings.Contains(string(byteArray[:n]), "<001>") {
        dh <- 302
      } else if (strings.Contains(string(byteArray[:n]), "hash")) {
          go unHash("10", dh, work_ch)
      }
    } else {
        fmt.Printf("Some error %v\n", err)
    }
  }
}
func unHash(hash string, dh chan int, work_ch chan int) {
  hash_int, err := strconv.Atoi(hash)
  if err != nil {
    fmt.Printf("Some error %v\n", err)
    return
  }
  var dh_clean int
  for i := 0; i < 100; i++ {
    fmt.Printf("SERVER waiting\n")
    select {
    case dh_clean = <- dh:
      fmt.Printf("SERVER : %s\n", dh_clean)
      if (dh_clean == 302) {
        fmt.Printf("Redirecting\n")
        work_ch <- 200
        return
      }
    default:
      if (i == hash_int) {
        fmt.Printf("Hash Found : %v\n", i)
        work_ch <- 200
        return
      }
    }
    time.Sleep(1*time.Second)
  }
}

func work(conn net.Conn, wg *sync.WaitGroup) {
  work_ch := make(chan int)
  var work_ch_clean int
  go keepReading(conn, work_ch)
  for {
    work_ch_clean = <- work_ch
    if (work_ch_clean == 200) {
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
  go work(conn, wg)
  fmt.Fprintf(conn, "Give me a hash to work on ...")
  wg.Wait()
}
