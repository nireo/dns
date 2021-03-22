package main

import (
	"flag"
	"log"
	"net"

	"golang.org/x/net/dns/dnsmessage"
)

// Service represents the dns server
type Service struct {
	forwards map[string][]Packet
	conn     *net.UDPConn
	cache    *Cache
}

// NewService creates a service instance appending the cache to it.
func NewService() *Service {
	return &Service{
		cache:    cache,
		forwards: make(map[string][]Packet),
	}
}

// Packet represents a single DNS packet
type Packet struct {
	message dnsmessage.Message
	addr    *net.UDPAddr
}

// Configuration options
var (
	ttl  = flag.Int64("ttl", 30, "the time to stay alive for domains")
	port = flag.Int("port", 8080, "the port to host the server on")
)

func main() {
	// take flags and init the cache
	flag.Parse()
	cache = &Cache{
		domains: make(map[string]cacheEntry),
	}
}

// Listen for requests and forward the incoming requst if possible (not fully implemented yet)
func (s *Service) Start() {
	var err error
	s.conn, err = net.ListenUDP("upd", &net.UDPAddr{Port: *port})
	if err != nil {
		log.Fatalf("could not start udp service")
	}
	defer s.conn.Close()

	log.Printf("running dns service on port: %d", *port)
	for {
		buffer := make([]byte, 512)
		_, addr, err := s.conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("erro reading udp connection: %s", err)
			// even though a error happened still listen for more connections
			continue
		}

		// parse the message from the buffer
		var msg dnsmessage.Message
		if err := msg.Unpack(buffer); err != nil {
			log.Printf("failed parings dns: %s", err)
			continue
		}

		if len(msg.Questions) == 0 {
			continue
		}

		go s.Query(Packet{addr: addr, message: msg})
	}
}

func (s *Service) Query(pac Packet) {
	// TODO: implement this
	// domain := pac.message.Questions[0].Name.String()
}

func (s *Service) Send(msg dnsmessage.Message, addr *net.UDPAddr) {
	// TODO
}

func (s *Service) Forward(msg dnsmessage.Message) {
	// TODO
}
