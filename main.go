package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
)

func getServicePort() (string, error) {
	servicePort, ok := os.LookupEnv("IP_SERVICE_PORT")
	if !ok {
		flagPort := flag.Int("port", 10059, "Port where application should serve.")
		flag.Parsed()
		servicePort = strconv.Itoa(*flagPort)
	}

	numericPort, _ := strconv.Atoi(servicePort)
	if 0 <= numericPort && numericPort <= 65536 {
		return ":" + servicePort, nil
	}

	return "", errors.New("port must be greater than 0 and less than or equal to 65536. provided: " + servicePort)
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	candidate := findIP(r)
	if strings.Contains(candidate, ":") {
		candidate, _, _ = net.SplitHostPort(candidate)
	}

	ip := net.ParseIP(candidate)

	if ip == nil {
		w.WriteHeader(400)
		_, _ = fmt.Fprint(w, "No ip found")
	} else {
		_, _ = fmt.Fprint(w, ip.To4().String())
	}
}

func signalHandler(server *http.Server, signalHandlingChan chan struct{}) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("HTTP Server Shutdown Error: %v", err)
	}
	close(signalHandlingChan)
}

func main() {
	log.SetFlags(0)

	servicePort, err := getServicePort()
	if err != nil {
		log.Fatal(err)
	}

	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", requestHandler)
	server := &http.Server{Addr: servicePort, Handler: serverMux}

	signalHandlingChan := make(chan struct{})
	go signalHandler(server, signalHandlingChan)

	log.Printf("Serving on %s\n", servicePort)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

	<-signalHandlingChan
	log.Printf("Cya!")
}
