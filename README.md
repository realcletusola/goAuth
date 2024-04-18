<h3>implementing jwt auth in golang</h3>

<h4>Dirs</h4>
<li><b>Database:</b> It contains the database init function and the migration function</li>
<li><b>Handler:</b>This contains all user auth http handler functions</li>
<li><b>Middleware:</b>The middleware.go file contains TokenBlackList middleware function(this function checks the gets token from header on signout request and check if it's blacklisted or not)</li>
<li><b>Router:</b>All routes</li>
<li><b>Validator:</b> This is where all field validaton takes place</li>  

<h4>Note:</h4>
<p>The JWT secret key is defined on line 24 in the handler.go file, this key can be set to anything(you can make it more secured by generating a long and strong random string).</p>
<p>It's good practice to keep all valuable keys in .env file</p>

<h4>Next(todo):</h4>
<p><b>Things that are yet to be done: </b></p>
<li>Unit Test</li>
<li>Crud operation on both User and Profile data</li>
<li>Dockerization</li>
<li>Deployment</li>

