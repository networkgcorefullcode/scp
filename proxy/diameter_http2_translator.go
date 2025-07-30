package proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/dict"
)

// StartDiameterTranslator initializes the Diameter to HTTP/2 translator
func StartDiameterTranslator() {
	// Start Diameter server
	go startDiameterServer()

	// Start HTTP/2 client
	fmt.Println("Diameter to HTTP/2 Translator is running...")
}

func startDiameterServer() {
	mux := diam.NewServeMux()
	mux.Handle("CER", handleCER())
	mux.Handle("AAR", handleAAR())

	// Start the Diameter server
	go func() {
		addr := ":3868"
		fmt.Printf("Starting Diameter server on %s\n", addr)
		if err := diam.ListenAndServe(addr, mux, dict.Default); err != nil {
			fmt.Printf("Failed to start Diameter server: %v\n", err)
		}
	}()
}

func handleCER() diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		fmt.Println("Received CER message")
		// Respond with CEA (Capabilities-Exchange-Answer)
		response := m.Answer(diam.Success)
		response.NewAVP(avp.OriginHost, diam.Mbit, 0, datatype.DiameterIdentity("example.com"))
		response.NewAVP(avp.OriginRealm, diam.Mbit, 0, datatype.DiameterIdentity("example.com"))
		response.NewAVP(diam.OriginRealm, diam.Mbit, 0, datatype.DiameterIdentity("example.com"))
		c.WriteMessage(response)
	}
}

func handleAAR() diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		fmt.Println("Received AAR message")
		// Extract AVPs from the Diameter message
		var sessionID datatype.UTF8String
		if err := m.Unmarshal(&sessionID); err != nil {
			fmt.Printf("Failed to unmarshal AAR: %v\n", err)
			return
		}

		// Translate to HTTP/2 request
		url := "https://5g-service.example.com/api"
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			fmt.Printf("Failed to create HTTP/2 request: %v\n", err)
			return
		}
		req.Header.Set("Session-ID", string(sessionID))

		// Send HTTP/2 request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Failed to send HTTP/2 request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		// Respond back to Diameter client
		response.NewAVP(avp.SessionID, diam.Mbit, 0, sessionID)
		response.NewAVP(diam.SessionID, diam.Mbit, 0, sessionID)
		c.WriteMessage(response)
	}
}

// StartHTTPToDiameterProxy inicia el proxy HTTP a Diameter
func StartHTTPToDiameterProxy() {
	// Inicia el servidor HTTP
	http.HandleFunc("/proxy", handleHTTPToDiameter)
	fmt.Println("HTTP to Diameter Proxy is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func handleHTTPToDiameter(w http.ResponseWriter, r *http.Request) {
	// Leer el cuerpo de la solicitud HTTP
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Crear un mensaje Diameter (AAR como ejemplo)
	m := diam.NewRequest(265, 0, dict.Default) // 265 es el código para AAR
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String("session-12345"))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("example.com"))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("example.com"))
	m.NewAVP(avp.DestinationRealm, avp.Mbit, 0, datatype.DiameterIdentity("4g-core.com"))
	m.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String(string(body)))

	// Enviar el mensaje Diameter
	conn, err := diam.Dial("tcp", "4g-core.example.com:3868", nil, dict.Default)
	if err != nil {
		http.Error(w, "Failed to connect to Diameter server", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	if _, err := m.WriteTo(conn); err != nil {
		http.Error(w, "Failed to send Diameter message", http.StatusInternalServerError)
		return
	}

	// Leer la respuesta Diameter
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		http.Error(w, "Failed to read Diameter response", http.StatusInternalServerError)
		return
	}

	// Enviar la respuesta al cliente HTTP
	w.WriteHeader(http.StatusOK)
	w.Write(buf[:n])
}
