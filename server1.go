package main
import (
    "fmt"
    "time"
    "net"
)


func sendResponse(conn *net.UDPConn, addr *net.UDPAddr, hash string) {
    _,err := conn.WriteToUDP([]byte("<name>New Job</name><code>001</code><hash>36</hash><start>0</start><end>50</end>"), addr)
    if err != nil {
        fmt.Printf("Couldn't send response %v", err)
    }
    time.Sleep(3*time.Second)
    _,err = conn.WriteToUDP([]byte("<name>New Job</name><code>001</code><hash>36</hash><start>0</start><end>50</end>"), addr)
    if err != nil {
        fmt.Printf("Couldn't send response %v", err)
    }
}


func main() {
    hash := "36";
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
    for {
        n, remoteaddr, err := ser.ReadFromUDP(p)
        fmt.Printf("CLIENT : %v : %s\n", remoteaddr, p[:n])
        if err != nil {
            fmt.Printf("Some error  %v", err)
            continue
        }
        go sendResponse(ser, remoteaddr, hash)
    }
}
