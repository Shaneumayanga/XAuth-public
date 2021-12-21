package mysqlrepo

import (
	"fmt"

	"github.com/Shaneumayanga/XAuth/api/models"
	"github.com/jmoiron/sqlx"
)

type codeRepo struct {
	db *sqlx.DB
}

var (
	insertcode       = `INSERT INTO Codes VALUES (?,?);`
	getcodebycode    = `SELECT * FROM Codes WHERE code = ?;`
	deletecodebycode = `DELETE FROM Codes WHERE code= ?;`
)

func NewCodeRepo(db *sqlx.DB) *codeRepo {
	return &codeRepo{
		db: db,
	}
}
func (cr *codeRepo) SaveCode(code *models.Code) string {
	result, err := cr.db.Exec(insertcode, code.Code, code.Userid)
	if err != nil {
		fmt.Println("My sql!")
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

	defer result.Close()

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
