package controller

import (
	"log"
	"net/http"
	"text/template"

	"forum/internal/models"
)

func (h *Handler) homepage(w http.ResponseWriter, r *http.Request) {
	h.Service.Posts.Create("email", "name")
	tmpl, err := template.ParseFiles("ui/index.html")
	if err != nil {
		h.ErrorLog.Println(err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handler) user(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	name := r.FormValue("name")
	password := r.FormValue("password")
	h.Service.Users.Create(models.User{0, email, name, password})
	user, err := h.Service.Users.Get(email, name)
	if err != nil {
		h.ErrorLog.Println(err)
		return
	}
	log.Printf("user created id %d email: %s, name: %s, password: %s", user.ID, user.Email, user.Name, user.Password)
}
