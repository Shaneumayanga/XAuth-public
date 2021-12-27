package oauth

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/Shaneumayanga/XAuth/api/models"
	"github.com/Shaneumayanga/XAuth/api/repositories"
	"github.com/Shaneumayanga/XAuth/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type OAuthHandler struct {
	userRepo repositories.UserRepo
	appRepo  repositories.AppRepo
	codeRepo repositories.CodeRepo
}

func NewOAuthHandler(userRepo repositories.UserRepo, appRepo repositories.AppRepo, codeRepo repositories.CodeRepo) *OAuthHandler {
	return &OAuthHandler{
		userRepo: userRepo,
		appRepo:  appRepo,
		codeRepo: codeRepo,
	}
}

var (
	authorizationLoginTemplate    = template.Must(template.ParseFiles("././static/oauth/authorize_login.html"))
	authorizationRegisterTemplate = template.Must(template.ParseFiles("././static/oauth/authorize_register.html"))
)

type AuthorizeLoginRequest struct {
	Email    string
	Password string
}

type AuthorizeRegisterRequest struct {
	Email    string
	Name     string
	Password string
}

//TODO  : scopes with the code :-))
func (oh *OAuthHandler) AuthorizeLogin(c *gin.Context) {
	authorizationReq := new(AuthorizeLoginRequest)
	if err := json.NewDecoder(c.Request.Body).Decode(authorizationReq); err != nil {
		c.JSON(http.StatusBadRequest, "Bad request")
		return
	}
	//sanitize and validate data after decode
	if !utils.IsValidvalid(authorizationReq.Email) {
		c.JSON(http.StatusBadGateway, "Bad Request")
		return
	}
	client_id := c.Request.URL.Query().Get("client_id")
	redirect_url := c.Request.URL.Query().Get("redirect_url")
	response_type := c.Request.URL.Query().Get("response_type")
	if client_id == "" || redirect_url == "" || response_type != "code" {
		c.JSON(200, map[string]string{
			"redirect_url": "",
			"error":        "INCORRECT_CONFIGURATION",
		})
		return
	}
	app := oh.appRepo.GetAppByClientId(client_id)
	if app == nil {
		c.JSON(200, map[string]string{
			"redirect_url": "",
			"error":        "INVALID_APP_CLIENT",
		})
		return
	}
	authorizedUser := oh.userRepo.GetUserByEmail(authorizationReq.Email)
	if authorizedUser == nil {
		c.JSON(200, map[string]string{
			"redirect_url": "",
			"error":        "INCORRECT_EMAIL_OR_PASSWORD",
		})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(authorizedUser.Password), []byte(authorizationReq.Password)); err != nil {
		c.JSON(200, map[string]string{
			"redirect_url": "",
			"error":        "INCORRECT_EMAIL_OR_PASSWORD",
		})
		return
	}
	if response_type == "code" {
		codeDB := new(models.Code)
		codeDB.Code = uuid.New().String()
		codeDB.Userid = authorizedUser.ID
		//the PK of the code table
		code := oh.codeRepo.SaveCode(codeDB)
		//Redirect to unknown error occured
		if code == "" {
			fmt.Println("Code was not saved in DB")
			return
		}

		//The JWT generated with the pk of the code table table := [code, userid]
		c.JSON(200, map[string]string{
			"redirect_url": redirect_url + "?code=" + GenerateJWTWithCode(code),
			"error":        "",
		})
	}
}

func (oh *OAuthHandler) AuthorizeRegister(c *gin.Context) {
	registerRequest := new(AuthorizeRegisterRequest)

	if err := json.NewDecoder(c.Request.Body).Decode(registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, "Bad request")
		return
	}

	if !utils.IsValidvalid(registerRequest.Email) {
		c.JSON(http.StatusBadGateway, "Bad Request")
		return
	}

	client_id := c.Request.URL.Query().Get("client_id")
	redirect_url := c.Request.URL.Query().Get("redirect_url")
	response_type := c.Request.URL.Query().Get("response_type")
	if client_id == "" || redirect_url == "" || response_type == "" || response_type != "code" {
		c.JSON(200, map[string]string{
			"redirect_url": "",
			"error":        "INCORRECT_CONFIGURATION",
		})
		return
	}
	app := oh.appRepo.GetAppByClientId(client_id)
	if app == nil {
		c.JSON(200, map[string]string{
			"redirect_url": "",
			"error":        "INVALID_APP_CLIENT",
		})
		return
	}
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(c.Request.FormValue("password")), 12)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "InternalServerError")
		return
	}
	userId := uuid.New().String()
	user := new(models.User)
	user.ID = userId
	user.Password = string(passwordBytes)
	user.Email = registerRequest.Email
	user.Name = registerRequest.Name
	if err := oh.userRepo.SaveUser(*user); err != nil {
		c.JSON(http.StatusInternalServerError, "InternalServerError")
		return
	}
	if response_type == "code" {
		codeDB := new(models.Code)
		codeDB.Userid = userId
		codeDB.Code = uuid.New().String()
		//the PK of the code table
		code := oh.codeRepo.SaveCode(codeDB)
		//Redirect to unknown error occured
		if code == "" {
			fmt.Println("Code was not saved in DB")
			return
		}
		//The JWT generated with the pk of the code table table := [code, userid]
		c.JSON(200, map[string]string{
			"redirect_url": redirect_url + "?code=" + GenerateJWTWithCode(code),
			"error":        "",
		})
	}

}
