package handler 

import (
	"encoding/json"
	"net/http"
	"strconv"
	"log"

	"github.com/cletushunsu/goAuth/Database"
	"github.com/cletushunsu/goAuth/Validator"
	"github.com/go-chi/chi/v5"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// jwt claims 
type JWTClaims struct {
	Username   string   `json:"username"`
	IsAdmin    bool     `json:"isAdmin"`
	jwt.StandardClaims	
}

// jwt secret key  
var JWTSecret = []byte("secretkey")


// user registration handler 
func UserRegistrationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// user registration  request struct from database.go 
		var request database.UserRegistrationRequest
		err := json.NewDecoder(r.Body).Decode(&request) // decode json request body into user registration struct
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return 
		}

		// validate username using the isValidUsername function in validator.go
		isUsernameValid, errMsg := validator.isValidUsername(request.Username)
		if !isUsernameValid {
			validator.errorMsg(w, "username", errMsg)  // returning error msg using the errorMsg function in validator.go
			return
		}

		// validate email 
		isEmailValid, errMsg := validator.isValidEmail(request.Email)
		if !isEmailValid {
			validator.errorMsg(w, "email", errMsg)
			return
		}

		// validate password 
		isPasswordValid, errMsg := validator.isValidPassword(request.Password, request.Password2)
		if !isPasswordValid {
			validator.errorMsg(w, "password". errMsg)
			return
		}

		// hash user's password before saving it to database 
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password2), bcrypt.DefaultCost)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "An error occured, please try again later", http.StatusInternalServerError)
			return 
		}

		// save user to database 
		_, err = database.db.Exec("INSERT INTO users (username, email, password, is_admin, is_active) VALUES ($1, $2, $3, $4)",
		request.Username, request.Email, string(hashedPassword), request.IsAdmin, request.IsActive)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Unable to create user, please try again later", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated) // return 201 created status if registration is successful
	}
}


// login handler  
func UserLoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request database.UserLoginRequest 
		err := json.NewDecoder(r.Body).Decode(&request) // decode json body into userLogin struct
		if err != nil {
			http.Error(w, "Invalid data, please check all fields and try again", http.StatusBadRequest)
			return
		}

		var dbUser database.User 
		// query user data from the database based on provided username or email 
		row := database.db.QueryRow()
	}
}