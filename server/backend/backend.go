package backend

import (
	"bytebox/server/middleware"
	"net/http"

	"gorm.io/gorm"
)

func StaticFileRouteRegister(url string, path string) {
	fs := http.FileServer(http.Dir(path))
	loggedFs := middleware.LoggingHandlerMiddleware(fs)
	http.Handle(url, http.StripPrefix(url, loggedFs))
}

func ApiRouteRegister(url string, apiHandlerFunc func(http.ResponseWriter, *http.Request, *gorm.DB), db *gorm.DB) {
	loggedHandlerFunc := middleware.LoggingHandlerFuncMiddleware(
		func(responseWriter http.ResponseWriter, request *http.Request) {
			apiHandlerFunc(responseWriter, request, db)
		},
	)
	http.HandleFunc(url, loggedHandlerFunc)
}
