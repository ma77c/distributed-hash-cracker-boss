package main

import (
	"encoding/json"
	"fmt"
	"net"
	"math"
	"crypto/md5"
	"encoding/hex"
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
	/*
		test data

		* password: 73,111,111
		* major range: 0 -> 99,999,999
		* increment: 1,000,000

	*/
	password := "73111111"
	// hash the password
	hasher := md5.New()
	hasher.Write([]byte(password))
	hash := hex.EncodeToString(hasher.Sum(nil))
	fmt.Printf("Hash: %s\n", hash)
	// limits
	numChars := len(password)
	lowerLimit := 0
	upperLimit := math.Pow(10, float64(numChars-1))
	for i := 0; i < numChars-1; i++ {
		upperLimit = upperLimit + (math.Pow(10, float64(i)))
	}
	upperLimit = upperLimit * 9
	majorRange := Range { Start: lowerLimit, End: int(upperLimit) }
	fmt.Printf("Major Range: %+v\n", majorRange)
	// increment set to  1,000,000
	increment := 1000000
    // prepare server address for listener
    sAddr := net.UDPAddr {
        Port: 1234,
        IP:   net.ParseIP("127.0.0.1"),
    }
	fmt.Printf("Server Address: %+v\n", sAddr)
    // listener for incoming clients
	conn, err := net.ListenUDP("udp", &sAddr)
    if err != nil {
        fmt.Printf("Error %v", err)
        return
    }
	// keep listening
	for {
		ba := make([]byte, 2048)
	    // blocking call
	    n, cAddr, err := conn.ReadFromUDP(ba)
	    if err != nil {
	        fmt.Printf("Error %v", err)
	        return
	    }
		fmt.Printf("Client: %s\n", cAddr)
	    inMessage := Message{}
		err = json.Unmarshal(ba[:n], &inMessage)
		if err != nil {
			fmt.Printf("Error %v\n", err)
			return
		}
		if inMessage.Code == 1 {
			minorRange := &Range {
				Start: majorRange.Start,
				End: majorRange.End,
			}
			r, err := json.Marshal(minorRange)
			if err != nil {
		        fmt.Printf("Error %v", err)
		        return
		    }
			outMessage := &Message {
				Code: 2,
				Payload: r,
			}
		}
		// json.Unmarshal(m.Payload, &r)
		// fmt.Printf("Range %+v", r)
	}
}
