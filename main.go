package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	mux := http.NewServeMux()

	getTodo_V := http.HandlerFunc(getTodo)
	// getTodoMW_V := getTodoMW(getTodo_V)
	// parentGetTodoMW_V := parentGetTodoMW(getTodoMW_V)

	mux.HandleFunc("GET /", test)
	// mux.HandleFunc("GET /get", parentGetTodoMW_V.ServeHTTP)
	mux.Handle("GET /get", parentGetTodoMW(getTodoMW(getTodo_V)))

	//run server
	err := http.ListenAndServe(":8000", mux)
	log.Fatal((err))
}

// to get
func getTodo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Todo is getted")
}

// to test
func test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "API WORKING Middleware")
}

// testing middleware for gettodo
func getTodoMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("MW WORKING")
		next.ServeHTTP(w, r)
	})
}

// parent MW of getTodoMW
func parentGetTodoMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("Parrent MW is called")
		next.ServeHTTP(w, r)
		// a:=A{
		// 	user:"sss",
		// 	id:1,
		// }
	})
}

// generate JWT TOKEN
func generateJWT(payload string) (string, error) {
	secret := []byte("this is secret of jwt practice")
	claims := jwt.MapClaims{
		"data": payload,
		"exp":  time.Now().Add(time.Minute * 3),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(secret)
}

type A interface{}
