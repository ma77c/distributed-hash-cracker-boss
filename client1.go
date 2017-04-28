package main
import (
    "fmt"
    "net"
    "strings"
    "strconv"
    "crypto/md5"
    "encoding/hex"
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
  startInt, err := strconv.Atoi(start)
  endInt, err := strconv.Atoi(end)

  fmt.Printf("starting new job = %s\n", hash)
  if err != nil {
    fmt.Printf("Some error %v\n", err)
    return
  }
  for i:= startInt; i <= endInt; i++ {
    select {
    case kch := <-killChannel:
      if strings.Contains(kch, "<code>003</code>") {
        fmt.Printf("stopping job\n")
        unHashChannel <- "<code>004</code><start>"+strconv.Itoa(i)+"</start><end>"+strconv.Itoa(endInt)+"</end>"
        return
      }
    default:
      iString := strconv.Itoa(i)
      hasher := md5.New()
      hasher.Write([]byte(iString))
      possibleHash := hex.EncodeToString(hasher.Sum(nil))
      // fmt.Printf("current hash = %s\n", possibleHash)
      if possibleHash == hash {
        unHashChannel <- "<code>002</code><password>"+strconv.Itoa(i)+"</password>"
        return
      }
    }
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
        hash := getValueByElementXXX(smch, "hash")
        start := getValueByElementXXX(smch, "start")
        end := getValueByElementXXX(smch, "end")
        go unHashXXX(hash, start, end, unHashChannel, killChannel)
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
