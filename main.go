package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", test)

	//run server
	err := http.ListenAndServe(":8000", mux)
	log.Fatal((err))
}

// to test
func test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "API WORKING Middleware")
}
