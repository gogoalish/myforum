package controller

import (
	"bytes"
	"html/template"
	"net/http"
	"path/filepath"
	"runtime/debug"

	"forum/internal/models"
)

type Data struct {
	User         models.User
	Content      any
	IsAuthorized bool
	ErrMsgs      map[string]string
}

type ErrorData struct {
	Status int
	Text   string
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
	if err != nil {
		h.ErrorLog.Printf("server error: %v", err)
	}
	msg := http.StatusText(status)
	errdata := ErrorData{status, msg}
	h.templaterender(w, status, "errors.html", errdata)
}

func (h *Handler) serverError(w http.ResponseWriter, err error) {
	h.ErrorLog.Printf("%s\n%s", err.Error(), debug.Stack())
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}
