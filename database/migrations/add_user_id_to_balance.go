package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddUserIDToBalance() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "adduseridtobalance",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`
				ALTER TABLE balances ADD COLUMN user_id INT REFERENCES users(id)
			`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec("ALTER TABLE balances DROP COLUMN user_id").Error
		},
	}
}
