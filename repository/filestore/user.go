package filestore

import (
	"encoding/json"
	"errors"
	"fmt"
	"go/entity"
	"os"
	"strconv"
	"strings"
)

const (
	path  = "user.txt"
	Alaki = "Alaki"
	Json  = "Json"
)

type FileStore struct {
	filePath          string
	serializationMode string
}

func New(path, serializationMode string) FileStore {
	return FileStore{filePath: path, serializationMode: serializationMode}
}

func (f FileStore) Save(u entity.User) {
	f.writeUserToFile(u)
}

func (f FileStore) Load() []entity.User {
	var uStore []entity.User
	file, err := os.Open(f.filePath)
	if err != nil {
		fmt.Println("Error", err)
	}
	var data []byte = make([]byte, 1024)
	_, rErr := file.Read(data)
	if rErr != nil {
		fmt.Println("Error", rErr)
	}

	var dataStr = string(data)
	// dataStr = strings.Trim(dataStr, "\n")
	userSlice := strings.Split(dataStr, "\n")
	for _, u := range userSlice {
		var userStruct = entity.User{}

		switch f.serializationMode {
		case Alaki:
			var dErr error
			userStruct, dErr = deserializeFromAlaki(u)
			if dErr != nil {
				fmt.Println("Cant deserialize")

				return nil
			}
		case Json:
			fmt.Println("Json mode", u)
			uErr := json.Unmarshal([]byte(u), &userStruct)
			if uErr != nil {
				fmt.Println("Cant unmarshall")

				return nil
			}

		default:
			fmt.Println("Invalid mode")

			return nil
		}
		uStore = append(uStore, userStruct)
	}
	return uStore
}

func (f FileStore) writeUserToFile(user entity.User) {
	var file *os.File
	file, err := os.OpenFile(f.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("File error", err)

		return
	}

	defer file.Close()
	var data []byte
	if f.serializationMode == Alaki {
		data = []byte(fmt.Sprintf("id:%d, name:%s, email:%s, password:%s\n", user.ID, user.Name, user.Email, user.Password))

	} else if f.serializationMode == Json {

		var jErr error
		data, jErr = json.Marshal(user)
		if err != nil {
			fmt.Println("File error", jErr)

			return
		}

		data = append(data, []byte("\n")...)
	} else {
		fmt.Println("Sm Error")

		return
	}

	numOfWrittenBytes, wErr := file.Write(data)
	if wErr != nil {
		fmt.Println("cant write to file")

		return
	}

	fmt.Println("numOfWrittenBytes", numOfWrittenBytes)
}

func deserializeFromAlaki(userStr string) (entity.User, error) {
	if userStr == "" {
		return entity.User{}, errors.New("Err")
	}
	var user = entity.User{}
	userFields := strings.Split(userStr, ",")
	for _, field := range userFields {
		values := strings.Split(field, ": ")
		if len(values) != 2 {
			continue
		}
		fieldName := strings.ReplaceAll(values[0], " ", "")
		fieldValue := values[1]

		switch fieldName {
		case "id":
			id, err := strconv.Atoi(fieldValue)
			if err != nil {
				fmt.Println(err)

				return entity.User{}, errors.New("Err")
			}
			user.ID = id
		case "name":
			user.Name = fieldValue
		case "email":
			user.Email = fieldValue
		case "password":
			user.Password = fieldValue
		}
	}
	// fmt.Printf("user: %+v\n", user)
	return user, nil
}
