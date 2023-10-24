package main

import (
	"encoding/json"
	"fmt"
	"go/delivery/param"
	"go/repository/memorystore"
	"go/service/task"
	"log"
	"net"
)

func main() {
	const (
		network = "tcp"
		address = ":9986"
	)

	//create new listener
	listener, err := net.Listen(network, address)

	if err != nil {
		log.Fatalln("Can't listen on given address", address, err)
	}

	fmt.Println("server listening on", listener.Addr())

	taskMemoryRepo := memorystore.NewTaskStore()
	taskService := task.NewService(taskMemoryRepo)

	for {
		//listen for new connection
		connection, aErr := listener.Accept()

		if aErr != nil {
			log.Fatalln("Can't listen to new connection", aErr)
		}

		//process request
		var rawRequest = make([]byte, 1024)
		numberOfReadBytes, readErr := connection.Read(rawRequest)

		if readErr != nil {
			fmt.Println("Can't read data from connection", readErr)

			continue
		}

		fmt.Printf("client address: %s,numberOfReadBytes:%d,data:%s\n ",
			connection.RemoteAddr(), numberOfReadBytes, string(rawRequest))

		req := &param.Request{}
		if uErr := json.Unmarshal(rawRequest[:numberOfReadBytes], req); uErr != nil {
			log.Println("Bad Request", uErr)

			continue
		}

		switch req.Command {
		case "create-task":
			response, cErr := taskService.Create(task.CreateRequest{
				Title:               req.CreateTaskRequest.Title,
				DueDate:             req.CreateTaskRequest.DueDate,
				CategoryId:          req.CreateTaskRequest.CategoryId,
				AuthenticatedUserId: 0,
			})

			if cErr != nil {
				_, wErr := connection.Write([]byte(cErr.Error()))

				if wErr != nil {
					log.Println("Can't write data to connection", wErr)

					continue
				}
			}

			data, mErr := json.Marshal(&response)
			if mErr != nil {
				_, wErr := connection.Write([]byte(mErr.Error()))

				if wErr != nil {
					log.Println("Can't marshal response", wErr)

					continue
				}
				continue
			}
			_, wErr := connection.Write([]byte(data))

			if wErr != nil {
				log.Println("Can't write data to connection", wErr)

				continue
			}
		}

		connection.Close()

	}
}
