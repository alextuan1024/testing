// +build integration persistence

package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-pg/pg"
	"github.com/khaiql/dbcleaner/engine"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/xuanit/testing/todo/pb"
	"gopkg.in/khaiql/dbcleaner.v2"
)

var cleaner = dbcleaner.New()

type ToDoRepositorySuite struct {
	db *pg.DB
	suite.Suite
	todoRep ToDoImpl
}

func (s *ToDoRepositorySuite) SetupSuite() {
	// Connect to PostgresQL
	opt := &pg.Options{
		User:                  "postgres",
		Password:              "example",
		Database:              "todo",
		Addr:                  "localhost" + ":" + "5433",
		RetryStatementTimeout: true,
		MaxRetries:            4,
		MinRetryBackoff:       250 * time.Millisecond,
	}
	s.db = pg.Connect(opt)

	// Create Table
	_ = s.db.CreateTable(&pb.Todo{}, nil)

	s.todoRep = ToDoImpl{DB: s.db}

	// setup cleaner
	dbe := engine.NewPostgresEngine(fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", opt.User, opt.Password, opt.Addr, opt.Database))
	cleaner.SetEngine(dbe)
}

func (s *ToDoRepositorySuite) TearDownSuite() {
	_ = s.db.DropTable(&pb.Todo{}, nil)
	_ = s.db.Close()
}

func (s *ToDoRepositorySuite) TestInsert() {
	item := &pb.Todo{Id: "new_item", Title: "meeting"}
	err := s.todoRep.Insert(item)

	s.Nil(err)

	newTodo, err := s.todoRep.Get(item.Id)
	s.Nil(err)
	s.Equal(item, newTodo)
}

func (s *ToDoRepositorySuite) BeforeTest(_, _ string) {
	cleaner.Acquire("todos")
}
func (s *ToDoRepositorySuite) AfterTest(_, _ string) {
	cleaner.Clean("todos")
}

func (s *ToDoRepositorySuite) TestGet() {
	item := &pb.Todo{Id: "star_platinum", Title: "Strongest stand"}
	err := s.todoRep.Insert(item)
	s.Nil(err)

	star, err := s.todoRep.Get(item.Id)
	s.Nil(err)
	s.Equal(item, star)
}

func (s *ToDoRepositorySuite) TestList() {
	tda := &pb.Todo{Id: "magician_s_red", Title: "Avdol"}
	err := s.todoRep.Insert(tda)
	s.Nil(err)

	tdb := &pb.Todo{Id: "silver_chariot", Title: "Polnareff"}
	err = s.todoRep.Insert(tdb)
	s.Nil(err)

	expected := []*pb.Todo{tda, tdb}
	todos, err := s.todoRep.List(2, false)
	s.Nil(err)
	s.Equal(expected, todos)

}

func (s *ToDoRepositorySuite) TestDelete() {
	td := &pb.Todo{Id: "last_mission", Title: "Dio"}
	err := s.todoRep.Insert(td)
	s.Nil(err)

	err = s.db.Delete(&pb.Todo{Id: td.Id})
	s.Nil(err)
}

func TestToDoRepository(t *testing.T) {
	suite.Run(t, new(ToDoRepositorySuite))
}
