package function

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	log.Printf("r.RemoteAddr: %s", r.RemoteAddr)

	ip, err := getIP(r)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("No valid ip"))
	}
	log.Printf("IP: %s", ip)

	w.WriteHeader(200)
	w.Write([]byte(ip))
}

func getIP(r *http.Request) (string, error) {
	//Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	//Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", fmt.Errorf("No valid ip found")
}
