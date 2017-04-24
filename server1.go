package main
import (
    "fmt"
    "net"
    "crypto/md5"
    "encoding/hex"
    "math"
)


func sendResponse(conn *net.UDPConn, addr *net.UDPAddr, hash string) {
    _,err := conn.WriteToUDP([]byte("<name>New Job</name><code>001</code><hash>"+string(hash[:])+"</hash><start>0</start><end>1000001</end>"), addr)
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
    numChars := 7
    hasher := md5.New()
    hasher.Write([]byte(password))
    hash := hex.EncodeToString(hasher.Sum(nil))
    p := make([]byte, 2048)
    addr := net.UDPAddr{
        Port: 1234,
        IP: net.ParseIP("127.0.0.1"),
    }
    ser, err := net.ListenUDP("udp", &addr)
    if err != nil {
        fmt.Printf("Some error %v\n", err)
        return
    }
    var clients []*net.UDPAddr
    for {
        n, remoteaddr, err := ser.ReadFromUDP(p)
        fmt.Printf("CLIENT : %v : %s\n", remoteaddr, p[:n])
        if err != nil {
            fmt.Printf("Some error  %v", err)
            continue
        }
        clients = append(clients, remoteaddr)
        numClients := len(clients)

        lowerLimit := 0
        upperLimit := math.Pow(10, float64(numChars - 1))
        for i := 0; i < numChars - 1; i++ {
          upperLimit = upperLimit + (math.Pow(10, float64(i)))
        }
        upperLimit = upperLimit * 9
        fmt.Printf("number of clients = %d\n", numClients)
        fmt.Printf("lower limit = %d\n", lowerLimit)
        fmt.Printf("upper limit = %v\n", upperLimit)

        // for i := 0; i < numClients; i++ {
        //
        // }

        go sendResponse(ser, remoteaddr, hash)
    }
}
