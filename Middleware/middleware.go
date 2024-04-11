package middleware

import(
	"net/http"

	"github.com/cletushunsu/goAuth/Database"
)



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