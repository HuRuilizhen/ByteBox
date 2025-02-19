package middleware

import (
	"bytebox/logger"
	"fmt"
	"net/http"
)

func logRequest(request *http.Request) {
	loggerInstence := logger.GetLoggerInstance()
	logMessage := fmt.Sprintf(
		"client request logging %s %s %s %s",
		request.RemoteAddr,
		request.Method,
		request.Proto,
		request.URL.Path,
	)
	loggerInstence.Info(logMessage)
}

func LoggingHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		logRequest(request)
		next.ServeHTTP(responseWriter, request)
	})
}

func LoggingHandlerFuncMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		logRequest(request)
		next(responseWriter, request)
	}
}
