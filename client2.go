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
	r := &Range { Start: 1, End: 10, }
	j, err := json.Marshal(r)
	m := &Message { Code: 4, Payload: j, }
	j, err = json.Marshal(m)
	fmt.Printf("JSON %+s\n", j)
	if err != nil {
		fmt.Printf("Error %v", err)
		return
	}
	fmt.Fprintf(conn, string(j))

}
