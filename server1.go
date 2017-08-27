package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"strings"
)


/*

	new work
		-

	stop work

	reporting

	has Password




*/



//// codes ////
//
//// 0000 - Request - New Job
//
//// 0001 - Response - New Job
//
//// 0010 - Request -
//
//// 0011 - Request -
//
//// 0020 - Request -
//
//// 0021 - Response -
//
//// 0030 - Request -
//
//// 0031 - Response -
//

type Range struct {
	Start    int
	End      int
	ClientID int
}

type Client struct {
	ID       int
	Address  *net.UDPAddr
	Start    int
	End      int
	Code     string
	Hash     string
	Password string
}

func send(conn *net.UDPConn, addr *net.UDPAddr, response []byte) {
	_, err := conn.WriteToUDP([]byte(response), addr)
	if err != nil {
		fmt.Printf("Error %v", err)
		return
	}
}

func removeRangeByID(rangeArray []Range, ClientID int) []Range {
	for _, r := range rangeArray {
		if r.ClientID == ClientID {
			r.ClientID = rangeArray[len(rangeArray)-1].ClientID
			r.Start = rangeArray[len(rangeArray)-1].Start
			r.End = rangeArray[len(rangeArray)-1].End
			return rangeArray[:len(rangeArray)-1]
		}
	}
	return nil
}

func main() {
	// hash
	password := "73111111"
	hasher := md5.New()
	hasher.Write([]byte(password))
	hash := hex.EncodeToString(hasher.Sum(nil))
	// limits
	numChars := 8
	lowerLimit := 0
	upperLimit := math.Pow(10, float64(numChars-1))
	for i := 0; i < numChars-1; i++ {
		upperLimit = upperLimit + (math.Pow(10, float64(i)))
	}
	upperLimit = upperLimit * 9
	// ranges
	minorRanges := []Range{}
	majorRange := Range{Start: lowerLimit, End: int(upperLimit)}

	// increment set to  1,000,000
	increment := 1000000

	activeClients := []Client{}

	//// prepare address for listener
	addr := net.UDPAddr{
		Port: 1234,
		IP:   net.ParseIP("127.0.0.1"),
	}
	//// listener for incoming clients
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Error %v", err)
		return
	}
	//// keep listening - cycle when receive message
	for {
		byteArray := make([]byte, 2048)
		//// blocks until reads a message
		n, inAddress, err := conn.ReadFromUDP(byteArray)
		if err != nil {
			fmt.Printf("Error %v", err)
			return
		}
		//// boom - message received - lets parse the message into a temporary client
		//// initiate temporary client
		inClient := Client{}
		if err == json.Unmarshal(byteArray[:n], &inClient) {
			if err != nil {
				fmt.Printf("Error %v", err)
				return
			}
		}
		inClient.Address = inAddress
		fmt.Printf("CLIENT : %v : %+v\n", inClient.Address, inClient)

		//// if brand new client
		if inClient.ID == -1 {
			//// assign client ID
			greatestID := -1
			for _, ac := range activeClients {
				if ac.ID > greatestID {
					greatestID = ac.ID
				}
			}
			if greatestID == -1 {
				inClient.ID = 0
			} else {
				inClient.ID = greatestID + 1
			}
			inClient.Start = majorRange.Start
			inClient.End = majorRange.Start + increment - 1
			majorRange.Start = inClient.End + 1
			minor := Range{Start: inClient.Start, End: inClient.End, ClientID: inClient.ID}
			minorRanges = append(minorRanges, minor)
			inClient.Code = "0001"
			inClient.Hash = hash
			fmt.Printf("CLIENT : %v : %v\n", inClient.Address, inClient)
			assignment, err := json.Marshal(inClient)
			if err != nil {
				fmt.Printf("Error %v", err)
				return
			}
			go send(conn, inClient.Address, assignment)
			activeClients = append(activeClients, inClient)

			//// if active client
		} else {
			if inClient.Code == "0000" {
				minorRanges = removeRangeByID(minorRanges, inClient.ID)
				inClient.Start = majorRange.Start
				inClient.End = majorRange.Start + increment - 1
				majorRange.Start = inClient.End + 1
				minor := Range{Start: inClient.Start, End: inClient.End, ClientID: inClient.ID}
				minorRanges = append(minorRanges, minor)
				inClient.Code = "0001"
				assignment, err := json.Marshal(inClient)
				if err != nil {
					fmt.Printf("Error %v", err)
					return
				}
				go send(conn, inClient.Address, assignment)
			} else if inClient.Code == "0091" {
				fmt.Printf("password = %s\n", inClient.Password)
				fmt.Printf("END!")
				return
			}
		}
	}
}
