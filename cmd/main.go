package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kapi1023/kingGivefivecoins-Authentication/internal/config"
	_ "github.com/kapi1023/kingGivefivecoins-Authentication/internal/docs"
	"github.com/kapi1023/kingGivefivecoins-Authentication/internal/handlers"
	"github.com/kapi1023/kingGivefivecoins-Authentication/internal/middlewares"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Authentication Service API
// @version 1.0
// @description API for user authentication, OAuth, and token management
// @contact.name Developer Support
// @contact.url kapi1023
// @contact.email kacper.pietrzak533@gmail.com
// @BasePath /
func main() {
	cfg := config.Load()
	r := mux.NewRouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	r.Use(middlewares.LoggingMiddleware)
	r.Use(middlewares.CORSMiddleware)

	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/refresh", handlers.RefreshTokenHandler).Methods("POST")
	r.HandleFunc("/oauth/callback", handlers.OAuthCallbackHandler).Methods("GET")

	log.Println("Server is running on port", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
