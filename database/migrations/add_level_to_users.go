package migrations

import(
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddLevelToUsers() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "addleveltousers",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`ALTER TABLE users ADD COLUMN IF NOT EXISTS level INT NOT NULL DEFAULT 1`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec(`ALTER TABLE users DROP COLUMN IF EXISTS level`).Error
		},
	}
}
