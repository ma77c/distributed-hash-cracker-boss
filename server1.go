package main
import (
    "fmt"
    "net"
    "crypto/md5"
    "encoding/hex"
    "math"
    "strconv"
)
//// codes ////
//
//// 0010 - Request - New Job
//
//// 0011 - Response - New Job
//
//// 0020 - Request - Stop Work Return Ranges
//
//// 0021 - Response - Stop Work Return Ranges
//
//// 0030 - Request -
//
//// 0031 - Response -
//
//// 0040 - Request -
//
//// 0041 - Response -
//
//// 0050 - Request -
//
//// 0051 - Response -
//
//// 0060 - Request -
//
//// 0061 - Response -
//
//// 0070 - Request -
//
//// 0071 - Response -
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
  address *net.UDPAddr
  start int
  end int
  checkIn bool
  id int
}

func (c Client) SetAddress(address *net.UDPAddr) {
  c.address = address
}

func (c Client) SetRange(start int, end int) {
  c.start = start
  c.end = end
}

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr, response string) {
    _,err := conn.WriteToUDP([]byte(response), addr)
    checkError(err)
}


func main() {
  //// hash
  password:= "44444444"
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
  upperLimit = int(upperLimit)


  activeClients := []Client{}
  waitList := []Client{}

  ranges := []Range{}
  ranges = append(ranges, Range{start: lowerLimit, end: upperLimit, clientID: nil})

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
    inMessage := byteArray[:n]
    inID := int(getValueByElementXXX(inMessage, "id"))
    inStart := int(getValueByElementXXX(inMessage, "start"))
    inEnd :=  int(getValueByElementXXX(inMessage, "end"))
    //// initiate temporary client
    inClient := Client{ id: inID, address: inAddress, start: inStart, end: inEnd, checkIn : false }
    fmt.Printf("CLIENT : %v : %s\n", inClient.address, inClient.message)

    //// if first client ever
    if len(activeClients) == 0 {
      inClient.id = 0;
      ranges[0].clientID = 0
      inClient.start = ranges[0].start
      inClient.end = ranges[0].end
      response := "<id>0</id><name>New Job</name><code>0011</code><start>"+strconv.Itoa(inClient.start)+"</start><end>"+strconv.Itoa(inClient.start)+"</end>"
      go sendResponse(conn, ac.address, response)
      activeClients = append(activeClients, inClient)
    }
    //// if brand new client
    if inClient.id == -1 {
      //// add to waitlist
      waitList = append(waitList, inClient)
      //// send message to active clients - stop work request range
      for _, ac := range activeClients {
        ac.checkIn = false
        response := "<name>Stop Work Request Range</name><code>0020</code>"
        go sendResponse(conn, ac.address, response)
      }
    //// if active client
    } else {
      inCode := getValueByElementXXX(inMessage, "code")
      if inCode == "0021" {
        for _, ac : = range activeClients {
          if ac.id == inClient.id {
            ac.start = inClient.start
            ac.end = inClient.end
            ac.checkIn = true
          }
        }
        allIn := true
        for _, ac := range activeClients {
          if ac.checkIn == false {
            allIn == false
          }
        }
        //// if all clients check in
        if allIn == true {
          for _, c := range clients {
            for _, r := range ranges {
              if r.clientID == c.id {
                
              }
            }
          }
        }
      }
    }
  }
}
