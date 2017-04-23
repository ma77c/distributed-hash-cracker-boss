package main
import (
    "fmt"
    "time"
    "net"
    "strings"
    "strconv"
)


func getValueByElementXXX(w string, element string) string {
  x := strings.Split(w, "<"+element+">")
  y := strings.Split(x[1], "</"+element+">")
  z := y[0]
  return z
}
func keepReadingXXX(conn net.Conn, serverMessageChannel chan string) {
  var byteArray = make([]byte, 2048)
  for {
    n, err := conn.Read(byteArray)
    if err == nil {
      fmt.Printf("got message\n")
      serverMessageChannel <- string(byteArray[:n])
    }
  }
}

func unHashXXX(hash string, start string, end string, unHashChannel chan string, killChannel chan string) {
  hashInt, err := strconv.Atoi(hash)
  startInt, err := strconv.Atoi(start)
  endInt, err := strconv.Atoi(end)

  fmt.Printf("starting new job = %v\n", hashInt)
  if err != nil {
    fmt.Printf("Some error %v\n", err)
    return
  }
  for i:= startInt; i <= endInt; i++ {
    select {
    case kch := <-killChannel:
      if strings.Contains(kch, "<code>001</code>") {
        fmt.Printf("stopping job\n")
        return
      }
    default:
      if i == hashInt {
        unHashChannel <- "<code>002</code><password>"+strconv.Itoa(i)+"</password>"
        return
      }
    }
    time.Sleep(1 * time.Second)
  }
}

func main() {
  conn, err := net.Dial("udp", "127.0.0.1:1234")
  if err != nil {
      fmt.Printf("Some error %v", err)
      return
  }
  //// send message to server
  fmt.Fprintf(conn, "Give me a hash to work on ...")
  serverMessageChannel := make(chan string)
  unHashChannel := make(chan string)
  killChannel := make(chan string)
  go keepReadingXXX(conn, serverMessageChannel)
  var working bool
  for {
    select {
    case smch := <-serverMessageChannel:
      fmt.Printf("SERVER : %s\n", smch)
      code := getValueByElementXXX(smch, "code")
      if code == "001" {
        if working == true {
          killChannel <- "<code>001</code>"
          working = false
        }
        start := getValueByElementXXX(smch, "start")
        end := getValueByElementXXX(smch, "end")
        go unHashXXX("10", start, end, unHashChannel, killChannel)
        working = true
      }
    case uhch := <-unHashChannel:
      code := getValueByElementXXX(uhch, "code")
      if code == "002" {
        password := getValueByElementXXX(uhch, "password")
        fmt.Printf("password = %s\n", password)
        fmt.Printf("END!")
        return
      }
    }
  }

}
