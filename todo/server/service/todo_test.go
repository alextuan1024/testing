package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/xuanit/testing/todo/pb"
	"github.com/xuanit/testing/todo/server/repository/mocks"
	"testing"
)

func TestGetToDo(t *testing.T) {
	mockToDoRep := &mocks.ToDo{}
	toDo := &pb.Todo{}
	req := &pb.GetTodoRequest{Id: "123"}
	mockToDoRep.On("Get", req.Id).Return(toDo, nil)
	service := ToDo{ToDoRepo: mockToDoRep}

	res, err := service.GetTodo(nil, req)

	expectedRes := &pb.GetTodoResponse{Item: toDo}

	assert.Nil(t, err)
	assert.Equal(t, expectedRes, res)
	mockToDoRep.AssertExpectations(t)
}

func TestCreateTodo(t *testing.T) {
	mockToDoRepo := &mocks.ToDo{}
	toDo := &pb.Todo{Title: "Star Platinum", Description: "The World"}
	req := &pb.CreateTodoRequest{Item: toDo}
	mockToDoRepo.On("Insert", toDo).Return(nil)

	service := ToDo{ToDoRepo: mockToDoRepo}
	res, err := service.CreateTodo(nil, req)

	expectedRes := &pb.CreateTodoResponse{Id: toDo.Id}

	assert.Nil(t, err)
	assert.Equal(t, expectedRes, res)
	mockToDoRepo.AssertExpectations(t)
}

func TestListTodo(t *testing.T) {
	// mocking DB
	mockTodoRepo := &mocks.ToDo{}
	toDos := []*pb.Todo{
		{
			Id:          "1",
			Title:       "Jonathan Joestar",
			Description: "Hamon user",
		}, {
			Id:          "3",
			Title:       "Jotaro Kujo",
			Description: "Stand user",
		},
	}
	req := &pb.ListTodoRequest{Limit: 2, NotCompleted: false}

	// set mock behavior
	mockTodoRepo.On("List", req.Limit, req.NotCompleted).Return(toDos, nil)

	// create service, using mocked db
	service := ToDo{ToDoRepo: mockTodoRepo}
	expectedRes := &pb.ListTodoResponse{Items: toDos}

	// call ListTodo
	res, err := service.ListTodo(nil, req)

	// assertion
	assert.Nil(t, err)
	assert.Equal(t, expectedRes, res)
	mockTodoRepo.AssertExpectations(t)

}

func TestDeleteTodo(t *testing.T) {
	mockTodoRepo := &mocks.ToDo{}

	req := &pb.DeleteTodoRequest{Id: "Mark43"}
	expectedRes := &pb.DeleteTodoResponse{}
	mockTodoRepo.On("Delete", req.Id).Return(nil)

	service := ToDo{ToDoRepo: mockTodoRepo}
	res, err := service.DeleteTodo(nil, req)

	assert.Nil(t, err)
	assert.Equal(t, expectedRes, res)
	mockTodoRepo.AssertExpectations(t)
}
