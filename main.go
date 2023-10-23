package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"go/contract"
	"go/entity"
	"go/repository/filestore"
	"go/repository/memorystore"
	"go/service/task"
	"os"
	"strconv"
)

// global variable

var (
	userStorage       []entity.User
	authenticatedUser *entity.User

	CategoryStorage []entity.Category

	serializationMode string
)

const (
	path  = "user.txt"
	Alaki = "Alaki"
	Json  = "Json"
)

func main() {
	taskMemoryRepo := memorystore.NewTaskStore()
	taskService := task.NewService(taskMemoryRepo)
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

	var userFileStore = filestore.New(path, serializationMode)
	users := userFileStore.Load()
	userStorage = append(userStorage, users...)

	for {
		runCommand(userFileStore, *command, &taskService)
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("please enter another command")
		scanner.Scan()
		*command = scanner.Text()
	}
}

func runCommand(store contract.UserWriteStore, command string, taskService *task.Service) {

	if command != "register" && command != "exit" && authenticatedUser == nil {
		fmt.Println("you must log in first")
		login()
		if authenticatedUser == nil {
			return
		}
	}

	switch command {
	case "create-task":
		createTask(taskService)
	case "task-list":
		listTask(taskService)
	case "create-category":
		createCategory()
	case "category-list":
		categoryList()
	case "register":
		register(store)
	case "exit":
		os.Exit(0)
	default:
		fmt.Println("Command not valid", command)
	}
}

func createTask(taskService *task.Service) {

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

	fmt.Println("please enter the dueDate")
	scanner.Scan()
	dueDate = scanner.Text()

	if authenticatedUser != nil {

	}

	res, cErr := taskService.Create(task.CreateRequest{
		Title:               title,
		DueDate:             dueDate,
		CategoryId:          categoryId,
		AuthenticatedUserId: authenticatedUser.ID,
	})
	if cErr != nil {
		fmt.Println("error", cErr)
		return
	}

	fmt.Println("create task", res.Task)

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

	category := entity.Category{
		ID:     len(CategoryStorage) + 1,
		Title:  title,
		Color:  color,
		UserId: authenticatedUser.ID,
	}

	CategoryStorage = append(CategoryStorage, category)

}

func register(s contract.UserWriteStore) {
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

	user := entity.User{
		ID:       len(userStorage) + 1,
		Name:     name,
		Email:    email,
		Password: hashPassword(password),
	}
	userStorage = append(userStorage, user)
	//writeUserToFile(user)
	s.Save(user)
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
		if user.Email == email && user.Password == hashPassword(password) {
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

func listTask(taskService *task.Service) {
	userTasks, err := taskService.List(task.ListRequest{UserID: authenticatedUser.ID})
	if err != nil {
		fmt.Println("error", err)
		return
	}

	fmt.Println("user tasks", userTasks.Tasks)
}

func categoryList() {
	for _, category := range CategoryStorage {
		if category.UserId == authenticatedUser.ID {
			fmt.Println(category)
		}
	}
}

func hashPassword(password string) string {
	hash := md5.Sum([]byte(password))

	return hex.EncodeToString(hash[:])
	// fmt.Println(string(hash))

}
