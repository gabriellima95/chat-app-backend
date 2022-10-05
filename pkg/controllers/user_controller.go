package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"msn/pkg/models"
	"msn/storage"
	"msn/websocket"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	userRepository storage.UserRepository
	socketNotifier websocket.Notifier
}

func NewUserController(userRepository storage.UserRepository, socketNotifier websocket.Notifier) UserController {
	return UserController{
		userRepository: userRepository,
		socketNotifier: socketNotifier,
	}
}

var jwtKey = []byte("my_secret_key")

type Claims struct {
	UserID string `json:"id"`
	jwt.StandardClaims
}

func (c UserController) Signup(w http.ResponseWriter, r *http.Request) {
	var userPayload UserRequestSchema
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err := json.NewDecoder(r.Body).Decode(&userPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if userPayload.Username == "" {
		http.Error(w, "Username should not be empty", http.StatusBadRequest)
		return
	}
	if userPayload.Password == "" {
		http.Error(w, "Password should not be empty", http.StatusBadRequest)
		return
	}

	user := models.User{
		Username: userPayload.Username,
		Password: hashAndSalt(userPayload.Password),
	}

	// user.ID = uuid.New()

	err = c.userRepository.Create(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// token, err := generateToken(user.ID.String())
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// userResponse := LoginResponseSchema{
	// 	Token: token,
	// }

	json.NewEncoder(w).Encode("")
}

func (c UserController) Login(w http.ResponseWriter, r *http.Request) {
	log.Println("UserController.Login")
	var userPayload UserRequestSchema
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err := json.NewDecoder(r.Body).Decode(&userPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if userPayload.Username == "" {
		http.Error(w, "Username should not be empty", http.StatusBadRequest)
		return
	}
	if userPayload.Password == "" {
		http.Error(w, "Password should not be empty", http.StatusBadRequest)
		return
	}

	user, err := c.userRepository.GetByUsername(userPayload.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userPayload.Password))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := generateToken(user.ID.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	loginResponse := LoginResponseSchema{
		Token: token,
	}

	json.NewEncoder(w).Encode(loginResponse)
}

func (u UserController) ConnectSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UserController: ConnectSocket")
	params := mux.Vars(r)
	ID := params["user_id"]

	userID, err := uuid.Parse(ID)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "chatID should be valid uuid", http.StatusBadRequest)
		return
	}

	u.socketNotifier.AddConnection(w, r, userID)
}

func hashAndSalt(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	return string(hash)
}

func generateToken(userID string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("Error generating token: %v", err.Error())
		return "", nil
	}

	return tokenString, nil
}
