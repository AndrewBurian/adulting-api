package data

import (
	"fmt"

	"github.com/go-pg/pg"
)

var ErrNotFound = fmt.Errorf("Entry not found")

type FullDAL interface {
	UserDAL
	ActivityDAL
}

type pgDal struct {
	db *pg.DB
}
