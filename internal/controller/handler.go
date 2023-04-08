package controller

import (
	"log"
	"net/http"

	"forum/internal/service"
)

type Handler struct {
	ErrorLog *log.Logger
	Service  *service.Service
}

func NewHandler(logger *log.Logger, service *service.Service) *Handler {
	return &Handler{logger, service}
}

func Routes(h *Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.homepage)
	mux.HandleFunc("/sign", h.user)
	return mux
}
