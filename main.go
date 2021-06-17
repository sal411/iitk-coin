package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/sal411/iitk-coin/routes"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	// routes.PrintABC()
	// fmt.Println("open to port " + port)

	http.Handle("/", routes.Handlers())

	log.Printf("Server is up on port '%s'", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
