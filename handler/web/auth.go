package web

import (
	"a21hc3NpZ25tZW50/client"
	"embed"
	"fmt"
	"net/http"
	"path"
	"text/template"
	"time"
)

type AuthWeb interface {
	Login(w http.ResponseWriter, r *http.Request)
	LoginProcess(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	RegisterProcess(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type authWeb struct {
	userClient client.UserClient
	embed      embed.FS
}

func NewAuthWeb(userClient client.UserClient, embed embed.FS) *authWeb {
	return &authWeb{userClient, embed}
}

func (a *authWeb) Login(w http.ResponseWriter, r *http.Request) {
	// TODO: answer here

	filepath := path.Join("views", "auth", "login.html")
	header := path.Join("views", "general", "header.html")

	t, err := template.ParseFS(a.embed, filepath, header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"Message": "",
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (a *authWeb) LoginProcess(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	userId, status, err := a.userClient.Login(email, password)

	if status == 200 {
		http.SetCookie(w, &http.Cookie{
			Name:     "user_id",
			Value:    fmt.Sprintf("%d", userId),
			Path:     "/",
			Expires:  time.Now().Add(1 * time.Hour),
			Domain:   "",
			HttpOnly: true,
		})

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	} else {
		filepath := path.Join("views", "auth", "login.html")
		header := path.Join("views", "general", "header.html")

		t, _ := template.ParseFS(a.embed, filepath, header)

		data := map[string]interface{}{
			"Message": err.Error(),
		}

		err = t.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func (a *authWeb) Register(w http.ResponseWriter, r *http.Request) {
	// TODO: answer here

	filepath := path.Join("views", "auth", "register.html")
	header := path.Join("views", "general", "header.html")

	t, err := template.ParseFS(a.embed, filepath, header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"Message": "",
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (a *authWeb) RegisterProcess(w http.ResponseWriter, r *http.Request) {
	fullname := r.FormValue("fullname")
	email := r.FormValue("email")
	password := r.FormValue("password")

	userId, status, err := a.userClient.Register(fullname, email, password)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	if status == 201 {
		http.SetCookie(w, &http.Cookie{
			Name:   "user_id",
			Value:  fmt.Sprintf("%d", userId),
			Path:   "/",
			MaxAge: 31536000,
			Domain: "",
		})

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	} else {
		filepath := path.Join("views", "auth", "register.html")
		header := path.Join("views", "general", "header.html")

		t, _ := template.ParseFS(a.embed, filepath, header)

		data := map[string]interface{}{
			"Message": err.Error(),
		}

		err = t.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (a *authWeb) Logout(w http.ResponseWriter, r *http.Request) {
	// TODO: answer here

	http.SetCookie(w, &http.Cookie{
		Name:   "user_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
		Domain: "",
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)

}
