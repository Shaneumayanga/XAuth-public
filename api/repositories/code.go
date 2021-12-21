package repositories

import (
	"fmt"

	"github.com/Shaneumayanga/XAuth/api/models"
	"github.com/jmoiron/sqlx"
)

type CodeRepo interface {
	SaveCode(code *models.Code) string
	GetCodeByCode(code string) *models.Code
	DeleteCodeByTable(code string)
}

type codeRepo struct {
	db *sqlx.DB
}

var (
	insertcode       = `INSERT INTO Codes VALUES ($1,$2);`
	getcodebycode    = `SELECT * FROM Codes WHERE code = $1;`
	deletecodebycode = `DELETE FROM Codes WHERE code= $1;`
)

func NewCodeRepo(db *sqlx.DB) *codeRepo {
	return &codeRepo{
		db: db,
	}
}

func (cr *codeRepo) SaveCode(code *models.Code) string {
	result, err := cr.db.Exec(insertcode, code.Code, code.Userid)
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		return ""
	}
	fmt.Println(result.RowsAffected())
	return code.Code
}

func (cr *codeRepo) GetCodeByCode(code string) *models.Code {
	result, err := cr.db.Queryx(getcodebycode, code)
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		return nil
	}
	for result.Next() {
		code := new(models.Code)
		result.StructScan(code)
		fmt.Printf("code.Code: %v\n", code.Code)
		fmt.Printf("code.Userid: %v\n", code.Userid)
		return code
	}
	return nil
}

func (cr *codeRepo) DeleteCodeByTable(code string) {
	result := cr.db.MustExec(deletecodebycode, code)
	fmt.Println(result.RowsAffected())
}
