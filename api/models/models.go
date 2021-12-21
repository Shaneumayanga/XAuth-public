package models

type User struct {
	ID       string
	Email    string
	Name     string
	Password string
}

func (user *User) ClearPassword() {
	user.Password = ""
}

//App has a clientID and a clientSecret
type App struct {
	ID             string //generated
	UserId         string //getting from the session
	Clientid       string //generated
	Clientsecret   string //generated
	Appname        string //user provided
	Appdescription string //user provided
	CallbackURL    string //user provided
}

type Code struct {
	Code       string
	Userid     string
	Created_at string
}

type UserWithApps struct {
	ID             string
	Email          string
	Name           string
	Password       string
	UserId         string
	Clientid       string
	Clientsecret   string
	Appname        string
	Appdescription string
	CallbackURL    string
}

func (user *UserWithApps) ClearPassword() {
	user.Password = ""
}
