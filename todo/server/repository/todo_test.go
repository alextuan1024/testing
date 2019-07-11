package repository

import (
	"github.com/go-pg/pg"
	"github.com/stretchr/testify/suite"
	"github.com/xuanit/testing/todo/pb" // Update
	"testing"
	"time"
)

type ToDoRepositorySuite struct {
	db *pg.DB
	suite.Suite
	todoRep ToDoImpl
}

func (s *ToDoRepositorySuite) SetupSuite()  {
	// Connect to PostgresQL
	s.db = pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "example",
		Database: "todo",
		Addr:     "localhost" + ":" + "5433",
		RetryStatementTimeout: true,
		MaxRetries:            4,
		MinRetryBackoff:       250 * time.Millisecond,
	})

	// Create Table
	s.db.CreateTable(&pb.Todo{}, nil)

	s.todoRep = ToDoImpl{DB: s.db}
}

func (s *ToDoRepositorySuite) TearDownSuite() {
	s.db.Close()
}

func (s *ToDoRepositorySuite)TestInsert()  {

	item := &pb.Todo{ Id: "new_item4"}
	err := s.todoRep.Insert(item)

	s.Nil(err)

	newTodo, err := s.todoRep.Get(item.Id)
	s.Nil(err)
	s.Equal(item, newTodo)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ToDoRepositorySuite))
}