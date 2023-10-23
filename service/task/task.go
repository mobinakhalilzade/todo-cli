package task

import (
	"fmt"
	"go/entity"
)

// ServiceRepository : It's a place where everything that task service is going to need it, we put it here
type ServiceRepository interface {
	//DoesThisUserHaveThisCategoryID(userID, categoryID int) bool
	CreateNewTask(t entity.Task) (entity.Task, error)
	ListUserTask(userID int) ([]entity.Task, error)
}

type Service struct {
	repository ServiceRepository
}

func NewService(rep ServiceRepository) Service {
	return Service{
		repository: rep,
	}
}

type CreateRequest struct {
	Title               string
	DueDate             string
	CategoryId          int
	AuthenticatedUserId int
}

type CreateResponse struct {
	Task entity.Task
}

func (s Service) Create(req CreateRequest) (CreateResponse, error) {

	//ok := s.repository.DoesThisUserHaveThisCategoryID(req.AuthenticatedUserId, req.CategoryId)
	//
	//if !ok {
	//	return CreateResponse{}, fmt.Errorf("user does not have this category: %d", req.CategoryId)
	//}

	createdTask, cErr := s.repository.CreateNewTask(entity.Task{
		Title:      req.Title,
		DueDate:    req.DueDate,
		CategoryId: req.CategoryId,
		IsDone:     false,
		UserId:     req.AuthenticatedUserId,
	})
	if cErr != nil {
		return CreateResponse{}, fmt.Errorf("cant create new task: %v", cErr)
	}

	return CreateResponse{Task: createdTask}, nil
}

type ListRequest struct {
	UserID int
}

type ListResponse struct {
	Tasks []entity.Task
}

func (s Service) List(req ListRequest) (ListResponse, error) {
	tasks, err := s.repository.ListUserTask(req.UserID)
	if err != nil {
		return ListResponse{}, fmt.Errorf("cant list user tasks: %v", err)
	}

	return ListResponse{Tasks: tasks}, nil
}

//func (t Task) editTask()                        {}
//func (t Task) deleteTask()                      {}
