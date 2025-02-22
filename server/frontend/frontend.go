package frontend

import (
	"bytebox/logger"
	"bytebox/server/middleware"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

func RenderTemplate(responseWriter http.ResponseWriter, templateName string, data map[string]interface{}) {
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

	if err = tmpl.ExecuteTemplate(responseWriter, "base.html", data); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		logger := logger.GetLoggerInstance()
		logger.Error(fmt.Sprintf("error executing template: %s", err))
	}
}

func MakeTemplateHandler(templateName string) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		RenderTemplate(responseWriter, templateName, nil)
	}
}

func TemplateRouteRegister(url string, templateName string) {
	http.HandleFunc(url, middleware.LoggingHandlerFuncMiddleware(MakeTemplateHandler(templateName)))
}
