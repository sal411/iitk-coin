package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

//....................................................
type UserData struct {
	Name     string `json:"name"`
	Rollno   string `json:"rollno"`
	Password string `json:"password"`
}

//.....................................................
type User struct {
	DB *sql.DB
}

func NewUser(db *sql.DB) *User {

	stmt, err := db.Prepare(`
			CREATE TABLE IF NOT EXISTS 
				data (rollno TEXT NOT NULL PRIMARY KEY UNIQUE, 
				name TEXT,
				password TEXT ) 
	`)
	if err != nil {
		log.Fatal(err)
	}

	stmt.Exec()

	return &User{
		DB: db,
	}
}

func (user *User) Add(userdata UserData) error {
	stmt, err := user.DB.Prepare(`
			INSERT INTO data 
				(rollno, name, password) VALUES(?, ?, ?)
	`)
	PrintError(err)

	stmt.Exec(userdata.Rollno, userdata.Name, userdata.Password)
	if err != nil {
		return err
	}
	return nil

}

//.....................................................

func PrintError(err error) {
	if err != nil {
		log.Fatal(err)
	}

}

func ConnectDB() *sql.DB {

	db, err := sql.Open("sqlite3", "./users.db")
	PrintError(err)

	return db
}

var db = ConnectDB()

//............................................................

func CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method == "POST" {
		var newuser UserData
		err := json.NewDecoder(r.Body).Decode(&newuser)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println(newuser)

		hashed_password, err := bcrypt.GenerateFromPassword([]byte(newuser.Password), bcrypt.DefaultCost)
		PrintError(err)

		item := NewUser(db)

		newUserData := UserData{
			Rollno:   newuser.Rollno,
			Name:     newuser.Name,
			Password: string(hashed_password),
		}

		err_in_write := item.Add(newUserData)
		if err_in_write != nil {
			log.Printf("Body read error, %v", err_in_write)
			w.WriteHeader(500) // Return 500 Internal Server Error.
			w.Write([]byte("Roll number must be unique "))
			return
		}
		fmt.Fprintf(w, "Congratulations! Your account has been successfully created")
		db.Close()

	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Only Post Methods are supported, please try again")

	}

}

//..................................................................

func GenerateToken(userRollNo string) (time.Time, string, error) {
	var err error
	//Creating Access Token

	errenv := godotenv.Load()
	if errenv != nil {
		log.Fatal("Error loading .env file")
	}

	secretKey := os.Getenv("secretKey")

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_roll_no"] = userRollNo
	expTime := time.Now().Add(time.Minute * 10)
	atClaims["exp"] = expTime.Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(secretKey))
	if err != nil {
		return time.Now(), "", err
	}

	return expTime, token, nil

}

func getHashedPassword(rollno string) string {
	db, err := sql.Open("sqlite3", "./users.db")
	PrintError(err)
	integerRollNo, _ := strconv.Atoi(rollno)

	sqlStatement := `SELECT password FROM data WHERE rollno= $1;`
	row := db.QueryRow(sqlStatement, integerRollNo)

	var hashedPassword string
	row.Scan(&hashedPassword)
	//fmt.Println("hey getting hashed password")
	//fmt.Println(hashedPassword)
	return (hashedPassword)

}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			w.WriteHeader(500)
			return
		}

		var newuser UserData
		err := json.NewDecoder(r.Body).Decode(&newuser)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		fmt.Println(newuser)
		//print ho rha hai

		hashedPassword := getHashedPassword(newuser.Rollno)
		//fmt.Fprintf(w, hashedPassword)
		//print ho rha hai hashed password
		if hashedPassword == "" {
			w.WriteHeader(500) // send server error
			w.Write([]byte("User does not exist, Please signup first"))
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(newuser.Password)); err != nil {
			w.WriteHeader(500) // send server error
			w.Write([]byte("wrong Password"))
			return
		}

		expirationTime, token, err := GenerateToken(newuser.Rollno)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(err.Error()))
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   token,
			Expires: expirationTime,
		})

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Password was correct!, You are logged in to localhost:8080 "))

	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Only Post Methods are supported, please try again")
	}

}

//................................................................................

func VerifyJWToken(requestedToken string) (*jwt.Token, error) {

	errenv := godotenv.Load()
	if errenv != nil {
		log.Fatal("Error in loading .env file")
	}

	secretKey := os.Getenv("secretKey")

	tokenString := requestedToken
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ExtractMetadata(userToken string) (string, error) {
	token, err := VerifyJWToken(userToken)
	if err != nil {
		return " ", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		rollNo, _ := claims["user_roll_no"].(string)

		return rollNo, err
	}
	return " ", err
}

func secretPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/secretpage" {
		w.WriteHeader(404)
		fmt.Fprint(w, "Error 404 Page not found")
		return
	}

	if r.Method == "GET" {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				// cookie is not set => return an unauthorized status
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			//  if any other type of error => return a bad request status
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tokenFromUser := c.Value
		//fmt.Fprintf(w, tokenFromUser)
		userRollNo, _ := ExtractMetadata(tokenFromUser)
		fmt.Println(userRollNo, "Hello")
		fmt.Fprint(w, "You are authorized and accessing the secret page, your roll number is "+userRollNo)

	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Only GET Methods are supported, please try again")
	}

}

//..............................................................................

func DataPrint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World, This is woking")

}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	http.HandleFunc("/signup", CreateUser)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/", DataPrint)
	http.HandleFunc("/secretpage", secretPage)

	log.Printf("Server is up on port '%s'", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
