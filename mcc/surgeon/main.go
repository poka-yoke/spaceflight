package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/poka-yoke/spaceflight/mcc/surgeon/surgeon"
	"log"
	"sort"
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
			"'postgres://[user[:password]@]hostname[:port]/database",
	)
	showInactive := flag.Bool("i", false, "Show inactive processes")

	flag.Parse()

	db, err := GetDB(*connection)
	defer db.Close()
	if err != nil {
		log.Panic(err.Error())
	}
	ps, err := surgeon.GetProcessLists(db)
	if err != nil {
		fmt.Print("In surgeon.GetProcessLists: ")
		log.Panic(err.Error())
	}
	fmt.Printf(" %5s %10s %19s %17s %13s %9s %s \n",
		"PID", "USER", "DB", "CLIENT", "DURATION", "WAITING", "QUERY")
	sort.Sort(ps)
	for _, process := range ps {
		if *showInactive || process.IsActive() {
			fmt.Printf(
				" %5d %12s %20s %17s %10.2f %9s %s \n",
				process.Pid,
				process.Usename,
				process.Datname,
				process.Client.String(),
				process.Duration.Duration,
				process.Waiting,
				process.Query,
			)
		}
	}
}
