package internal

import (
	"encoding/json"
	"github.com/recoilme/pudge"
	"html/template"
	"path"
	"log"
	"net/http"
	"strconv"

	"collator-module/internal/session"
)

var (

	tmpl = template.Must(template.ParseFiles(path.Join(".", "ui", "templates", "login.html"),
		path.Join(".", "ui", "templates", "admin.html"),
		path.Join(".", "ui", "templates", "scan.html"),
		path.Join(".", "ui", "templates", "report.html")))
)

type UserSession struct {
	ApiKey      string `json:"apiKey"`
	CurrentUser string `json:"currentUser"`
}

// root will always redirect to login page. is it right or ?
func handlerRoot(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "/login", http.StatusFound)
}

func handlerLogin(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		return
	}
	err := tmpl.ExecuteTemplate(writer, "login.html", nil)
	if err != nil {
		responseInternalError(writer, err)
		return
	}
}

func handlerLoginCheck(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		http.Redirect(writer, request, "/login", http.StatusFound)
	} else {
		// trying to load a session from user's request
		var userSession UserSession
		err := session.Load(&userSession, writer, request)
		if err != nil {
			responseInternalError(writer, err)
			return
		}

		currentUser := request.FormValue("login")
		passw := request.FormValue("password")
		address, apiKey, err := authorizeUser(currentUser, passw)
		if err != nil {
			// incorrect login or password
			responseForbidden(writer, err)
			return
		}

		// seems good. we can save the session now
		userSession.CurrentUser = currentUser
		userSession.ApiKey = apiKey
		err = session.Save(userSession, writer, request)
		if err != nil {
			responseInternalError(writer, err)
			return
		}

		http.Redirect(writer, request, "/"+address, http.StatusFound)
	}
}

func handlerAdmin(writer http.ResponseWriter, request *http.Request) {
	// trying to load a session from user's request
	var userSession UserSession
	err := session.Load(&userSession, writer, request)
	if err != nil {
		responseInternalError(writer, err)
		return
	}

	if userSession.ApiKey == "admin" {
		data, err := FormData()
		if err != nil {
			responseInternalError(writer, err)
			return
		}

		err = tmpl.ExecuteTemplate(writer, "admin.html", data)
		if err != nil {
			responseInternalError(writer, err)
			return
		}
	} else {
		http.Redirect(writer, request, "/login", http.StatusFound)
	}
}

func handlerScan(writer http.ResponseWriter, request *http.Request) {
	// trying to load a session from user's request
	var userSession UserSession
	err := session.Load(&userSession, writer, request)
	if err != nil {
		responseInternalError(writer, err)
		return
	}

	if userSession.ApiKey == "user" {
		err := tmpl.ExecuteTemplate(writer, "scan.html", userSession.CurrentUser)
		if err != nil {
			responseInternalError(writer, err)
			return
		}
		return
	} else {
		http.Redirect(writer, request, "/login", http.StatusFound)
	}
}

func handlerLogout(writer http.ResponseWriter, request *http.Request) {
	// trying to load a session from user's request
	var userSession UserSession
	err := session.Load(&userSession, writer, request)
	if err != nil {
		responseInternalError(writer, err)
		return
	}
	userSession.ApiKey = ""
	userSession.CurrentUser = ""
	//  don't forget to save it!

	err = session.Save(userSession, writer, request)
	if err != nil {
		responseInternalError(writer, err)
		return
	}
	http.Redirect(writer, request, "login", http.StatusFound)
}

func handlerReceiveReport(writer http.ResponseWriter, request *http.Request) {
	var userSession UserSession
	err := session.Load(&userSession, writer, request)
	if err != nil {
		responseInternalError(writer, err)
		return
	}
	if request.Method == "POST" {
		if userSession.ApiKey == "user" {
			err := jsonToReportObject(request)
			if err != nil {
				responseInternalError(writer, err)
				return
			}
		}
	}
}

