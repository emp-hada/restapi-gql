package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/emp/restapi-gql/data"

	"github.com/emp/restapi-gql/model"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/emp/restapi-gql/graph"
	"github.com/emp/restapi-gql/graph/generated"
	"github.com/gorilla/mux"
)

var mySigningKey = []byte("key")

// Get all books.
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data.Books)
}

// Get single book by id.
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, item := range data.Books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&model.Book{})
}

// Create a new book.
func createBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book model.Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(10000000)) // Mock ID - not safe
	data.Books = append(data.Books, &book)

	json.NewEncoder(w).Encode(book)
}

// Update an book by id.
func updateBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range data.Books {
		if item.ID == params["id"] {
			data.Books = append(data.Books[:index], data.Books[index+1:]...)
			var book model.Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]
			data.Books = append(data.Books, &book)
			break
		}
	}
	json.NewEncoder(w).Encode(data.Books)
}

// Delete a book by id.
func deleteBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range data.Books {
		if item.ID == params["id"] {
			data.Books = append(data.Books[:index], data.Books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(data.Books)
}

func secretTest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Secret test")
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] != nil {
			token, err := jwt.Parse(r.Header["Authorization"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return mySigningKey, nil
			})

			if err != nil {
				fmt.Fprintf(w, err.Error())
				return
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Not authorized"))
		}
	})
}

func isAuthorizedMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//fmt.Println(r.URL.Path)
			if r.URL.Path == "/gql" {
				next.ServeHTTP(w, r)
			}
			if r.Header["Authorization"] != nil {
				token, err := jwt.Parse(r.Header["Authorization"][0], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return mySigningKey, nil
				})

				if err != nil {
					fmt.Fprintf(w, err.Error())
					return
				}

				if token.Valid {
					next.ServeHTTP(w, r)
					//endpoint(w, r)
				}
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Not authorized"))
			}
		})
	}
}

func main() {
	// Init Mux Router
	r := mux.NewRouter()

	r.Use(isAuthorizedMiddleware())

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	r.Handle("/gql", playground.Handler("GraphQL playground", "/gql/query"))
	r.Handle("/gql/query", srv)

	// Mock data - @todo - implement DB
	data.Books = append(data.Books, &model.Book{
		ID: "1", Isbn: "445445", Title: "A Best Book",
		Author: &model.Author{Firstname: "First1", Lastname: "Last1"}})

	data.Books = append(data.Books, &model.Book{
		ID: "2", Isbn: "445446", Title: "A Better Book",
		Author: &model.Author{Firstname: "First2", Lastname: "Last2"}})

	r.HandleFunc("/", secretTest).Methods("GET")
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/book/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBooks).Methods("POST")
	r.HandleFunc("/api/book/{id}", updateBooks).Methods("PUT")
	r.HandleFunc("/api/book/{id}", deleteBooks).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", r))
}
