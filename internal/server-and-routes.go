package internal

import (
	"encoding/json"
	"github.com/recoilme/pudge"
	"html/template"
	"log"
	"net/http"

	"collator-module/internal/session"
)

var (
	tmpl = template.Must(template.ParseFiles("./ui/templates/login.html", "./ui/templates/admin.html", "./ui/templates/scan.html"))
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
	}

	if userSession.ApiKey == "admin" {
		data, err := FormData()
		if err != nil {
			responseInternalError(writer, err)
		}

		err = tmpl.ExecuteTemplate(writer, "admin.html", data)
		if err != nil {
			responseInternalError(writer, err)
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
	}

	if userSession.ApiKey == "user" {
		err := tmpl.ExecuteTemplate(writer, "scan.html", userSession.CurrentUser)
		if err != nil {
			responseInternalError(writer, err)
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

func handlerReceiveReport(writer http.ResponseWriter, request *http.Request){
	if request.Method == "POST" {
		err := jsonToReportObject(request)
		if err != nil {
			responseInternalError(writer, err)
		}
	}
}

func handlerAdminDeleteUser(writer http.ResponseWriter, request *http.Request){
	if request.Method == "POST" {
		decoder := json.NewDecoder(request.Body)
		login := ""
		decoder.Decode(&login)
		err := DeleteUser(login)
		if err != nil {
			responseInternalError(writer, err)
		}
		http.Redirect(writer, request, "/admin", http.StatusFound)
	}
}

func handlerAdminDeleteReport(writer http.ResponseWriter, request *http.Request){
	if request.Method == "POST" {
		decoder := json.NewDecoder(request.Body)
		orderNumber := ""
		decoder.Decode(&orderNumber)
		 err := DeleteReport(orderNumber)
		 if err != nil {
		 	responseInternalError(writer, err)
		 }
		http.Redirect(writer, request, "/admin", http.StatusFound)
	}
}

func handlerAdminNewUser(writer http.ResponseWriter, request *http.Request) {
	// trying to load a session from user's request
	var userSession UserSession
	err := session.Load(&userSession, writer, request)
	if err != nil {
		responseInternalError(writer, err)
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

func createFirstAdmin() {
	defer closeAllDB()
	u := User{
		Login:    "admin",
		Password: "admin",
		Admin:    true,
	}
	err := pudge.Set("./db/users", u.Login, u)
	if err != nil {
		// do something
	}
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
	http.HandleFunc("/admin/delete_report", handlerAdminDeleteReport)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("ui/"))))
	log.Fatal(http.ListenAndServe(":9090", nil))
}
