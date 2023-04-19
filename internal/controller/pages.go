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
		err := h.Service.SignUp(form)
		if err != nil {
			h.ErrorLog.Println(err)
		}
		log.Printf("user created email: %s, name: %s, password: %s", form.Email, form.Name, form.Password)
	}
}

func (h *Handler) signin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := h.Tempcache.ExecuteTemplate(w, "signin.html", nil)
		if err != nil {
			h.ErrorLog.Println(err)
		}
	case http.MethodPost:
		email := r.FormValue("email")
		password := r.FormValue("password")
		user, err := h.Service.SignIn(email, password)
		if err != nil {
			h.ErrorLog.Println(err)
		}
		log.Println(user)
		// http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
