// Package mdns provides mDNS service registration for local network discovery.
package mdns

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/mdns"
)

// Server wraps the mDNS server.
type Server struct {
	server *mdns.Server
}

// Start registers the LaNotifica service via mDNS.
// The service will be discoverable as "lanotifica.local" on the local network.
func Start(port string) (*Server, error) {
	portNum, err := parsePort(port)
	if err != nil {
		return nil, fmt.Errorf("parsing port: %w", err)
	}

	info := []string{"LaNotifica notification forwarder"}

	service, err := mdns.NewMDNSService(
		"lanotifica",       // Instance name
		"_lanotifica._tcp", // Service type
		"",                 // Domain (empty = .local)
		"",                 // Host name (empty = use system hostname)
		portNum,            // Port
		nil,                // IPs (nil = all interfaces)
		info,               // TXT records
	)
	if err != nil {
		return nil, fmt.Errorf("creating mDNS service: %w", err)
	}

	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		return nil, fmt.Errorf("starting mDNS server: %w", err)
	}

	log.Printf("mDNS: registered as lanotifica.local:%d", portNum)

	return &Server{server: server}, nil
}

// Stop shuts down the mDNS server.
func (s *Server) Stop() error {
	if s.server != nil {
		return s.server.Shutdown()
	}
	return nil
}

func parsePort(port string) (int, error) {
	p := strings.TrimPrefix(port, ":")
	return strconv.Atoi(p)
}
