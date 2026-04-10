package account

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Gwilides/finance-tracker/configs"
	"github.com/Gwilides/finance-tracker/pkg/middleware"
	"github.com/Gwilides/finance-tracker/pkg/req"
	"github.com/Gwilides/finance-tracker/pkg/res"
	"gorm.io/gorm"
)

type AccountHandlerDeps struct {
	Service *AccountService
	Config  *configs.AuthConfig
}

type AccountHandler struct {
	service *AccountService
}

func NewAccountHandler(router *http.ServeMux, deps *AccountHandlerDeps) {
	handler := &AccountHandler{
		service: deps.Service,
	}
	router.Handle("POST /account", middleware.IsAuthed(handler.create(), deps.Config))
	router.Handle("GET /account/{id}", middleware.IsAuthed(handler.goTo(), deps.Config))
	router.Handle("PATCH /account/{id}", middleware.IsAuthed(handler.update(), deps.Config))
	router.Handle("DELETE /account/{id}", middleware.IsAuthed(handler.delete(), deps.Config))
}

func (handler *AccountHandler) extractEmail(w http.ResponseWriter, r *http.Request) (string, bool) {
	email, ok := r.Context().Value(middleware.ContextEmailKey).(string)
	if !ok {
		res.Json(w, "unauthorized", http.StatusUnauthorized)
		return "", ok
	}
	return email, ok
}

func (handler *AccountHandler) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := handler.extractEmail(w, r)
		if !ok {
			return
		}
		body, err := req.HandleBody[AccountCreateRequest](w, r)
		if err != nil {
			return
		}
		account, err := handler.service.Create(email, body)
		if err != nil {
			res.Json(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Json(w, account, http.StatusCreated)
	}
}

func (handler *AccountHandler) goTo() http.HandlerFunc {
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
		account, err := handler.service.GetById(email, uint(id))
		if errors.Is(err, ErrForbidden) {
			res.Json(w, err.Error(), http.StatusForbidden)
			return
		}
		if err != nil {
			res.Json(w, err.Error(), http.StatusNotFound)
			return
		}
		res.Json(w, account, http.StatusOK)
	}
}

func (handler *AccountHandler) update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := handler.extractEmail(w, r)
		if !ok {
			return
		}
		body, err := req.HandleBody[AccountUpdateRequest](w, r)
		if err != nil {
			return
		}
		stringId := r.PathValue("id")
		id, err := strconv.Atoi(stringId)
		if err != nil {
			res.Json(w, err.Error(), http.StatusBadRequest)
			return
		}
		account, err := handler.service.Update(email, &Account{
			Model: gorm.Model{ID: uint(id)},
			Type:  body.Type,
			Title: body.Title,
		})
		if errors.Is(err, ErrForbidden) {
			res.Json(w, err.Error(), http.StatusForbidden)
			return
		}
		if err != nil {
			res.Json(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Json(w, account, http.StatusOK)
	}
}

func (handler *AccountHandler) delete() http.HandlerFunc {
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
		if errors.Is(err, ErrAccountNotFound) {
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
