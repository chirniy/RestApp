package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"serversTest2/internal/domain"
	"serversTest2/internal/jwt"
	m "serversTest2/internal/middleware"
	"serversTest2/internal/usecase"
)

type UserHandler struct {
	usecase *usecase.UserUsecase
}

func NewUserHandler(u *usecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: u}
}

func (h *UserHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var userInput domain.UserInput

	err := json.NewDecoder(r.Body).Decode(&userInput)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	if err != nil {
		// handle error
	}
	userInput.Password = string(hashedPassword)
	log.Println(userInput)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = h.usecase.Create(r.Context(), userInput)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) handleGetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.usecase.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Error fetching users", http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, "Error encoding users", http.StatusBadRequest)
		return
	}
}

func (h *UserHandler) handleGetUser(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	user, err := h.usecase.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Error fetching user", http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Error encoding users", http.StatusBadRequest)
		return
	}
}

func (h *UserHandler) handleUpdateUser(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	var userPartial domain.PartialUser
	err := json.NewDecoder(r.Body).Decode(&userPartial)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = h.usecase.Update(r.Context(), id, userPartial)
	if err != nil {
		http.Error(w, "Error updating user", http.StatusBadRequest)
		return
	}

	h.handleGetUser(w, r, id)
}

func (h *UserHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	err := h.usecase.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, "Error deleting user", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

func parseUUIDFromRequest(r *http.Request) (uuid.UUID, error) {
	vars := mux.Vars(r)
	return uuid.Parse(vars["id"])
}

func (h *UserHandler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	m.WithHeaders(w)

	switch r.Method {
	case http.MethodGet:
		h.handleGetAllUsers(w, r)
		return
	case http.MethodPost:
		h.handleCreateUser(w, r)
		return
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (h *UserHandler) UsersHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	m.WithHeaders(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	parsedId, err := parseUUIDFromRequest(r)
	if err != nil {
		http.Error(w, "Неверный UUID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetUser(w, r, parsedId)
		return
	case http.MethodPut, http.MethodPatch:
		h.handleUpdateUser(w, r, parsedId)
		return
	case http.MethodDelete:
		h.handleDeleteUser(w, r, parsedId)
		return
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var loginUser domain.LoginInput

	err := json.NewDecoder(r.Body).Decode(&loginUser)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	user, err := h.usecase.GetByEmail(r.Context(), loginUser.Email)
	if err != nil {
		http.Error(w, "Error fetching user", http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusBadRequest)
		return
	}

	tokenString, err := jwt.GenerateToken(user.ID.String())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}
