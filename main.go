package main

import (
	"encoding/json"
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
	mux.Handle("GET /get", JWTMW(parentGetTodoMW(getTodoMW(getTodo_V))))
	mux.HandleFunc("POST /login", login)

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

var secret = []byte("this is secret of jwt practice")

// generate JWT TOKEN
func generateJWT(payload string) (string, error) {
	claims := jwt.MapClaims{
		"data": payload,
		"exp":  time.Now().Add(time.Minute * 3).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString(secret)
	if err != nil {
		print("errrrrrrrrrr")
	}
	return s, err

}

// dummy login
func login(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		var loginModel_V loginModel
		json.NewDecoder(r.Body).Decode(&loginModel_V)
		print(loginModel_V.Username)
		token, err := generateJWT(loginModel_V.Username)
		if err != nil {
			print("errrrrrrrrrr22")
		}
		var loginRes_V loginRes
		loginRes_V.Jwt_token = "Bearrer " + token
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(&loginRes_V)
	}
}

type loginModel struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type loginRes struct {
	Jwt_token string `json:"jwt_token"`
}

// JWT Middleware
func JWTMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" || len(header) < 8 {
			http.Error(w, "No Token Found", http.StatusUnauthorized)
			return
		}
		jwt_token := header[len("bearer "):]
		var payload_data string
		JWTParser(jwt_token, w, &payload_data)
		println("payload = ", payload_data)
		if payload_data != "" {
			next.ServeHTTP(w, r)
		} else {
			print("wwwwwwwwwww")
		}
	})
}

// JWT parser
func JWTParser(jwt_token string, w http.ResponseWriter, payload_data *string) {
	token, err := jwt.Parse(jwt_token, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	claim := token.Claims.(jwt.MapClaims)
	*payload_data = claim["data"].(string)
}
