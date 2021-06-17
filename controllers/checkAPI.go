package controllers

import (
	"fmt"
	"net/http"
)

func CheckAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World, This is woking")

}
