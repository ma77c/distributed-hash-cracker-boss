package main

import (
	"encoding/json"
	"fmt"
	"net"
)

type Message struct {
	Code	int
	Payload	json.RawMessage
}
type Range struct {
	Start	int
	End		int
}
type User struct {
	ID		int
	Pass	int
}


func main() {
    // prepare server address for listener
    sAddr := net.UDPAddr {
        Port: 1234,
        IP:   net.ParseIP("127.0.0.1"),
    }
    // listener for incoming clients
	conn, err := net.ListenUDP("udp", &sAddr)
    if err != nil {
        fmt.Printf("Error %v", err)
        return
    }
    ba := make([]byte, 2048)
    // blocking
    n, cAddr, err := conn.ReadFromUDP(ba)
    if err != nil {
        fmt.Printf("Error %v", err)
        return
    }
	fmt.Printf("Message from address: %s\n", cAddr)
    m := Message{}
	err = json.Unmarshal(ba[:n], &m)
	if err != nil {
		fmt.Printf("Error %v", err)
		return
	}
	r := Range{}
	json.Unmarshal(m.Payload, &r)
	fmt.Printf("Range %+v", r)
}
