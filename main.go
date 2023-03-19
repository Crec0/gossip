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
	"syscall"
	"time"
)

func main() {
	log.SetFlags(0)

	servicePort, err := getServicePort()
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{Addr: servicePort, Handler: http.HandlerFunc(requestHandler)}
	go func() {
		log.Printf("Serving on %s\n", servicePort)
		if err = server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe Error: %v", err)
		}
	}()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err = server.Shutdown(ctx); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
	}()

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
	log.Printf("Cya!")
}

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
	ip := net.ParseIP(findIP(r))

	if ip == nil {
		w.WriteHeader(400)
		_, _ = fmt.Fprint(w, "No ip found")
	} else {
		_, _ = fmt.Fprint(w, ip.String())
	}
}
