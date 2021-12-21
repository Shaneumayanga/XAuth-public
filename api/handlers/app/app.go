package apphandlers

import (
	"html/template"
	"net/http"

	"github.com/Shaneumayanga/XAuth/api/models"
	"github.com/Shaneumayanga/XAuth/api/repositories"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sethvargo/go-password/password"
)

type Apphandlers struct {
	apprepo repositories.AppRepo
}

func NewAppHandler(apprepo repositories.AppRepo) *Apphandlers {
	return &Apphandlers{
		apprepo: apprepo,
	}
}

var (
	registerappTemplate = template.Must(template.ParseFiles("././static/app/register_app.html"))
)

func GenerateClientsecret() string {
	password, err := password.Generate(20, 5, 0, false, false)
	if err != nil {
		panic(err)
	}
	return password
}

func (ah *Apphandlers) GetRegisterApp(c *gin.Context) {
	registerappTemplate.Execute(c.Writer, nil)
}

func (ah *Apphandlers) RegisterApp(c *gin.Context) {
	currentUser := c.Request.Context().Value("user").(*models.User)
	app := new(models.App)
	app.UserId = currentUser.ID
	app.ID = uuid.New().String()
	app.Appname = c.Request.FormValue("appname")
	app.Appdescription = c.Request.FormValue("appdescription")
	app.CallbackURL = c.Request.FormValue("callbackurl")
	app.Clientid = uuid.New().String()
	app.Clientsecret = GenerateClientsecret()
	if err := ah.apprepo.SaveApp(app); err != nil {
		c.JSON(http.StatusConflict, err.Error())
		return
	}
	http.Redirect(c.Writer, c.Request, "/", http.StatusSeeOther)

}
