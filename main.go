package main

import (
	database "bab6/databases"
	"bab6/routers"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "nickytirta54"
	dbName   = "mcs_bab_6"
)

var (
	DB  *sql.DB
	err error
)

func main() {
	var PORT = ":50000"

	psqlInfo := fmt.Sprintf(
		`host=%s port=%d user=%s password=%s dbname=%s sslmode=disable`,
		host, port, user, password, dbName,
	)
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Error Open DB: ", err)
	}

	database.DBMigrate(DB)

	defer DB.Close()

	routers.StartServer().Run(PORT)
	fmt.Printf("Connect Successfully")
}
