package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"msn/pkg/controllers"
	"msn/storage"
	"msn/storage/postgres"

	ws "msn/websocket"
)

func handleRequests(userController controllers.UserController, chatController controllers.ChatController, messageController controllers.MessageController, dataController controllers.DataController) {

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/app-info", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("OK")
		json.NewEncoder(w).Encode("OK")
	}).Methods("GET")
	myRouter.HandleFunc("/signup", userController.Signup).Methods("POST")
	myRouter.HandleFunc("/login", userController.Login).Methods("POST")
	myRouter.HandleFunc("/users/{user_id}/ws", userController.ConnectSocket).Methods("GET")
	myRouter.HandleFunc("/{user_id}/chats", controllers.AuthMiddleware(chatController.ListChats)).Methods("GET")
	myRouter.HandleFunc("/{user_id}/generic_chats", chatController.ListGenericChats).Methods("GET")
	myRouter.HandleFunc("/{chat_id}/messages", messageController.ListMessages).Methods("GET")
	myRouter.HandleFunc("/messages", messageController.CreateMessage).Methods("POST")
	myRouter.HandleFunc("/data", dataController.Populate).Methods("POST")
	myRouter.HandleFunc("/data", dataController.Clear).Methods("DELETE")
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	fmt.Println("Starting server on 8080")

	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(originsOk, headersOk, methodsOk)(myRouter)))
}

func Serve() {

	db := postgres.SetupDatabase()
	userRepository := storage.NewUserRepository(db)
	chatRepository := storage.NewChatRepository(db)
	genericChatRepository := storage.NewGenericChatRepository(db)
	messageRepository := storage.NewMessageRepository(db)
	socketNotifier := ws.NewSocketNotifier()

	userController := controllers.NewUserController(userRepository, socketNotifier)
	chatController := controllers.NewChatController(chatRepository, genericChatRepository)
	messageController := controllers.NewMessageController(messageRepository, chatRepository, genericChatRepository, socketNotifier)

	dataController := controllers.NewDataController(chatRepository, userRepository, messageRepository, genericChatRepository)
	handleRequests(userController, chatController, messageController, dataController)
}
