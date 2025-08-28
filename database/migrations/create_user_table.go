package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateUserTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "createusertable",
		Migrate: func(d *gorm.DB) error {
			return d.Exec(`
				CREATE TABLE IF NOT EXISTS users (
					id SERIAL PRIMARY KEY,
					username VARCHAR(255) NOT NULL,
					password VARCHAR(255) NOT NULL,
					created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
					updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
					deleted_at TIMESTAMPTZ
				)
			`).Error
		},
		Rollback: func(d *gorm.DB) error {
			return d.Exec("DROP TABLE IF EXISTS users").Error
		},
	}
}
