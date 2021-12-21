//app is basically the client

package repositories

import (
	"errors"
	"fmt"

	"github.com/Shaneumayanga/XAuth/api/models"
	"github.com/jmoiron/sqlx"
)

type AppRepo interface {
	SaveApp(app *models.App) error
	GetAppByClientId(clientID string) *models.App
}

type appRepo struct {
	db *sqlx.DB
}

var (
	insertApp        = `INSERT INTO Apps VALUES ($1,$2,$3,$4,$5,$6,$7);`
	getappbyclientid = `SELECT * FROM Apps WHERE clientid = $1;`
	// addusertoapp     = `` TODO
)

func NewAppRepo(db *sqlx.DB) *appRepo {
	return &appRepo{
		db: db,
	}
}

func (ar *appRepo) SaveApp(app *models.App) error {
	result, err := ar.db.Exec(insertApp, app.ID, app.UserId, app.Appname, app.Appdescription, app.CallbackURL, app.Clientid, app.Clientsecret)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "apps_appname_key"` {
			return errors.New("app name already exists")
		} else {
			return errors.New("unknown error occured please try again")
		}
	}
	fmt.Println(result.LastInsertId())
	return nil
}

///Think about some better error handling over here :-)
func (ar *appRepo) GetAppByClientId(clientID string) *models.App {
	result, err := ar.db.Queryx(getappbyclientid, clientID)

	if err != nil {
		return nil
	}
	for result.Next() {
		app := new(models.App)
		result.StructScan(app)
		fmt.Printf("app.Appname: %v\n", app.Appname)
		return app
	}
	return nil
}
