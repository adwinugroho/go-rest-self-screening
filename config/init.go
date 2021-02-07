package config

import (
	"context"
	"log"
	"os"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

// create global variable, we use environment variable.
var (
	// create a variable for return connection struct (pointer)
	Instance *Connection
	// DB URL
	DBURL string = os.Getenv("DBURL")
	// DB username
	DBUsername string = os.Getenv("DBUSERNAME")
	// DB password
	DBPassword string = os.Getenv("DBPASSWORD")
	// DB Name
	DBName string = os.Getenv("DBNAME")
	// DB Log Name
	DBLogName string = os.Getenv("DBLOGNAME")
	// DB Log URL
	DBLOGURL string = os.Getenv("DBLOGURL")
)

// create a struct for save DB (driver.Database)
type Connection struct {
	DBLive driver.Database
	DBLog  driver.Database
}

func init() {
	// create a connection to DB
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{DBURL}, //DB Url
	})
	if err != nil {
		panic(err)
	}
	// create a new connection to DB client
	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(DBUsername, DBPassword),
	})
	if err != nil {
		panic(err)
	}
	// create a variable to initial context background (for connection between client and DB)
	ctx := context.Background()
	// connect client and DB
	db, err := client.Database(ctx, DBName)
	if err != nil {
		log.Printf("Error connecting to database, cause: %+v\n", err)
		panic(err)
	}
	// same as DBLive, we create client for log
	connLog, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{DBLOGURL}, //DB Url
	})
	if err != nil {
		panic(err)
	}
	// create a new connection to DB client
	clientLog, err := driver.NewClient(driver.ClientConfig{
		Connection:     connLog,
		Authentication: driver.BasicAuthentication(DBUsername, DBPassword),
	})
	if err != nil {
		panic(err)
	}
	// connect client and DB
	dbLog, err := clientLog.Database(ctx, DBLogName)
	if err != nil {
		log.Printf("Error connecting to database, cause: %+v\n", err)
		panic(err)
	}

	//in the end we assign DB client to variable instance (*Connection)
	Instance = &Connection{
		DBLive: db,
		DBLog:  dbLog,
	}
}

// create function to get variable instance (*Connection)
func GetInstance() *Connection {
	return Instance
}
