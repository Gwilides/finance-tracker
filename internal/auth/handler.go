package auth

import (
	"net/http"

	"github.com/Gwilides/finance-tracker/configs"
	"github.com/Gwilides/finance-tracker/pkg/jwt"
	"github.com/Gwilides/finance-tracker/pkg/req"
	"github.com/Gwilides/finance-tracker/pkg/res"
)

type AuthHandlerDeps struct {
	AuthService *AuthService
	Config      *configs.AuthConfig
}

type AuthHandler struct {
	authService *AuthService
	config      *configs.AuthConfig
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		authService: deps.AuthService,
		config:      deps.Config,
	}
	router.HandleFunc("POST /auth/login", handler.login())
	router.HandleFunc("POST /auth/register", handler.register())
}

func (handler *AuthHandler) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LoginRequest](w, r)
		if err != nil {
			return
		}
		email, err := handler.authService.Login(body)
		if err != nil {
			res.Json(w, err.Error(), http.StatusUnauthorized)
			return
		}
		token, err := jwt.NewJWT(handler.config.Secret).Create(&jwt.JWTData{
			Email: email,
		})
		if err != nil {
			res.Json(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := LoginResponse{
			Token: token,
		}
		res.Json(w, data, http.StatusOK)
	}
}

func (handler *AuthHandler) register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[RegisterRequest](w, r)
		if err != nil {
			return
		}
		email, err := handler.authService.Register(body)
		if err != nil {
			res.Json(w, err.Error(), http.StatusInternalServerError)
			return
		}
		token, err := jwt.NewJWT(handler.config.Secret).Create(&jwt.JWTData{
			Email: email,
		})
		if err != nil {
			res.Json(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := RegisterResponse{
			Token: token,
		}
		res.Json(w, data, http.StatusCreated)
	}
}
