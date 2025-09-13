package mysqlrepo

import (
	"errors"
	"fmt"

	"github.com/Shaneumayanga/XAuth/api/models"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

var (
	insertuser        = `INSERT INTO Users VALUES (? , ? , ? , ?);`
	queryuserbyemail  = `SELECT * FROM Users WHERE email = ?;`
	queryuserbyuserid = `SELECT * FROM Users WHERE id = ?;`
	getuserwithapps   = `SELECT * FROM Apps INNER JOIN Users on Users.id = Apps.userid  WHERE Users.id = ?;`
)

func NewUserRepo(db *sqlx.DB) *userRepo {
	return &userRepo{
		db: db,
	}
}

func (ur *userRepo) SaveUser(user models.User) error {
	result, err := ur.db.Exec(insertuser, user.ID, user.Email, user.Name, user.Password)
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		return errors.New("unknown error occured please try again later")
	}
	fmt.Println(result.RowsAffected())
	return nil
}

func (ur *userRepo) GetUserByEmail(email string) *models.User {
	result, err := ur.db.Queryx(queryuserbyemail, email)
	if err != nil {
		return nil
	}
	defer result.Close()

	for result.Next() {
		user := new(models.User)
		result.StructScan(user)
		fmt.Printf("user.Name: %v\n", user.Name)
		return user
	}
	return nil
}

func (ur *userRepo) GetUserByUserId(userid string) *models.User {
	result, err := ur.db.Queryx(queryuserbyuserid, userid)
	if err != nil {
		return nil
	}
	defer result.Close()
	for result.Next() {
		user := new(models.User)
		result.StructScan(user)
		fmt.Printf("user.Name: %v\n", user.Name)
		return user
	}

	return nil
}

func (ur *userRepo) GetUserWithApps(userid string) []models.UserWithApps {
	apps := make([]models.UserWithApps, 0)
	result, err := ur.db.Queryx(getuserwithapps, userid)
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		return nil
	}
	defer result.Close()
	for result.Next() {
		userwithApps := new(models.UserWithApps)
		result.StructScan(userwithApps)
		userwithApps.ClearPassword()
		apps = append(apps, *userwithApps)
	}
	if len(apps) == 0 {
		return nil
	}
	return apps
}
