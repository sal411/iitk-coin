package routes

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sal411/iitk-coin/controllers"
)

func PrintABC() {
	fmt.Println("a b c")
}

func Handlers() *mux.Router {

	r := mux.NewRouter().StrictSlash(true)
	r.Use(CommonMiddleware)

	r.HandleFunc("/", controllers.CheckAPI).Methods("GET")
	r.HandleFunc("/api", controllers.CheckAPI).Methods("GET")
	r.HandleFunc("/signup", controllers.CreateUser).Methods("POST")
	r.HandleFunc("/login", controllers.Login).Methods("POST")
	r.HandleFunc("/secretpage", controllers.SecretPage).Methods("GET")

	return r

}

func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}
