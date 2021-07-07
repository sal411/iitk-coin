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

	http.Handle("/", routes.Handlers())

	// a, err := database.GetUserFromRollNo("190951")
	// fmt.Println(a)

	//b, err := database.GetCoinsFromRollno("190951")

	//fmt.Println(b)

	log.Printf("Server is up on port '%s'", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
