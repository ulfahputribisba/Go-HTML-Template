package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
)

const STATIC_URL string = "/assets/"
const STATIC_ROOT string = "assets/"

type Context struct {
	Title  string
	Static string
}

type students struct {
	Idstudent   int
	Nim         sql.NullString
	Name        sql.NullString
	Dateofbirth sql.NullString
	City        sql.NullString
	Major       sql.NullString
	Year_in     sql.NullString
}

func konekKeDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./celerates_unv.db")

	if err != nil {
		panic(err.Error)
	}

	return db
}

func Home(w http.ResponseWriter, req *http.Request) {
	context := Context{Title: "Welcome!"}
	render(w, "index", context)
}

func About(w http.ResponseWriter, req *http.Request) {
	context := Context{Title: "About"}
	render(w, "about", context)
}

func render(w http.ResponseWriter, tmpl string, context Context) {
	context.Static = STATIC_URL
	tmpl_list := []string{"templates/base.html",
		"templates/footer.html",
		"templates/header.html",
		"templates/sidebar.html",
		"templates/topassets.html",
		"templates/bottomassets.html",
		fmt.Sprintf("templates/%s.html", tmpl)}
	t, err := template.ParseFiles(tmpl_list...)
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, context)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

func StaticHandler(w http.ResponseWriter, req *http.Request) {
	static_file := req.URL.Path[len(STATIC_URL):]
	if len(static_file) != 0 {
		f, err := http.Dir(STATIC_ROOT).Open(static_file)
		if err == nil {
			content := io.ReadSeeker(f)
			http.ServeContent(w, req, static_file, time.Now(), content)
			return
		}
	}
	http.NotFound(w, req)
}

func checkErr(err error, args ...string) {
	if err != nil {
		fmt.Println("Error Nya : ")
		fmt.Println("%q: %s", err, args)
	}
}

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/about/", About)
	http.HandleFunc(STATIC_URL, StaticHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
