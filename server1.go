package main
import (
    "fmt"
    "net"
    "crypto/md5"
    "encoding/hex"
    "math"
)

type Client struct {
  address *net.UDPAddr
  startRange int
  endRange int
}
func (c Client) SetAddress(address *net.UDPAddr) {
  c.address = address
}
func (c Client) SetRange(startRange int, endRange int) {
  c.startRange = startRange
  c.endRange = endRange
}

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr, response string) {
    _,err := conn.WriteToUDP([]byte(response), addr)
    if err != nil {
        fmt.Printf("Couldn't send response %v", err)
    }
    // time.Sleep(3*time.Second)
    // _,err = conn.WriteToUDP([]byte("<name>New Job</name><code>001</code><hash>"+string(hash[:])+"</hash><start>0</start><end>50</end>"), addr)
    // if err != nil {
    //     fmt.Printf("Couldn't send response %v", err)
    // }
}


func main() {
    password:= "1000000"
    hasher := md5.New()
    hasher.Write([]byte(password))
    hash := hex.EncodeToString(hasher.Sum(nil))

    numChars := 7
    lowerLimit := 0
    upperLimit := math.Pow(10, float64(numChars - 1))
    for i := 0; i < numChars - 1; i++ {
      upperLimit = upperLimit + (math.Pow(10, float64(i)))
    }
    upperLimit = upperLimit * 9

    clients := []Client{}

    addr := net.UDPAddr{
        Port: 1234,
        IP: net.ParseIP("127.0.0.1"),
    }
    ser, err := net.ListenUDP("udp", &addr)
    if err != nil {
        fmt.Printf("Some error %v\n", err)
        return
    }
    ranges := [][]int{}
    for {
        p := make([]byte, 2048)
        n, remoteaddr, err := ser.ReadFromUDP(p)
        fmt.Printf("CLIENT : %v : %s\n", remoteaddr, p[:n])
        if err != nil {
            fmt.Printf("Some error  %v", err)
            continue
        }


        isNewClient := true
        for _, v := range clients {
          if (v.address == remoteaddr) {
            isNewClient = false
          }
        }

        if isNewClient {
          if len(clients) < 1 {
            newRange := []int{int(lowerLimit), int(upperLimit)}
            c := Client{address: remoteaddr, startRange: int(lowerLimit), endRange: int(upperLimit)}
            ranges = append(ranges, []int{int(lowerLimit), int(upperLimit)})
            clients = append(clients, c)

            response := "<name>New Job</name><code>001</code><hash>"+string(hash[:])+"</hash><start>"++"</start><end>1000001</end>"
            go sendResponse(ser, remoteaddr, response)

          } else {

              for _, c range clients {
                response := "<name>Halt Request Current Range</name><code>003</code>"
              }

          }


          for _, r := range ranges {
            fmt.Printf("range %v %v\n", r[0], r[1])
          }
          fmt.Printf("number of clients = %d\n", len(clients))
          fmt.Println(clients)
          fmt.Printf("lower limit = %d\n", lowerLimit)
          fmt.Printf("upper limit = %v\n", upperLimit)

        } else {

        }


        // for i := 0; i < numClients; i++ {
        //
        // }
    }
}
