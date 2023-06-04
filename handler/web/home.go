package web

import (
	"embed"
	"net/http"
	"path"
	"text/template"
)

type HomeWeb interface {
	Index(w http.ResponseWriter, r *http.Request)
}

type homeWeb struct {
	embed embed.FS
}

func NewHomeWeb(embed embed.FS) *homeWeb {
	return &homeWeb{embed}
}

func (h *homeWeb) Index(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("user_id")
	userId := true
	if err != nil {
		userId = false
	} else if c.Value == "" {
		userId = false
	}
	var filepath = path.Join("views", "main", "index.html")
	var header = path.Join("views", "general", "header.html")

	var data = map[string]interface{}{
		"userId": userId,
	}
	var tmpl = template.Must(template.ParseFS(h.embed, filepath, header))

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
