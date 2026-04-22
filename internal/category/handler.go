package category

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Gwilides/finance-tracker/configs"
	"github.com/Gwilides/finance-tracker/pkg/middleware"
	"github.com/Gwilides/finance-tracker/pkg/req"
	"github.com/Gwilides/finance-tracker/pkg/res"
)

type CategoryHandlerDeps struct {
	Service *CategoryService
	Config  *configs.AuthConfig
}

type CategoryHandler struct {
	service *CategoryService
}

func NewCategoryHandler(router *http.ServeMux, deps *CategoryHandlerDeps) {
	handler := &CategoryHandler{
		service: deps.Service,
	}
	router.Handle("POST /category", middleware.IsAuthed(handler.create(), deps.Config))
	router.Handle("GET /category", middleware.IsAuthed(handler.getAll(), deps.Config))
	router.Handle("PATCH /category/{id}", middleware.IsAuthed(handler.update(), deps.Config))
	router.Handle("DELETE /category/{id}", middleware.IsAuthed(handler.delete(), deps.Config))
}

func (handler *CategoryHandler) extractEmail(w http.ResponseWriter, r *http.Request) (string, bool) {
	email, ok := r.Context().Value(middleware.ContextEmailKey).(string)
	if !ok {
		res.Json(w, "unauthorized", http.StatusUnauthorized)
		return "", ok
	}
	return email, ok

}

func (handler *CategoryHandler) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := handler.extractEmail(w, r)
		if !ok {
			return
		}
		body, err := req.HandleBody[CategoryCreateRequest](w, r)
		if err != nil {
			return
		}
		category, err := handler.service.Create(email, body)
		if err != nil {
			res.Json(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Json(w, category, http.StatusCreated)
	}
}

func (handler *CategoryHandler) getAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := handler.extractEmail(w, r)
		if !ok {
			return
		}
		categories, err := handler.service.GetAll(email)
		if err != nil {
			res.Json(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Json(w, categories, http.StatusOK)
	}
}

func (handler *CategoryHandler) update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := handler.extractEmail(w, r)
		if !ok {
			return
		}
		body, err := req.HandleBody[CategoryUpdateRequest](w, r)
		if err != nil {
			return
		}
		stringId := r.PathValue("id")
		id, err := strconv.Atoi(stringId)
		if err != nil {
			res.Json(w, err.Error(), http.StatusBadRequest)
			return
		}
		category, err := handler.service.Update(email, uint(id), body)
		if errors.Is(err, ErrCategoryNotFound) {
			res.Json(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, ErrForbidden) {
			res.Json(w, err.Error(), http.StatusForbidden)
			return
		}
		if err != nil {
			res.Json(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Json(w, category, http.StatusOK)
	}
}

func (handler *CategoryHandler) delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := handler.extractEmail(w, r)
		if !ok {
			return
		}
		stringId := r.PathValue("id")
		id, err := strconv.Atoi(stringId)
		if err != nil {
			res.Json(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = handler.service.Delete(email, uint(id))
		if errors.Is(err, ErrCategoryNotFound) {
			res.Json(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, ErrForbidden) {
			res.Json(w, err.Error(), http.StatusForbidden)
			return
		}
		if err != nil {
			res.Json(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
