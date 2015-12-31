/*
Package storage defines the interface and functions for database crud operations.
*/

package storage

import (
	"log"
	"time"

	r "github.com/dancannon/gorethink"
	e "github.com/gabesullice/accountman/errors"
)

var sess Session

// Session type aliases the rethinkDb Session
type Session *r.Session

// Storage is a data model for storing a Session and any connection options.
type Storage struct {
	Session     Session
	ConnectOpts r.ConnectOpts
}

// StorageConfig provides an interface for creating sets of configuration able
// to be passes to NewStorage()
type StorageConfig interface {
	Configure(*Storage)
}

func Insert(db string, table string, data interface{}) (string, error) {
	res, err := r.DB(db).
		Table(table).
		Insert(data).
		RunWrite(getSession())

	if err != nil {
		log.Println(err)
		return "", e.ErrDatabase
	}

	return res.GeneratedKeys[0], nil
}

func GetAll(db string, table string, result interface{}) error {
	res, err := r.DB(db).Table(table).Run(getSession())

	if err != nil {
		return err
	}

	return res.All(result)
}

func GetOne(db string, table string, id string, result interface{}) error {
	if res, err := r.DB(db).Table(table).Get(id).Run(getSession()); err != nil {
		return err
	} else {
		return res.One(result)
	}
}

func GetByName(db string, table string, name string, result interface{}) error {
	if res, err := r.DB(db).Table(table).GetAllByIndex("name", name).Run(getSession()); err != nil {
		return err
	} else {
		return res.One(result)
	}
}

func Update(db string, table string, name string, data interface{}) error {
	res, err := r.DB(db).
		Table(table).
		GetAllByIndex("name", name).
		Update(data).
		RunWrite(getSession())

	if err != nil {
		log.Println(err)
		return e.ErrDatabase
	}

	if res.Updated != 1 {
		log.Printf("Unexpected Update Result\nOrganization: %+v\nDB Response: %+v", data, res)
	}

	return nil
}

func Delete(db string, table string, id string) error {
	res, err := r.DB(db).Table(table).Get(id).Delete().RunWrite(getSession())
	if err != nil {
		return err
	}

	if res.Deleted == 0 {
		return e.ErrEntityNotFound
	}

	return nil
}

func DeleteByName(db string, table string, name string) error {
	res, err := r.DB(db).Table(table).GetAllByIndex("name", name).RunWrite(getSession())
	if err != nil {
		return err
	}

	if res.Deleted == 0 {
		return e.ErrEntityNotFound
	}

	return nil
}

// UniqueField returns whether a value is  unique accross fields and tables.
func UniqueField(value interface{}, lookups map[string]string, db string) (bool, error) {
	for table, field := range lookups {
		res, err := r.DB(db).
			Table(table).
			Field(field).
			Contains(value).
			Run(getSession())

		if err != nil {
			log.Println(err)
			return false, e.ErrDatabase
		}

		var exists bool
		if err = res.One(&exists); err != nil {
			log.Println(err)
			return false, e.ErrDatabase
		}

		if exists {
			return false, nil
		}
	}

	return true, nil
}

// GetSession exports the getSession function.
func GetSession() Session {
	return getSession()
}

// getSession returns a db session that can be passed to the database package.
func getSession() Session {
	return sess
}

// NewStorage is an initializer function for the database.
//
// It takes 0 or more StorageConfig options, which it applies in order to the
// database connection options prior to connecting to the database.
func OpenSession(settings ...StorageConfig) {
	var s Storage

	s.ConnectOpts = r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "",
		MaxIdle:  10,
		MaxOpen:  10,
		Timeout:  time.Second * 10,
	}

	for _, setting := range settings {
		setting.Configure(&s)
	}

	session, err := r.Connect(s.ConnectOpts)

	if err != nil {
		log.Fatalln(err.Error())
	}

	sess = session
}
