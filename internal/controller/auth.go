package controller

import (
	"net/http"

	"forum/internal/models"
)

type Data struct {
	User         models.User
	Content      any
	IsAuthorized bool
	Validator    map[string]string
}

type ErrorData struct {
	Status int
	Text   string
}

func (h *Handler) homepage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.errorpage(w, http.StatusNotFound, nil)
		return
	}
	posts, err := h.Service.Posts.GetAll()
	if err != nil {
		h.errorpage(w, http.StatusInternalServerError, err)
		return
	}
	data := r.Context().Value(ctxKey).(*Data)
	data.Content = posts
	h.templaterender(w, http.StatusOK, "index.html", data)
}

func (h *Handler) signup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := r.Context().Value(ctxKey).(*Data)
		h.templaterender(w, http.StatusOK, "signup.html", data)
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			h.errorpage(w, http.StatusInternalServerError, err)
			return
		}
		form := models.User{
			Email:    r.PostForm.Get("email"),
			Name:     r.PostForm.Get("name"),
			Password: r.PostForm.Get("password"),
		}
		err = h.Service.SignUp(form)
		if err != nil {
			h.errorpage(w, http.StatusBadRequest, err)
			return
		}
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
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)
	err := h.Service.LogOut(*user.Token)
	if err != nil {
		h.ErrorLog.Println(err)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
