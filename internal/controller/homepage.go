package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"forum/internal/models"
)

func (h *Handler) homepage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.errorpage(w, http.StatusNotFound, nil)
		return
	}
	data := r.Context().Value(ctxKey).(*Data)
	var err error
	switch r.Method {
	case http.MethodGet:
		query := r.URL.Query()
		var catID []int
		for _, value := range query["filter"] {
			number, _ := strconv.Atoi(value)
			if number > 0 {
				catID = append(catID, number)
			}
		}
		if catID != nil {
			data.Content, err = h.Service.Posts.GetFiltered(catID)
			if err != nil {
				h.errorpage(w, http.StatusInternalServerError, err)
				return
			}
		} else {
			data.Content, err = h.Service.Posts.GetAll()
			if err != nil {
				h.errorpage(w, http.StatusInternalServerError, err)
				return
			}
		}
		data.IsEmpty = (len(data.Content.([]*models.Post)) == 0)
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
			fmt.Println(err)
			h.errorpage(w, http.StatusBadRequest, nil)
			return
		}
		err = h.Service.Posts.React(postID, userID, reaction)
		if err != nil {
			h.errorpage(w, http.StatusInternalServerError, err)
			return
		}
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	}
}
