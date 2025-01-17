package handler

import (
	"net/http"

	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jinzhu/copier"
	"github.com/tendant/simple-user/user"
	"golang.org/x/exp/slog"
)

type Handler struct {
	UserService *user.UserService
}

func (h *Handler) Routes(r *chi.Mux) {
	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, http.StatusText(http.StatusOK))
	})

	r.With(httpin.NewInput(UserInput{})).Post("/api/users", h.handleCreateUser)
	r.Get("/api/users", h.handleFindUsers)
}

type UserParams struct {
	Email string `json:"email"`
}

type UserInput struct {
	Payload *UserParams `in:"body=json"`
}

func (h *Handler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	slog.Debug("debug ***")
	input := r.Context().Value(httpin.Input).(*UserInput)
	slog.Debug("input ****:", "input", input, "user params", input.Payload)
	params := input.Payload
	userParams := user.UserParams{}
	copier.Copy(&userParams, params)
	// svc := user.UserService{}
	user, err := h.UserService.Create(r.Context(), userParams)
	if err != nil {
		slog.Error("Falied creating users", "err", err)
		render.PlainText(w, r, http.StatusText(http.StatusInternalServerError))
	}
	render.JSON(w, r, user)
}

func (h *Handler) handleFindUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.UserService.FindUsers(r.Context())
	if err != nil {
		slog.Error("Failed finding users:", "err", err)
		render.PlainText(w, r, http.StatusText(http.StatusInternalServerError))
	}
	slog.Debug("users:", "users", users)
	render.JSON(w, r, users)
}
