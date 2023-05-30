package web

import (
	"a21hc3NpZ25tZW50/client"
	"embed"
	"fmt"
	"log"
	"net/http"
	"path"
	"text/template"
	"time"
)

type DashboardWeb interface {
	Dashboard(w http.ResponseWriter, r *http.Request)
}

type dashboardWeb struct {
	categoryClient client.CategoryClient
	embed          embed.FS
}

func NewDashboardWeb(catClient client.CategoryClient, embed embed.FS) *dashboardWeb {
	return &dashboardWeb{catClient, embed}
}

func (d *dashboardWeb) Dashboard(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id")

	categories, err := d.categoryClient.GetCategories(userId.(string))
	if err != nil {
		log.Println("error get cat: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var dataTemplate = map[string]interface{}{
		"categories": categories,
	}

	var funcMap = template.FuncMap{
		"categoryInc": func(catId int) int {
			return catId + 1
		},
		"categoryDec": func(catId int) int {
			return catId - 1
		},
		"CorrectTime": func(t time.Time) string {
			t = t.UTC()

			if t.IsZero() {
				return ""
			}

			return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
				t.Year(), t.Month(), t.Day(),
				t.Hour(), t.Minute(), t.Second())
		},
		"Reminder": func(t time.Time) bool {
			return !t.IsZero()
		},
	}

	// ignore this
	_ = dataTemplate
	_ = funcMap
	//

	// TODO: answer here
	var header = path.Join("views", "general", "header.html")
	var dashboard = path.Join("views", "main", "dashboard.html")
	tmpl, err := template.New("dashboard.html").Funcs(funcMap).ParseFS(d.embed, dashboard, header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, dataTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
