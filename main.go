package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	getTodo_V := http.HandlerFunc(getTodo)
	getTodoMW_V := getTodoMW(getTodo_V)

	mux.HandleFunc("GET /", test)
	mux.HandleFunc("GET /get", getTodoMW_V.ServeHTTP)

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
