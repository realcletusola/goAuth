package validator  

import (
	"encoding/json" 
	"net/http"
	"strings"
	"regexp"

	"github.com/cletushunsu/goAuth/Database"
)

// custom error message function 
func errorMsg(w http.ResponseWriter, field string, message string) {
	error := Error{Field: field, Message: message}
	w.Header().Set("Content-Type", "application/json") // set http header 
	w.WriteHeader(http.StatusBadRequest) // set http status code 
	json.NewEncoder(w).Encode(error) // write out error message in json 
}

// username validation 
func isValidUsername(username string) (bool, string) {
	var count int // declare count variale to check if username already exist

	if len(username) < 3 || len(username) > 30 { // check username length
		return false, "Username length must be between 3 and 35 characters"
	}
	if len(strings.TrimSpace(username)) == 0 { // use trimspace to remove white space and check if username field is empty
		return false, "Username cannot be empty"
	}
	// Define a regular expression pattern that matches any of the specified characters
	regex := regexp.MustCompile(`[!@#$%^&*()_+={}\[\]|\\:;"'<>,.?/]`)
	if regex.MatchString(username) { // if username contain special character
		return false, "Username cannot contain special characters"
	}
	// query database for username if it already exists 
	err := database.db.QueryRow("SELECT COUNT(*) FROM users WHERE email=$1", username).Scan(&count)
	if err != nil {
		log.Println(err)	
	}
	if count > 0 {
		return false, "Username already exists"
	}
	
	return true,"" // username is valid 
}

// email validation  
func isValidEmail(email string) (bool, string) {
	var count int // declare count variale to check if email already exist

	if len(email) < 4 || len(email) > 160 { // check email length
		return false, "Email length must be between 4 to 150 characters"
	}
	if len(strings.TrimSpace(email)) == 0 { // use trimspace to remove white space and check if email field is empty
		return false, "Email cannot be empty"
	}
	// Regular expression for basic email validation
    regex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !regex.MatchString(email){ // use regex to check if email is in a valid format
		return false, "Invalid Email"
	}
	// query database for email if it already exists 
	err := database.db.QueryRow("SELECT COUNT(*) FROM users WHERE email=$1", email).Scan(&count)
	if err != nil {
		log.Println(err)
		
	}
	if count > 0 { // if email already exist
		return false, "Email already exists"
	}

	return true,"" // email is valid 

}

// password validator 
func isValidPassword(password string, password2 string) (bool, string) {
	if len(password) < 8 || or len(password2) < 8 { // check password length 
		return false, "Password must be at least 8 characters"
	} 
	if len(strings.TrimSpace(password)) == 0 || len(strings.TrimSpace(password2)) == 0 { // check if password field is empty
		return false, "Password cannot be empty"
	}
	if password != password2 {
		return false, "Both passwords must match"
	}
	// regular expression for checking password strength
	regex := regexp.MustCompile(`^(?=.*[A-Z])(?=.*[a-z])(?=.*\d)(?=.*[^A-Za-z0-9]).+$`)
	if !regex.MatchString(password) || !regex.MatchString(password2) { // check password strength
		return false, "Password must contain at least one uppercase, lowercase, digit, and special character"
	}
	return true,"" // password is valid 
}