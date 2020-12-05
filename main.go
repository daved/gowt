package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	xlog "log"
	"net/http"
	"os"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

var log = xlog.New(ioutil.Discard, "", xlog.LstdFlags)

// Tool struct
type Tool struct {
	Id       int
	Name     string
	Category string
	URL      string
	Rating   int
	Notes    string
}

func dbConn() (*sql.DB, error) {
	efmt := "dbconn: %w"

	dbDriver := "mysql"
	dbUser := os.Getenv("DATABASE_USERNAME")
	dbPass := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE_NAME")
	dbServer := os.Getenv("DATABASE_SERVER")
	dbPort := os.Getenv("DATABASE_PORT")
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbServer+":"+dbPort+")/"+dbName)
	if err != nil {
		return nil, fmt.Errorf(efmt, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf(efmt, err)
	}

	return db, nil
}

var noDBMsg = "no db conn"

type toolSvc struct {
	db *sql.DB
	t  *template.Template
}

func newToolSvc(db *sql.DB, t *template.Template) *toolSvc {
	return &toolSvc{
		db: db,
		t:  t,
	}
}

//Index handler
func (s *toolSvc) Index(w http.ResponseWriter, r *http.Request) {
	selDB, err := s.db.Query("SELECT * FROM tools ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}

	tool := Tool{}
	res := []Tool{}

	for selDB.Next() {
		var id, rating int
		var name, category, url, notes string
		err := selDB.Scan(&id, &name, &category, &url, &rating, &notes)
		if err != nil {
			panic(err.Error())
		}
		log.Println("Listing Row: Id " + string(id) + " | name " + name + " | category " + category + " | url " + url + " | rating " + string(rating) + " | notes " + notes)

		tool.Id = id
		tool.Name = name
		tool.Category = category
		tool.URL = url
		tool.Rating = rating
		tool.Notes = notes
		res = append(res, tool)
	}

	_ = s.t.ExecuteTemplate(w, "Index", res)
}

//Show handler
func (s *toolSvc) Show(w http.ResponseWriter, r *http.Request) {
	nId := r.URL.Query().Get("id")
	selDB, err := s.db.Query("SELECT * FROM tools WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	defer selDB.Close() //nolint

	tool := Tool{}

	for selDB.Next() {
		var id, rating int
		var name, category, url, notes string
		err := selDB.Scan(&id, &name, &category, &url, &rating, &notes)
		if err != nil {
			panic(err.Error())
		}

		log.Println("Showing Row: Id " + string(id) + " | name " + name + " | category " + category + " | url " + url + " | rating " + string(rating) + " | notes " + notes)

		tool.Id = id
		tool.Name = name
		tool.Category = category
		tool.URL = url
		tool.Rating = rating
		tool.Notes = notes
	}

	_ = s.t.ExecuteTemplate(w, "Show", tool) //nolint
}

func (s *toolSvc) New(w http.ResponseWriter, r *http.Request) {
	_ = s.t.ExecuteTemplate(w, "New", nil) //nolint
}

func (s *toolSvc) Edit(w http.ResponseWriter, r *http.Request) {
	nId := r.URL.Query().Get("id")
	selDB, err := s.db.Query("SELECT * FROM tools WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}

	tool := Tool{}

	for selDB.Next() {
		var id, rating int
		var name, category, url, notes string
		err := selDB.Scan(&id, &name, &category, &url, &rating, &notes)
		if err != nil {
			panic(err.Error())
		}

		tool.Id = id
		tool.Name = name
		tool.Category = category
		tool.URL = url
		tool.Rating = rating
		tool.Notes = notes
	}

	_ = s.t.ExecuteTemplate(w, "Edit", tool) //nolint
}

func (s *toolSvc) Insert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		category := r.FormValue("category")
		url := r.FormValue("url")
		rating := r.FormValue("rating")
		notes := r.FormValue("notes")
		insForm, err := s.db.Prepare("INSERT INTO tools (name, category, url, rating, notes) VALUES (?, ?, ?, ?, ?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, category, url, rating, notes)
		log.Println("Insert Data: name " + name + " | category " + category + " | url " + url + " | rating " + rating + " | notes " + notes)
	}
	http.Redirect(w, r, "/", 301)
}

func (s *toolSvc) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		category := r.FormValue("category")
		url := r.FormValue("url")
		rating := r.FormValue("rating")
		notes := r.FormValue("notes")
		id := r.FormValue("uid")
		insForm, err := s.db.Prepare("UPDATE tools SET name=?, category=?, url=?, rating=?, notes=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, category, url, rating, notes, id)
		log.Println("UPDATE Data: name " + name + " | category " + category + " | url " + url + " | rating " + rating + " | notes " + notes)
	}
	http.Redirect(w, r, "/", 301)
}

func (s *toolSvc) Delete(w http.ResponseWriter, r *http.Request) {
	tool := r.URL.Query().Get("id")
	delForm, err := s.db.Prepare("DELETE FROM tools WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(tool)
	log.Println("DELETE " + tool)
	http.Redirect(w, r, "/", 301)
}

func main() {
	var (
		debug bool
	)

	flag.BoolVar(&debug, "debug", debug, "turn on debug logging")
	flag.Parse()

	if debug {
		log = xlog.New(os.Stdout, "", xlog.LstdFlags)
	}

	log.Println("Server started on: http://localhost:8080")

	db, err := dbConn()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	t, err := template.ParseGlob("templates/*")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	tSvc := newToolSvc(db, t)

	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(tSvc.Index))
	mux.Handle("/show", http.HandlerFunc(tSvc.Show))
	mux.Handle("/new", http.HandlerFunc(tSvc.New))
	mux.Handle("/edit", http.HandlerFunc(tSvc.Edit))
	mux.Handle("/insert", http.HandlerFunc(tSvc.Insert))
	mux.Handle("/update", http.HandlerFunc(tSvc.Update))
	mux.Handle("/delete", http.HandlerFunc(tSvc.Delete))

	http.ListenAndServe(":8080", mux)
}
