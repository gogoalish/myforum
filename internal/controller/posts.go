package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"forum/internal/models"
)

func (h *Handler) postcreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.templaterender(w, http.StatusOK, "create.html", nil)
	case http.MethodPost:
		data := r.Context().Value(ctxKey).(*Data)
		err := r.ParseForm()
		if err != nil {
			h.errorpage(w, http.StatusInternalServerError, err)
			return
		}

		var catid []int

		for _, value := range r.PostForm["cat"] {
			number, _ := strconv.Atoi(value)
			catid = append(catid, number)
		}
		post := &models.Post{
			UserID:  data.User.ID,
			Title:   r.PostForm.Get("title"),
			Content: r.PostForm.Get("content"),
			CatID:   catid,
		}
		if post.Title == "" && post.Content == "" {
			h.errorpage(w, http.StatusBadRequest, nil)
			return
		}
		id, err := h.Service.Posts.Create(post)
		if err != nil {
			h.errorpage(w, http.StatusInternalServerError, err)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/posts/%v", id), http.StatusSeeOther)
	}
}

func (h *Handler) postview(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value(ctxKey).(*Data)
	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if err != nil {
		h.errorpage(w, http.StatusNotFound, err)
		return
	}
	switch r.Method {
	case http.MethodGet:
		post, err := h.Service.Posts.GetById(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				h.errorpage(w, http.StatusNotFound, err)
				return
			}
			h.errorpage(w, http.StatusInternalServerError, err)
			return
		}
		post.Comments, err = h.Service.Comments.Fetch(post.ID)
		if err != nil && !errors.Is(err, models.ErrNoRecord) {
			h.errorpage(w, http.StatusInternalServerError, err)
			return
		}
		data.Content = post
		h.templaterender(w, http.StatusOK, "post.html", data)
	case http.MethodPost:
		if data.User == (models.User{}) {
			h.templaterender(w, http.StatusUnauthorized, "index.html", nil)
			return
		}
		if err := r.ParseForm(); err != nil {
			h.errorpage(w, http.StatusInternalServerError, err)
			return
		}
		comment := &models.Comment{
			PostID:  id,
			UserID:  data.User.ID,
			Content: r.PostForm.Get("content"),
		}
		comment.ParentID, err = strconv.Atoi(r.PostForm.Get("parent"))
		if comment.Content == "" && err != nil {
			h.errorpage(w, http.StatusBadRequest, err)
			return
		}
		h.Service.Comments.Create(comment)
		http.Redirect(w, r, fmt.Sprintf("/posts/%v", id), http.StatusSeeOther)
	}
}

func (h *Handler) postreaction(w http.ResponseWriter, r *http.Request) {
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
	reaction := r.PostForm.Get("reaction")
	userID := data.User.ID
	postID, err := strconv.Atoi(r.PostForm.Get("post"))
	if err != nil || (reaction != "like" && reaction != "dislike") {
		h.errorpage(w, http.StatusBadRequest, nil)
		return
	}
	err = h.Service.Posts.React(postID, userID, reaction)
	if err != nil {
		h.errorpage(w, http.StatusInternalServerError, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/posts/%v", postID), http.StatusSeeOther)
}

func (h *Handler) likedposts(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value(ctxKey).(*Data)
	var err error
	data.Content, err = h.Service.Posts.GetUserLiked(data.User.ID)
	if err != nil {
		h.errorpage(w, http.StatusInternalServerError, fmt.Errorf("controller - liked - GetUserLiked - %w", err))
		return
	}
	h.templaterender(w, http.StatusOK, "index.html", data)
}

func (h *Handler) createdposts(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value(ctxKey).(*Data)
	var err error
	data.Content, err = h.Service.Posts.GetUserCreated(data.User.ID)
	if err != nil {
		h.errorpage(w, http.StatusInternalServerError, fmt.Errorf("controller - liked - GetUserLiked - %w", err))
		return
	}
	h.templaterender(w, http.StatusOK, "index.html", data)
}
