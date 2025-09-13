package repositories

import (
	"errors"
	"fmt"

	"github.com/Shaneumayanga/XAuth/api/models"
	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	SaveUser(user models.User) error
	GetUserByEmail(email string) *models.User
	GetUserByUserId(userid string) *models.User
	GetUserWithApps(userid string) []models.UserWithApps
}

type userRepo struct {
	db *sqlx.DB
}

var (
	insertuser        = `INSERT INTO Users VALUES ($1 , $2 , $3 , $4);`
	queryuserbyemail  = `SELECT * FROM Users WHERE email = $1;`
	queryuserbyuserid = `SELECT * FROM Users WHERE id = $1;`
	getuserwithapps   = `SELECT * FROM Apps INNER JOIN Users on Users.id = Apps.userid  WHERE Users.id = $1;`
)

func NewUserRepo(db *sqlx.DB) *userRepo {
	return &userRepo{
		db: db,
	}
}

func (ur *userRepo) SaveUser(user models.User) error {
	result, err := ur.db.Exec(insertuser, user.ID, user.Email, user.Name, user.Password)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
			return errors.New("email already exists try again")
		} else {
			fmt.Printf("err.Error(): %v\n", err.Error())
			return errors.New("unknown error occured please try again later")
		}
	}
	fmt.Println(result.RowsAffected())
	return nil
}

func (ur *userRepo) GetUserByEmail(email string) *models.User {
	result, err := ur.db.Queryx(queryuserbyemail, email)
	if err != nil {
		panic(err)
	}

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
		panic(err)
	}
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
