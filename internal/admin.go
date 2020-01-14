package internal

import (
	"github.com/recoilme/pudge"
	"path"
)

type AdminData struct {
	Users   []User   `json:"users"`
	Reports []Report `json:"reports"`
}

func getUsers() ([]User, error) {
	defer closeAllDB()
	keys, err := pudge.Keys(path.Join(".", "db", "users"), 0, 0, 0, true)
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(keys))
	for _, key := range keys {
		var u User
		err := pudge.Get(path.Join(".", "db", "users"), key, &u)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

func getReports() ([]Report, error) {
	defer closeAllDB()
	keys, err := pudge.Keys(path.Join(".", "db", "reports"), 0, 0, 0, true)
	if err != nil {
		return nil, err
	}
	reports := make([]Report, 0, len(keys))
	for _, key := range keys {
		var r Report
		err := pudge.Get(path.Join(".", "db", "reports"), key, &r)
		if err != nil {
			return nil, err
		}
		reports = append(reports, r)
	}
	return reports, nil
}

func FormData() (AdminData, error) {
	users, err := getUsers()
	if err != nil {
		return AdminData{}, err
	}
	reports, err := getReports()
	if err != nil {
		return AdminData{}, err
	}
	return AdminData{
		Users:   users,
		Reports: reports,
	}, nil
}

func AddUser(login, pass string, adm bool) (AdminData, error) {
	defer closeAllDB()
	u := &User{
		Login:    login,
		Password: pass,
		Admin:    adm,
	}
	err := pudge.Set(path.Join(".", "db", "users"), u.Login, u)
	if err != nil {
		return AdminData{}, err
	}

	return FormData()
}

func DeleteUser(login string) error {
	defer pudge.CloseAll()
	err := pudge.Delete(path.Join(".", "db", "users"), login)
	if err != nil {
		return err
	}
	return nil
}
