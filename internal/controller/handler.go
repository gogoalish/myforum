package controller

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"runtime/debug"

	"forum/internal/repository"
	"forum/internal/service"
)

type Handler struct {
	ErrorLog  *log.Logger
	Service   *service.Service
	Tempcache *template.Template
}

func NewHandler(logger *log.Logger, DB *sql.DB) (*Handler, error) {
	repo := repository.New(DB)
	service := service.New(repo)
	tempcache, err := template.ParseGlob("ui/*.html")
	return &Handler{logger, service, tempcache}, err
}

func Routes(h *Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.homepage)
	mux.HandleFunc("/signup", h.signup)
	mux.HandleFunc("/signin", h.signin)
	mux.HandleFunc("/create", h.CheckAuth(h.create))
	return SecureHeaders(mux)
}

func Newtscache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./ui/*.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}

// func (h *Handler) render(w http.ResponseWriter, status int, page string, data string) {
// 	ts, ok := h.Tempcache[page]
// 	if !ok {
// 		err := fmt.Errorf("the template %s does not exist", page)
// 		h.serverError(w, err)
// 		return
// 	}
// 	buf := new(bytes.Buffer)
// 	err := ts.ExecuteTemplate(buf, "index", data)
// 	if err != nil {
// 		h.serverError(w, err)
// 		return
// 	}
// 	w.WriteHeader(status)
// 	buf.WriteTo(w)
// }

func (app *Handler) serverError(w http.ResponseWriter, err error) {
	log.Printf("%s\n%s", err.Error(), debug.Stack())
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

// unused yet
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}

	return cache, nil
}
