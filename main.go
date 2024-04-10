package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
)

type registrationData struct {
	Name string
	Text string
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
	
	_, err := req.Cookie("session")
	if err != nil {
		if req.Method == http.MethodPost {
			// parse the file
			if req.FormValue("email") == "raphaelhebert18@gmail.com" && req.FormValue("password") == "1234" {
				c := &http.Cookie{
					Name: "session",
					Value: "sessionValue",
				}
				http.SetCookie(w, c)
				http.Redirect(w, req, "/welcome", http.StatusSeeOther)
			}
		}
		tpl.ExecuteTemplate(w, "login.gohtml", nil)
		return
	}
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func signUp(w http.ResponseWriter, req *http.Request){
	// TODO: check if already logged in
	if IsAlreadyLoggedIn(req){
		http.Redirect(w, req, "/welcome", http.StatusSeeOther)
	}

	// TODO parse form data
		if req.Method == http.MethodPost {
			// parse the form value
			// create new user
			nu := User{req.FormValue("firstName"), req.FormValue("lastName"), req.FormValue("email"), []byte(req.FormValue("password"))}
			udb[req.FormValue("email")] = nu
			sID := "someRandomSid"
			c := &http.Cookie{
				Name: "session",
				Value: sID,
			}
			// open session for new user
			sdb[sID] = nu.email
			http.SetCookie(w, c)
			http.Redirect(w, req, "/welcome", http.StatusSeeOther)
		}
		tpl.ExecuteTemplate(w, "signup.gohtml", nil)
}

func logout(w http.ResponseWriter, req *http.Request){
	c, _ := req.Cookie("session")
	c.MaxAge = -1
	http.SetCookie(w, c)
	http.Redirect(w, req, "/login", http.StatusSeeOther)
}

func welcome(w http.ResponseWriter, req *http.Request){
	var data registrationData

	// check cookie 
	c, err := req.Cookie("session")
	if err != nil {
		// TODO: redirect to login when auth is available
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		// create cookie
		c = &http.Cookie{
			Name: "session",
			Value: "sessionValue",
		}
		http.SetCookie(w, c)
	}

	
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

