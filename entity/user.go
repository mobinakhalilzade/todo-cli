package entity

import "fmt"

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

func (u User) Print() {
	fmt.Println("user", u.ID, u.Name, u.Email)
}
