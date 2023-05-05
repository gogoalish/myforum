package controller

import (
	"net/http"
	"strconv"
)

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
	switch r.Method {
	case http.MethodGet:
		h.templaterender(w, http.StatusOK, "index.html", data)
	case http.MethodPost:
		err = r.ParseForm()
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		reaction := r.PostForm.Get("reaction")
		userID := data.User.ID
		postID, err := strconv.Atoi(r.PostForm.Get("post"))
		if err != nil || (reaction != "like" && reaction != "dislike") {
			h.errorpage(w, http.StatusBadRequest, nil)
			return
		}
		err = h.Service.React(postID, userID, reaction)
		if err != nil {
			h.errorpage(w, http.StatusInternalServerError, err)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (h *Handler) filter(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var catID []int
	for _, value := range r.Form["filter"] {

		number, _ := strconv.Atoi(value)
		catID = append(catID, number)
	}
	posts, err := h.Service.Posts.GetFiltered(catID)
	if err != nil {
		h.errorpage(w, http.StatusInternalServerError, err)
		return
	}
	data := r.Context().Value(ctxKey).(*Data)
	data.Content = posts

	h.templaterender(w, http.StatusOK, "index.html", data)
	// switch r.Method {
	// case http.MethodGet:
	// case http.MethodPost:
	// 	err = r.ParseForm()
	// 	if err != nil {
	// 		http.Redirect(w, r, "/filter", http.StatusSeeOther)
	// 		return
	// 	}
	// 	reaction := r.PostForm.Get("reaction")
	// 	userID := data.User.ID
	// 	postID, err := strconv.Atoi(r.PostForm.Get("post"))
	// 	if err != nil || (reaction != "like" && reaction != "dislike") {
	// 		h.errorpage(w, http.StatusBadRequest, nil)
	// 		return
	// 	}
	// 	err = h.Service.React(postID, userID, reaction)
	// 	if err != nil {
	// 		h.errorpage(w, http.StatusInternalServerError, err)
	// 		return
	// 	}
	// 	http.Redirect(w, r, "/filter", http.StatusSeeOther)
	// }
}
