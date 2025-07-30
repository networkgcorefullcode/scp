package proxy

import (
	"fmt"
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/networkgcorefullcode/scp/logger"
)

// This code creates a new HTTP server that listens on port 8080 for default and uses the `handleRequest`
// function to handle incoming requests. It also logs any errors
// that occur while starting the server.
func Start_Proxy_Server(httpPort int) {
	// Create a new HTTP server with the handleRequest function as the handler
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	logger.AppLog.Infof("Starting proxy server on :%d", httpPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), proxy)
	if err != nil {
		logger.AppLog.Error("Error starting proxy server: ", err)
	}
	logger.AppLog.Infof("Proxy server on :%d stopped", httpPort)
}
