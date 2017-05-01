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
      if strings.Contains(kch, "<code>0020</code>") {
        fmt.Printf("stopping job\n")
        unHashChannel <- "<code>0021</code><start>"+strconv.Itoa(i)+"</start><end>"+strconv.Itoa(endInt)+"</end>"
        return
      }
    default:
      iString := strconv.Itoa(i)
      hasher := md5.New()
      hasher.Write([]byte(iString))
      possibleHash := hex.EncodeToString(hasher.Sum(nil))
      // fmt.Printf("current hash = %s\n", possibleHash)
      if possibleHash == hash {
        unHashChannel <- "<code>0091</code><password>"+strconv.Itoa(i)+"</password>"
        return
      }
    }
  }
  fmt.Printf("finished job - need new work\n")
  unHashChannel <- "<name>New Job</name><code>0000</code>"
  return
}

func main() {
  var id int
  conn, err := net.Dial("udp", "10.10.10.2:1234")
  if err != nil {
      fmt.Printf("Some error %v", err)
      return
  }
  //// send message to server
  fmt.Fprintf(conn, "<id>-1</id><name>New Job</name><type>Request</type><code>0000</code>")
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
      if code == "0001" {
          id, err = strconv.Atoi(getValueByElementXXX(smch, "id"))
          hash := getValueByElementXXX(smch, "hash")
          start := getValueByElementXXX(smch, "start")
          end := getValueByElementXXX(smch, "end")
          go unHashXXX(hash, start, end, unHashChannel, killChannel)
          working = true
      } else if code == "0020" {
          if working == true {
            killChannel <- smch
            working = false
          }
      }
    case uhch := <-unHashChannel:
      code := getValueByElementXXX(uhch, "code")
      //// response - return ranges
      if (code) == "0021" {
        //// forward ranges to server
        fmt.Fprintf(conn, uhch+"<id>"+strconv.Itoa(id)+"</id>")
      } else if code == "0091" {
        password := getValueByElementXXX(uhch, "password")
        fmt.Fprintf(conn, uhch+"<id>"+strconv.Itoa(id)+"</id>")
        fmt.Printf("password = %s\n", password)
        fmt.Printf("END!")
        return
      } else if code == "0000" {
        response := uhch+"<id>"+strconv.Itoa(id)+"</id>"
        fmt.Fprintf(conn, response)
      }
    }
  }

}
