package bongo

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"strings"
)

type Config struct {
	ConnectionString string
	Database         string
}

// var EncryptionKey [32]byte
// var EnableEncryption bool

type Connection struct {
	Config  *Config
	Session *mgo.Session
	// collection []Collection
}

// Create a new connection and run Connect()
func Connect(config *Config) (*Connection, error) {
	conn := &Connection{
		Config: config,
	}

	err := conn.Connect()

	return conn, err
}

// Connect to the database using the provided config
func (m *Connection) Connect() (err error) {
	defer func() {
		if r := recover(); r != nil {
			// panic(r)
			// return
			if e, ok := r.(error); ok {
				err = e
			} else if e, ok := r.(string); ok {
				err = errors.New(e)
			} else {
				err = errors.New(fmt.Sprint(r))
			}

		}
	}()

	var session *mgo.Session
	var sessError error

	if strings.ContainsRune(m.Config.ConnectionString, '@') {
		parse := strings.Split(m.Config.ConnectionString, string('@'))
		auth := strings.Split(parse[0], string(':'))
		user := auth[0]
		password := auth[1]

		session, sessError = mgo.DialWithInfo(&mgo.DialInfo{
			Database: m.Config.Database,
			Username: user,
			Password: password,
			Addrs: []string{parse[1]},
		})
	} else {
		session, sessError = mgo.Dial(m.Config.ConnectionString)
	}

	if sessError != nil {
		return sessError
	}

	m.Session = session

	m.Session.SetMode(mgo.Monotonic, true)
	return nil
}

func (m *Connection) Collection(name string) *Collection {

	// Just create a new instance - it's cheap and only has name
	return &Collection{
		Connection: m,
		Name:       name,
	}
}
