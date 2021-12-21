package session

import (
	"net/http"
	"os"

	"github.com/Shaneumayanga/XAuth/api/models"
	"github.com/Shaneumayanga/XAuth/api/repositories"
	"github.com/gorilla/securecookie"
)

var hashKey = os.Getenv("HASH_KEY")
var blockKey = os.Getenv("BLOCK_KEY")
var sc = securecookie.New([]byte(hashKey), []byte(blockKey))

var cookieName = "session"

type Session struct {
	userRepo repositories.UserRepo
}

func NewSession(userRepo repositories.UserRepo) *Session {
	return &Session{
		userRepo: userRepo,
	}
}

func (s *Session) SaveSessionUser(rw http.ResponseWriter, r *http.Request, user *models.User) {
	values := map[string]string{
		"userID":    user.ID,
		"userEmail": user.Email,
	}
	encoded, err := sc.Encode(cookieName, values)

	if err == nil {
		cookie := &http.Cookie{
			Name:     cookieName,
			Value:    encoded,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			MaxAge:   60 * 60 * 72,
		}
		http.SetCookie(rw, cookie)
	} else {
		panic(err)
	}

}

func (s *Session) GetUserBySession(rw http.ResponseWriter, r *http.Request) *models.User {
	values := make(map[string]string)
	cookie, err := r.Cookie(cookieName)

	if err != nil {
		return nil
	}
	err = sc.Decode(cookieName, cookie.Value, &values)
	if err == nil {
		email := values["userEmail"]
		return s.userRepo.GetUserByEmail(email)
	}
	return nil
}
