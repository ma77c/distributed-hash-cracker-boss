package main
import (
    "fmt"
    "net"
    "crypto/md5"
    "encoding/hex"
    "math"
    "strconv"
    "strings"
)
//// codes ////
//
//// 0000 - Request - New Job
//
//// 0001 - Response - New Job
//
//// 0010 - Request -
//
//// 0011 - Request -
//
//// 0020 - Request -
//
//// 0021 - Response -
//
//// 0030 - Request -
//
//// 0031 - Response -
//
func checkError(err error) {
  if err != nil {
      fmt.Printf("Error %v", err)
      return
  }
}
func getValueByElementXXX(w string, element string) string {
  x := strings.Split(w, "<"+element+">")
  y := strings.Split(x[1], "</"+element+">")
  z := y[0]
  return z
}

type Range struct {
  start int
  end int
  clientID int
}
type Client struct {
  id int
  address *net.UDPAddr
}

func (c Client) SetAddress(address *net.UDPAddr) {
  c.address = address
}

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr, response string) {
    _,err := conn.WriteToUDP([]byte(response), addr)
    checkError(err)
}

func removeRangeByID(rangeArray []Range, clientID int) []Range {
  for _, r := range rangeArray {
    if r.clientID == clientID {
      r.clientID = rangeArray[len(rangeArray)-1].clientID
      r.start = rangeArray[len(rangeArray)-1].start
      r.end = rangeArray[len(rangeArray)-1].end
      return rangeArray[:len(rangeArray)-1]
    }
  }
  return nil
}


func main() {
  //// hash
  password:= "73111111"
  hasher := md5.New()
  hasher.Write([]byte(password))
  hash := hex.EncodeToString(hasher.Sum(nil))
  //// the big range
  numChars := 8
  lowerLimit := 0
  upperLimit := math.Pow(10, float64(numChars - 1))
  for i := 0; i < numChars - 1; i++ {
    upperLimit = upperLimit + (math.Pow(10, float64(i)))
  }
  upperLimit = upperLimit * 9

  increment := 1000000

  activeClients := []Client{}
  minorRanges := []Range{}
  majorRange := Range{start: lowerLimit, end: int(upperLimit)}

  addr := net.UDPAddr {
      Port: 1234,
      IP: net.ParseIP("127.0.0.1"),
  }
  //// listener for incoming clients
  conn, err := net.ListenUDP("udp", &addr)
  checkError(err)
  //// keep listening
  for {
    byteArray := make([]byte, 2048)
    //// blocks until reads a message
    n, inAddress, err := conn.ReadFromUDP(byteArray)
    checkError(err)
    //// boom - message received - lets parse the message
    inMessage := string(byteArray[:n])
    inID, err := strconv.Atoi(getValueByElementXXX(inMessage, "id"))
    inCode := getValueByElementXXX(inMessage, "code")
    //// initiate temporary client
    inClient := Client{ id: inID, address: inAddress }
    fmt.Printf("CLIENT : %v : %s\n", inClient.address, inMessage)

    //// if brand new client
    if inClient.id == -1 {
      greatestID := -1
      for _, ac := range activeClients {
        if ac.id > greatestID {
          greatestID = ac.id
        }
      }
      if greatestID == -1 {
        inClient.id = 0
      } else {
        inClient.id = greatestID + 1
      }
      start := majorRange.start
      end := majorRange.start + increment - 1
      majorRange.start = end + 1
      minor := Range{start: start, end: end, clientID: inClient.id}
      minorRanges = append(minorRanges, minor)

      response := "<id>0</id><name>New Job</name><code>0001</code><hash>"+hash+"</hash><start>"+strconv.Itoa(minor.start)+"</start><end>"+strconv.Itoa(minor.end)+"</end>"
      go sendResponse(conn, inClient.address, response)

    //// if active client
    } else {
      if inCode == "0000" {
        minorRanges = removeRangeByID(minorRanges, inClient.id)
        start := majorRange.start
        end := majorRange.start + increment - 1
        majorRange.start = end + 1
        minor := Range{start: start, end: end, clientID: inClient.id}
        minorRanges = append(minorRanges, minor)
        response := "<id>"+strconv.Itoa(inClient.id)+"</id><name>New Job</name><code>0001</code><hash>"+hash+"</hash><start>"+strconv.Itoa(minor.start)+"</start><end>"+strconv.Itoa(minor.end)+"</end>"
        go sendResponse(conn, inClient.address, response)
      } else if inCode == "0091" {
        password := getValueByElementXXX(inMessage, "password")
        fmt.Printf("password = %s\n", password)
        fmt.Printf("END!")
        return
      }
    }
  }
}
