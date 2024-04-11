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
	r := router.NewRouter()

	// define database connection string 
	connectionString := "postgres://postgres:postgres@localhost/goauth?sslmode=disable"

	// database instance 
    err := database.InitDB(connectionString)
    if err != nil {
        log.Fatal("Error initializing database: ", err)
    }

	// make migrations
	err = database.CreateMigrations()
    if err != nil {
        log.Fatal("Error running database migrations: ", err)
    }

	// start server 
	log.Println("server listening on port 8080")
	http.ListenAndServe(":8080", r)	
}
