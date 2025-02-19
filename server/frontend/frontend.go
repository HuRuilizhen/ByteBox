package frontend

import (
	"bytebox/logger"
	"bytebox/server/midlleware"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

func renderTemplate(responseWriter http.ResponseWriter, templateName string) {
	tmpl, err := template.ParseFiles(
		filepath.Join("template", "base.html"),
		filepath.Join("template", templateName),
	)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		loggerInstance := logger.GetLoggerInstance()
		loggerInstance.Error(fmt.Sprintf("error parsing template %s", err))
		return
	}

	if err = tmpl.ExecuteTemplate(responseWriter, "base.html", nil); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		logger := logger.GetLoggerInstance()
		logger.Error(fmt.Sprintf("error executing template: %s", err))
	}
}

func MakeTemplateHandler(templateName string) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		renderTemplate(responseWriter, templateName)
	}
}

func TemplateRouteRegister(url string, templateName string) {
	http.HandleFunc(url, middleware.LoggingHandlerFuncMiddleware(MakeTemplateHandler(templateName)))
}
