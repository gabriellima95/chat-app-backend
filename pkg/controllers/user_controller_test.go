package controllers

import (
	"bytes"
	"encoding/json"
	"msn/pkg/models"
	"msn/storage"
	"msn/storage/sqlite"
	"msn/websocket"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func TestUserController(t *testing.T) {
	db := sqlite.SetupDatabase()
	userRepository := storage.NewUserRepository(db)
	notifierMock := &websocket.NotifierMock{}
	userController := NewUserController(userRepository, notifierMock)

	t.Run("case=signup-must-save-new-user", func(t *testing.T) {
		sqlite.DB.Exec("DELETE FROM users")
		username := "abc"
		password := "def"
		jsonMap := map[string]string{"username": username, "password": password}
		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(jsonMap)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(http.MethodPost, "/users", &b)
		w := httptest.NewRecorder()

		userController.Signup(w, req)

		// responseBody := map[string]interface{}{}
		if w.Code != 200 {
			t.Errorf("Error creating user: response status code is not 200")
		}
		// if err := json.Unmarshal([]byte(w.Body.String()), &responseBody); err != nil {
		// 	t.Errorf("Response body is not valid json")
		// }
		// if len(fmt.Sprint(responseBody["id"])) != 0 {
		// 	t.Errorf("ID field should be populated")
		// }
		// if fmt.Sprint(responseBody["username"]) != username {
		// 	t.Errorf("Username field should be %v", username)
		// }
	})

	// t.Run("case=signup-must-return-with-valid-token", func(t *testing.T) {
	// 	sqlite.DB.Exec("DELETE FROM users")
	// 	username := "abc"
	// 	password := "def"
	// 	jsonMap := map[string]string{"username": username, "password": password}
	// 	var b bytes.Buffer
	// 	err := json.NewEncoder(&b).Encode(jsonMap)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	req := httptest.NewRequest(http.MethodPost, "/users", &b)
	// 	w := httptest.NewRecorder()

	// 	userController.Signup(w, req)

	// 	responseBody := map[string]interface{}{}
	// 	if err := json.Unmarshal([]byte(w.Body.String()), &responseBody); err != nil {
	// 		t.Errorf("Response body is not valid json: %v", err.Error())
	// 	}
	// 	tokenStr := responseBody["token"].(string)
	// 	claims := &Claims{}
	// 	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
	// 		return jwtKey, nil
	// 	})
	// 	if err == jwt.ErrSignatureInvalid {
	// 		t.Errorf("Token should have correct signature")
	// 	}
	// 	if err != nil {
	// 		t.Errorf("Token should be parsed correctly")
	// 	}
	// 	if !token.Valid {
	// 		t.Errorf("Token should be valid")
	// 	}
	// })

	t.Run("case=signup-must-not-save-new-user-with-empty-username", func(t *testing.T) {
		sqlite.DB.Exec("DELETE FROM users")
		var b bytes.Buffer
		user := models.User{
			Password: "abc",
		}
		err := json.NewEncoder(&b).Encode(user)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(http.MethodPost, "/users", &b)
		w := httptest.NewRecorder()

		userController.Signup(w, req)

		if w.Code != 400 {
			t.Errorf("Status code should be 400")
		}
	})

	t.Run("case=signup-must-not-save-new-user-with-empty-password", func(t *testing.T) {
		sqlite.DB.Exec("DELETE FROM users")
		var b bytes.Buffer
		user := models.User{
			Username: "abc",
		}
		err := json.NewEncoder(&b).Encode(user)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(http.MethodPost, "/users", &b)
		w := httptest.NewRecorder()

		userController.Signup(w, req)

		if w.Code != 400 {
			t.Errorf("Status code should be 400")

		}
	})

	t.Run("case=signup-must-save-new-user-with-hashed-password", func(t *testing.T) {
		sqlite.DB.Exec("DELETE FROM users")
		username := "abc"
		password := "def"
		jsonMap := map[string]string{"username": username, "password": password}
		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(jsonMap)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(http.MethodPost, "/users", &b)
		w := httptest.NewRecorder()

		userController.Signup(w, req)

		savedUser, err := userRepository.GetByUsername(username)
		if err != nil {
			t.Errorf("Error fetching user: %v", err.Error())
		}
		err = bcrypt.CompareHashAndPassword([]byte(savedUser.Password), []byte(password))
		if err != nil {
			t.Errorf("Saved password should be hashed: %v", err.Error())
		}
	})

	t.Run("case=login-must-generate-token-for-valid-user", func(t *testing.T) {
		sqlite.DB.Exec("DELETE FROM users")
		username := "abc"
		password := "def"
		user := &models.User{
			ID:       uuid.New(),
			Username: username,
			Password: hashAndSalt(password),
		}
		err := userRepository.Create(user)
		if err != nil {
			t.Errorf("Error saving user: %v", err)
		}
		jsonMap := map[string]string{"username": username, "password": password}
		var b bytes.Buffer
		err = json.NewEncoder(&b).Encode(jsonMap)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(http.MethodPost, "/login", &b)
		w := httptest.NewRecorder()

		userController.Login(w, req)

		responseBody := map[string]interface{}{}
		if w.Code != 200 {
			t.Errorf("Error getting token: response status code is not 200")
		}
		if err := json.Unmarshal([]byte(w.Body.String()), &responseBody); err != nil {
			t.Errorf("Response body is not valid json: %v", err.Error())
		}
		tokenStr := responseBody["token"].(string)
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err == jwt.ErrSignatureInvalid {
			t.Errorf("Token should have correct signature")
		}
		if err != nil {
			t.Errorf("Token should be parsed correctly")
		}
		if !token.Valid {
			t.Errorf("Token should be valid")
		}
	})

	t.Run("case=login-must-return-unauthorized-if-passwords-are-different", func(t *testing.T) {
		sqlite.DB.Exec("DELETE FROM users")
		username := "abc"
		user := &models.User{
			ID:       uuid.New(),
			Username: "abc",
			Password: "def",
		}
		err := userRepository.Create(user)
		if err != nil {
			t.Errorf("Error saving user: %v", err)
		}
		jsonMap := map[string]string{"username": username, "password": "xyz"}
		var b bytes.Buffer
		err = json.NewEncoder(&b).Encode(jsonMap)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(http.MethodPost, "/login", &b)
		w := httptest.NewRecorder()

		userController.Login(w, req)

		if w.Code != 401 {
			t.Errorf("Response status code should be 401")
		}
	})

	t.Run("case=login-must-salt-and-hash-password-before-comparing", func(t *testing.T) {
		sqlite.DB.Exec("DELETE FROM users")
		username := "abc"
		password := "def"
		user := &models.User{
			ID:       uuid.New(),
			Username: username,
			Password: hashAndSalt(password),
		}
		err := userRepository.Create(user)
		if err != nil {
			t.Errorf("Error saving user: %v", err)
		}
		jsonMap := map[string]string{"username": username, "password": password}
		var b bytes.Buffer
		err = json.NewEncoder(&b).Encode(jsonMap)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(http.MethodPost, "/login", &b)
		w := httptest.NewRecorder()

		userController.Login(w, req)

		if w.Code != 200 {
			t.Errorf("Error getting token: response status code is not 200")
		}
	})

	t.Run("case=login-must-return-bad-request-when-username-is-empty", func(t *testing.T) {
		sqlite.DB.Exec("DELETE FROM users")
		var b bytes.Buffer
		user := models.User{
			Password: "abc",
		}
		err := json.NewEncoder(&b).Encode(user)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(http.MethodPost, "/login", &b)
		w := httptest.NewRecorder()

		userController.Login(w, req)

		if w.Code != 400 {
			t.Errorf("Status code should be 400")
		}
	})

	t.Run("case=login-must-return-bad-request-when-password-is-empty", func(t *testing.T) {
		sqlite.DB.Exec("DELETE FROM users")
		var b bytes.Buffer
		user := models.User{
			Username: "abc",
		}
		err := json.NewEncoder(&b).Encode(user)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(http.MethodPost, "/login", &b)
		w := httptest.NewRecorder()

		userController.Login(w, req)

		if w.Code != 400 {
			t.Errorf("Status code should be 400")
		}
	})
}
