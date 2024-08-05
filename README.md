<h3>implementing jwt auth in golang</h3>

<h4>Dirs</h4>  
<li><b>Database:</b> It contains the database initialization and migration function</li>
<li><b>Handler:</b>This contains all user auth http handler functions</li>
<li><b>Middleware:</b>The middleware.go file contains TokenBlackList middleware function(this function checks the gets token from header on signout request and check if it's blacklisted or not)</li>
<li><b>Router:</b>All routes</li>
<li><b>Validator:</b> This is where all field validaton takes place</li>  

<h4>Note:</h4>
<p>The JWT secret key is defined on line 24 in the handler.go file, this key can be set to anything(you can make it more secured by generating a long and strong random string).You can generate a secretkey with the below code: </p>
    
	import (
	     "crypto/rand"
	     "encoding/base64"
	     "fmt"  
        )
	
        func generateSecretKey(length int) (string, error) {
	        // Create a byte slice to hold the random bytes
	        key := make([]byte, length)
	        
	        // Read random bytes from the crypto/rand package
	        _, err := rand.Read(key)
	        if err != nil {
		        return "", err
	        }

	        // Encode the bytes to a base64 string
	        secretKey := base64.URLEncoding.EncodeToString(key)

	        return secretKey, nil
        }

<h4>Next(todo):</h4>
<p><b>Things that are yet to be done: </b></p>
<li>Unit Test</li>
<li>Crud operation on both User and Profile data</li>
<li>Dockerization</li>  
<li>Deployment</li>

