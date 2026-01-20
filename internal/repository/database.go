package repository

import (
	"fmt"
	"globe/internal/repository/entities"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	// load .env variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Couldn't load .env veriables %v\n", err)
	}
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s port=%s sslmode=disable dbname=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// connect to database
	var DB *gorm.DB
	var err error
	if DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		log.Fatalf("Couldn't start database %v\n", err)
	}

	// create trigram extention
	DB.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm;")

	// create index from product name
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_products_name_trgm ON products USING gin (name gin_trgm_ops);")

	// migrations
	if err := DB.AutoMigrate(
		&entities.User{},
		&entities.UnverifiedUser{},
		&entities.RefreshToken{},
		&entities.Product{},
	); err != nil {
		log.Fatalf("Couldn't migrate tables %v\n", err)
	}

	// sequence for unverified users
	DB.Exec("ALTER SEQUENCE unverified_users_id_seq "+
	"INCREMENT BY 1 "+
	"START WITH 100001 "+
	"CYCLE "+
	"CACHE 20")

	// sequence for users
	DB.Exec("ALTER SEQUENCE users_id_seq "+
	"INCREMENT BY 1 "+
	"START WITH 100001 "+
	"CYCLE "+
	"CACHE 20")

	// sequence for refresh tokens
	DB.Exec("ALTER SEQUENCE refresh_tokens_id_seq "+
	"INCREMENT BY 1 "+
	"START WITH 100001 "+
	"CYCLE "+
	"CACHE 20")

	// sequence for products
	DB.Exec("ALTER SEQUENCE products_id_seq "+
	"INCREMENT BY 1 "+
	"START WITH 100001 "+
	"CYCLE "+
	"CACHE 20")
	return DB
}