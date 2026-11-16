package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/WhiteRain7/GoTestCase/pkg/utils/coded-errors"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/users/domain"
	"github.com/WhiteRain7/GoTestCase/wallets/pkg/decoder"
)

type Handler struct {
	uc     UserUseCase
	auth   JWTGen
	logger Logger
}

type Logger interface {
	Info(string)
}

type UserUseCase interface {
	Create(ctx context.Context, username string, password string, email string) (*domain.User, error)
	Validate(ctx context.Context, username string, password string) error
	GetID(ctx context.Context, username string) (int, error)
}

type JWTGen interface {
	GenerateJWT(id int, username string) (string, error)
}

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	Message string `json:"message"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

func NewHandler(uc UserUseCase, auth JWTGen, logger Logger) *Handler {
	logger.Info("Handler init - Wallets")
	return &Handler{uc, auth, logger}
}

// SignUp
// @Summary      Register new user
// @Description  Creates new user with provided email, password and username.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request body SignUpRequest true "Body"
// @Success      200 {object} SignUpResponse
// @Failure      400 {string} string "Username or email already exists"
// @Router       /api/v1/users/register [post]
func (h *Handler) SignUp(res http.ResponseWriter, req *http.Request) {
	body := &SignUpRequest{}
	if err := decoder.Decode(req, body); err != nil {
		codederrors.ErrorResponse(err, res)
	}

	fmt.Print(body)

	_, err := h.uc.Create(req.Context(), body.Username, body.Password, body.Email)
	if err != nil {
		codederrors.ErrorResponse(err, res)
	} else {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(201)
		json.NewEncoder(res).Encode(SignUpResponse{"User registered successfully"})
	}
}

// SignIn
// @Summary      Login
// @Description  Authorize for token.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request body SignInRequest true "Body"
// @Success      200 {object} SignInResponse
// @Failure      400 {string} string "Username not exists or password didn't match"
// @Router       /api/v1/users/login [post]
func (h *Handler) SignIn(res http.ResponseWriter, req *http.Request) {
	body := &SignInRequest{}
	if err := decoder.Decode(req, body); err != nil {
		codederrors.ErrorResponse(err, res)
	}

	err := h.uc.Validate(req.Context(), body.Username, body.Password)
	if err != nil {
		codederrors.ErrorResponse(err, res)
		return
	}

	id, err := h.uc.GetID(req.Context(), body.Username)
	if err != nil {
		codederrors.ErrorResponse(err, res)
		return
	}

	token, err := h.auth.GenerateJWT(id, body.Username)
	if err != nil {
		codederrors.ErrorResponse(err, res)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(200)
	json.NewEncoder(res).Encode(SignInResponse{token})
}
