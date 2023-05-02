package controller

import (
	"bytes"
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
	mux.HandleFunc("/", (h.homepage))
	mux.HandleFunc("/signup", (h.signup))
	mux.HandleFunc("/signin", h.signin)
	mux.HandleFunc("/posts/create", (h.create))
	mux.HandleFunc("/logout", (h.logout))
	mux.HandleFunc("/posts/", (h.showpost))
	return h.CheckAuth(SecureHeaders(mux))
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

func (h *Handler) templaterender(w http.ResponseWriter, status int, page string, data any) {
	buf := new(bytes.Buffer)
	err := h.Tempcache.ExecuteTemplate(buf, page, data)
	if err != nil {
		h.serverError(w, err)
		return
	}
	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (h *Handler) errorpage(w http.ResponseWriter, status int, err error) {
	msg := http.StatusText(status)
	if err != nil {
		h.ErrorLog.Printf("server error: %v", err)
	}
	errdata := ErrorData{status, msg}
	h.templaterender(w, status, "errors.html", errdata)
}

func (h *Handler) serverError(w http.ResponseWriter, err error) {
	h.ErrorLog.Printf("%s\n%s", err.Error(), debug.Stack())
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}
