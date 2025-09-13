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
		CREATE TABLE IF NOT EXISTS Users (
		id VARCHAR(255) PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		name VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS Apps (
		id VARCHAR(255) PRIMARY KEY,
		userid VARCHAR(255),
		appname VARCHAR(255) UNIQUE NOT NULL,
		appdescription VARCHAR(255) NOT NULL, 
		callbackurl VARCHAR(255) NOT NULL,
		clientid VARCHAR(255) UNIQUE NOT NULL,
		clientsecret VARCHAR(255) NOT NULL,
		CONSTRAINT fk_user FOREIGN KEY (userid) REFERENCES Users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS Codes (
		code VARCHAR(255) PRIMARY KEY,
		userid VARCHAR(255) NOT NULL,
		CONSTRAINT fk_code_user FOREIGN KEY (userid) REFERENCES Users(id) ON DELETE CASCADE
	);


	`
	db.MustExec(sqlMake)
	fmt.Println("Created tables successfully")
}
