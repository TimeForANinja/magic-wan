package gui

import (
	"embed"
	"html/template"
	"magic-wan/rest/shared"
	"net/http"
)

//go:embed login_page.html
var loginPageFS embed.FS

func LoginHandler(w http.ResponseWriter, _ *http.Request, _ *shared.User) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	tmpl := template.Must(template.ParseFS(loginPageFS, "login_page.html"))
	data := struct {
	}{}

	_ = tmpl.Execute(w, data)
}
