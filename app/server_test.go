package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPostOrder(t *testing.T) {
	db, mock, _ := sqlmock.New()

	// Create a mock Config struct
	// type mockConfig struct {
	// 	DB *sql.DB
	// }
	// type MockDB struct {
	// 	ExecFunc func(query string, args ...interface{}) (result sql.Result, err error)
	// }
	// mockDB := &MockDB{
	// 	ExecFunc: func(query string, args ...interface{}) (result sql.Result, err error) {
	// 		return nil, nil
	// 	},
	// }
	mockConfig := &Config{
		DB: db,
	}

	// Create a mock request
	order := Order{
		Email:  "test@example.com",
		Amount: "100.00",
	}
	jsonOrder, _ := json.Marshal(order)
	req, _ := http.NewRequest("POST", "/api/app/order", bytes.NewReader(jsonOrder))

	// Create a mock response writer
	var buf bytes.Buffer
	w := &httptest.ResponseRecorder{
		Body: &buf,
	}

	// Call the function with the mock Config and request
	PostOrder(mockConfig, w, req)

	// Check that the response status code is 200
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	// Check that the response body is "Order Succeeded"
	if w.Body.String() != "Order Succeeded" {
		t.Errorf("Expected response body 'Order Succeeded', got '%s'", w.Body.String())
	}

	mock.ExpectBegin()
	expectedQuery := `INSERT INTO orders (user_id, amount) VALUES ((SELECT id FROM accounts WHERE email = 'test@example.com'), '100.00')`
	mock.ExpectExec(expectedQuery).WithoutArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// // Check that the database query was executed correctly
	// expectedQuery := `INSERT INTO orders (user_id, amount) VALUES ((SELECT id FROM accounts WHERE email = 'test@example.com'), '100.00')`
	// if mockDB.ExecFunc == nil {
	// 	t.Error("Expected ExecFunc to be set, but it was not")
	// } else if query := mockDB.ExecFunc[0]; query != expectedQuery {
	// 	t.Errorf("Expected ExecFunc to be called with query '%s', but it was called with '%s'", expectedQuery, query)
	// }
}
