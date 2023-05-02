package controller

import "net/http"

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
