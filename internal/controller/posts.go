package controller

import (
	"net/http"

	"forum/internal/models"
)

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
