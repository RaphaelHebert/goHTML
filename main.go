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
	http.HandleFunc("/logout", logout)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	// listen to port 8080 and use the default mux handler
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func login(w http.ResponseWriter, req *http.Request){
	_, err := req.Cookie("session")
	if err != nil {
		if req.Method == http.MethodPost {
			// parse the file
			if req.FormValue("name") == "Raphael" && req.FormValue("password") == "1234" {
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

