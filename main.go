package main

import (
	"fmt"
	"net/http"

	"github.com/RivaldoHardiantoWibowo/backend-test/database"
	"github.com/RivaldoHardiantoWibowo/backend-test/handlers"
	"github.com/RivaldoHardiantoWibowo/backend-test/middleware"
	"github.com/RivaldoHardiantoWibowo/backend-test/response"
	"github.com/RivaldoHardiantoWibowo/backend-test/routes"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	db := database.Connect()

	// r.Get("/", routes.GetBalance(db))
	// r.Get("/init-balance", routes.InitializeBalance(db))
	// r.Post("/take-balance", routes.PostTakeBalance(db))

	r.Post("/login", handlers.Login(db))
	r.Post("/register", handlers.Register(db))
	r.Get("/list-user", routes.ListUser(db))

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Get("/secure-route", func(w http.ResponseWriter, r *http.Request) {
			userID := middleware.GetUserID(r.Context())
			construct := map[string]any{
				"user_id": userID,
			}
			response.JSONSuccess(w, construct, nil, nil)
		})

		r.Post("/take-balance", routes.PostTakeBalance(db))
		r.Get("/balance", routes.GetBalance(db))
		r.Put("/add-balance", routes.AddBalance(db))
		r.Delete("/delete-user/{id}", routes.DeleteUser(db))
	})

	fmt.Println("Server running on port 3000")
	http.ListenAndServe(":3000", r)
}
