package controller

import (
	"log"
	"net/http"

	"forum/internal/models"
)

func (h *Handler) homepage(w http.ResponseWriter, r *http.Request) {
	h.Tempcache.ExecuteTemplate(w, "index.html", nil)
}

func (h *Handler) signup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := h.Tempcache.ExecuteTemplate(w, "signup.html", nil)
		if err != nil {
			h.ErrorLog.Println(err)
		}
	case http.MethodPost:
		form := models.User{
			Email:    r.FormValue("email"),
			Name:     r.FormValue("name"),
			Password: r.FormValue("password"),
		}
		h.Service.SignUp(form)
		user, err := h.Service.Get(form.Email, form.Name)
		if err != nil {
			h.ErrorLog.Println(err)
			return
		}
		log.Printf("user created id %d email: %s, name: %s, password: %s", user.ID, user.Email, user.Name, user.Password)
	}
}
