package api

import (
	"embed"
	"fmt"
	"html/template"
	"magic-wan/pkg/frr"
	"magic-wan/pkg/wg"
	"magic-wan/rest/shared"
	"net/http"
)

//go:embed debug_page.html
var debugPageFS embed.FS

func DebugV1Handler(w http.ResponseWriter, _ *http.Request, _ *shared.User) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// collect data
	frrContent, err := frr.Debug()
	if err != nil {
		frrContent = fmt.Sprintf("Error: %v", err)
	}

	wgContent, err := wg.Debug()
	if err != nil {
		wgContent = fmt.Sprintf("Error: %v", err)
	}

	tmpl := template.Must(template.ParseFS(debugPageFS, "debug_page.html"))
	data := struct {
		FrrContent string
		WgContent  string
	}{
		FrrContent: frrContent,
		WgContent:  wgContent,
	}

	_ = tmpl.Execute(w, data)
}
