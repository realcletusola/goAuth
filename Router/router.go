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
	apiRouter.Use(middleware.Logger)
	apiRouter.Use(middleware.AllowContentType("application/json"))
	apiRouter.Use(middleware.CleanPath)
	apiRouter.Use(middleware.AllowContentEncoding("deflate","gzip"))

	// define routes 
	router.Post("/signup", handler.UserRegistrationHandler) 
	router.Post("/signin", handler.UserLoginHandler)
	// apply middleware to signout route
	router.With(middleware.BlacklistMiddleware).Post("/signout", handler.LogoutHandler)


}