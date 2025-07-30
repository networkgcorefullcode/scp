package proxy

import (
	"fmt"
	"io"
	"net/http"

	"github.com/networkgcorefullcode/scp/logger"
)

var customTransport = http.DefaultTransport

func init() {
	// Here, you can customize the transport, e.g., set timeouts or enable/disable keep-alive
}

// This function takes an incoming HTTP request, creates a new request with the
// same method, URL, and body, sends the new request using the custom transport,
// and forwards the response back to the client. It also handles copying headers
// between the original and proxy requests and responses.
func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Create a new HTTP request with the same method, URL, and body as the original request
	targetURL := r.URL

	logger.AppLog.Debugf("Handling request for URL: %s", targetURL)
	// Only log and use the path and raw query from the URL
	path := targetURL.Path
	if targetURL.RawQuery != "" {
		path = path + "?" + targetURL.RawQuery
	}
	logger.AppLog.Debugf("Request path: %s", path)

	proxyReq, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
	if err != nil {
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		return
	}

	// Copy the headers from the original request to the proxy request
	for name, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(name, value)
		}
	}

	// Send the proxy request using the custom transport
	resp, err := customTransport.RoundTrip(proxyReq)
	if err != nil {
		http.Error(w, "Error sending proxy request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the headers from the proxy response to the original response
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// Set the status code of the original response to the status code of the proxy response
	w.WriteHeader(resp.StatusCode)

	// Copy the body of the proxy response to the original response
	io.Copy(w, resp.Body)
}

// This code creates a new HTTP server that listens on port 8080 for default and uses the `handleRequest`
// function to handle incoming requests. It also logs any errors
// that occur while starting the server.
func Start_Proxy_Server(httpPort int) {
	// Create a new HTTP server with the handleRequest function as the handler

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: http.HandlerFunc(handleRequest),
	}

	// Start the server and log any errors
	logger.AppLog.Infof("Starting proxy server on :%d", httpPort)
	err := server.ListenAndServe()
	if err != nil {
		logger.AppLog.Error("Error starting proxy server: ", err)
	}
	logger.AppLog.Infof("Proxy server on :%d was started successfully", httpPort)
}
