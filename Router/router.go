package router 

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/cletushunsu/goAuth/Handler"
	"github.com/cletushunsu/goAuth/Middleware"
)

// router initiation 
func NewRouter() http.Handler {
	// define router 
	router := chi.NewRouter()

	// set middlewares
	router.Use(middleware.Logger)
	router.Use(middleware.AllowContentType("application/json"))
	router.Use(middleware.CleanPath)
	router.Use(middleware.AllowContentEncoding("deflate","gzip"))

	// define routes 
	router.Post("/signup", handler.UserRegistrationHandler) 
	router.Post("/signin", handler.UserLoginHandler)
	// apply middleware to signout route
	router.With(auth_middleware.BlacklistMiddleware).Post("/signout", handler.LogoutHandler)

	return router
}
