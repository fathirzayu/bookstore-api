package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"bookstore-api/config"
	"bookstore-api/routes"
)

func main() {
	config.InitDB()
	defer config.DB.Close()

	r := gin.Default()
	routes.Register(r)

	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	log.Println("listening on :" + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
