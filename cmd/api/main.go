package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/it0dan/golang-studies/internal/server"
	"github.com/it0dan/golang-studies/pkg/db/postgres"
)

func main() {
	psqlDB, err := postgres.NewDBConnect()
	if err != nil {
		log.Fatal("Something is wrong. Connection not established: ", err)
	} else {
		log.Println("Postgres connection established")
	}

	defer psqlDB.Close()

	godotenv.Load()
	serverAdd := os.Getenv("SERVER_ADDR")
	println("Server: ", serverAdd)

	router := gin.Default()
	router.GET("/customers", server.GetCustomers)
	router.GET("/customers/:id", server.GetCustomersById)
	router.POST("/customers", server.CreateCustomer)
	router.PUT("/customers/:id", server.UpdateCustomer)
	router.DELETE("/customers/:id", server.RemoveCustomersById)

	router.Run(serverAdd)
}
