package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"forum/internal/models"
)

func (h *Handler) reaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errorpage(w, http.StatusMethodNotAllowed, nil)
		return
	}
	data := r.Context().Value(ctxKey).(*Data)
	if data.User == (models.User{}) {
		h.errorpage(w, http.StatusUnauthorized, nil)
		return
	}
	if err := r.ParseForm(); err != nil {
		h.errorpage(w, http.StatusInternalServerError, err)
		return
	}
	object := r.PostForm.Get("object")
	reaction := r.PostForm.Get("reaction")
	userID := data.User.ID

	id, err := strconv.Atoi(r.PostForm.Get("id"))
	if err != nil || (reaction != "like" && reaction != "dislike") {
		h.errorpage(w, http.StatusBadRequest, nil)
		return
	}
	if object == "comment" {
		err = h.Service.Comments.React(id, userID, reaction)
		if err != nil {
			h.errorpage(w, http.StatusInternalServerError, err)
			return
		}
		comment, err := h.Service.Comments.GetByID(id)
		if err != nil && !errors.Is(err, models.ErrNoRecord) {
			h.errorpage(w, http.StatusInternalServerError, err)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/posts/%v", comment.PostID), http.StatusSeeOther)
		return
	}
	err = h.Service.Posts.React(id, userID, reaction)
	if err != nil {
		h.errorpage(w, http.StatusInternalServerError, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/posts/%v", id), http.StatusSeeOther)
}
