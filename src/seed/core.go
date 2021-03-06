package seed

import (
	"time"
)

// seed specific types
type User struct {
	Username string
	Name string
	Email string
	IsDriver bool
	Balance string
	Picture string
	Age int32
	Gender string
	Insurer string
	Phone string
	CC int64
	Address string
	About string
	Registered time.Time
	Latitude float64
	Longitude float64
	Tags []string
	Friends []string
	FavoriteItemArray []string
}

type DbLayer interface {
	OpenDB(constr string) (interface{}, error)
	CreateUsers(users []User) ([]uint64, error)
	CreateCol(col string) (error)
	Update(col string, docid uint64, data string) error
	Delete(col string, docid uint64) error
	Query(col, querystr string) ([]byte, error)
	Get(col string, docid uint64 ) (interface{}, error)
}

