package main

import (
	"html/template"
	"io"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type textData struct {
	Name string
	Text string
	User User
}

var tpl *template.Template

// init 
func init(){
	tpl = template.Must(tpl.ParseGlob("templates/*"))
}

func main (){
	http.HandleFunc("/", welcome)
	http.HandleFunc("/login", login)
	http.HandleFunc("/sign-up", signUp)
	http.HandleFunc("/authors", authors)

	http.HandleFunc("/logout", logout)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	// listen to port 8080 and use the default mux handler
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func login(w http.ResponseWriter, req *http.Request){
	if IsAlreadyLoggedIn(req){
		http.Redirect(w, req, "/welcome", http.StatusSeeOther)
	}
	
	if req.Method == http.MethodPost {
		// parse the file
		u, ok := udb[req.FormValue("email")]
		if !ok {
			http.Error(w, "user's email and password do not match", 403)
			return
		}
		// TODO use bcrypt to decode password
		if err := bcrypt.CompareHashAndPassword(u.Password, []byte(req.FormValue("password"))); err != nil{
			http.Error(w, "user's email and password do not match", 403)
			return
		}
		c := makeSessionCookie()
		sdb[c.Value] = u.Email
		http.SetCookie(w, c)
		http.Redirect(w, req, "/welcome", http.StatusSeeOther)
	}
	tpl.ExecuteTemplate(w, "login.gohtml", nil)
}

func signUp(w http.ResponseWriter, req *http.Request){
	if IsAlreadyLoggedIn(req){
		http.Redirect(w, req, "/welcome", http.StatusSeeOther)
	}

	// TODO parse form data
	if req.Method == http.MethodPost {
		// parse the form value
		// create new user
		// TODO use bcrypt to encode password
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
		c.MaxAge = 15
		// open session for new user
		sdb[c.Value] = nu.Email
		http.SetCookie(w, c)
		http.Redirect(w, req, "/welcome", http.StatusSeeOther)
	}
	tpl.ExecuteTemplate(w, "signup.gohtml", nil)
}

func logout(w http.ResponseWriter, req *http.Request){
	c, _ := req.Cookie("session")

	// close session
	delete(sdb, c.Value)

	//remove cookie
	c.MaxAge = -1
	c.Value = ""
	http.SetCookie(w, c)

	http.Redirect(w, req, "/login", http.StatusSeeOther)
}

func welcome(w http.ResponseWriter, req *http.Request){
	if !IsAlreadyLoggedIn(req){
		http.Redirect(w, req, "/login", http.StatusSeeOther)
	}

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

func authors(w http.ResponseWriter, req *http.Request){
	if !IsAlreadyLoggedIn(req) || !IsAdmin(req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
	}

	tpl.ExecuteTemplate(w, "authors.gohtml", udb)
}


