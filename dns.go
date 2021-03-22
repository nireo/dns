package main

import (
	"flag"
	"log"
	"net"
	"sync"

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
	port = flag.Int("port", 1053, "the port to host the service")

	rw sync.RWMutex
)

func main() {
	// take flags and init the cache
	flag.Parse()
	cache = &Cache{
		domains: make(map[string]cacheEntry),
	}

	s := NewService()
	s.Start()
}

// Listen for requests and forward the incoming requst if possible (not fully implemented yet)
func (s *Service) Start() {
	var err error
	s.conn, err = net.ListenUDP("udp", &net.UDPAddr{Port: *port})
	if err != nil {
		log.Fatalf("could not start udp service: %s", err)
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
	domain := pac.message.Questions[0].Name.String()
	if pac.message.Response {
		// set the domain in the cache
		s.cache.Set(domain, pac.message)

		// get the client address
		for index, val := range s.forwards[domain] {
			if val.message.ID == pac.message.ID {
				go s.Send(pac.message, val.addr)

				if len(s.forwards)-1 == index {
					s.forwards[domain] = s.forwards[domain][:len(s.forwards[domain])-1]
				} else {
					s.forwards[domain] = append(s.forwards[domain][:index], s.forwards[domain][index+1:]...)
				}

				break
			}
		}

		return
	}

	log.Printf("%s %s %d", domain, pac.addr.IP.String(), pac.message.ID)

	// check cache
	if msg, ok := s.cache.Get(domain); ok {
		msg.ID = pac.message.ID

		go s.Send(pac.message, pac.addr)
	}

	// if a cache entry doesn't exist, add to the forwarders and forward the request
	rw.Lock()
	s.forwards[domain] = append(s.forwards[domain], Packet{
		addr:    pac.addr,
		message: pac.message,
	})
	rw.Unlock()

	// forward the request
	s.Forward(pac.message)
}

func (s *Service) Send(msg dnsmessage.Message, addr *net.UDPAddr) {
	// pack full message
	pack, err := msg.Pack()
	if err != nil {
		log.Printf("packing dns message failed, ID: %d, err: %s", msg.ID, err)
		return
	}

	if _, err := s.conn.WriteToUDP(pack, addr); err != nil {
		log.Printf("error responding to client: %s", err)
	}
}

func (s *Service) Forward(msg dnsmessage.Message) {
	// pack full message
	pack, err := msg.Pack()
	if err != nil {
		log.Printf("packing dns message failed, ID: %d, err: %s", msg.ID, err)
		return
	}

	resolv := net.UDPAddr{IP: net.IP{114, 114, 114, 114}, Port: 53}
	if _, err := s.conn.WriteToUDP(pack, &resolv); err != nil {
		log.Printf("error responding to client: %s", err)
	}
}
