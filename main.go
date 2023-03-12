package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	portToServe, ok := os.LookupEnv("PORT")
	if !ok {
		portToServe = "10059"
	}
	log.Printf("Serving on %s\n", portToServe)
	log.Fatal(http.ListenAndServe(":"+portToServe, http.HandlerFunc(serveIP)))
}

func serveIP(writer http.ResponseWriter, request *http.Request) {
	candidate := grabIP(request)
	if strings.Contains(candidate, ":") {
		candidate, _, _ = net.SplitHostPort(candidate)
	}
	ip := net.ParseIP(candidate)

	if ip == nil {
		writer.WriteHeader(400)
		_, _ = fmt.Fprint(writer, "No ip found")
	} else {
		_, _ = fmt.Fprint(writer, ip.To4().String())
	}
}

var headersInOrder = []string{
	"x-client-ip",
	"cf-connecting-ip",
	"fastly-client-ip",
	"true-client-ip",
	"x-real-ip",
	"x-cluster-client-ip",
	"x-forwarded",
	"forwarded-for",
	"forwarded",
	"x-appengine-user-ip",
	"cf-pseudo-ipv4",
}

func grabIP(r *http.Request) string {
	marshalled, err := json.MarshalIndent(r, " ", "  ")
	if err != nil {
		log.Println("ERROR: failed to marshal request")
	}
	log.Println(string(marshalled))

	if ip := r.Header.Get("x-forwarded-for"); ip != "" {
		candidates := strings.Split(ip, ",")
		if len(candidates) > 0 {
			ip = strings.TrimSpace(candidates[0])
		}
		log.Printf("Found ip %s in header X-Forwarded-For", ip)
		return ip
	}

	for _, header := range headersInOrder {
		if ip := r.Header.Get(header); ip != "" {
			log.Printf("Found ip %s in header %s", ip, header)
			return ip
		}
	}

	log.Printf("Fallback to remote header ip %s", r.RemoteAddr)
	return r.RemoteAddr
}
