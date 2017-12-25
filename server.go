package main

import (
	"encoding/json"
	"fmt"
	"net"
	"crypto/md5"
	"encoding/hex"
	"github.com/m4ttclendenen/basen"
)

type Message struct {
	Code	int
	Payload	json.RawMessage
}
type Range struct {
	Start	[]byte
	End		[]byte
}
type User struct {
	ID		int
	Pass	int
}
type Job struct {
	Hash	string
	Range	Range
}
type Password struct {
	Value	string
}

func main() {
	////////////////////////////////////////////////////////////////////////////
	/////////////////////////////// Test Data //////////////////////////////////
	////////////////////////////////////////////////////////////////////////////
	password := "ax3"
	// hash the password
	hasher := md5.New()
	hasher.Write([]byte(password))
	hash := hex.EncodeToString(hasher.Sum(nil))
	fmt.Printf("\nHash: %s\n", hash)

	// base 62 arithmetic
	bn := basen.NewBaseN([]byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"), 62)
	// Starting string
	currBase := []byte("0")
	// Increment set to 4C92 = 1,000,000 in base10
	increment := []byte("4C92")
    // Prepare server address for listener

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
		read := make([]byte, 2048)
	    // Blocking call
	    n, cAddr, err := conn.ReadFromUDP(read)
	    if err != nil {
	        fmt.Printf("Error %v", err)
	        return
	    }
	    inMessage := Message{}
		err = json.Unmarshal(read[:n], &inMessage)
		if err != nil {
			fmt.Printf("Error %v\n", err)
			return
		}
		fmt.Printf("\nMessage Received!\nFrom Address: %s\nCode: %v\n", cAddr, inMessage.Code)
		if inMessage.Code == 1 {
			newBase := bn.Add(currBase, increment)
			minorRange := Range {
				Start: currBase,
				End: newBase,
			}
			job := Job {
				Hash: hash,
				Range: minorRange,
			}
			// Update current base
			currBase = newBase
			// Marshall data
			j, err := json.Marshal(job)
			if err != nil {
		        fmt.Printf("Error %v", err)
		        return
		    }
			outMessage := Message {
				Code: 2,
				Payload: j,
			}
			om, err := json.Marshal(outMessage)
			if err != nil {
		        fmt.Printf("Error %v", err)
		        return
		    }
			_, err = conn.WriteToUDP(om, cAddr)
			if err != nil {
				fmt.Printf("Error %v", err)
				return
			}
			fmt.Printf("\nMessage Sent!\nTo Address: %s\nCode: %v\nRange: %v - %v\n",
						cAddr, outMessage.Code, minorRange.Start, minorRange.End)
		} else if inMessage.Code == 3 {
			password := Password{}
			err = json.Unmarshal(inMessage.Payload, &password)
			fmt.Printf("Client Found the Password! %s\n", password.Value)
			return
		}
	}
}
