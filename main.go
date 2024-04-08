package main

import (
	"html/template"
	"log"
	"net/http"
)

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
	tpl.ExecuteTemplate(w, "welcome.gohtml", nil)
}