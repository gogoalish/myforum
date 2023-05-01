package controller

import (
	"net/http"
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
		data := r.Context().Value(ctxKey).(Data)
		h.Service.Create(data.User.ID, title, content)
	}
}
