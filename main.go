package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"

	"./models"
)

var articlesCollection *mgo.Collection

func main() {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	articlesCollection = session.DB("testblog").C("posts")

	router := mux.NewRouter().StrictSlash(true)

	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/articles", indexHandler).Methods("GET")

	router.HandleFunc("/articles/new", newHandler).Methods("GET")
	router.HandleFunc("/create", postCreateHandler).Methods("POST")

	router.HandleFunc("/articles/{id}", showHandler).Methods("GET")

	router.HandleFunc("/edit/{id}", editHandler).Methods("GET")
	router.HandleFunc("/update", postUpdateHandler).Methods("PUT")

	router.HandleFunc("/delete/{id}", deleteHandler)

	log.Fatal(http.ListenAndServe(":3000", router))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html",
		"templates/header.html",
		"templates/footer.html")

	if err != nil {
		panic(err)
	}

	//article := models.NewArticle("cat", "tit", "subt", "cont")
	//articlesCollection.Insert(article)

	articles := []models.Article{}
	err = articlesCollection.Find(nil).All(&articles)
	if err != nil {
		panic(err)
	}

	t.ExecuteTemplate(w, "index", articles)
}

func showHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/show.html",
		"templates/header.html",
		"templates/footer.html")

	if err != nil {
		panic(err)
	}

	params := mux.Vars(r)
	id := params["id"]
	art := models.Article{}

	articlesCollection.FindId(id).One(&art)

	t.ExecuteTemplate(w, "show", art)
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/new.html",
		"templates/header.html",
		"templates/footer.html")

	if err != nil {
		panic(err)
	}

	t.ExecuteTemplate(w, "new", nil)
}

func postCreateHandler(w http.ResponseWriter, r *http.Request) {
	category := r.FormValue("category")
	title := r.FormValue("title")
	content := r.FormValue("content")

	article := models.NewArticle(category, title, content)
	err := articlesCollection.Insert(article)

	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/", 302)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	articlesCollection.RemoveId(id)

	http.Redirect(w, r, "/", 302)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/edit.html",
		"templates/header.html",
		"templates/footer.html")

	if err != nil {
		panic(err)
	}

	params := mux.Vars(r)
	id := params["id"]
	article := models.Article{}
	articlesCollection.FindId(id).One(&article)

	t.ExecuteTemplate(w, "edit", article)
}

func postUpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	category := r.FormValue("category")
	title := r.FormValue("title")
	content := r.FormValue("content")

	article := models.NewArticle(category, title, content)

	articlesCollection.UpdateId(id, article)

	http.Redirect(w, r, "/", 302)
}
