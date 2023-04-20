package controller

import (
	"log"
	"net/http"

	"forum/internal/models"
)

func (h *Handler) homepage(w http.ResponseWriter, r *http.Request) {
	posts, _ := h.Service.All()
	h.Tempcache.ExecuteTemplate(w, "index.html", posts)
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
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (h *Handler) signin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := h.Tempcache.ExecuteTemplate(w, "signin.html", nil)
		if err != nil {
			h.ErrorLog.Println(err)
			return
		}
	case http.MethodPost:
		email := r.FormValue("email")
		password := r.FormValue("password")
		user, err := h.Service.SignIn(email, password)
		if err != nil {
			h.ErrorLog.Println(err)
			return
		}
		cookie := &http.Cookie{
			Name:  "session",
			Value: *user.Token,
		}
		http.SetCookie(w, cookie)
		log.Println(user)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)
	err := h.Service.LogOut(*user.Token)
	if err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := h.Tempcache.ExecuteTemplate(w, "create.html", nil)
		if err != nil {
			h.ErrorLog.Println(err)
		}
	case http.MethodPost:
		title := r.FormValue("title")
		content := r.FormValue("content")
		user := r.Context().Value("user").(models.User)
		h.Service.Create(user.ID, title, content)
	}
}
