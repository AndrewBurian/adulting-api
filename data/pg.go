package data

import "github.com/go-pg/pg"

type FullDAL interface {
	UserDAL
	ActivityDAL
}

type pgDal struct {
	db *pg.DB
}
