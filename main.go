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
    database.InitDB(connectionString)
	
	defer database.DB.Close()
	
	//make migrations 
	database.CreateMigrations()

	// start server 
	log.Println("server listening on port 8080")
	http.ListenAndServe(":8080", r)	
}
