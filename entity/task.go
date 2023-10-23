package entity

type Task struct {
	ID         int
	Title      string
	DueDate    string
	CategoryId int
	IsDone     bool
	UserId     int
}
