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
type Job struct {
	Hash	string
	Range	Range
}
type Password struct {
	Value	string
}

func main() {
	/*
		* for this proof of concept program, we assume passwords only use chars [0-9][a-f][A-F]
		* so, 62 possible chars for a single place
		* for a 5 character password, this is 62 ^ 8
	*/
	passwordX := "abcde"
	numCharsX := 5


	/*
		test data

		* password: 73,111,111
		* major range: 0 -> 99,999,999
		* increment: 1,000,000

	*/
	password := "73111111"
	// password := "5000000"
	// hash the password
	hasher := md5.New()
	hasher.Write([]byte(password))
	hash := hex.EncodeToString(hasher.Sum(nil))
	fmt.Printf("\nHash: %s\n", hash)
	// limits
	numChars := len(password)
	lowerLimit := 0
	upperLimit := math.Pow(10, float64(numChars-1))
	for i := 0; i < numChars - 1; i++ {
		upperLimit = upperLimit + (math.Pow(10, float64(i)))
	}
	upperLimit = upperLimit * 9
	majorRange := Range {
		Start: lowerLimit,
		End: int(upperLimit),
	}
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
	    inMessage := Message{}
		err = json.Unmarshal(ba[:n], &inMessage)
		if err != nil {
			fmt.Printf("Error %v\n", err)
			return
		}
		fmt.Printf("\nMessage Received!\nFrom Address: %s\nCode: %v\n", cAddr, inMessage.Code)
		if inMessage.Code == 1 {
			minorRange := Range {
				Start: majorRange.Start,
				End: majorRange.Start + increment - 1,
			}
			job := Job {
				Hash: hash,
				Range: minorRange,
			}
			// update major range
			majorRange.Start = majorRange.Start + increment
			//
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
