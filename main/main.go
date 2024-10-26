package main

import (
	"fmt"
	"net/http"

	orderapp "go-pq8/app"
	authapp "go-pq8/middleware"
)

func main() {
	fmt.Println("!!! Project begins !!!")
	// config init
	cfg, err := orderapp.InitConfig()
	if err != nil {
		fmt.Errorf("init orderapp failed: (%w)", err)
	}

	http.HandleFunc("/api/v1/signup", handleSignup(cfg))
	http.HandleFunc("/api/v1/login", handleLogin(cfg))
	http.HandleFunc("/api/v1/order", authapp.Auth(cfg, handleOrder(cfg)))

	fmt.Println("!!! Server listening on port :8090")
	http.ListenAndServe(":8090", nil)
}

func handleSignup(cfg *orderapp.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// server signup method
		orderapp.Signup(cfg, w, r)
	}
}

func handleLogin(cfg *orderapp.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// server login method
		orderapp.Login(cfg, w, r)
	}
}

func handleOrder(cfg *orderapp.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method)
		if r.Method == "GET" {
			orderapp.CheckMyOrder(cfg, w, r)
		} else if r.Method == "POST" {
			orderapp.PostOrder(cfg, w, r)
		}
	}
}
