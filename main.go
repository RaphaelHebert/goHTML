package main

// Package main provides HTTP handlers for serving various routes.
// This package is written for education purposes.
//
// All handlers in this package follow a similar pattern:
//   - They handle HTTP requests to specific routes.
//   - They retrieve user information from the request.
//   - If it's a POST request, they process uploaded text data.
//   - They execute corresponding templates with user data.
//
// Parameters:
//   - w: ResponseWriter for constructing HTTP responses.
//   - req: Request representing an HTTP request received by the server.
//
// Notes:
//   - Assumes "tpl" represents an initialized template object.
//   - Sends appropriate HTTP error responses if file parsing or reading fails.

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// TextData represents a text extracted from a text file submitted by a user.
type textData struct {
	Name string
	Text string
	User User
}

var tpl *template.Template

// ExpireTime sets the lifespan of a session.
const expireTime int = 600 // session lifespan in sec

// init
func init() {
	tpl = template.Must(tpl.ParseGlob("templates/*"))
}

// Main is the entry point of the application. 
// It sets up HTTP routes with corresponding handlers.
// It starts the HTTP server to listen for incoming requests on port 8080.
func main() {
	http.HandleFunc("/", authorized(welcome))
	http.HandleFunc("/login", guestGuard(login))
	http.HandleFunc("/sign-up", guestGuard(signUp))
	http.HandleFunc("/authors", admin(authors))

	http.HandleFunc("/logout", logout)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	// listen to port 8080 and use the default mux handler
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Login handles HTTP requests to the "/login" route. 
// If the request method is POST, it validates the user's credentials,
// It creates a session for the user, and redirects them to the "/welcome" page.
// If the credentials are invalid, it returns a forbidden error.
func login(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		// parse the file
		u, ok := udb[req.FormValue("email")]
		if !ok {
			http.Error(w, "user's email and password do not match", 403)
			return
		}
		// TODO use bcrypt to decode password
		if err := bcrypt.CompareHashAndPassword(u.Password, []byte(req.FormValue("password"))); err != nil {
			http.Error(w, "user's email and password do not match", 403)
			return
		}
		c := makeSessionCookie()
		c.MaxAge = expireTime
		sdb[c.Value] = Session{u.Email, time.Now()}
		http.SetCookie(w, c)
		http.Redirect(w, req, "/welcome", http.StatusSeeOther)
	}
	tpl.ExecuteTemplate(w, "login.gohtml", nil)
}

// SignUp handles HTTP requests to the "/signup" route. 
// It processes form data submitted via POST request to create a new user account. 
// Passwords are hashed using bcrypt for security.
// If the email provided already exists in the user database, it returns an error response.
// Upon successful signup, it creates a session for the new user and redirects them to the "/welcome" page.
func signUp(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		if _, ok := udb[req.FormValue("email")]; ok {
			http.Error(w, "This email is already linked to an account", http.StatusForbidden)
		}
		password, err := bcrypt.GenerateFromPassword([]byte(req.FormValue("password")), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "Could not proceed password", http.StatusInternalServerError)
		}
		nu := User{req.FormValue("firstName"), req.FormValue("lastName"), req.FormValue("email"), password, req.FormValue("role")}
		udb[req.FormValue("email")] = nu
		c := makeSessionCookie()
		c.MaxAge = expireTime
		// open session for new user
		sdb[c.Value] = Session{nu.Email, time.Now()}
		http.SetCookie(w, c)
		http.Redirect(w, req, "/welcome", http.StatusSeeOther)
	}
	tpl.ExecuteTemplate(w, "signup.gohtml", nil)
}

// Logout handles the HTTP requests to the "/logout" route. 
// Deletes the session associated with the user.
// Removes the session cookie.
// Redirects the user to the "/login" page.
func logout(w http.ResponseWriter, req *http.Request) {
	c, _ := req.Cookie("session")

	// close session
	delete(sdb, c.Value)

	//remove cookie
	c.MaxAge = -1
	c.Value = ""
	http.SetCookie(w, c)

	// TODO: move to CRON for production
	go CleanUpSession()

	http.Redirect(w, req, "/login", http.StatusSeeOther)
}

// Welcome handles HTTP requests to "/welcome". 
// Welcome is the home page for logged in users.
// Retrieve users information and execute welcome template with the user data.
func welcome(w http.ResponseWriter, req *http.Request) {
	u := GetUser(req)

	var data textData

	data.User = u

	if req.Method == http.MethodPost {
		// parse the file
		data.Name = req.FormValue("name")
		f, _, err := req.FormFile("textFile")
		if err != nil {
			http.Error(w, "could not find file", http.StatusInternalServerError)
		}
		c, err := io.ReadAll(f)
		if err != nil {
			http.Error(w, "could not find file", http.StatusInternalServerError)
		}
		data.Text = string(c)
	}
	tpl.ExecuteTemplate(w, "welcome.gohtml", data)
}

// Authors handles HTTP requests to the "/authors" route. 
// It renders the "authors.gohtml" template, passing the user database (udb) to be displayed.
func authors(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "authors.gohtml", udb)
}
