package main

import (
	"encoding/json"
	"githubSearch/gitAPISearch"
	"html/template"
	"log"
	"net/http"
)

type getSearchResult struct {
	/*
	 *Structure that will be encoded and sent to the client as JSON.
	 */
	Repos []gitAPISearch.GitRepos
}

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
		keys, errQuery := r.URL.Query()["search"]
		var key = ""
		if !errQuery || len(keys[0]) < 1 {
			log.Println("Url Param 'search' is missing or invalid")
			return
		}
		key = keys[0]
		var sr = getSearchResult{}
		var totalRepo = 20
		sr.Repos = gitAPISearch.GetSearch(key, totalRepo)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sr)
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
