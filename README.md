# resty

resty is a wrapper over [[Gorilla Mux|http://www.gorillatoolkit.org/pkg/mux]] which provides mapping for REST resources.

Example usage:

```go
r := resty.NewRouter()	// Create a new resty router
r.EnableDebug() 	// This will log incoming requests to stdout

r.Resource([]string{"users"}, controller.User{}).Only().	// Only generate 'create' method: PUT /users
	Collection("login", []shttp.Method{shttp.Post}).	// Add POST '/users/login' 
	Collection("logout", []shttp.Method{shttp.Post}).	// Add POST '/users/logout'
	Collection("register", []shttp.Method{shttp.Post})	// Add POST '/users/register'

r.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) { // Handle root path
	app := reactor.New("HelloWorld")		// Create a new single-page React app called 'HelloWorld'
	app.MapJavascriptFolder("public/js", "js")	// Load all Javascript files in /public/js using <script> tags
	app.MapCssFolder("public/css", "css")		// Load all CSS files in /public/css using <link> tags

	io.WriteString(response, app.Html().String()) 	// Generate HTML for the single-page app
})

r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/"))) // Serve all files in /public

r.PrintRoutes(os.Stdout)		// This will displasy all registered REST routes
r.HandleRoot()				// This will call http.HandleRoot to register the router
http.ListenAndServe(":5000", nil)	// Start listening on port 5000
```
