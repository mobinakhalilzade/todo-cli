package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

type Task struct {
	ID         int
	Title      string
	DueDate    string
	categoryId int
	IsDone     bool
	UserId     int
}

type Category struct {
	ID     int
	Title  string
	Color  string
	UserId int
}

// global variable

var (
	userStorage       []User
	authenticatedUser *User

	taskStorage     []Task
	CategoryStorage []Category

	serializationMode string
)

const (
	path  = "user.txt"
	Alaki = "alaki"
	Json  = "json"
)

func main() {
	loadUserStorageFromFile()
	fmt.Println("Hello ...")
	command := flag.String("command", "no command", "create a new task")
	serializeMode := flag.String("serialize", Alaki, "serialize")
	flag.Parse()

	switch *serializeMode {
	case Alaki:
		serializationMode = Alaki
	default:
		serializationMode = Json
	}
	for {
		runCommand(*command)
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("please enter another command")
		scanner.Scan()
		*command = scanner.Text()
	}
}

func runCommand(command string) {

	if command != "register" && command != "exit" && authenticatedUser == nil {
		fmt.Println("you must log in first")
		login()
		if authenticatedUser == nil {
			return
		}
	}

	switch command {
	case "create-task":
		createTask()
	case "task-list":
		taskList()
	case "create-category":
		createCategory()
	case "category-list":
		categoryList()
	case "register":
		register()
	case "exit":
		os.Exit(0)
	default:
		fmt.Println("Command not valid", command)
	}
}

func (u User) print() {
	fmt.Println("user", u.ID, u.Name, u.Email)
}

func createTask() {

	if authenticatedUser != nil {
		authenticatedUser.print()
	}
	scanner := bufio.NewScanner(os.Stdin)

	var title, dueDate, category string
	fmt.Println("please enter the title")
	scanner.Scan()
	title = scanner.Text()

	fmt.Println("please enter the category")
	scanner.Scan()
	category = scanner.Text()

	categoryId, err := strconv.Atoi(category)
	if err != nil {
		fmt.Println("category id is not valid", err)
		return
	}

	isFound := false
	for _, c := range CategoryStorage {
		if c.ID == categoryId && c.UserId == authenticatedUser.ID {
			isFound = true

			break
		}
	}

	if !isFound {
		fmt.Println("category id is not valid")

		return
	}

	fmt.Println("please enter the dueDate")
	scanner.Scan()
	dueDate = scanner.Text()

	if authenticatedUser != nil {
		task := Task{
			ID:         len(taskStorage) + 1,
			Title:      title,
			categoryId: categoryId,
			DueDate:    dueDate,
			IsDone:     false,
			UserId:     authenticatedUser.ID,
		}
		taskStorage = append(taskStorage, task)
	}
	fmt.Println("task", title, dueDate, category)
}

func createCategory() {
	scanner := bufio.NewScanner(os.Stdin)

	var title, color string
	fmt.Println("please enter the category title")
	scanner.Scan()
	title = scanner.Text()

	fmt.Println("please enter the color")
	scanner.Scan()
	color = scanner.Text()

	fmt.Println("category", title, color)

	category := Category{
		ID:     len(CategoryStorage) + 1,
		Title:  title,
		Color:  color,
		UserId: authenticatedUser.ID,
	}

	CategoryStorage = append(CategoryStorage, category)

}

func register() {
	scanner := bufio.NewScanner(os.Stdin)

	var name, email, password string

	fmt.Println("please enter the email")
	scanner.Scan()
	email = scanner.Text()

	fmt.Println("please enter your name")
	scanner.Scan()
	name = scanner.Text()

	fmt.Println("please enter the password")
	scanner.Scan()
	password = scanner.Text()
	user := User{
		ID:       len(userStorage) + 1,
		Name:     name,
		Email:    email,
		Password: password,
	}
	userStorage = append(userStorage, user)
	writeUserToFile(user)
}

func login() {

	fmt.Println("login process")
	scanner := bufio.NewScanner(os.Stdin)

	var email, password string
	fmt.Println("please enter the email")
	scanner.Scan()
	email = scanner.Text()

	fmt.Println("please enter the password")
	scanner.Scan()
	password = scanner.Text()

	fmt.Println("user", email, password)

	for _, user := range userStorage {
		if user.Email == email && user.Password == password {
			fmt.Println("You are logged in")
			authenticatedUser = &user

			break
		}

	}

	if authenticatedUser == nil {
		fmt.Println("Your data is wrong")

		return
	}
}

func taskList() {
	for _, task := range taskStorage {
		if task.UserId == authenticatedUser.ID {
			fmt.Println(task)
		}
	}
}

func categoryList() {
	for _, category := range CategoryStorage {
		if category.UserId == authenticatedUser.ID {
			fmt.Println(category)
		}
	}
}
func loadUserStorageFromFile() {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error", err)
	}
	var data []byte = make([]byte, 10240)
	_, rErr := file.Read(data)
	if rErr != nil {
		fmt.Println("Error", rErr)
	}

	var dataStr = string(data)
	//dataStr = strings.Trim(dataStr, "\n")
	userSlice := strings.Split(dataStr, "\n")
	for _, u := range userSlice {
		if u == "" {
			continue
		}
		userFields := strings.Split(u, ",")
		for _, field := range userFields {
			values := strings.Split(field, ": ")
			if len(values) != 2 {
				continue
			}
			fieldName := strings.ReplaceAll(values[0], " ", "")
			fieldValue := values[1]

			var user = User{}
			switch fieldName {
			case "id":
				id, err := strconv.Atoi(fieldValue)
				if err != nil {
					fmt.Println(err)

					return
				}
				user.ID = id
			case "name":
				user.Name = fieldValue
			case "email":
				user.Email = fieldValue
			case "password":
				user.Password = fieldValue
			}
			fmt.Printf("user: %+v\n", user)
		}
	}

}
func writeUserToFile(user User) {
	var file *os.File
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("File error", err)

		return
	}

	defer file.Close()
	var data []byte
	if serializationMode == Alaki {
		data = []byte(fmt.Sprintf("id:%s, name:%s, email:%s, password:%s\n", user.ID, user.Name, user.Email, user.Password))

	} else if serializationMode == Json {

		var jErr error
		data, jErr = json.Marshal(user)
		if err != nil {
			fmt.Println("File error", jErr)

			return
		}
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
