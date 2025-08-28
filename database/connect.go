package database

import (
	"fmt"
	"log"

	"github.com/Shercosta/digi-wallet/database/migrations"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect() *gorm.DB {
	dsn := "host=127.0.0.1 user=postgres password=postgres dbname=digidb port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	fmt.Println("connected to database")

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		migrations.CreateBalanceTable(),
		migrations.CreateUserTable(),
		migrations.AddUserIDToBalance(),
		migrations.AddLevelToUsers(),
	})

	if err := m.Migrate(); err != nil {
		log.Fatal("failed to migrate database", err)
	}

	log.Println("database migrated")

	return db
}
