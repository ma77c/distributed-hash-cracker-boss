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
	// dial connection
	conn, err := net.Dial("udp", "127.0.0.1:1234")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	// build marshalled data
	outMessage := &Message { Code: 1, Payload: nil, }
	om, err := json.Marshal(outMessage)
	if err != nil {
		fmt.Printf("Error %v", err)
		return
	}
	fmt.Fprintf(conn, string(om))
	fmt.Printf("Message Sent!\nCode: 1\n")
	// lets keep reading
	for {
		ba := make([]byte, 2048)
		n, err := conn.Read(ba)
		if err != nil {
			fmt.Printf("Error %v", err)
			return
		}
		inMessage := Message{}
		err = json.Unmarshal(ba[:n], &inMessage)
		if err != nil {
			fmt.Printf("Error %v", err)
			return
		}
		if inMessage.Code == 2 {
			minorRange := Range{}
			err = json.Unmarshal(inMessage.Payload, &minorRange)
			if err != nil {
				fmt.Printf("Error %v", err)
				return
			}
			fmt.Printf("\nMessage Received!\nCode: %v\nRange: %+v\n", inMessage.Code, minorRange)
		}
	}
}
