package auth

import (
	"html/template"
	"net/http"

	"github.com/Shaneumayanga/XAuth/api/models"
	"github.com/Shaneumayanga/XAuth/api/repositories"
	"github.com/Shaneumayanga/XAuth/api/session"
	"github.com/Shaneumayanga/XAuth/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo repositories.UserRepo
	session  *session.Session
}

var (
	loginTemplate    = template.Must(template.ParseFiles("././static/auth/login.html"))
	registerTemplate = template.Must(template.ParseFiles("././static/auth/register.html"))
)

func NewAuthHandler(userRepo repositories.UserRepo, session *session.Session) *AuthHandler {
	return &AuthHandler{
		userRepo: userRepo,
		session:  session,
	}
}

func (ah *AuthHandler) Login(c *gin.Context) {
	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")

	if !utils.IsValidvalid(email) {
		c.JSON(http.StatusBadRequest, "Bad request")
		return
	}
	user := ah.userRepo.GetUserByEmail(email)
	if user == nil {
		c.JSON(http.StatusUnauthorized, "Incorrect Email or password")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(http.StatusUnauthorized, "Incorrect Email or password")
		return
	}
	//saves session here
	ah.session.SaveSessionUser(c.Writer, c.Request, user)
	http.Redirect(c.Writer, c.Request, "/", http.StatusSeeOther)
}

func (ah *AuthHandler) Register(c *gin.Context) {
	user := new(models.User)
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(c.Request.FormValue("password")), 12)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "InternalServerError")
		return
	}
	user.ID = uuid.New().String()
	user.Email = c.Request.FormValue("email")
	user.Name = c.Request.FormValue("name")
	user.Password = string(passwordBytes)
	if !utils.IsValidvalid(user.Email) {
		c.JSON(http.StatusBadRequest, "Bad request")
		return
	}
	if err := ah.userRepo.SaveUser(*user); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	//save session here
	ah.session.SaveSessionUser(c.Writer, c.Request, user)

	http.Redirect(c.Writer, c.Request, "/", http.StatusSeeOther)

}

func (ah *AuthHandler) GetLogin(c *gin.Context) {
	loginTemplate.Execute(c.Writer, nil)
}

func (ah *AuthHandler) GetRegister(c *gin.Context) {
	registerTemplate.Execute(c.Writer, nil)
}
