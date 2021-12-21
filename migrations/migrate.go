package migrations

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func Migrate(db *sqlx.DB) {
	sqlMake := `
		CREATE TABLE IF NOT EXISTS Users(
			id varchar(255) PRIMARY KEY,
			email varchar(255) UNIQUE NOT NULL,
			name varchar(255) NOT NULL,
			password varchar(255) NOT NULL
		);

		CREATE TABLE IF NOT EXISTS Apps(
			id varchar(255) PRIMARY KEY,
			userid varchar(255) REFERENCES Users(id),
			appname varchar(255) UNIQUE NOT NULL,
			appdescription varchar(255) NOT NULL, 
			callbackurl varchar(255) NOT NULL,
			clientid varchar(255) UNIQUE NOT NULL,
			clientsecret varchar(255) NOT NULL
		);
		CREATE TABLE IF NOT EXISTS Codes(
			code varchar(255) PRIMARY KEY,
			userid varchar(255) NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		);
	`
	db.MustExec(sqlMake)
	fmt.Println("Created tables successfully")
}

func MigrateMYSQL(db *sqlx.DB) {
	sqlMake := `
	CREATE TABLE IF NOT EXISTS Users(
		id varchar(255),
		email varchar(255) UNIQUE NOT NULL,
		name varchar(255) NOT NULL,
		password varchar(255) NOT NULL,
		PRIMARY KEY (id)
	);
	CREATE TABLE IF NOT EXISTS Apps(
		id varchar(255),
		userid varchar(255),
		appname varchar(255) UNIQUE NOT NULL,
		appdescription varchar(255) NOT NULL, 
		callbackurl varchar(255) NOT NULL,
		clientid varchar(255) UNIQUE NOT NULL,
		clientsecret varchar(255) NOT NULL,
		PRIMARY KEY (id),
		FOREIGN KEY (userid) REFERENCES Users(id)
	);
	CREATE TABLE IF NOT EXISTS Codes(
		code varchar(255),
		userid varchar(255) NOT NULL,
		PRIMARY KEY (code)
	);

	`
	db.MustExec(sqlMake)
	fmt.Println("Created tables successfully")
}
