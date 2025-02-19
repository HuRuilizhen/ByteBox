package backend

import (
	"bytebox/server/midlleware"
	"net/http"
)

func StaticFileRouteRegist(url string, path string) {
	fs := http.FileServer(http.Dir(path))
	loggedFs := middleware.LoggingHandlerMiddleware(fs)
	http.Handle(url, http.StripPrefix(url, loggedFs))
}
