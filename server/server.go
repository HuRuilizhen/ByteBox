package server

import (
	"bytebox/logger"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	serverInstance     *http.Server
	serverInstanceOnce sync.Once
)

type ServerConfig struct {
	addrIp        string
	addrPort      string
	maxHeaderByte int
	readTimeout   time.Duration
	writeTimout   time.Duration
}

var (
	serverConfigInstance     ServerConfig
	serverConfigInstanceOnce sync.Once
)

func GetServerConfigInstance() ServerConfig {
	serverConfigInstanceOnce.Do(func() {
		serverConfigInstance.addrIp = "0.0.0.0"
		serverConfigInstance.addrPort = "4000"
		serverConfigInstance.maxHeaderByte = 1 << 20
		serverConfigInstance.readTimeout = 10 * time.Second
		serverConfigInstance.writeTimout = 10 * time.Second
	})
	return serverConfigInstance
}

func GetServerInstance() *http.Server {
	serverInstanceOnce.Do(func() {
		serverConfig := GetServerConfigInstance()
		serverInstance = &http.Server{
			Addr:           serverConfigInstance.addrIp + ":" + serverConfigInstance.addrPort,
			MaxHeaderBytes: serverConfig.maxHeaderByte,
			ReadTimeout:    serverConfig.readTimeout,
			WriteTimeout:   serverConfig.writeTimout,
		}
	})
	return serverInstance
}

/* set and get method of server config */

func (serverConfig *ServerConfig) SetAddrIp(addrIp string) {
	serverConfig.addrIp = addrIp
}

func (serverConfig *ServerConfig) GetAddrIp() string {
	return serverConfig.addrIp
}

func (serverConfig *ServerConfig) SetAddrPort(addrPort string) {
	serverConfig.addrPort = addrPort
}

func (serverConfig *ServerConfig) GetAddrPort() string {
	return serverConfig.addrPort
}

/* server start up info log */

func getLocalIPs() ([]string, error) {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	if len(ips) == 0 {
		return ips, fmt.Errorf("cannot find any local IP addresses")
	}
	return ips, nil
}

func LogServerStartUpInfo() {
	protocol := "http://"
	localhostIp := "127.0.0.1"
	localAreaNetIps, _ := getLocalIPs()
	port := serverConfigInstance.GetAddrPort()

	loggerInstance := logger.GetLoggerInstance()
	loggerInstance.Info("server starting...")
	loggerInstance.Info(fmt.Sprintf("server listen at %s%s:%s", protocol, localhostIp, port))
	for _, localAreaNetIp := range localAreaNetIps {
		loggerInstance.Info(fmt.Sprintf("server listen at %s%s:%s", protocol, localAreaNetIp, port))
	}
}
