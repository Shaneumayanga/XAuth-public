package api

import (
	"context"
	"html/template"
	"net/http"
	"time"

	"github.com/Shaneumayanga/XAuth/api/cleanup"
	apphandlers "github.com/Shaneumayanga/XAuth/api/handlers/app"
	"github.com/Shaneumayanga/XAuth/api/handlers/auth"
	"github.com/Shaneumayanga/XAuth/api/handlers/oauth"
	"github.com/Shaneumayanga/XAuth/api/models"
	"github.com/Shaneumayanga/XAuth/api/ratelimit"

	mysqlrepo "github.com/Shaneumayanga/XAuth/api/repositories/mysql"
	"github.com/Shaneumayanga/XAuth/api/session"
	"github.com/Shaneumayanga/XAuth/db"
	helmet "github.com/danielkov/gin-helmet"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var homeTemplate = template.Must(template.ParseFiles("././static/index.html"))
var myappstemplate = template.Must(template.ParseFiles("././static/app/my_apps.html"))

//error templates
var (
	invalidappclient = template.Must(template.ParseFiles("././static/app/invalid_app_client.html"))
)

func AuthMiddleware(session *session.Session) func(c *gin.Context) {
	return func(c *gin.Context) {
		userCtx := c.Request.Context().Value("user")
		user := session.GetUserBySession(c.Writer, c.Request)
		if user == nil {
			http.Redirect(c.Writer, c.Request, "/login", http.StatusSeeOther)
			c.Abort()
			return
		} else {
			if userCtx == nil {
				c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "user", user))
			}
			c.Next()
		}
	}
}

func AuthMiddlewareWithOutloginRedirect(session *session.Session) func(c *gin.Context) {
	return func(c *gin.Context) {
		userCtx := c.Request.Context().Value("user")
		user := session.GetUserBySession(c.Writer, c.Request)
		if userCtx == nil && user != nil {
			c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "user", user))
		}
		c.Next()
	}
}

func SetXPoweredBy(c *gin.Context) {
	c.Writer.Header().Set("X-Powered-By", "XAuthHTTP")
	c.Next()
}

func Start() *gin.Engine {
	handler := gin.New()
	handler.Use(helmet.Default())
	handler.Use(SetXPoweredBy)
	handler.Static("/assets", "././static/assets/")
	limiter := ratelimit.NewRateLimit(15)

	go limiter.Start()

	db := db.GetDB()
	go cleanup.Run(db)
	handler.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Authorization", "Origin", "Cookie"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	//repos
	userRepo := mysqlrepo.NewUserRepo(db)
	appRepo := mysqlrepo.NewAppRepo(db)
	codeRepo := mysqlrepo.NewCodeRepo(db)

	//sessions
	session := session.NewSession(userRepo)

	oauthhandler := oauth.NewOAuthHandler(userRepo, appRepo, codeRepo)
	authHandler := auth.NewAuthHandler(userRepo, session)
	apphandler := apphandlers.NewAppHandler(appRepo)

	//get the user from the session and load with the hometemplate :-)
	handler.GET("/", AuthMiddlewareWithOutloginRedirect(session), func(c *gin.Context) {
		user := c.Request.Context().Value("user")
		if user != nil {
			userF := user.(*models.User)
			userF.ClearPassword()
			homeTemplate.Execute(c.Writer, user)
		} else {
			homeTemplate.Execute(c.Writer, nil)
		}
	})

	handler.GET("/login", authHandler.GetLogin)
	handler.GET("/register", authHandler.GetRegister)
	handler.POST("/login", limiter.Use(), authHandler.Login)
	handler.POST("/register", limiter.Use(), authHandler.Register)

	handler.GET("/app/myapps", limiter.Use(), AuthMiddleware(session), func(c *gin.Context) {
		currentUser := c.Request.Context().Value("user").(*models.User)
		myappstemplate.Execute(c.Writer, map[string][]models.UserWithApps{
			"Apps": userRepo.GetUserWithApps(currentUser.ID),
		})
	})

	handler.GET("/app/register", AuthMiddleware(session), apphandler.GetRegisterApp)
	handler.POST("/app/register", limiter.Use(), AuthMiddleware(session), apphandler.RegisterApp)

	//OAuth routes
	handler.GET("/login/oauth/authorize", oauthhandler.GetAuthorization)
	handler.POST("/login/oauth/authorize", limiter.UseREST(), oauthhandler.AuthorizeLogin)

	handler.GET("/register/oauth/authorize", oauthhandler.GetAuthorizationRegister)
	handler.POST("/register/oauth/authorize", limiter.UseREST(), oauthhandler.AuthorizeRegister)

	handler.POST("/login/oauth/access_token", limiter.UseREST(), oauthhandler.GetAccesToken)
	handler.POST("/login/oauth/user", limiter.UseREST(), oauthhandler.VerifyAccessTokenAndGetUserInfo)

	//Error routes
	handler.GET("/app/invalidapp", func(c *gin.Context) {
		invalidappclient.Execute(c.Writer, nil)
	})

	handler.GET("/incorrectlogin", func(c *gin.Context) {
		c.JSON(200, "Incorrect username or password")
	})

	handler.GET("/error/unknownerror", func(c *gin.Context) {
		c.JSON(200, "Unknown error occured")
	})

	handler.GET("/error/incorrectconfiguration", func(c *gin.Context) {
		c.JSON(200, "Incorrect configuration")
	})

	handler.GET("/onblock", func(c *gin.Context) {
		c.JSON(200, "You were blocked for 5 mins or less, but thank you for caring about me :-))))")
	})

	return handler
}
