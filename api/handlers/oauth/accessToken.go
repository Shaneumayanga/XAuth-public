package oauth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (oh *OAuthHandler) GetAccesToken(c *gin.Context) {
	client_id := c.Request.URL.Query().Get("client_id")
	client_secret := c.Request.URL.Query().Get("client_secret")
	codeJWT := c.Request.URL.Query().Get("code")
	clientDB := oh.appRepo.GetAppByClientId(client_id)
	if clientDB == nil {
		c.JSON(http.StatusUnauthorized, map[string]string{
			"Error": "Invalid client secret or client id",
		})
		return
	}
	if client_secret != clientDB.Clientsecret {
		c.JSON(http.StatusUnauthorized, map[string]string{
			"Error": "Invalid client secret or client id",
		})
		return
	}
	//the pk of the codes table
	code := ValidateTokenAndGetCode(codeJWT)
	if code == "" {
		c.JSON(http.StatusUnauthorized, map[string]string{
			"Error": "Invalid code",
		})
		return
	}
	if codeDB := oh.codeRepo.GetCodeByCode(code); codeDB != nil {
		accessToken := GenerateJWTAccessTokenWithUserId(codeDB.Userid)
		oh.codeRepo.DeleteCodeByTable(codeDB.Code)
		//delete the code from the db
		c.JSON(200, map[string]string{
			"accessToken": accessToken,
		})
		return
	}
	c.JSON(http.StatusUnauthorized, map[string]string{
		"Error": "code expired",
	})
}

func (oh *OAuthHandler) VerifyAccessTokenAndGetUserInfo(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	access_token := strings.Split(authHeader, " ")[1]

	userID := ValidateTokenAndGetUserID(access_token)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, "AccessToken invalid or expired")
		return
	}
	userDB := oh.userRepo.GetUserByUserId(userID)
	if userDB == nil {
		c.JSON(404, "User not found")
		return
	}
	userDB.ClearPassword()
	c.JSON(200, userDB)
}

//Get handlers
func (oh *OAuthHandler) GetAuthorization(c *gin.Context) {
	authorizationLoginTemplate.Execute(c.Writer, nil)
}

func (oh *OAuthHandler) GetAuthorizationRegister(c *gin.Context) {
	authorizationRegisterTemplate.Execute(c.Writer, nil)
}
