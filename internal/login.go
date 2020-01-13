package internal

import "github.com/recoilme/pudge"

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Admin    bool   `json:"admin"`
}

// depends on user role. if admin - admin panel. otherwise - scan panel
func authorizeUser(login, pass string) (string, string, error) {
	defer closeAllDB()

	// authorization and authentication logic
	user := User{}
	err := pudge.Get("./db/users", login, &user)
	if err != nil {
		return "", "", err
	}

	var addr string
	var apiKey string

	if user.Password == pass {
		if user.Admin {
			apiKey = "admin"
			addr = "admin"
		} else {
			apiKey = "user"
			addr = "scan"
		}
	} else {
		addr = "login"
	}

	return addr, apiKey, nil
}
