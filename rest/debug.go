package rest

import (
	"fmt"
	"html/template"
	"magic-wan/pkg/frr"
	"magic-wan/pkg/wg"
	"net/http"
	"time"
)

func StartRest() error {
	http.HandleFunc("/api/v1/debug", DebugV1Handler)
	port := 80

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadHeaderTimeout: 3 * time.Second,
	}

	fmt.Printf("Starting server on :%d...\n", port)
	err := server.ListenAndServe()
	return err
}

const debugPageTemplate = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Debug Information</title>
	<style>
		body {
			display: flex;
		}
		.panel {
			flex: 1;
			border: 1px solid #ccc;
			padding: 10px;
			overflow: auto;
		    white-space: pre-wrap;
		}
		#leftPanel {
			background-color: #f9f9f9;
		}
		#rightPanel {
			background-color: #f1f1f1;
		}
	</style>
	</head>
	<body>
		<div id="leftPanel" class="panel">{{.FrrContent}}</div>
		<div id="rightPanel" class="panel">{{.WgContent}}</div>
	</body>
	</html>`

func DebugV1Handler(w http.ResponseWriter, r *http.Request) {
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

	tmpl := template.Must(template.New("debugPage").Parse(debugPageTemplate))
	data := struct {
		FrrContent string
		WgContent  string
	}{
		FrrContent: frrContent,
		WgContent:  wgContent,
	}

	_ = tmpl.Execute(w, data)
}
