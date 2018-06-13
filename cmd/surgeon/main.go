package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"

	"github.com/Devex/spaceflight/pkg/surgeon"
)

// GetDB returns a sql.DB object to be used when running queries.
func GetDB(connection string) (db *sql.DB, err error) {
	engine := strings.Split(connection, ":")[0]

	db, err = sql.Open(engine, connection)
	err = db.Ping()
	return
}

// In returns true if the first argument is in the second.
func In(item string, list []string) bool {
	for _, element := range list {
		if element == item {
			return true
		}
	}
	return false
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
	function := func(surgeon.XODB, map[string]bool) (fmt.Stringer, error) {
		return nil, nil
	}
	switch {
	case In("ps", flag.Args()):
		function = surgeon.GetProcesses
	case In("locks", flag.Args()):
		function = surgeon.GetBlocks
	default:
		flag.Usage()
	}

	db, err := GetDB(*connection)
	defer db.Close()
	if err != nil {
		log.Panic(err.Error())
	}
	result, err := function(db, options)
	if err != nil {
		log.Panic(err.Error())
	}
	fmt.Print(result)
}
