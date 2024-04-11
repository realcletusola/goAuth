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
	IsActive    bool     `json:"isActive"`
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
			http.Error(w, "Invalid data format, check form and try again", http.StatusBadRequest)
			return 
		}

		// use a chanel to receive the result from field validation
		resultCh := make(chan bool)
		errorCh := make(chan string)

		// validate username with goroutine using the username validator func in validators.go
		go func(){
			isUsernameValid, errMsg := validator.isValidUsername(request.Username)
			resultCh <- isUsernameValid
			errorCh <- errMsg

			// wait for the result from the goroutine
			isUsernameValid := <- resultCh
			errMsg := <- errorCh

			if !isUsernameValid {
				validator.errorMsg(w, "username", errMsg)
				return 
			}
		}

		// validate email with goroutine using the email validator func in validators.go
		go func(){
			isEmailValid, errMsg := validator.isValidEmail(request.Email)
			resultCh <- isEmailValid
			errorCh <- errMsg

			// wait for the result from the goroutine
			isEmailValid := <- resultCh
			errMsg := <- errorCh

			if !isEmailValid {
				validator.errorMsg(w, "email", errMsg)
				return 
			}
		}
		
		// validate password with goroutine using the password validator func in validators.go
		go func(){
			isPasswordValid, errMsg := validator.isValidPassword(request.Password, request.Password2)
			resultCh <- isPasswordValid
			errorCh <- errMsg

			// wait for the result from the goroutine
			isPasswordValid := <- resultCh
			errMsg := <- errorCh

			if !isPasswordValid {
				validator.errorMsg(w, "password", errMsg)
				return 
			}
		}

		// hash user's password before saving it to database 
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password2), bcrypt.DefaultCost)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "An error occured, please try again later", http.StatusInternalServerError)
			return 
		}

		// save user to database 
		var userID int 
		_, err = database.db.QueryRow("INSERT INTO users (username, email, password, is_admin, is_active) VALUES ($1, $2, $3, $4) RETURNING id",
		request.Username, request.Email, string(hashedPassword), request.IsAdmin, request.IsActive).Scan(&userID)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Unable to create user, please try again later", http.StatusInternalServerError)
			return
		}

		// create profile for user  
		_, err = database.db.Exec("INSERT INTO profile (username, email, user_id) VALUES ($1, $2, $3)",
		request.Username, request.Email, userID)
		if err != nil {
			
			// ensure data consistency by deleting user data from database if profile fails to be created
			_, err := database.db.Exec("DELETE FROM users WHERE id = $1 ", userID) 
			if err != nil {
				http.Error(w, "Unable to rollback data, please try again later", http.StatusInternalServerError)
				return 
			}
			// if there is error creating profile data 
			log.Println(err.Error())
			http.Error(w, "Unable to create profile, please try again later", http.StatusInternalServerError)
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
		row := database.db.QueryRow("SELECT id, username, email, password, is_admin, is_active FROM users WHERE username = $1 OR email = $1", request.LoginId)
		// scan the retrieved row nto dbUser struct 
		err =  row.Scan(&dbUser.ID, &dbUser.Username, &dbUser.Email, &dbUser.Password, &dbUser.IsAdmin, &dbUser.IsActive)
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusBadRequest)
			return
		}

		// compare hashed password from user database with password provided by user 
		err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(request.Password))
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusBadRequest)
			return 
		}

		if !dbUser.IsActive {
			http.Error(w, "User account is not active, contact our support") // check if account is active
			return 
		} 
		// create jwt token 
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
			Username: dbUser.Username,
			IsActive: dbUser.isActive,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // token expires in 24 hours 
			}
		})

		// sign the token with jwtsecret key 
		tokenString, err := token.SignedString(JWTSecret)
		if err != nil {
			http.Error(w, "An error occured, please try again later", http.StatusInternalServerError)
			return
		}

		// encode token to json and send response  
		json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
	}


	// logout function ( the logout function will be wraped into the BlaclistMiddleware)
	func LogoutHandler(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		// add token to blacklist 
		err := addToBlacklist(token)
		if err != nil {
			http.Error(w, "An error occured, please try again later", http.StatusInternalServerError)
			return 
		}

		w.WriteHeader(http.StatusOK) // return success message if no error
	}

	// blacklist middleware is to check if token is blacklisted 
	func BlacklistMiddleware(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			// extract token from header 
			token := r.Header.Get("Authorization")

			// check if token is blaclisted using the isBlacklisted function 
			blacklist := isBlacklisted(token)
			
			// return unauthorized is token is blacklisted 
			if blacklist {
				http.Error(w, "Token is blacklisted, please login again", http.StatusUnauthorized)
				return
			}

			// call the next handler is token is not blacklisted
			next.ServeHTTP(w, r)
		}) 
	}

	// function to check the blacklist database if token is blacklisted 
	func isBlacklisted(token string) bool {
		var count int 
		row := database.db.QueryRow("SELECT COUNT(*) FROM blacklisted_token WHERE token = $1", token)
		err := row.Scan(&count)
		if err != nil {
			return false // return false if token is not found in blacklist 
		}
		return count > 0 // return true if token is found in blacklist 
	}

	// function to add token to blacklist
	func addToBlacklist(token string) error {
		// insert token into blacklist table in the database 
		_, err := database.db.Exec("INSERT INTO blacklisted_token (token) VALUES ($1)", token)
		return err
	}

}