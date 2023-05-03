package controller

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

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
	mux.HandleFunc("/posts/create", h.requireauth(h.postcreate))
	mux.HandleFunc("/logout", h.requireauth(h.logout))
	mux.HandleFunc("/posts/", (h.postview))
	return h.middleware(SecureHeaders(mux))
}
