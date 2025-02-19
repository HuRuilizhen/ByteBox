package main

import (
	"bytebox/server"
	"bytebox/server/backend"
	"bytebox/server/frontend"
	"os"
)

func main() {
	backend.StaticFileRouteRegist("/static/", "static")

	frontend.TemplateRouteRegister("/", "upload.html")
	frontend.TemplateRouteRegister("/upload", "upload.html")
	frontend.TemplateRouteRegister("/download", "download.html")

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	serverConfigInstance := server.GetServerConfigInstance()
	serverConfigInstance.SetAddrPort(port)

	server.LogServerStartUpInfo()
	serverInstance := server.GetServerInstance()
	serverInstance.ListenAndServe()
}
