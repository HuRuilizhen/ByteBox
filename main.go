package main

import (
	"bytebox/database"
	"bytebox/handler"
	"bytebox/logger"
	"bytebox/server"
	"bytebox/server/backend"
	"bytebox/server/frontend"
)

func main() {
	logger.LoadLoggerConfig()
	database.LoadDatabaseConfig()
	server.LoadServerConfig()
	handler.LoadDatabaseConfig()

	db := database.GetDatabaseInstance()

	backend.StaticFileRouteRegister("/static/", "static")
	backend.ApiRouteRegister("/api/upload", handler.UploadHandler, db)
	backend.ApiRouteRegister("/api/download", handler.DownloadHandler, db)

	frontend.TemplateRouteRegister("/", "upload.html")
	frontend.TemplateRouteRegister("/upload", "upload.html")
	frontend.TemplateRouteRegister("/download", "download.html")

	server.LogServerStartUpInfo()
	serverInstance := server.GetServerInstance()
	serverInstance.ListenAndServe()
}