func handlerAdminDeleteUser(writer http.ResponseWriter, request *http.Request) {
	var userSession UserSession
	err := session.Load(&userSession, writer, request)
	if err != nil {
		responseInternalError(writer, err)
		return
	}
	if request.Method == "POST" {
		if userSession.ApiKey == "admin" {
			decoder := json.NewDecoder(request.Body)
			login := ""
			decoder.Decode(&login)
			err := DeleteUser(login)
			if err != nil {
				responseInternalError(writer, err)
				return
			}
			http.Redirect(writer, request, "/admin", http.StatusFound)
		}
	}
}

func handlerAdminGetReport(writer http.ResponseWriter, request *http.Request) {
	var userSession UserSession
	err := session.Load(&userSession, writer, request)
	if err != nil {
		responseInternalError(writer, err)
		return
	}
	if request.Method == "GET" {
		if userSession.ApiKey == "admin" {
			orderNumber := request.URL.Query().Get("order_number")
			report, err := getReportFromDB(orderNumber)
			if err != nil {
				responseInternalError(writer, err)
				return
			}
			type Data struct {
				Data    Report   `json:"data"`
				Headers []string `json:"headers"`
			}
			layers, err := strconv.Atoi(report.ScansAmount)
			if err != nil {
				responseInternalError(writer, err)
				return
			}
			tmp := make([]string, 0, layers+2)
			tmp = append(tmp, "№ п./п.")
			for i := 1; i < layers+1; i++ {
				tmp = append(tmp, strconv.Itoa(i))
			}
			tmp = append(tmp, "Статус")
			data := Data{
				Data:    report,
				Headers: tmp,
			}
			err = tmpl.ExecuteTemplate(writer, "report.html", data)
			if err != nil {
				responseInternalError(writer, err)
				return
			}
		}
	}
}

func handlerAdminDeleteReport(writer http.ResponseWriter, request *http.Request) {
	var userSession UserSession
	err := session.Load(&userSession, writer, request)
	if err != nil {
		responseInternalError(writer, err)
		return
	}
	if request.Method == "POST" {
		if userSession.ApiKey == "admin" {
			decoder := json.NewDecoder(request.Body)
			orderNumber := ""
			decoder.Decode(&orderNumber)
			err := DeleteReport(orderNumber)
			if err != nil {
				responseInternalError(writer, err)
				return
			}
			http.Redirect(writer, request, "/admin", http.StatusFound)
		}
	}
}

func handlerAdminNewUser(writer http.ResponseWriter, request *http.Request) {
	// trying to load a session from user's request
	var userSession UserSession
	err := session.Load(&userSession, writer, request)
	if err != nil {
		responseInternalError(writer, err)
		return
	}

	if request.Method == "POST" {
		if userSession.ApiKey == "admin" {
			var isAdmin bool
			if request.FormValue("admin") == "admin" {
				isAdmin = true
			} else {
				isAdmin = false
			}
			_, err := AddUser(request.FormValue("login"), request.FormValue("password"), isAdmin)
			if err != nil {
				responseInternalError(writer, err)
				return
			}
			http.Redirect(writer, request, "/admin", http.StatusFound)
		} else {
			http.Redirect(writer, request, "/admin", http.StatusFound)
		}
	}
}

func createFirstAdmin() error {
	defer closeAllDB()
	u := User{
		Login:    "admin",
		Password: "admin",
		Admin:    true,
	}
	err := pudge.Set(path.Join(".", "db", "users"), u.Login, u)
	if err != nil {
		return err
	}
	return nil
}

func ApplicationStart() {
	createFirstAdmin()
	http.HandleFunc("/", handlerRoot)
	http.HandleFunc("/scan", handlerScan)
	http.HandleFunc("/scan/send_report", handlerReceiveReport)
	http.HandleFunc("/login", handlerLogin)
	http.HandleFunc("/login/enter", handlerLoginCheck)
	http.HandleFunc("/logout", handlerLogout)
	http.HandleFunc("/admin", handlerAdmin)
	http.HandleFunc("/admin/new_user", handlerAdminNewUser)
	http.HandleFunc("/admin/delete_user", handlerAdminDeleteUser)
	http.HandleFunc("/admin/open_report/", handlerAdminGetReport)
	http.HandleFunc("/admin/delete_report", handlerAdminDeleteReport)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("ui/"))))
	log.Fatal(http.ListenAndServe(":9090", nil))
}
