package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Shercosta/digi-wallet/models"
	"github.com/Shercosta/digi-wallet/response"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtSecre = []byte("Shercosta")

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			response.JSONError(w, http.StatusBadRequest, err.Error(), nil)
			return
		}

		var user models.User
		if err := db.Where("username = ?", body.Username).First(&user).Error; err != nil {
			response.JSONError(w, http.StatusUnauthorized, err.Error(), nil)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
			response.JSONError(w, http.StatusUnauthorized, err.Error(), nil)
			return
		}

		// generate JWT
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": user.ID,
			"exp":     time.Now().Add(time.Hour * 72).Unix(),
		})
		tokenString, _ := token.SignedString(jwtSecre)

		constructResponse := map[string]string{
			"user_id": strconv.Itoa(int(user.ID)),
			"token":   tokenString,
		}

		response.JSONSuccess(w, constructResponse, nil, nil)
	}
}

func Register(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			response.JSONError(w, http.StatusBadRequest, err.Error(), nil)
			return
		}

		var exist models.User
		if err := db.Where("username = ?", body.Username).First(&exist).Error; err == nil {
			response.JSONError(w, http.StatusBadRequest, "username already exists", nil)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		user := models.User{
			Username: body.Username,
			Password: string(hashedPassword),
		}

		if err := db.Create(&user).Error; err != nil {
			response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		// create balance for user
		balance := models.Balance{
			Amount: 100000,
			UserID: user.ID,
		}

		if err := db.Create(&balance).Error; err != nil {
			response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		response.JSONSuccess(w, user, nil, nil)
	}
}
