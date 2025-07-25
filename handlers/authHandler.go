package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	//"github.com/D:/Golang/GoNotes-Backend/models"
	"github.com/Sunny914/GoNotes-Backend/models"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthHandler struct {
	userModel *models.UserModel
	jwtSecret []byte
}

func NewAuthHandler(userModel *models.UserModel, jwtSecret []byte) *AuthHandler {
	return &AuthHandler{
		userModel: userModel,
		jwtSecret: jwtSecret,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email     string  `json:"email"`
		Password  string  `json:"password"` 
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status" :  false,
			"message":  err.Error(),
			"token"  :  "",
		})
		return
	}

	user, err := h.userModel.Create(input.Email, input.Password)
	if err != nil {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status" :  false,
			"message":  err.Error(),
			"token"  :  "",
		})
		return 
	}

	tokenString, err := h.generateJWT(user.ID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status" :  false,
			"message":  "Failed to generate token",
			"token"  :  "",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		   "status" :  true,
			"message":  "User Registered Successfully",
			"token"  :  tokenString,
	})

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email     string  `json:"email"`
		Password  string  `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":   false,
			"message":  err.Error(),
			"token":    "",
		})	
		return 
	}

	user, err := h.userModel.GetByEmail(input.Email)
	if err != nil || !h.userModel.VerifyPassword(user, input.Password) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":   false,
			"message":  "Invalid credentials",
			"token":    "",
		})
		return 
	}

	tokenString, err := h.generateJWT(user.ID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":   false,
			"message":  "Failed to Generate token",
			"token":    "",
		})
		return 
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		   "status":   true,
			"message":  "Login Successfull",
			"token":    tokenString,
	})

}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// JWT is stateless : Just ask client to Discard the Token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   true,
		"message":  "Logout Successfull",
		"token":    "",
	})
}

func (h *AuthHandler) generateJWT(userID primitive.ObjectID) (string, error) {
	claims := jwt.MapClaims {
		"user_id": userID.Hex(),
		"iat":     time.Now().Unix(), // gives new token at every login
		                              // No Expiration - Token Never Expires
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(h.jwtSecret)
}

