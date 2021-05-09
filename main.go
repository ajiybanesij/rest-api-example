package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

var mySigningKey = []byte("SECRET_KET_!123")

type Login struct {
	Username string
	Password string
}

type Token struct {
	Token   string
	Message string
}

// Book Struct(Model)
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// Init books var as a slice Book struct
var books []Book

// Handlers
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)

}
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//Get params
	params := mux.Vars(r)

	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})

}
func createBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(1000000))
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}
func updateBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	json.NewEncoder(w).Encode(books)
}
func deleteBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return mySigningKey, nil
			})
			if err != nil {
				fmt.Fprintf(w, err.Error())
			}
			if token.Valid {
				endpoint(w, r)
			}

		} else {

			fmt.Fprintf(w, "Not Authorized", http.StatusForbidden)
		}

	})
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user Login
	_ = json.NewDecoder(r.Body).Decode(&user)
	token, err := GenerateJWT(user)
	var token_struct Token
	if err != nil {
		token_struct.Message = "Error"
		token_struct.Token = ""
	}
	token_struct.Token = token
	token_struct.Message = "Success"

	json.NewEncoder(w).Encode(token_struct)
}
func GenerateJWT(user Login) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = user.Username
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something went wrong : %s", err.Error())
		return "", err
	}
	return tokenString, nil

}

func main() {
	// Init Routher
	r := mux.NewRouter()

	// Mock Data @todo - implement DB
	books = append(books, Book{ID: "1", Isbn: "12345", Title: "Book One", Author: &Author{Firstname: "John", Lastname: "Doe"}})
	books = append(books, Book{ID: "2", Isbn: "12346", Title: "Book Two", Author: &Author{Firstname: "John2", Lastname: "Doe2"}})

	// Route Handlers / Endpoints
	r.HandleFunc("/api/login", login).Methods("POST")
	r.Handle("/api/books", isAuthorized(getBooks)).Methods("GET")
	r.Handle("/api/books/{id}", isAuthorized(getBook)).Methods("GET")
	r.Handle("/api/books", isAuthorized(createBooks)).Methods("POST")
	r.Handle("/api/books/{id}", isAuthorized(updateBooks)).Methods("PUT")
	r.Handle("/api/books", isAuthorized(deleteBooks)).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", r))

}
