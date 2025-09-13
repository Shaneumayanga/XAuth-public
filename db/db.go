package db

import (
	"fmt"
	"os"
	"time"

	"github.com/Shaneumayanga/XAuth/migrations"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DBConfig struct {
	DBUser     string
	DBName     string
	DBPassword string
	Host       string
	Port       string
	SSLmode    string
}

func GetDB() *sqlx.DB {
	dbconfig := DBConfig{}
	if os.Getenv("PRODUCTION") == "TRUE" {
		dbconfig.DBUser = os.Getenv("POSTGRES_USER")
		dbconfig.DBPassword = os.Getenv("POSTGRES_PASSWORD")
		dbconfig.DBName = os.Getenv("POSTGRES_DBNAME")
		dbconfig.Port = os.Getenv("DATABASE_PORT")
		dbconfig.Host = os.Getenv("DATABASE_HOST")
		dbconfig.SSLmode = "require"
	} else {
		dbconfig.DBUser = "postgres"
		dbconfig.DBPassword = "mossmoss"
		dbconfig.DBName = "xauth"
		dbconfig.Host = "localhost"
		dbconfig.Port = "5432"
		dbconfig.SSLmode = "disable"
	}

	if os.Getenv("POSTGRES") == "TRUE" {

		connection := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=%s password=%s port=%s", dbconfig.Host, dbconfig.DBUser, dbconfig.DBName, dbconfig.SSLmode, dbconfig.DBPassword, dbconfig.Port)
		db, err := sqlx.Connect("postgres", connection)
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)
		if err != nil {
			panic(err)
		}
		migrations.Migrate(db)
		return db
	} else {

		dsn := os.Getenv("DSN_MYSQL")

		db, err := sqlx.Connect("mysql", dsn)
		if err != nil {
			panic(err)
		}
		db.SetMaxIdleConns(100)
		db.SetMaxOpenConns(100)
		db.SetConnMaxLifetime(1 * time.Minute)

		migrations.MigrateMYSQL(db)

		return db
	}

}
