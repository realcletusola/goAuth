package main 

import (
	"log"
	"net/http"

	"github.com/cletushunsu/goAuth/Router"
	"github.com/cletushunsu/goAuth/Database"
)

// main function 
func main(){
	// router instance 
	r := routes.NewRouter()

	// define database connection string 
	connectionString := "postgres://postgres:postgres@localhost/goauth?sslmode=disable"

	// database instance 
	conn, err := database.InitDb(connectionString)
	if err != nil {
		log.Panic(err)
	}

	// make migrations
	migration, err := database.createMigrations() 
	if err != nil {
		log.Println("Unable to make migrations")
		log.Panic(err)
	}

	// start server 
	log.Println("server listening on port 8080")
	http.ListenAndServe(":8080", r)	
}
