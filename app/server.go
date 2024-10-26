package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Signup(cfg *Config, w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST /api/v1/signup")

	account := Order{}
	json.NewDecoder(r.Body).Decode(&account)

	hash, _ := bcrypt.GenerateFromPassword([]byte(account.Password), 10)

	insertStr := fmt.Sprintf(`INSERT INTO accounts (email, password) VALUES ('%s', '%s')`, account.Email, string(hash))
	_, err := cfg.DB.Exec(insertStr)
	if err != nil {
		fmt.Errorf("Signup failed: (%w)", err)
		w.WriteHeader(http.StatusNotImplemented)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode("Signup failed")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": account.Email,
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, _ := token.SignedString([]byte(cfg.JwtKey))

	w.Header().Set("Jwt-Token", tokenString)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode("Signup Passed")
}

func Login(cfg *Config, w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST /api/v1/login")

	account := Order{}
	json.NewDecoder(r.Body).Decode(&account)

	queryStr := fmt.Sprintf(`SELECT password FROM accounts WHERE email = '%s'`, account.Email)
	rows, err := cfg.DB.Query(queryStr)
	if err != nil {
		fmt.Errorf("query by email failed: (%w)", err)
		w.WriteHeader(http.StatusNotImplemented)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode("query by email failed")
		return
	}
	defer rows.Close()

	var pwdFromDB string
	for rows.Next() {
		rows.Scan(&pwdFromDB)
	}

	err = bcrypt.CompareHashAndPassword([]byte(pwdFromDB), []byte(account.Password))
	if err != nil {
		fmt.Errorf("invalid email or password: (%w)", err)
		w.WriteHeader(http.StatusNotImplemented)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode("invalid email or password")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": account.Email,
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, _ := token.SignedString([]byte(cfg.JwtKey))

	w.Header().Set("Jwt-Token", tokenString)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode("Login Passed")
}

func PostOrder(cfg *Config, w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST /api/app/order")

	order := Order{}
	json.NewDecoder(r.Body).Decode(&order)

	insertStr := fmt.Sprintf(`INSERT INTO orders (user_id, amount) VALUES ((SELECT id FROM accounts WHERE email = '%s'), '%s')`, order.Email, order.Amount)
	_, err := cfg.DB.Exec(insertStr)
	if err != nil {
		fmt.Errorf("make order failed: (%w)", err)
		w.WriteHeader(http.StatusNotImplemented)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode("make order failed")
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode("Order Succeeded")
}

func CheckMyOrder(cfg *Config, w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /api/app/order")

	// order := Order{}
	// json.NewDecoder(r.Body).Decode(&order)
	params := r.URL.Query()
	email := params["email"][0]

	queryStr := fmt.Sprintf(`SELECT SUM(amount) FROM orders JOIN accounts ON orders.user_id = accounts.id WHERE accounts.email = '%s' GROUP BY accounts.id`, email)
	rows, err := cfg.DB.Query(queryStr)
	if err != nil {
		fmt.Errorf("check order failed: (%w)", err)
		w.WriteHeader(http.StatusNotImplemented)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode("check order failed")
		return
	}
	defer rows.Close()

	var amountFromDB string
	for rows.Next() {
		rows.Scan(&amountFromDB)
	}

	res := []Order{{Amount: amountFromDB}}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(res)
}
