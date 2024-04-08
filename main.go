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
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	// listen to port 8080 and use the default mux handler
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func welcome(w http.ResponseWriter, req *http.Request){
	var data registrationData
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