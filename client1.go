package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Client struct {
	ID       int
	Address  *net.UDPAddr
	Start    int
	End      int
	Code     string
	Hash     string
	Password int
}

type Message struct {
	Code string
}

func keepReadingXXX(conn net.Conn, serverMessageChannel chan []byte) {
	var byteArray = make([]byte, 2048)
	for {
		n, err := conn.Read(byteArray)
		if err == nil {
			fmt.Printf("got message\n")
			serverMessageChannel <- byteArray[:n]
		}
	}
}

func unHashXXX(inClient Client, unHashChannel chan []byte, killChannel chan []byte) {

	fmt.Printf("starting new job = %s\n", inClient.Hash)
	for i := inClient.Start; i <= inClient.End; i++ {
		select {
			case kch := <-killChannel:
				fmt.Printf("wtf %v\n", kch)
				continue

			default:
				iString := strconv.Itoa(i)
				hasher := md5.New()
				hasher.Write([]byte(iString))
				possibleHash := hex.EncodeToString(hasher.Sum(nil))
				// fmt.Printf("current hash = %s\n", possibleHash)
				if possibleHash == inClient.Hash {
					inClient.Password = i
					inClient.Code = "0091"
					assignment, err := json.Marshal(inClient)
					if err != nil {
						fmt.Printf("Error %v", err)
						return
					}
					unHashChannel <- assignment
					return
				}
		}
	}
	fmt.Printf("finished job - need new work\n")
	inClient.Code = "0000"
	assignment, err := json.Marshal(inClient)

	if err != nil {
		fmt.Printf("Error %v", err)
		return
	}
	unHashChannel <- assignment
	return
}

func main() {
	// dial connection
	conn, err := net.Dial("udp", "127.0.0.1:1234")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	// send message to server
	inClient := &Client{ID: -1, Code: "0000"}
	fmt.Printf("%+v\n", inClient)
	assignment, err := json.Marshal(inClient)
	if err != nil {
		fmt.Printf("Error %v", err)
		return
	}
	fmt.Println(string(assignment))
	fmt.Fprintf(conn, string(assignment))
	// channels
	serverMessageChannel := make(chan []byte)
	unHashChannel := make(chan []byte)
	killChannel := make(chan []byte)
	go keepReadingXXX(conn, serverMessageChannel)
	var working bool
	for {
		select {
			case smch := <-serverMessageChannel:
				fmt.Printf("SERVER : %s\n", smch)
				inClient := Client{}
				if err == json.Unmarshal(smch, &inClient) {
					fmt.Printf("Error %v", err)
					return
				}
				if inClient.Code == "0001" {
					go unHashXXX(inClient, unHashChannel, killChannel)
					working = true
				} else if inClient.Code == "0020" {
					if working == true {
						killChannel <- smch
						working = false
					}
				}
			case uhch := <-unHashChannel:
				inClient := Client{}
				if err == json.Unmarshal(uhch, &inClient) {
					if err != nil {
						fmt.Printf("Error %v", err)
						return
					}
				}
				//// response - return ranges
				if inClient.Code == "0021" {
					//// forward ranges to server
					fmt.Fprintf(conn, string(uhch))
				} else if inClient.Code == "0091" {
					fmt.Fprintf(conn, string(uhch))
					fmt.Printf("password = %s\n", inClient.Password)
					fmt.Printf("END!")
					return
				} else if inClient.Code == "0000" {
					fmt.Fprintf(conn, string(uhch))
				}
		}
	}

}
