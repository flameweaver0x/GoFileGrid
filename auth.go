package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Struct for JWT claims
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Load .env file variables
func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
}

// Generate JWT token
func generateJWT(username string) (string, error) {
	var mySigningKey = []byte(os.Getenv("SECRET_KEY"))
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Middleware to validate JWT token
func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] != nil {
			token, err := jwt.ParseWithClaims(r.Header["Authorization"][0], &Claims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("SECRET_KEY")), nil
			})

			if err == nil {
				if token.Valid {
					endpoint(w, r)
				}
			} else {
				fmt.Fprintf(w, "Not Authorized")
			}
		} else {
			fmt.Fprintf(w, "Not Authorized")
		}
	})
}

// Protected endpoint example
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Super Secret Information")
}

// Login simulation
func login(w http.ResponseWriter, r *http.Request) {
	var username, password string
	// Normally you'd get username and password from the request, validate it against your user store
	username = "admin" // This is a stub. In real scenarios, get this from request body
	password = "password" // This is a stub. Replace with real validation

	// Here, we're simulating an always successful login for demonstration
	if username == "admin" && password == "password" {
		validToken, err := generateJWT(username)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}

		fmt.Fprintf(w, validToken)
	} else {
		fmt.Fprintf(w, "Wrong Username or Password")
	}
}

// Main function - Set up routes and start the server
func main() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.Handle("/", isAuthorized(homePage))
	myRouter.HandleFunc("/login", login)
	http.ListenAndServe(":8080", myRouter)
}