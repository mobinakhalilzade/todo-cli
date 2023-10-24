package main

import (
	"encoding/json"
	"fmt"
	"go/delivery/param"
	"log"
	"net"
	"os"
)

func main() {
	fmt.Println("command", os.Args[0])

	if len(os.Args) < 2 {
		log.Fatalln("You should set the ip address of server")
	}

	serverAddress := os.Args[1]

	message := "default"
	if len(os.Args) > 2 {
		message = os.Args[2]
	}

	connection, err := net.Dial("tcp", serverAddress)

	if err != nil {
		log.Fatalln("Can't dial the given address", err)
	}
	fmt.Println("local address", connection.LocalAddr())

	req := param.Request{Command: message}

	if req.Command == "create-task" {
		req.CreateTaskRequest = param.CreateTaskRequest{
			Title:      "test",
			DueDate:    "test",
			CategoryId: 1,
		}
	}

	serializedData, sErr := json.Marshal(&req)

	if sErr != nil {
		log.Fatalln("Can't write data", sErr)
	}

	numberOfWrittenByte, wErr := connection.Write(serializedData)

	if wErr != nil {
		log.Fatalln("Can't write data", wErr)
	}
	fmt.Println("numberOfWrittenByte", numberOfWrittenByte)

	var data = make([]byte, 1024)
	_, readErr := connection.Read(data)

	if readErr != nil {
		fmt.Println("Can't read data from connection", readErr)

	}

	fmt.Println("server response", string(data))
}
