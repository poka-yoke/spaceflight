package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/poka-yoke/spaceflight/mcc/surgeon/surgeon"
	"log"
	"strings"
)

// GetDB returns a sql.DB object to be used when running queries.
func GetDB(connection string) (db *sql.DB, err error) {
	engine := strings.Split(connection, ":")[0]

	db, err = sql.Open(engine, connection)
	err = db.Ping()
	return
}

func main() {
	connection := flag.String(
		"connection",
		"postgres://localhost/postgres",
		"Connection string: "+
			"'postgres://[user[:password]@]hostname[:port]/"+
			"database",
	)
	showInactive := flag.Bool("i", false, "Show inactive processes")

	flag.Parse()
	options := map[string]bool{
		"showInactive": *showInactive,
	}

	db, err := GetDB(*connection)
	defer db.Close()
	if err != nil {
		log.Panic(err.Error())
	}
	result, err := surgeon.GetProcessLists(db, options)
	if err != nil {
		log.Panic(err.Error())
	}
	fmt.Print(result)
}
