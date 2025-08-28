package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Shercosta/digi-wallet/middleware"
	"github.com/Shercosta/digi-wallet/models"
	"github.com/Shercosta/digi-wallet/request"
	"github.com/Shercosta/digi-wallet/response"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	response.JSONSuccess(w, "hello from routes", nil, nil)
}

func ListUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sortOrder := r.URL.Query().Get("sort_amount")

		var users []struct {
			ID       uint    `json:"id"`
			Username string  `json:"username"`
			Amount   float64 `json:"amount"`
		}

		query := db.Table("users").
			Select("users.id, users.username, COALESCE(balances.amount, 0) as amount").
			Joins("LEFT JOIN balances ON users.id = balances.user_id")

		if sortOrder == "asc" {
			query = query.Order("balances.amount ASC")
		} else if sortOrder == "desc" {
			query = query.Order("balances.amount DESC")
		}

		if err := query.Scan(&users).Error; err != nil {
			response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		response.JSONSuccess(w, users, nil, nil)
	}
}

func DeleteUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		currentUserID := middleware.GetUserID(ctx)

		idParam := chi.URLParam(r, "id")
		targetID64, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			response.JSONError(w, http.StatusBadRequest, "invalid user id", nil)
			return
		}
		targetID := uint(targetID64)

		if uint(currentUserID) == targetID {
			response.JSONError(w, http.StatusForbidden, "you cannot delete yourself", nil)
			return
		}

		var currentUser models.User
		if err := db.WithContext(ctx).First(&currentUser, currentUserID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				response.JSONError(w, http.StatusUnauthorized, "current user not found", nil)
				return
			}
			response.JSONError(w, http.StatusInternalServerError, "failed to load current user", nil)
			return
		}

		var targetUser models.User
		if err := db.WithContext(ctx).First(&targetUser, targetID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				response.JSONError(w, http.StatusNotFound, "target user not found", nil)
				return
			}
			response.JSONError(w, http.StatusInternalServerError, "failed to load target user", nil)
			return
		}

		if currentUser.Level <= targetUser.Level {
			response.JSONError(w, http.StatusForbidden, "you cannot delete this user", nil)
			return
		}

		if err := db.WithContext(ctx).Delete(&targetUser).Error; err != nil {
			response.JSONError(w, http.StatusInternalServerError, "failed to delete user", nil)
			return
		}

		response.JSONSuccess(w, map[string]any{
			"deleted_user_id": targetUser.ID,
			"by_user_id":      currentUser.ID,
		}, nil, nil)
	}
}

func GetBalance(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var balance models.Balance

		userID := middleware.GetUserID(r.Context())

		fmt.Println("pass")

		err := db.Where("user_id = ?", userID).First(&balance).Error
		if err != nil {
			response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		response.JSONSuccess(w, balance, nil, nil)
	}
}

type AddBalanceRequest struct {
	Amount float64 `json:"amount"`
}

func AddBalance(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AddBalanceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.JSONError(w, http.StatusBadRequest, "Invalid request body", nil)
			return
		}

		if req.Amount <= 0 {
			response.JSONError(w, http.StatusBadRequest, "Amount must be greater than zero", nil)
			return
		}

		userID := middleware.GetUserID(r.Context())

		var balance models.Balance
		if err := db.Where("user_id = ?", userID).First(&balance).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				balance = models.Balance{
					UserID: uint(userID),
					Amount: req.Amount,
				}
				if err := db.Create(&balance).Error; err != nil {
					response.JSONError(w, http.StatusInternalServerError, "Failed to create balance", nil)
					return
				}
			} else {
				response.JSONError(w, http.StatusInternalServerError, "Failed to get balance", nil)
				return
			}
		} else {
			balance.Amount += req.Amount
			if err := db.Save(&balance).Error; err != nil {
				response.JSONError(w, http.StatusInternalServerError, "Failed to update balance", nil)
				return
			}
		}

		var user models.User
		if err := db.First(&user, userID).Error; err == nil {
			if balance.Amount > 3000000 {
				user.Level = 4
			} else if balance.Amount > 2000000 {
				user.Level = 3
			} else if balance.Amount > 1000000 {
				user.Level = 2
			} else {
				user.Level = 1
			}
			db.Save(&user)
		}

		response.JSONSuccess(w, map[string]any{
			"new_balance": balance.Amount,
			"level":       user.Level,
		}, nil, nil)
	}
}

func InitializeBalance(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var balance models.Balance

		err := db.First(&balance).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				balance = models.Balance{Amount: 100000}
				if err := db.Create(&balance).Error; err != nil {
					response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
					return
				}
				response.JSONSuccess(w, "Balance initialized", nil, nil)
				return
			}

			response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		balance.Amount = 100000
		if err := db.Save(&balance).Error; err != nil {
			response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		response.JSONSuccess(w, "Balance reset", nil, nil)
	}
}

func PostTakeBalance(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var balance models.Balance

		var req request.TakeRequest
		req.AssignFormValues(r)

		userID := middleware.GetUserID(r.Context())

		err := db.Where("user_id = ?", userID).First(&balance).Error
		if err != nil {
			response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		balance.Amount -= *req.Amount
		if err := db.Save(&balance).Error; err != nil {
			response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		response.JSONSuccess(w, balance, nil, nil)
	}
}
