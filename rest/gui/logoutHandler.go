package gui

import (
	"embed"
	"html/template"
	"magic-wan/rest/shared"
	"net/http"
)

//go:embed logout_page.html
var logoutPageFS embed.FS

func LogoutHandler(w http.ResponseWriter, _ *http.Request, user *shared.User) {
	http.SetCookie(w, &http.Cookie{Name: "Authorization", Value: "", MaxAge: -1, Path: "/"})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	tmpl := template.Must(template.ParseFS(logoutPageFS, "logout_page.html"))
	data := struct {
		User string
	}{
		User: user.Name,
	}

	_ = tmpl.Execute(w, data)
}
