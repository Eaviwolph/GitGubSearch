package main

import (
	"githubSearch/gitAPISearch"
	"html/template"
	"log"
	"net/http"
)

/*
 *That call the function gitSearch and send the result to the client
 *if the method is GET and the url path is /search.
 *for all other case, the function send the index page
 */
var handleFunc = func(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Println(err)
	}
	if r.Method == "GET" && r.URL.Path == "/search" {
		gitAPISearch.GetSearch(w, r)
	} else {
		tpl.ExecuteTemplate(w, "index.html", nil)
	}
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("templates/index.html"))
}

/*
 *Main function that start the server and redirect source files like css of js.
 */
func main() {
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/web/", http.StripPrefix("/web", fs))
	http.HandleFunc("/", handleFunc)
	http.ListenAndServe(":8080", nil)
}
